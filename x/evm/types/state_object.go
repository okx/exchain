package types

import (
	"bytes"
	"fmt"
	types2 "github.com/ethereum/go-ethereum/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"io"
	"math/big"
	"sync"

	"github.com/VictoriaMetrics/fastcache"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

const keccak256HashSize = 100000

var (
	_ StateObject = (*stateObject)(nil)

	emptyCodeHash          = ethcrypto.Keccak256(nil)
	keccak256HashCache, _  = lru.NewARC(keccak256HashSize)
	keccak256HashFastCache = fastcache.New(128 * keccak256HashSize) // 32 + 20 + 32

	keccakStatePool = &sync.Pool{
		New: func() interface{} {
			return ethcrypto.NewKeccakState()
		},
	}
)

func keccak256HashWithSyncPool(data ...[]byte) (h ethcmn.Hash) {
	d := keccakStatePool.Get().(ethcrypto.KeccakState)
	defer keccakStatePool.Put(d)
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	return h
}

func keccak256HashWithLruCache(compositeKey []byte) ethcmn.Hash {
	if value, ok := keccak256HashCache.Get(string(compositeKey)); ok {
		return value.(ethcmn.Hash)
	}
	value := keccak256HashWithSyncPool(compositeKey)
	keccak256HashCache.Add(string(compositeKey), value)
	return value
}

func keccak256HashWithFastCache(compositeKey []byte) (hash ethcmn.Hash) {
	if _, ok := keccak256HashFastCache.HasGet(hash[:0], compositeKey); ok {
		return
	}
	hash = keccak256HashWithSyncPool(compositeKey)
	keccak256HashFastCache.Set(compositeKey, hash[:])
	return
}

func Keccak256HashWithCache(compositeKey []byte) ethcmn.Hash {
	// if length of compositeKey + hash size is greater than 128, use lru cache
	if len(compositeKey) > 128-ethcmn.HashLength {
		return keccak256HashWithLruCache(compositeKey)
	} else {
		return keccak256HashWithFastCache(compositeKey)
	}
}

// StateObject interface for interacting with state object
type StateObject interface {
	GetCommittedState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash
	GetState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash
	SetState(db ethstate.Database, key, value ethcmn.Hash)

	Code(db ethstate.Database) []byte
	SetCode(codeHash ethcmn.Hash, code []byte)
	CodeHash() []byte

	AddBalance(amount *big.Int)
	SubBalance(amount *big.Int)
	SetBalance(amount *big.Int)

	Balance() *big.Int
	ReturnGas(gas *big.Int)
	Address() ethcmn.Address

	SetNonce(nonce uint64)
	Nonce() uint64
}

// stateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call CommitTrie to write the modified storage trie into a database.
type stateObject struct {
	trie      ethstate.Trie // storage trie, which becomes non-nil on first access
	stateRoot ethcmn.Hash   // merkle root of the storage trie

	code types.Code // contract bytecode, which gets set when code is loaded
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	originStorage  ethstate.Storage // Storage cache of original entries to dedup rewrites
	dirtyStorage   ethstate.Storage // Storage entries that need to be flushed to disk
	pendingStorage ethstate.Storage // Storage entries that need to be flushed to disk, at the end of an entire block
	fakeStorage    ethstate.Storage // Fake storage which constructed by caller for debugging purpose.

	// DB error
	dbErr   error
	stateDB *CommitStateDB
	account *types.EthAccount

	address  ethcmn.Address
	addrHash ethcmn.Hash

	// cache flags
	//
	// When an object is marked suicided it will be delete from the trie during
	// the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

func newStateObject(db *CommitStateDB, accProto authexported.Account, stateRoot ethcmn.Hash) *stateObject {
	ethermintAccount, ok := accProto.(*types.EthAccount)
	if !ok {
		panic(fmt.Sprintf("invalid account type for state object: %T", accProto))
	}

	// set empty code hash
	if ethermintAccount.CodeHash == nil {
		ethermintAccount.CodeHash = emptyCodeHash
	}
	if stateRoot == (ethcmn.Hash{}) {
		stateRoot = types2.EmptyRootHash
	}

	ethAddr := ethermintAccount.EthAddress()
	return &stateObject{
		stateDB:        db,
		stateRoot:      stateRoot,
		account:        ethermintAccount,
		address:        ethAddr,
		addrHash:       ethcrypto.Keccak256Hash(ethAddr[:]),
		originStorage:  make(ethstate.Storage),
		pendingStorage: make(ethstate.Storage),
		dirtyStorage:   make(ethstate.Storage),
	}
}

// ----------------------------------------------------------------------------
// Setters
// ----------------------------------------------------------------------------

// SetState updates a value in account storage. Note, the key will be prefixed
// with the address of the state object.
func (so *stateObject) SetState(db ethstate.Database, key, value ethcmn.Hash) {
	// If the fake storage is set, put the temporary state update here.
	if so.fakeStorage != nil {
		so.fakeStorage[key] = value
		return
	}
	// If the new value is the same as old, don't set
	prev := so.GetState(db, key)
	if prev == value {
		return
	}

	// New value is different, update and journal the change
	so.stateDB.journal.append(storageChange{
		account:   &so.address,
		key:       key,
		prevValue: prev,
	})
	so.setState(key, value)
}

// setState sets a state with a prefixed key and value to the dirty storage.
func (so *stateObject) setState(key, value ethcmn.Hash) {
	so.dirtyStorage[key] = value
}

// SetCode sets the state object's code.
func (so *stateObject) SetCode(codeHash ethcmn.Hash, code []byte) {
	prevCode := so.Code(so.stateDB.db)
	so.stateDB.journal.append(codeChange{
		account:  &so.address,
		prevHash: so.CodeHash(),
		prevCode: prevCode,
	})
	so.setCode(codeHash, code)
}

func (so *stateObject) setCode(codeHash ethcmn.Hash, code []byte) {
	so.code = code
	so.account.CodeHash = codeHash.Bytes()
	so.dirtyCode = true
}

// AddBalance adds an amount to a state object's balance. It is used to add
// funds to the destination account of a transfer.
func (so *stateObject) AddBalance(amount *big.Int) {
	amt := sdk.NewDecFromBigIntWithPrec(amount, sdk.Precision) // int2dec
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.

	// NOTE: this will panic if amount is nil
	if amt.IsZero() {
		if so.empty() {
			so.touch()
		}
		return
	}

	newBalance := so.account.GetCoins().AmountOf(sdk.DefaultBondDenom).Add(amt)
	so.SetBalance(newBalance.BigInt())
}

// SubBalance removes an amount from the stateObject's balance. It is used to
// remove funds from the origin account of a transfer.
func (so *stateObject) SubBalance(amount *big.Int) {
	amt := sdk.NewDecFromBigIntWithPrec(amount, sdk.Precision) // int2dec
	if amt.IsZero() {
		return
	}
	newBalance := so.account.GetCoins().AmountOf(sdk.DefaultBondDenom).Sub(amt)
	so.SetBalance(newBalance.BigInt())
}

// SetBalance sets the state object's balance.
func (so *stateObject) SetBalance(amount *big.Int) {
	amt := sdk.NewDecFromBigIntWithPrec(amount, sdk.Precision) // int2dec

	so.stateDB.journal.append(balanceChange{
		account: &so.address,
		prev:    so.account.GetCoins().AmountOf(sdk.DefaultBondDenom), // int2dec
	})

	so.setBalance(sdk.DefaultBondDenom, amt)
}

func (so *stateObject) setBalance(denom string, amount sdk.Dec) {
	so.account.SetBalance(denom, amount)
}

// SetNonce sets the state object's nonce (i.e sequence number of the account).
func (so *stateObject) SetNonce(nonce uint64) {
	so.stateDB.journal.append(nonceChange{
		account: &so.address,
		prev:    so.account.Sequence,
	})

	so.setNonce(nonce)
}

func (so *stateObject) setNonce(nonce uint64) {
	if so.account == nil {
		panic("state object account is empty")
	}
	so.account.Sequence = nonce
}

// setError remembers the first non-nil error it is called with.
func (so *stateObject) setError(err error) {
	if so.dbErr == nil {
		so.dbErr = err
	}
}

func (so *stateObject) markSuicided() {
	so.suicided = true
}

// commitState commits all dirty storage to a KVStore and resets
// the dirty storage slice to the empty state.
func (so *stateObject) commitState() {
	// Make sure all dirty slots are finalized into the pending storage area
	so.finalise(false) // Don't prefetch any more, pull directly if need be
	if len(so.pendingStorage) == 0 {
		return
	}

	ctx := so.stateDB.ctx
	store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), AddressStoragePrefix(so.Address()))
	for key, value := range so.pendingStorage {
		// Skip noop changes, persist actual changes
		if value == so.originStorage[key] {
			continue
		}
		so.originStorage[key] = value

		prefixKey := so.GetStorageByAddressKey(key.Bytes())
		if (value == ethcmn.Hash{}) {
			store.Delete(prefixKey.Bytes())
			so.stateDB.ctx.Cache().UpdateStorage(so.address, prefixKey, value.Bytes(), true)
			if !so.stateDB.ctx.IsCheckTx() {
				if so.stateDB.Watcher.Enabled() {
					so.stateDB.Watcher.SaveState(so.Address(), prefixKey.Bytes(), ethcmn.Hash{}.Bytes())
				}
			}
		} else {
			store.Set(prefixKey.Bytes(), value.Bytes())
			so.stateDB.ctx.Cache().UpdateStorage(so.address, prefixKey, value.Bytes(), true)
			if !so.stateDB.ctx.IsCheckTx() {
				if so.stateDB.Watcher.Enabled() {
					so.stateDB.Watcher.SaveState(so.Address(), prefixKey.Bytes(), value.Bytes())
				}
			}
		}
	}

	if len(so.pendingStorage) > 0 {
		so.pendingStorage = make(ethstate.Storage)
	}

	return
}

// commitCode persists the state object's code to the KVStore.
func (so *stateObject) commitCode() {
	ctx := so.stateDB.ctx
	store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), KeyPrefixCode)
	store.Set(so.CodeHash(), so.code)
	ctx.Cache().UpdateCode(so.CodeHash(), so.code, true)
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

// Address returns the address of the state object.
func (so stateObject) Address() ethcmn.Address {
	return so.address
}

// Balance returns the state object's current balance.
func (so *stateObject) Balance() *big.Int {
	balance := so.account.Balance(sdk.DefaultBondDenom).BigInt()
	if balance == nil {
		return zeroBalance
	}
	return balance
}

// CodeHash returns the state object's code hash.
func (so *stateObject) CodeHash() []byte {
	if so.account == nil || len(so.account.CodeHash) == 0 {
		return emptyCodeHash
	}
	return so.account.CodeHash
}

// Nonce returns the state object's current nonce (sequence number).
func (so *stateObject) Nonce() uint64 {
	if so.account == nil {
		return 0
	}
	return so.account.Sequence
}

// Code returns the contract code associated with this object, if any.
func (so *stateObject) Code(db ethstate.Database) []byte {
	if !tmtypes.HigherThanMars(so.stateDB.ctx.BlockHeight()) {
		if len(so.code) > 0 {
			return so.code
		}

		if bytes.Equal(so.CodeHash(), emptyCodeHash) {
			return nil
		}

		code := make([]byte, 0)
		ctx := so.stateDB.ctx
		if data, ok := ctx.Cache().GetCode(so.CodeHash()); ok {
			code = data
		} else {
			store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), KeyPrefixCode)
			code = store.Get(so.CodeHash())
			ctx.Cache().UpdateCode(so.CodeHash(), code, false)
		}

		if len(code) == 0 {
			so.setError(fmt.Errorf("failed to get code hash %x for address %s", so.CodeHash(), so.Address().String()))
		} else {
			so.code = code
		}

		return code
	} else {
		return so.CodeMpt(db)
	}
}

// GetState retrieves a value from the account storage trie. Note, the key will
// be prefixed with the address of the state object.
func (so *stateObject) GetState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash {
	// If the fake storage is set, only lookup the state here(in the debugging mode)
	if so.fakeStorage != nil {
		return so.fakeStorage[key]
	}
	// if we have a dirty value for this state entry, return it
	value, dirty := so.dirtyStorage[key]
	if dirty {
		return value
	}

	// otherwise return the entry's original value
	return so.GetCommittedState(db, key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
//
// NOTE: the key will be prefixed with the address of the state object.
func (so *stateObject) GetCommittedState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash {
	if !tmtypes.HigherThanMars(so.stateDB.ctx.BlockHeight()) {
		// If the fake storage is set, only lookup the state here(in the debugging mode)
		if so.fakeStorage != nil {
			return so.fakeStorage[key]
		}

		// If we have a pending write or clean cached, return that
		if value, pending := so.pendingStorage[key]; pending {
			return value
		}
		if value, cached := so.originStorage[key]; cached {
			return value
		}

		// otherwise load the value from the KVStore
		state := NewState(key, ethcmn.Hash{})

		ctx := so.stateDB.ctx
		rawValue := make([]byte, 0)
		var ok bool

		prefixKey := so.GetStorageByAddressKey(key.Bytes())
		rawValue, ok = ctx.Cache().GetStorage(so.address, prefixKey)
		if !ok {
			store := so.stateDB.dbAdapter.NewStore(ctx.KVStore(so.stateDB.storeKey), AddressStoragePrefix(so.Address()))
			rawValue = store.Get(prefixKey.Bytes())
			ctx.Cache().UpdateStorage(so.address, prefixKey, rawValue, false)
		}

		if len(rawValue) > 0 {
			state.Value.SetBytes(rawValue)
		}

		so.originStorage[key] = state.Value
		return state.Value
	} else {
		return so.GetCommittedStateMpt(db, key)
	}
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

// ReturnGas returns the gas back to the origin. Used by the Virtual machine or
// Closures. It performs a no-op.
func (so *stateObject) ReturnGas(gas *big.Int) {}

func (so *stateObject) deepCopy(db *CommitStateDB) *stateObject {
	if !tmtypes.HigherThanMars(so.stateDB.ctx.BlockHeight()) {
		newAccount := types.ProtoAccount().(*types.EthAccount)
		jsonAccount, err := so.account.MarshalJSON()
		if err != nil {
			return nil
		}
		err = newAccount.UnmarshalJSON(jsonAccount)
		if err != nil {
			return nil
		}
		newStateObj := newStateObject(db, newAccount, so.stateRoot)

		newStateObj.code = make(types.Code, len(so.code))
		copy(newStateObj.code, so.code)
		newStateObj.dirtyStorage = so.dirtyStorage.Copy()
		newStateObj.originStorage = so.originStorage.Copy()
		newStateObj.suicided = so.suicided
		newStateObj.dirtyCode = so.dirtyCode
		newStateObj.deleted = so.deleted

		return newStateObj
	} else {
		return so.deepCopyMpt(db)
	}
}

// empty returns whether the account is considered empty.
func (so *stateObject) empty() bool {
	balace := so.account.Balance(sdk.DefaultBondDenom)
	return so.account == nil ||
		(so.account != nil &&
			so.account.Sequence == 0 &&
			(balace.BigInt() == nil || balace.IsZero()) &&
			bytes.Equal(so.account.CodeHash, emptyCodeHash))
}

// EncodeRLP implements rlp.Encoder.
func (so *stateObject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, so.account)
}

func (so *stateObject) touch() {
	so.stateDB.journal.append(touchChange{
		account: &so.address,
	})

	if so.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		so.stateDB.journal.dirty(so.address)
	}
}

// GetStorageByAddressKey returns a hash of the composite key for a state
// object's storage prefixed with it's address.
func (so stateObject) GetStorageByAddressKey(key []byte) ethcmn.Hash {
	prefix := so.Address().Bytes()
	compositeKey := make([]byte, len(prefix)+len(key))

	copy(compositeKey, prefix)
	copy(compositeKey[len(prefix):], key)
	return Keccak256HashWithCache(compositeKey)
}

// stateEntry represents a single key value pair from the StateDB's stateObject mappindg.
// This is to prevent non determinism at genesis initialization or export.
type stateEntry struct {
	// address key of the state object
	address     ethcmn.Address
	stateObject *stateObject
}
