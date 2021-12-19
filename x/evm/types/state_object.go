package types

import (
	"bytes"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	types2 "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"io"
	"math/big"
)

var (
	_ StateObject = (*stateObject)(nil)

	emptyCodeHash = ethcrypto.Keccak256(nil)
)

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
	address   ethcmn.Address
	addrHash  ethcmn.Hash
	stateDB   *CommitStateDB
	account   *types.EthAccount

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	trie ethstate.Trie // storage trie, which becomes non-nil on first access
	code types.Code    // contract bytecode, which gets set when code is loaded

	//// State objects are used by the consensus core and VM which are
	//// unable to deal with database-level errors. Any error that occurs
	//// during a database read is memoized here and will eventually be returned
	//// by StateDB.Commit.
	originStorage  ethstate.Storage // Storage cache of original entries to dedup rewrites, reset for every transaction
	pendingStorage ethstate.Storage // Storage entries that need to be flushed to disk, at the end of an entire block
	dirtyStorage   ethstate.Storage // Storage entries that have been modified in the current transaction execution
	fakeStorage    ethstate.Storage // Fake storage which constructed by caller for debugging purpose.

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

func newStateObject(db *CommitStateDB, accProto authexported.Account) *stateObject {
	// func newStateObject(db *CommitStateDB, accProto authexported.Account, balance sdk.Int) *stateObject {
	ethermintAccount, ok := accProto.(*types.EthAccount)
	if !ok {
		panic(fmt.Sprintf("invalid account type for state object: %T", accProto))
	}

	// set empty code hash
	if ethermintAccount.CodeHash == nil {
		ethermintAccount.CodeHash = emptyCodeHash
	}
	if ethermintAccount.StateRoot == (ethcmn.Hash{}) {
		ethermintAccount.StateRoot = types2.EmptyRootHash
	}

	ethAddr := ethermintAccount.EthAddress()
	return &stateObject{
		stateDB:        db,
		account:        ethermintAccount.Copy(),
		address:        ethAddr,
		addrHash:       ethcrypto.Keccak256Hash(ethAddr[:]),
		originStorage:  make(ethstate.Storage),
		pendingStorage: make(ethstate.Storage),
		dirtyStorage:   make(ethstate.Storage),
	}
}

func (s *stateObject) getTrie(db ethstate.Database) ethstate.Trie {
	if s.trie == nil {
		var err error
		s.trie, err = db.OpenStorageTrie(s.addrHash, s.account.StateRoot)
		if err != nil {
			s.trie, _ = db.OpenStorageTrie(s.addrHash, ethcmn.Hash{})
			s.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return s.trie
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
		account:  &so.address,
		key:      key,
		prevalue: prev,
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
		prevhash: so.CodeHash(),
		prevcode: prevCode,
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

// UpdateRoot sets the trie root to the current root hash of
func (s *stateObject) updateRoot(db ethstate.Database) {
	// If nothing changed, don't bother with hashing anything
	if s.updateTrie(db) == nil {
		return
	}
	s.account.StateRoot = s.trie.Hash()
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
	if so.code != nil {
		return so.code
	}
	if bytes.Equal(so.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(so.addrHash, ethcmn.BytesToHash(so.CodeHash()))
	if err != nil {
		so.setError(fmt.Errorf("can't load code hash %x: %v", so.CodeHash(), err))
	}
	so.code = code
	return code
}

// GetState retrieves a value from the account storage trie. Note, the key will
// be prefixed with the address of the state object.
func (so *stateObject) GetState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash {
	// If the fake storage is set, only lookup the state here(in the debugging mode)
	if so.fakeStorage != nil {
		return so.fakeStorage[key]
	}
	// If we have a dirty value for this state entry, return it
	value, dirty := so.dirtyStorage[key]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return so.GetCommittedState(db, key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
//
// NOTE: the key will be prefixed with the address of the state object.
func (so *stateObject) GetCommittedState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash {
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

	var (
		enc []byte
		err error
	)

	if enc, err = so.getTrie(db).TryGet(key.Bytes()); err != nil {
		so.setError(err)
		return ethcmn.Hash{}
	}

	var value ethcmn.Hash
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			so.setError(err)
		}
		value.SetBytes(content)
	}
	so.originStorage[key] = value
	return value
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

// ReturnGas returns the gas back to the origin. Used by the Virtual machine or
// Closures. It performs a no-op.
func (so *stateObject) ReturnGas(gas *big.Int) {}

func (so *stateObject) deepCopy(db *CommitStateDB) *stateObject {
	acc := db.accountKeeper.NewAccountWithAddress(db.ctx, so.account.Address)
	newStateObj := newStateObject(db, acc)
	if so.trie != nil {
		newStateObj.trie = db.db.CopyTrie(so.trie)
	}

	newStateObj.code = so.code
	newStateObj.dirtyStorage = so.dirtyStorage.Copy()
	newStateObj.originStorage = so.originStorage.Copy()
	newStateObj.pendingStorage = so.pendingStorage.Copy()
	newStateObj.suicided = so.suicided
	newStateObj.dirtyCode = so.dirtyCode
	newStateObj.deleted = so.deleted

	return newStateObj
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

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (s *stateObject) CommitTrie(db ethstate.Database) error {
	// If nothing changed, don't bother with hashing anything
	if s.updateTrie(db) == nil {
		return nil
	}
	if s.dbErr != nil {
		return s.dbErr
	}

	root, err := s.trie.Commit(nil)
	if err == nil {
		s.account.StateRoot = root
	}
	return err
}

// updateTrie writes cached storage modifications into the object's storage trie.
// It will return nil if the trie has not been loaded and no changes have been made
func (s *stateObject) updateTrie(db ethstate.Database) ethstate.Trie {
	// Make sure all dirty slots are finalized into the pending storage area
	s.finalise(false) // Don't prefetch any more, pull directly if need be
	if len(s.pendingStorage) == 0 {
		return s.trie
	}

	// Insert all the pending updates into the trie
	tr := s.getTrie(db)
	for key, value := range s.pendingStorage {
		// Skip noop changes, persist actual changes
		if value == s.originStorage[key] {
			continue
		}
		s.originStorage[key] = value

		var v []byte
		if (value == ethcmn.Hash{}) {
			s.setError(tr.TryDelete(key[:]))
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ = rlp.EncodeToBytes(ethcmn.TrimLeftZeroes(value[:]))
			s.setError(tr.TryUpdate(key[:], v))
		}
	}

	if len(s.pendingStorage) > 0 {
		s.pendingStorage = make(ethstate.Storage)
	}
	return tr
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (s *stateObject) finalise(prefetch bool) {
	for key, value := range s.dirtyStorage {
		s.pendingStorage[key] = value
	}

	if len(s.dirtyStorage) > 0 {
		s.dirtyStorage = make(ethstate.Storage)
	}
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (s *stateObject) CodeSize(db ethstate.Database) int {
	if s.code != nil {
		return len(s.code)
	}
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return 0
	}
	size, err := db.ContractCodeSize(s.addrHash, ethcmn.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code size %x: %v", s.CodeHash(), err))
	}
	return size
}

// SetStorage replaces the entire state storage with the given one.
//
// After this function is called, all original state will be ignored and state
// lookup only happens in the fake state storage.
//
// Note this function should only be used for debugging purpose.
func (s *stateObject) SetStorage(storage map[ethcmn.Hash]ethcmn.Hash) {
	// Allocate fake storage if it's nil.
	if s.fakeStorage == nil {
		s.fakeStorage = make(ethstate.Storage)
	}
	for key, value := range storage {
		s.fakeStorage[key] = value
	}
	// Don't bother journal since this function should only be used for
	// debugging and the `fake` storage won't be committed to database.
}

func (s *stateObject) UpdateAccInfo() {
	accProto := s.stateDB.accountKeeper.GetAccount(s.stateDB.ctx, s.account.Address)
	if accProto != nil {
		ethAccount, ok := accProto.(*types.EthAccount)
		if !ok {
			return
		}

		// only need to update these field
		s.account.Coins = ethAccount.Coins
		s.account.Sequence = ethAccount.Sequence
	}
}
