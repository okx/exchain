package types

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/VictoriaMetrics/fastcache"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/okx/okbchain/app/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
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
	cacheKey := string(compositeKey)
	if value, ok := keccak256HashCache.Get(cacheKey); ok {
		return value.(ethcmn.Hash)
	}
	value := keccak256HashWithSyncPool(compositeKey)
	keccak256HashCache.Add(cacheKey, value)
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

// Keccak256HashWithCache returns the Keccak256 hash of the given data.
// this function should not keep the reference of the input data after return.
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

	SetStorage(storage map[ethcmn.Hash]ethcmn.Hash)
}

// stateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call CommitTrie to write the modified storage trie into a database.
type stateObject struct {
	trie ethstate.Trie // storage trie, which becomes non-nil on first access

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

func newStateObject(db *CommitStateDB, accProto authexported.Account) *stateObject {
	ethermintAccount, ok := accProto.(*types.EthAccount)
	if !ok {
		panic(fmt.Sprintf("invalid account type for state object: %T", accProto))
	}

	// set empty code hash
	if ethermintAccount.CodeHash == nil {
		ethermintAccount.CodeHash = emptyCodeHash
	}
	if ethermintAccount.StateRoot == (ethcmn.Hash{}) {
		ethermintAccount.StateRoot = ethtypes.EmptyRootHash
	}

	ethAddr := ethermintAccount.EthAddress()
	return &stateObject{
		stateDB:        db,
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
	if err != nil {
		so.stateDB.Logger().Debug("stateObject", "error", err)
	}
	if so.dbErr == nil {
		so.dbErr = err
	}
}

func (so *stateObject) markSuicided() {
	so.suicided = true
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

// Address returns the address of the state object.
func (so *stateObject) Address() ethcmn.Address {
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
	return so.CodeInRawDB(db)
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
	return so.GetCommittedStateMpt(db, key)
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

// ReturnGas returns the gas back to the origin. Used by the Virtual machine or
// Closures. It performs a no-op.
func (so *stateObject) ReturnGas(gas *big.Int) {}

func (so *stateObject) deepCopy(db *CommitStateDB) *stateObject {
	return so.deepCopyMpt(db)
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

// stateEntry represents a single key value pair from the StateDB's stateObject mappindg.
// This is to prevent non determinism at genesis initialization or export.
type stateEntry struct {
	// address key of the state object
	address     ethcmn.Address
	stateObject *stateObject
}
