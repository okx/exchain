package types

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	types2 "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/x/common/analyzer"
	"github.com/pkg/errors"
	"math/big"
	"sort"
)

var (
	_ ethvm.StateDB = (*CommitStateDB)(nil)

	zeroBalance = sdk.ZeroInt().BigInt()
)

type revision struct {
	id           int
	journalIndex int
}

type CommitStateDBParams struct {
	StoreKey      sdk.StoreKey
	ParamSpace    Subspace
	AccountKeeper AccountKeeper
	SupplyKeeper  SupplyKeeper
	Watcher       Watcher
	BankKeeper    BankKeeper
	Ada           DbAdapter
}

type CacheCode struct {
	CodeHash []byte
	Code     []byte
}

// CommitStateDB implements the Geth state.StateDB interface. Instead of using
// a trie and database for querying and persistence, the Keeper uses KVStores
// and an account mapper is used to facilitate state transitions.
//
// TODO: This implementation is subject to change in regards to its statefull
// manner. In otherwords, how this relates to the keeper in this module.
type CommitStateDB struct {
	db   ethstate.Database
	//trie ethstate.Trie // only storage addr -> storageMptRoot in this mpt tree

	// TODO: We need to store the context as part of the structure itself opposed
	// to being passed as a parameter (as it should be) in order to implement the
	// StateDB interface. Perhaps there is a better way.
	ctx sdk.Context

	storeKey      sdk.StoreKey
	paramSpace    Subspace
	accountKeeper AccountKeeper
	supplyKeeper  SupplyKeeper
	Watcher       Watcher
	bankKeeper    BankKeeper

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjects        map[ethcmn.Address]*stateObject
	stateObjectsPending map[ethcmn.Address]struct{} // State objects finalized but not yet written to the mpt tree
	stateObjectsDirty   map[ethcmn.Address]struct{} // State objects modified in the current execution

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// The refund counter, also used by state transitioning.
	refund uint64

	bhash, thash ethcmn.Hash
	txIndex      int
	logs         map[ethcmn.Hash][]*ethtypes.Log
	logSize      uint

	preimages map[ethcmn.Hash][]byte

	// Per-transaction access list
	accessList *accessList

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionID int

	params    *Params
	codeCache map[ethcmn.Address]CacheCode
	dbAdapter DbAdapter

	updatedAccount map[ethcmn.Address]struct{} // will destroy every block
}

func NewCommitStateDB(csdbParams CommitStateDBParams) *CommitStateDB {
	csdb := &CommitStateDB{
		db: types.InstanceOfEvmStore(),

		storeKey:      csdbParams.StoreKey,
		paramSpace:    csdbParams.ParamSpace,
		accountKeeper: csdbParams.AccountKeeper,
		supplyKeeper:  csdbParams.SupplyKeeper,
		bankKeeper:    csdbParams.BankKeeper,
		Watcher:       csdbParams.Watcher,

		stateObjects:        make(map[ethcmn.Address]*stateObject),
		stateObjectsPending: make(map[ethcmn.Address]struct{}),
		stateObjectsDirty:   make(map[ethcmn.Address]struct{}),
		preimages:           make(map[ethcmn.Hash][]byte),
		journal:             newJournal(),
		validRevisions:      []revision{},
		accessList:          newAccessList(),
		logSize:             0,
		logs:                make(map[ethcmn.Hash][]*ethtypes.Log),
		codeCache:           make(map[ethcmn.Address]CacheCode, 0),
		dbAdapter:           csdbParams.Ada,
		updatedAccount:      make(map[ethcmn.Address]struct{}),
	}

	//latestHeight := csdb.GetLatestBlockHeight()
	//lastRootHash := csdb.GetRootMptHash(latestHeight)
	//csdb.OpenTrie(lastRootHash)

	return csdb
}

func CreateEmptyCommitStateDB(csdbParams CommitStateDBParams, ctx sdk.Context) *CommitStateDB {
	return NewCommitStateDB(csdbParams).WithContext(ctx)
}

func (csdb *CommitStateDB) SetInternalDb(dba DbAdapter) {
	csdb.dbAdapter = dba
}

// WithContext returns a Database with an updated SDK context
func (csdb *CommitStateDB) WithContext(ctx sdk.Context) *CommitStateDB {
	csdb.ctx = ctx
	return csdb
}

// GetParams returns the total set of evm parameters.
func (csdb *CommitStateDB) GetParams() Params {
	if csdb.params == nil {
		var params Params
		csdb.paramSpace.GetParamSet(csdb.ctx, &params)
		csdb.params = &params
	}
	return *csdb.params
}

// SetParams sets the evm parameters to the param space.
func (csdb *CommitStateDB) SetParams(params Params) {
	csdb.params = &params
	csdb.paramSpace.SetParamSet(csdb.ctx, &params)
}

// ----------------------------------------------------------------------------
// StateDB Interface
// ----------------------------------------------------------------------------

// CreateAccount explicitly creates a state object. If a state object with the
// address already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might
// arise that a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (csdb *CommitStateDB) CreateAccount(addr ethcmn.Address) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	newobj, prevobj := csdb.createObject(addr)
	if prevobj != nil {
		newobj.setBalance(sdk.DefaultBondDenom, sdk.NewDecFromBigIntWithPrec(prevobj.Balance(), sdk.Precision)) // int2dec
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (csdb *CommitStateDB) SubBalance(addr ethcmn.Address, amount *big.Int) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

// AddBalance adds amount to the account associated with addr.
func (csdb *CommitStateDB) AddBalance(addr ethcmn.Address, amount *big.Int) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

// SetBalance sets the balance of an account.
func (csdb *CommitStateDB) SetBalance(addr ethcmn.Address, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetBalance(amount)
	}
}

// GetBalance retrieves the balance from the given address or 0 if object not
// found.
func (csdb *CommitStateDB) GetBalance(addr ethcmn.Address) *big.Int {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Balance()
	}
	return ethcmn.Big0
}

// GetNonce returns the nonce (sequence number) for a given account.
func (csdb *CommitStateDB) GetNonce(addr ethcmn.Address) uint64 {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}

	return 0
}

// SetNonce sets the nonce (sequence number) of an account.
func (csdb *CommitStateDB) SetNonce(addr ethcmn.Address, nonce uint64) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

// GetCodeHash returns the code hash for a given account.
func (csdb *CommitStateDB) GetCodeHash(addr ethcmn.Address) ethcmn.Hash {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject == nil {
		return ethcmn.Hash{}
	}
	return ethcmn.BytesToHash(stateObject.CodeHash())
}

// GetCode returns the code for a given account.
func (csdb *CommitStateDB) GetCode(addr ethcmn.Address) []byte {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	// check for the contract calling from blocked list if contract blocked list is enabled
	if csdb.GetParams().EnableContractBlockedList && csdb.IsContractInBlockedList(addr.Bytes()) {
		panic(addr)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Code(csdb.db)
	}
	return nil
}

// SetCode sets the code for a given account.
func (csdb *CommitStateDB) SetCode(addr ethcmn.Address, code []byte) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		hash := crypto.Keccak256Hash(code)
		stateObject.SetCode(hash, code)

		csdb.codeCache[addr] = CacheCode{
			CodeHash: hash.Bytes(),
			Code:     code,
		}
	}
}

// GetCodeSize returns the code size for a given account.
func (csdb *CommitStateDB) GetCodeSize(addr ethcmn.Address) int {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.CodeSize(csdb.db)
	}
	return 0
}

// AddRefund adds gas to the refund counter.
func (csdb *CommitStateDB) AddRefund(gas uint64) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	csdb.journal.append(refundChange{prev: csdb.refund})
	csdb.refund += gas
}

// SubRefund removes gas from the refund counter. It will panic if the refund
// counter goes below zero.
func (csdb *CommitStateDB) SubRefund(gas uint64) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	csdb.journal.append(refundChange{prev: csdb.refund})
	if gas > csdb.refund {
		panic(fmt.Sprintf("Refund counter below zero (gas: %d > refund: %d)", gas, csdb.refund))
	}
	csdb.refund -= gas
}

// GetRefund returns the current value of the refund counter.
func (csdb *CommitStateDB) GetRefund() uint64 {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	return csdb.refund
}

// GetCommittedState retrieves a value from the given account's committed
// storage.
func (csdb *CommitStateDB) GetCommittedState(addr ethcmn.Address, hash ethcmn.Hash) ethcmn.Hash {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetCommittedState(csdb.db, hash)
	}
	return ethcmn.Hash{}
}

// GetState retrieves a value from the given account's storage store.
func (csdb *CommitStateDB) GetState(addr ethcmn.Address, hash ethcmn.Hash) ethcmn.Hash {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.GetState(csdb.db, hash)
	}
	return ethcmn.Hash{}
}

// SetState sets the storage state with a key, value pair for an account.
func (csdb *CommitStateDB) SetState(addr ethcmn.Address, key, value ethcmn.Hash) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetState(csdb.db, key, value)
	}
}

// Suicide marks the given account as suicided and clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (csdb *CommitStateDB) Suicide(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject == nil {
		return false
	}
	csdb.journal.append(suicideChange{
		account:     &addr,
		prev:        stateObject.suicided,
		prevbalance: sdk.NewDecFromBigIntWithPrec(stateObject.Balance(), sdk.Precision), // int2dec,
	})
	stateObject.markSuicided()
	stateObject.SetBalance(new(big.Int))

	return true
}

// HasSuicided returns if the given account for the specified address has been
// killed.
func (csdb *CommitStateDB) HasSuicided(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	stateObject := csdb.getStateObject(addr)
	if stateObject != nil {
		return stateObject.suicided
	}
	return false
}

// Exist reports whether the given account address exists in the state. Notably,
// this also returns true for suicided accounts.
func (csdb *CommitStateDB) Exist(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	return csdb.getStateObject(addr) != nil
}

// Empty returns whether the state object is either non-existent or empty
// according to the EIP161 specification (balance = nonce = code = 0).
func (csdb *CommitStateDB) Empty(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	return so == nil || so.empty()
}

func (csdb *CommitStateDB) PrepareAccessList(sender ethcmn.Address, dst *ethcmn.Address, precompiles []ethcmn.Address, list ethtypes.AccessList) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	csdb.AddAddressToAccessList(sender)
	if dst != nil {
		csdb.AddAddressToAccessList(*dst)
		// If it's a create-tx, the destination will be added inside evm.create
	}
	for _, addr := range precompiles {
		csdb.AddAddressToAccessList(addr)
	}
	for _, el := range list {
		csdb.AddAddressToAccessList(el.Address)
		for _, key := range el.StorageKeys {
			csdb.AddSlotToAccessList(el.Address, key)
		}
	}
}

// AddressInAccessList returns true if the given address is in the access list.
func (csdb *CommitStateDB) AddressInAccessList(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	return csdb.accessList.ContainsAddress(addr)
}

// SlotInAccessList returns true if the given (address, slot)-tuple is in the access list.
func (csdb *CommitStateDB) SlotInAccessList(addr ethcmn.Address, slot ethcmn.Hash) (bool, bool) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	return csdb.accessList.Contains(addr, slot)
}

// AddAddressToAccessList adds the given address to the access list
func (csdb *CommitStateDB) AddAddressToAccessList(addr ethcmn.Address) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	if csdb.accessList.AddAddress(addr) {
		csdb.journal.append(accessListAddAccountChange{&addr})
	}
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (csdb *CommitStateDB) AddSlotToAccessList(addr ethcmn.Address, slot ethcmn.Hash) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	addrMod, slotMod := csdb.accessList.AddSlot(addr, slot)
	if addrMod {
		// In practice, this should not happen, since there is no way to enter the
		// scope of 'address' without having the 'address' become already added
		// to the access list (via call-variant, create, etc).
		// Better safe than sorry, though
		csdb.journal.append(accessListAddAccountChange{&addr})
	}
	if slotMod {
		csdb.journal.append(accessListAddSlotChange{
			address: &addr,
			slot:    &slot,
		})
	}
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (csdb *CommitStateDB) RevertToSnapshot(revid int) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(csdb.validRevisions), func(i int) bool {
		return csdb.validRevisions[i].id >= revid
	})
	if idx == len(csdb.validRevisions) || csdb.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := csdb.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	csdb.journal.revert(csdb, snapshot)
	csdb.validRevisions = csdb.validRevisions[:idx]
}

// Snapshot returns an identifier for the current revision of the state.
func (csdb *CommitStateDB) Snapshot() int {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	id := csdb.nextRevisionID
	csdb.nextRevisionID++
	csdb.validRevisions = append(csdb.validRevisions, revision{id, csdb.journal.length()})
	return id
}

// AddLog adds a new log to the state and sets the log metadata from the state.
func (csdb *CommitStateDB) AddLog(log *ethtypes.Log) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	csdb.journal.append(addLogChange{txhash: csdb.thash})

	log.TxHash = csdb.thash
	log.BlockHash = csdb.bhash
	log.TxIndex = uint(csdb.txIndex)
	log.Index = csdb.logSize
	csdb.logs[csdb.thash] = append(csdb.logs[csdb.thash], log)
	csdb.logSize++
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (csdb *CommitStateDB) AddPreimage(hash ethcmn.Hash, preimage []byte) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	if _, ok := csdb.preimages[hash]; !ok {
		csdb.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		csdb.preimages[hash] = pi
	}
}

// ForEachStorage iterates over each storage items, all invoke the provided
// callback on each key, value pair.
func (csdb *CommitStateDB) ForEachStorage(addr ethcmn.Address, cb func(key, value ethcmn.Hash) (stop bool)) error {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so == nil {
		return nil
	}
	it := trie.NewIterator(so.getTrie(csdb.db).NodeIterator(nil))

	for it.Next() {
		key := ethcmn.BytesToHash(so.trie.GetKey(it.Key))
		if value, dirty := so.dirtyStorage[key]; dirty {
			if !cb(key, value) {
				return nil
			}
			continue
		}

		if len(it.Value) > 0 {
			_, content, _, err := rlp.Split(it.Value)
			if err != nil {
				return err
			}
			if !cb(key, ethcmn.BytesToHash(content)) {
				return nil
			}
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
// code related
// ----------------------------------------------------------------------------

func (csdb *CommitStateDB) GetCacheCode(addr ethcmn.Address) *CacheCode {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	code, ok := csdb.codeCache[addr]
	if ok {
		return &code
	}

	return nil
}

// GetCode returns the code for a given code hash.
func (csdb *CommitStateDB) GetCodeByHash(hash ethcmn.Hash) []byte {
	code, err := csdb.db.ContractCode(ethcmn.Hash{}, hash)
	if err != nil {
		return nil
	}

	return code
}

func (csdb *CommitStateDB) IteratorCode(cb func(addr ethcmn.Address, c CacheCode) bool) {
	for addr, v := range csdb.codeCache {
		cb(addr, v)
	}
}

// ----------------------------------------------------------------------------
// cosmos hash related
// ----------------------------------------------------------------------------

// SetHeightHash sets the block header hash associated with a given height.
func (csdb *CommitStateDB) SetHeightHash(height uint64, hash ethcmn.Hash) {
	if !csdb.ctx.IsCheckTx() {
		funcName := analyzer.RunFuncName()
		analyzer.StartTxLog(funcName)
		defer analyzer.StopTxLog(funcName)
	}

	store := csdb.dbAdapter.NewStore(csdb.ctx.KVStore(csdb.storeKey), KeyPrefixHeightHash)
	key := HeightHashKey(height)
	store.Set(key, hash.Bytes())
}

// GetHeightHash returns the block header hash associated with a given block height and chain epoch number.
func (csdb *CommitStateDB) GetHeightHash(height uint64) ethcmn.Hash {
	store := csdb.dbAdapter.NewStore(csdb.ctx.KVStore(csdb.storeKey), KeyPrefixHeightHash)
	key := HeightHashKey(height)
	bz := store.Get(key)
	if len(bz) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(bz)
}

// BlockHash returns the current block hash set by Prepare.
func (csdb *CommitStateDB) BlockHash() ethcmn.Hash {
	return csdb.bhash
}

func (csdb *CommitStateDB) SetBlockHash(hash ethcmn.Hash) {
	csdb.bhash = hash
}

// ----------------------------------------------------------------------------
// Transaction logs
// Required for upgrade logic or ease of querying.
// NOTE: we use BinaryLengthPrefixed since the tx logs are also included on Result data,
// which can't use BinaryBare.
// ----------------------------------------------------------------------------

// SetLogs sets the logs for a transaction in the KVStore.
func (csdb *CommitStateDB) SetLogs(hash ethcmn.Hash, logs []*ethtypes.Log) error {
	csdb.logs[hash] = logs
	return nil
}

// DeleteLogs removes the logs from the KVStore. It is used during journal.Revert.
func (csdb *CommitStateDB) DeleteLogs(hash ethcmn.Hash) {
	delete(csdb.logs, hash)
}

// GetLogs returns the current logs for a given transaction hash from the KVStore.
func (csdb *CommitStateDB) GetLogs(hash ethcmn.Hash) ([]*ethtypes.Log, error) {
	return csdb.logs[hash], nil
}

func (csdb *CommitStateDB) SetLogSize(logSize uint) {
	csdb.logSize = logSize
}

func (csdb *CommitStateDB) GetLogSize() uint {
	return csdb.logSize
}

// ----------------------------------------------------------------------------
// Persistence
// ----------------------------------------------------------------------------

// Commit writes the state to the underlying in-memory trie database.
func (csdb *CommitStateDB) Commit(deleteEmptyObjects bool) (ethcmn.Hash, error) {
	//if csdb.dbErr != nil {
	//	return ethcmn.Hash{}, fmt.Errorf("commit aborted due to earlier error: %v", csdb.dbErr)
	//}

	// Finalize any pending changes and merge everything into the tries
	csdb.IntermediateRoot(deleteEmptyObjects)

	// Commit objects to the trie, measuring the elapsed time
	codeWriter := csdb.db.TrieDB().DiskDB().NewBatch()
	for addr := range csdb.stateObjectsDirty {
		if obj := csdb.stateObjects[addr]; !obj.deleted {
			// Write any contract code associated with the state object
			if obj.code != nil && obj.dirtyCode {
				rawdb.WriteCode(codeWriter, ethcmn.BytesToHash(obj.CodeHash()), obj.code)
				obj.dirtyCode = false
			}

			// Write any storage changes in the state object to its storage trie
			if err := obj.CommitTrie(csdb.db); err != nil {
				return ethcmn.Hash{}, err
			}
			csdb.UpdateAccountInfo(obj.account)
		}
	}

	if len(csdb.stateObjectsDirty) > 0 {
		csdb.stateObjectsDirty = make(map[ethcmn.Address]struct{})
	}

	if codeWriter.ValueSize() > 0 {
		if err := codeWriter.Write(); err != nil {
			log.Crit("Failed to commit dirty codes", "error", err)
		}
	}

	return ethcmn.Hash{}, nil
}

// Finalise finalises the state by removing the s destructed objects and clears
// the journal as well as the refunds. Finalise, however, will not push any updates
// into the tries just yet. Only IntermediateRoot or Commit will do that.
func (csdb *CommitStateDB) Finalise(deleteEmptyObjects bool) {
	for addr := range csdb.journal.dirties {
		obj, exist := csdb.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `s.journal.dirties` but not in `s.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}
		if obj.suicided || (deleteEmptyObjects && obj.empty()) {
			obj.deleted = true
		} else {
			obj.finalise(true) // Prefetch slots in the background
		}
		csdb.stateObjectsPending[addr] = struct{}{}
		csdb.stateObjectsDirty[addr] = struct{}{}
	}

	// Invalidate journal because reverting across transactions is not allowed.
	csdb.clearJournalAndRefund()
}

// IntermediateRoot computes the current root hash of the state trie.
// It is called in between transactions to get the root hash that
// goes into transaction receipts.
func (csdb *CommitStateDB) IntermediateRoot(deleteEmptyObjects bool) ethcmn.Hash {
	// Finalise all the dirty storage states and write them into the tries
	csdb.Finalise(deleteEmptyObjects)

	// Although naively it makes sense to retrieve the account trie and then do
	// the contract storage and account updates sequentially, that short circuits
	// the account prefetcher. Instead, let's process all the storage updates
	// first, giving the account prefeches just a few more milliseconds of time
	// to pull useful data from disk.
	for addr := range csdb.stateObjectsPending {
		if obj := csdb.stateObjects[addr]; !obj.deleted {
			obj.updateRoot(csdb.db)
		}
	}

	for addr := range csdb.stateObjectsPending {
		if obj := csdb.stateObjects[addr]; obj.deleted {
			csdb.deleteStateObject(obj)
		} else {
			csdb.updateStateObject(obj)
		}
	}

	if len(csdb.stateObjectsPending) > 0 {
		csdb.stateObjectsPending = make(map[ethcmn.Address]struct{})
	}

	return ethcmn.Hash{}
}

// ----------------------------------------------------------------------------
// State object related
// ----------------------------------------------------------------------------

// updateStateObject writes the given state object to the store.
func (csdb *CommitStateDB) updateStateObject(so *stateObject) error {
	// NOTE: we don't use sdk.NewCoin here to avoid panic on test importer's genesis
	newBalance := sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDecFromBigIntWithPrec(so.Balance(), sdk.Precision)} // int2dec
	if !newBalance.IsValid() {
		return fmt.Errorf("invalid balance %s", newBalance)
	}

	//checking and reject tx if address in blacklist
	if csdb.bankKeeper.BlacklistedAddr(so.account.GetAddress()) {
		return fmt.Errorf("address <%s> in blacklist is not allowed", so.account.GetAddress().String())
	}

	coins := so.account.GetCoins()
	balance := coins.AmountOf(newBalance.Denom)
	if balance.IsZero() || !balance.Equal(newBalance.Amount) {
		coins = coins.Add(newBalance)
	}

	if err := so.account.SetCoins(coins); err != nil {
		return err
	}

	csdb.UpdateAccountInfo(so.account)

	return nil
}

func (csdb *CommitStateDB) UpdateAccountInfo(acc *types2.EthAccount) {
	csdb.accountKeeper.SetAccount(csdb.ctx, acc)
	if !csdb.ctx.IsCheckTx() {
		if csdb.Watcher.Enabled() {
			csdb.Watcher.SaveAccount(acc, false)
		}
	}
}

// deleteStateObject removes the given state object from the state store.
func (csdb *CommitStateDB) deleteStateObject(so *stateObject) {
	so.deleted = true
	csdb.accountKeeper.RemoveAccount(csdb.ctx, so.account)
}

// ClearStateObjects clears cache of state objects to handle account changes outside of the EVM
func (csdb *CommitStateDB) ClearStateObjects() {
	csdb.stateObjects = make(map[ethcmn.Address]*stateObject)
	csdb.stateObjectsPending = make(map[ethcmn.Address]struct{})
	csdb.stateObjectsDirty = make(map[ethcmn.Address]struct{})
}

// GetOrNewStateObject retrieves a state object or create a new state object if
// nil.
func (csdb *CommitStateDB) GetOrNewStateObject(addr ethcmn.Address) *stateObject {
	so := csdb.getStateObject(addr)
	if so == nil || so.deleted {
		so, _ = csdb.createObject(addr)
	}

	return so
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (csdb *CommitStateDB) createObject(addr ethcmn.Address) (newObj, prevObj *stateObject) {
	prevObj = csdb.getStateObject(addr)

	acc := csdb.accountKeeper.NewAccountWithAddress(csdb.ctx, sdk.AccAddress(addr.Bytes()))
	newObj = newStateObject(csdb, acc)
	newObj.setNonce(0) // sets the object to dirty

	if prevObj == nil {
		csdb.journal.append(createObjectChange{account: &addr})
	} else {
		csdb.journal.append(resetObjectChange{prev: prevObj})
	}
	csdb.setStateObject(newObj)

	if prevObj != nil && !prevObj.deleted {
		return newObj, prevObj
	}
	return newObj, nil
}

// getStateObject retrieves a state object given by the address, returning nil if
// the object is not found or was deleted in this execution context. If you need
// to differentiate between non-existent/just-deleted, use getDeletedStateObject.
func (csdb *CommitStateDB) getStateObject(addr ethcmn.Address) *stateObject {
	if obj := csdb.getDeletedStateObject(addr); obj != nil && !obj.deleted {
		return obj
	}
	return nil
}

// getDeletedStateObject is similar to getStateObject, but instead of returning
// nil for a deleted state object, it returns the actual object with the deleted
// flag set. This is needed by the state journal to revert to the correct s-
// destructed object instead of wiping all knowledge about the state object.
func (csdb *CommitStateDB) getDeletedStateObject(addr ethcmn.Address) *stateObject {
	// Prefer live objects if any is available
	if obj := csdb.stateObjects[addr]; obj != nil {
		if _, ok := csdb.updatedAccount[addr]; ok {
			obj.UpdateAccInfo()
			delete(csdb.updatedAccount, addr)
		}
		return obj
	}

	// otherwise, attempt to fetch the account from the account mapper
	acc := csdb.accountKeeper.GetAccount(csdb.ctx, sdk.AccAddress(addr.Bytes()))
	if acc == nil {
		csdb.setError(fmt.Errorf("no account found for address: %s", addr.String()))
		return nil
	}

	// insert the state object into the live set
	so := newStateObject(csdb, acc)
	csdb.setStateObject(so)

	return so
}

func (csdb *CommitStateDB) setStateObject(object *stateObject) {
	csdb.stateObjects[object.Address()] = object
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

// TxIndex returns the current transaction index set by Prepare.
func (csdb *CommitStateDB) TxIndex() int {
	return csdb.txIndex
}

// Database retrieves the low level database supporting the lower level trie
// ops. It is not used in Ethermint, so it returns nil.
func (csdb *CommitStateDB) Database() ethstate.Database {
	return csdb.db
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (csdb *CommitStateDB) Preimages() map[ethcmn.Hash][]byte {
	return csdb.preimages
}

// GetStateByKey retrieves a value from the given account's storage store.
func (csdb *CommitStateDB) GetStateByKey(addr ethcmn.Address, key ethcmn.Hash) ethcmn.Hash {
	var (
		enc []byte
		err error
	)

	if enc, err = csdb.StorageTrie(addr).TryGet(key.Bytes()); err != nil {
		return ethcmn.Hash{}
	}

	var value ethcmn.Hash
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			return ethcmn.Hash{}
		}
		value.SetBytes(content)
	}

	return value
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying account mapper and store keys to avoid reloading data for the
// next operations.
func (csdb *CommitStateDB) Reset(_ ethcmn.Hash) error {
	csdb.stateObjects = make(map[ethcmn.Address]*stateObject)
	csdb.stateObjectsPending = make(map[ethcmn.Address]struct{})
	csdb.stateObjectsDirty = make(map[ethcmn.Address]struct{})
	csdb.thash = ethcmn.Hash{}
	csdb.bhash = ethcmn.Hash{}
	csdb.txIndex = 0
	csdb.logSize = 0
	csdb.preimages = make(map[ethcmn.Hash][]byte)
	csdb.accessList = newAccessList()
	csdb.params = nil

	csdb.clearJournalAndRefund()
	return nil
}

func (csdb *CommitStateDB) clearJournalAndRefund() {
	if len(csdb.journal.entries) > 0 {
		csdb.journal = newJournal()
		csdb.refund = 0
	}
	csdb.validRevisions = csdb.validRevisions[:0] // Snapshots can be created without journal entires
}

// Prepare sets the current transaction hash and index and block hash which is
// used when the EVM emits new state logs.
func (csdb *CommitStateDB) Prepare(thash, bhash ethcmn.Hash, txi int) {
	csdb.thash = thash
	csdb.bhash = bhash
	csdb.txIndex = txi
}

// setError remembers the first non-nil error it is called with.
func (csdb *CommitStateDB) setError(err error) {
	if csdb.dbErr == nil {
		csdb.dbErr = err
	}
}

// Error returns the first non-nil error the StateDB encountered.
func (csdb *CommitStateDB) Error() error {
	return csdb.dbErr
}

// RawDump returns a raw state dump.
//
// TODO: Implement if we need it, especially for the RPC API.
func (csdb *CommitStateDB) RawDump() ethstate.Dump {
	return ethstate.Dump{}
}

func (csdb *CommitStateDB) MarkUpdatedAcc(addList []ethcmn.Address) {
	for _, addr := range addList {
		csdb.updatedAccount[addr] = struct{}{}
	}
}

// ----------------------------------------------------------------------------
// Proof related
// ----------------------------------------------------------------------------

type proofList [][]byte

func (n *proofList) Put(key []byte, value []byte) error {
	*n = append(*n, value)
	return nil
}

func (n *proofList) Delete(key []byte) error {
	panic("not supported")
}

// StorageTrie returns nil as the state in Ethermint does not use a direct
// storage trie.
func (csdb *CommitStateDB) StorageTrie(addr ethcmn.Address) ethstate.Trie {
	stateObject := csdb.getStateObject(addr)
	if stateObject == nil {
		return nil
	}
	cpy := stateObject.deepCopy(csdb)
	cpy.updateTrie(csdb.db)
	return cpy.getTrie(csdb.db)
}

// GetProof returns the Merkle proof for a given account.
func (csdb *CommitStateDB) GetProof(addr ethcmn.Address) ([][]byte, error) {
	var proof proofList
	accTrie := csdb.StorageTrie(addr)
	if accTrie == nil {
		return proof, errors.New("storage trie for requested address does not exist")
	}

	addrHash := crypto.Keccak256Hash(addr.Bytes())
	err := accTrie.Prove(addrHash[:], 0, &proof)
	return proof, err

}

// GetStorageProof returns the Merkle proof for given storage slot.
func (csdb *CommitStateDB) GetStorageProof(a ethcmn.Address, key ethcmn.Hash) ([][]byte, error) {
	var proof proofList
	addrTrie := csdb.StorageTrie(a)
	if addrTrie == nil {
		return proof, errors.New("storage trie for requested address does not exist")
	}
	err := addrTrie.Prove(crypto.Keccak256(key.Bytes()), 0, &proof)
	return proof, err
}
