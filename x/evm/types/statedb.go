package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/VictoriaMetrics/fastcache"
	ethermint "github.com/okx/okbchain/app/types"
	"github.com/tendermint/go-amino"
	"math/big"
	"sort"
	"sync"

	"github.com/okx/okbchain/libs/system/trace"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/prefix"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
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
	BankKeeper    BankKeeper
	Ada           DbAdapter
	// Amino codec
	Cdc *codec.Codec

	StateCache *fastcache.Cache

	DB ethstate.Database
}

type Watcher interface {
	SaveAccount(account auth.Account, isDirectly bool)
	AddDelAccMsg(account auth.Account, isDirectly bool)
	SaveState(addr ethcmn.Address, key, value []byte)
	Enabled() bool
	SaveContractBlockedListItem(addr sdk.AccAddress)
	SaveContractDeploymentWhitelistItem(addr sdk.AccAddress)
	DeleteContractBlockedList(addr sdk.AccAddress)
	DeleteContractDeploymentWhitelist(addr sdk.AccAddress)
	SaveContractMethodBlockedListItem(addr sdk.AccAddress, methods []byte)
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
// Warning!!! If you change CommitStateDB.member you must be careful ResetCommitStateDB contract BananaLF.
type CommitStateDB struct {
	db         ethstate.Database
	prefetcher *mpt.TriePrefetcher

	// TODO: We need to store the context as part of the structure itself opposed
	// to being passed as a parameter (as it should be) in order to implement the
	// StateDB interface. Perhaps there is a better way.
	ctx sdk.Context

	storeKey      sdk.StoreKey
	paramSpace    Subspace
	accountKeeper AccountKeeper
	supplyKeeper  SupplyKeeper
	bankKeeper    BankKeeper

	// array that hold 'live' objects, which will get modified while processing a
	// state transition
	stateObjects        map[ethcmn.Address]*stateObject
	stateObjectsPending map[ethcmn.Address]struct{} // State objects finalized but not yet written to the mpt tree
	stateObjectsDirty   map[ethcmn.Address]struct{} // State objects modified in the current execution

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash ethcmn.Hash
	txIndex      int
	logSize      uint
	logs         map[ethcmn.Hash][]*ethtypes.Log

	preimages map[ethcmn.Hash][]byte

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionID int

	// Per-transaction access list
	accessList *accessList

	// mutex for state deep copying
	lock sync.Mutex

	params    *Params
	codeCache map[ethcmn.Address]CacheCode
	dbAdapter DbAdapter

	// Amino codec
	cdc *codec.Codec

	updatedAccount map[ethcmn.Address]struct{} // will destroy every block

	GuFactor sdk.Dec
}

// Warning!!! If you change CommitStateDB.member you must be careful ResetCommitStateDB contract BananaLF.

type StoreProxy interface {
	Set(key, value []byte)
	Get(key []byte) []byte
	Delete(key []byte)
	Has(key []byte) bool
}

type DbAdapter interface {
	NewStore(parent types.KVStore, prefix []byte) StoreProxy
}

type DefaultPrefixDb struct {
}

func (d DefaultPrefixDb) NewStore(parent types.KVStore, Prefix []byte) StoreProxy {
	return prefix.NewStore(parent, Prefix)
}

// NewCommitStateDB returns a reference to a newly initialized CommitStateDB
// which implements Geth's state.StateDB interface.
//
// CONTRACT: Stores used for state must be cache-wrapped as the ordering of the
// key/value space matters in determining the merkle root.
func NewCommitStateDB(csdbParams CommitStateDBParams) *CommitStateDB {
	csdb := &CommitStateDB{
		db: csdbParams.DB,

		storeKey:      csdbParams.StoreKey,
		paramSpace:    csdbParams.ParamSpace,
		accountKeeper: csdbParams.AccountKeeper,
		supplyKeeper:  csdbParams.SupplyKeeper,
		bankKeeper:    csdbParams.BankKeeper,
		cdc:           csdbParams.Cdc,

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
		GuFactor:            DefaultGuFactor,
	}

	return csdb
}

func ResetCommitStateDB(csdb *CommitStateDB, csdbParams CommitStateDBParams, ctx *sdk.Context) {
	csdb.db = csdbParams.DB

	csdb.storeKey = csdbParams.StoreKey
	csdb.paramSpace = csdbParams.ParamSpace
	csdb.accountKeeper = csdbParams.AccountKeeper
	csdb.supplyKeeper = csdbParams.SupplyKeeper
	csdb.bankKeeper = csdbParams.BankKeeper
	csdb.cdc = csdbParams.Cdc

	if csdb.stateObjects != nil {
		for k := range csdb.stateObjects {
			delete(csdb.stateObjects, k)
		}
	} else {
		csdb.stateObjects = make(map[ethcmn.Address]*stateObject)
	}

	if csdb.stateObjectsPending != nil {
		for k := range csdb.stateObjectsPending {
			delete(csdb.stateObjectsPending, k)
		}
	} else {
		csdb.stateObjectsPending = make(map[ethcmn.Address]struct{})
	}

	if csdb.stateObjectsDirty != nil {
		for k := range csdb.stateObjectsDirty {
			delete(csdb.stateObjectsDirty, k)
		}
	} else {
		csdb.stateObjectsDirty = make(map[ethcmn.Address]struct{})
	}

	if csdb.preimages != nil {
		for k := range csdb.preimages {
			delete(csdb.preimages, k)
		}
	} else {
		csdb.preimages = make(map[ethcmn.Hash][]byte)
	}

	if csdb.journal != nil {
		csdb.journal.entries = nil
		if csdb.journal.dirties != nil {
			for k := range csdb.journal.dirties {
				delete(csdb.journal.dirties, k)
			}
		} else {
			csdb.journal.dirties = make(map[ethcmn.Address]int)
		}
	} else {
		csdb.journal = newJournal()
	}

	if csdb.validRevisions != nil {
		csdb.validRevisions = csdb.validRevisions[:0]
	} else {
		csdb.validRevisions = []revision{}
	}

	if csdb.accessList != nil {
		if csdb.accessList.addresses != nil {
			for k := range csdb.accessList.addresses {
				delete(csdb.accessList.addresses, k)
			}
		} else {
			csdb.accessList.addresses = make(map[ethcmn.Address]int)
		}
		csdb.accessList.slots = nil
	} else {
		csdb.accessList = newAccessList()
	}

	csdb.logSize = 0

	if csdb.logs != nil {
		for k := range csdb.logs {
			delete(csdb.logs, k)
		}
	} else {
		csdb.logs = make(map[ethcmn.Hash][]*ethtypes.Log)
	}

	if csdb.codeCache != nil {
		for k := range csdb.codeCache {
			delete(csdb.codeCache, k)
		}
	} else {
		csdb.codeCache = make(map[ethcmn.Address]CacheCode, 0)
	}

	csdb.dbAdapter = csdbParams.Ada

	if csdb.updatedAccount != nil {
		for k := range csdb.updatedAccount {
			delete(csdb.updatedAccount, k)
		}
	} else {
		csdb.updatedAccount = make(map[ethcmn.Address]struct{})
	}

	csdb.prefetcher = nil
	csdb.ctx = *ctx
	csdb.refund = 0
	csdb.thash = ethcmn.Hash{}
	csdb.bhash = ethcmn.Hash{}
	csdb.txIndex = 0
	csdb.dbErr = nil
	csdb.nextRevisionID = 0
	csdb.params = nil
	csdb.GuFactor = DefaultGuFactor
}

func CreateEmptyCommitStateDB(csdbParams CommitStateDBParams, ctx sdk.Context) *CommitStateDB {
	csdb := NewCommitStateDB(csdbParams).WithContext(ctx)
	return csdb
}

func (csdb *CommitStateDB) SetInternalDb(dba DbAdapter) {
	csdb.dbAdapter = dba
}

// WithContext returns a Database with an updated SDK context
func (csdb *CommitStateDB) WithContext(ctx sdk.Context) *CommitStateDB {
	csdb.ctx = ctx
	return csdb
}

func (csdb *CommitStateDB) GetCacheCode(addr ethcmn.Address) *CacheCode {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetCacheCode"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	code, ok := csdb.codeCache[addr]
	if ok {
		return &code
	}

	return nil
}

// IteratorCode is iterator code cache，it can't be used in onchain execution
func (csdb *CommitStateDB) IteratorCode(cb func(addr ethcmn.Address, c CacheCode) bool) {
	for addr, v := range csdb.codeCache {
		cb(addr, v)
	}
}

// ----------------------------------------------------------------------------
// Setters
// ----------------------------------------------------------------------------

// SetHeightHash sets the block header hash associated with a given height.
func (csdb *CommitStateDB) SetHeightHash(height uint64, hash ethcmn.Hash) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "SetHeightHash"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	csdb.setHeightHashInRawDB(height, hash)
}

// SetParams sets the evm parameters to the param space.
func (csdb *CommitStateDB) SetParams(params Params) {
	csdb.params = &params
	csdb.paramSpace.SetParamSet(csdb.ctx, &params)
	GetEvmParamsCache().SetNeedParamsUpdate()
}

// SetStorage replaces the entire storage for the specified account with given
// storage. This function should only be used for debugging.
func (csdb *CommitStateDB) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	stateObject := csdb.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetStorage(storage)
	}
}

// SetBalance sets the balance of an account.
func (csdb *CommitStateDB) SetBalance(addr ethcmn.Address, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetBalance(amount)
	}
}

// AddBalance adds amount to the account associated with addr.
func (csdb *CommitStateDB) AddBalance(addr ethcmn.Address, amount *big.Int) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "AddBalance"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.AddBalance(amount)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (csdb *CommitStateDB) SubBalance(addr ethcmn.Address, amount *big.Int) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "SubBalance"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SubBalance(amount)
	}
}

// SetNonce sets the nonce (sequence number) of an account.
func (csdb *CommitStateDB) SetNonce(addr ethcmn.Address, nonce uint64) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "SetNonce"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetNonce(nonce)
	}
}

// SetState sets the storage state with a key, value pair for an account.
func (csdb *CommitStateDB) SetState(addr ethcmn.Address, key, value ethcmn.Hash) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "SetState"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetState(csdb.db, key, value)
	}
}

// SetCode sets the code for a given account.
func (csdb *CommitStateDB) SetCode(addr ethcmn.Address, code []byte) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "SetCode"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		hash := Keccak256HashWithCache(code)
		so.SetCode(hash, code)
		csdb.codeCache[addr] = CacheCode{
			CodeHash: hash.Bytes(),
			Code:     code,
		}
	}
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

// AddLog adds a new log to the state and sets the log metadata from the state.
func (csdb *CommitStateDB) AddLog(log *ethtypes.Log) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "AddLog"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	csdb.journal.append(addLogChange{txhash: csdb.thash})

	log.TxHash = csdb.thash
	log.BlockHash = csdb.bhash
	log.TxIndex = uint(csdb.txIndex)
	log.Index = csdb.logSize

	csdb.logSize = csdb.logSize + 1
	csdb.logs[csdb.thash] = append(csdb.logs[csdb.thash], log)
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (csdb *CommitStateDB) AddPreimage(hash ethcmn.Hash, preimage []byte) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "AddPreimage"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	if _, ok := csdb.preimages[hash]; !ok {
		csdb.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		csdb.preimages[hash] = pi
	}
}

// AddRefund adds gas to the refund counter.
func (csdb *CommitStateDB) AddRefund(gas uint64) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "AddRefund"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	csdb.journal.append(refundChange{prev: csdb.refund})
	csdb.refund += gas
}

// SubRefund removes gas from the refund counter. It will panic if the refund
// counter goes below zero.
func (csdb *CommitStateDB) SubRefund(gas uint64) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "SubRefund"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	csdb.journal.append(refundChange{prev: csdb.refund})
	if gas > csdb.refund {
		panic("refund counter below zero")
	}

	csdb.refund -= gas
}

// AddAddressToAccessList adds the given address to the access list
func (csdb *CommitStateDB) AddAddressToAccessList(addr ethcmn.Address) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "AddAddressToAccessList"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	if csdb.accessList.AddAddress(addr) {
		csdb.journal.append(accessListAddAccountChange{&addr})
	}
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (csdb *CommitStateDB) AddSlotToAccessList(addr ethcmn.Address, slot ethcmn.Hash) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "AddSlotToAccessList"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
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
func (csdb *CommitStateDB) PrepareAccessList(sender ethcmn.Address, dest *ethcmn.Address, precompiles []ethcmn.Address, txAccesses ethtypes.AccessList) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "PrepareAccessList"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	csdb.AddAddressToAccessList(sender)
	if csdb != nil {
		csdb.AddAddressToAccessList(*dest)
		// If it's a create-tx, the destination will be added inside evm.create
	}
	for _, addr := range precompiles {
		csdb.AddAddressToAccessList(addr)
	}
	for _, el := range txAccesses {
		csdb.AddAddressToAccessList(el.Address)
		for _, key := range el.StorageKeys {
			csdb.AddSlotToAccessList(el.Address, key)
		}
	}
}

// AddressInAccessList returns true if the given address is in the access list.
func (csdb *CommitStateDB) AddressInAccessList(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := "AddressInAccessList"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	return csdb.accessList.ContainsAddress(addr)
}

// SlotInAccessList returns true if the given (address, slot)-tuple is in the access list.
func (csdb *CommitStateDB) SlotInAccessList(addr ethcmn.Address, slot ethcmn.Hash) (bool, bool) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "SlotInAccessList"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	return csdb.accessList.Contains(addr, slot)
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

// GetHeightHash returns the block header hash associated with a given block height and chain epoch number.
func (csdb *CommitStateDB) GetHeightHash(height uint64) ethcmn.Hash {
	return csdb.getHeightHashInRawDB(height)
}

// GetParams returns the total set of evm parameters.
func (csdb *CommitStateDB) GetParams() Params {
	if csdb.params == nil {
		var params Params
		if csdb.ctx.UseParamCache() {
			if GetEvmParamsCache().IsNeedParamsUpdate() {
				csdb.paramSpace.GetParamSet(csdb.ctx, &params)
				GetEvmParamsCache().UpdateParams(params, csdb.ctx.IsCheckTx())
			} else {
				params = GetEvmParamsCache().GetParams()
			}
		} else {
			csdb.paramSpace.GetParamSet(csdb.ctx, &params)
		}
		csdb.params = &params
	}
	return *csdb.params
}

// GetBalance retrieves the balance from the given address or 0 if object not
// found.
func (csdb *CommitStateDB) GetBalance(addr ethcmn.Address) *big.Int {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetBalance"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Balance()
	}

	return zeroBalance
}

// GetNonce returns the nonce (sequence number) for a given account.
func (csdb *CommitStateDB) GetNonce(addr ethcmn.Address) uint64 {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetNonce"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Nonce()
	}

	return 0
}

// TxIndex returns the current transaction index set by Prepare.
func (csdb *CommitStateDB) TxIndex() int {
	return csdb.txIndex
}

// BlockHash returns the current block hash set by Prepare.
func (csdb *CommitStateDB) BlockHash() ethcmn.Hash {
	return csdb.bhash
}

func (csdb *CommitStateDB) SetBlockHash(hash ethcmn.Hash) {
	csdb.bhash = hash
}

// GetCode returns the code for a given account.
func (csdb *CommitStateDB) GetCode(addr ethcmn.Address) []byte {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetCode"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	// check for the contract calling from blocked list if contract blocked list is enabled
	if csdb.GetParams().EnableContractBlockedList && csdb.IsContractInBlockedList(addr.Bytes()) {
		err := ErrContractBlockedVerify{fmt.Sprintf("failed. the contract %s is not allowed to invoke", addr.Hex())}
		panic(err)
	}

	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Code(csdb.db)
	}
	return nil
}

// GetCode returns the code for a given code hash.
func (csdb *CommitStateDB) GetCodeByHash(hash ethcmn.Hash) []byte {
	return csdb.GetCodeByHashInRawDB(hash)
}

// GetCodeSize returns the code size for a given account.
func (csdb *CommitStateDB) GetCodeSize(addr ethcmn.Address) int {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetCodeSize"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so != nil {
		return so.CodeSize(csdb.db)
	}
	return 0
}

// GetCodeHash returns the code hash for a given account.
func (csdb *CommitStateDB) GetCodeHash(addr ethcmn.Address) ethcmn.Hash {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetCodeHash"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so == nil {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(so.CodeHash())
}

// GetState retrieves a value from the given account's storage store.
func (csdb *CommitStateDB) GetState(addr ethcmn.Address, hash ethcmn.Hash) ethcmn.Hash {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetState"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so != nil {
		return so.GetState(csdb.db, hash)
	}

	return ethcmn.Hash{}
}

// GetStateByKey retrieves a value from the given account's storage store.
func (csdb *CommitStateDB) GetStateByKey(addr ethcmn.Address, key ethcmn.Hash) ethcmn.Hash {
	return csdb.GetStateByKeyMpt(addr, key)
}

// GetCommittedState retrieves a value from the given account's committed
// storage.
func (csdb *CommitStateDB) GetCommittedState(addr ethcmn.Address, hash ethcmn.Hash) ethcmn.Hash {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetCommittedState"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so != nil {
		return so.GetCommittedState(csdb.db, hash)
	}

	return ethcmn.Hash{}
}

// GetLogs returns the current logs for a given transaction hash from the KVStore.
func (csdb *CommitStateDB) GetLogs(hash ethcmn.Hash) ([]*ethtypes.Log, error) {
	return csdb.logs[hash], nil
}

// GetRefund returns the current value of the refund counter.
func (csdb *CommitStateDB) GetRefund() uint64 {
	if !csdb.ctx.IsCheckTx() {
		funcName := "GetRefund"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	return csdb.refund
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (csdb *CommitStateDB) Preimages() map[ethcmn.Hash][]byte {
	return csdb.preimages
}

// HasSuicided returns if the given account for the specified address has been
// killed.
func (csdb *CommitStateDB) HasSuicided(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := "HasSuicided"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so != nil {
		return so.suicided
	}

	return false
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

// ----------------------------------------------------------------------------
// Persistence
// ----------------------------------------------------------------------------

// Commit writes the state to the appropriate KVStores. For each state object
// in the cache, it will either be removed, or have it's code set and/or it's
// state (storage) updated. In addition, the state object (account) itself will
// be written. Finally, the root hash (version) will be returned.
func (csdb *CommitStateDB) Commit(deleteEmptyObjects bool) (ethcmn.Hash, error) {
	// Finalize any pending changes and merge everything into the tries
	csdb.IntermediateRoot(deleteEmptyObjects)

	// If there was a trie prefetcher operating, it gets aborted and irrevocably
	// modified after we start retrieving tries. Remove it from the statedb after
	// this round of use.
	//
	// This is weird pre-byzantium since the first tx runs with a prefetcher and
	// the remainder without, but pre-byzantium even the initial prefetcher is
	// useless, so no sleep lost.
	prefetcher := csdb.prefetcher
	if csdb.prefetcher != nil {
		defer func() {
			csdb.prefetcher.Close()
			csdb.prefetcher = nil
		}()
	}

	return csdb.CommitMpt(prefetcher)
}

// Finalise finalizes the state objects (accounts) state by setting their state,
// removing the csdb destructed objects and clearing the journal as well as the
// refunds.
func (csdb *CommitStateDB) Finalise(deleteEmptyObjects bool) {
	addressesToPrefetch := make([][]byte, 0, len(csdb.journal.dirties))
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

		// At this point, also ship the address off to the precacher. The precacher
		// will start loading tries, and when the change is eventually committed,
		// the commit-phase will be a lot faster
		addressesToPrefetch = append(addressesToPrefetch, ethcmn.CopyBytes(addr[:])) // Copy needed for closure
	}
	//TODO need to prefecth to acc trie

	// Invalidate journal because reverting across transactions is not allowed.
	csdb.clearJournalAndRefund()
}

// IntermediateRoot returns the current root hash of the state. It is called in
// between transactions to get the root hash that goes into transaction
// receipts.
//
// NOTE: The SDK has not concept or method of getting any intermediate merkle
// root as commitment of the merkle-ized tree doesn't happen until the
// BaseApps' EndBlocker.
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

	//usedAddrs := make([][]byte, 0, len(csdb.stateObjectsPending))
	for addr := range csdb.stateObjectsPending {
		if obj := csdb.stateObjects[addr]; obj.deleted {
			csdb.deleteStateObject(obj)
		} else {
			csdb.updateStateObject(obj)
		}
		//usedAddrs = append(usedAddrs, ethcmn.CopyBytes(addr[:])) // Copy needed for closure
	}
	//if csdb.prefetcher != nil {
	//	csdb.prefetcher.used(csdb.originalRoot, usedAddrs)
	//}

	if len(csdb.stateObjectsPending) > 0 {
		csdb.stateObjectsPending = make(map[ethcmn.Address]struct{})
	}

	return ethcmn.Hash{}
}

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

	csdb.accountKeeper.SetAccount(csdb.ctx, so.account)
	if !csdb.ctx.IsCheckTx() {
		if csdb.ctx.GetWatcher().Enabled() {
			csdb.ctx.GetWatcher().SaveAccount(so.account)
		}
	}

	return nil
}

// deleteStateObject removes the given state object from the state store.
func (csdb *CommitStateDB) deleteStateObject(so *stateObject) {
	csdb.accountKeeper.RemoveAccount(csdb.ctx, so.account)
}

// ----------------------------------------------------------------------------
// Snapshotting
// ----------------------------------------------------------------------------

// Snapshot returns an identifier for the current revision of the state.
func (csdb *CommitStateDB) Snapshot() int {
	if !csdb.ctx.IsCheckTx() {
		funcName := "Snapshot"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	id := csdb.nextRevisionID
	csdb.nextRevisionID++

	csdb.validRevisions = append(
		csdb.validRevisions,
		revision{
			id:           id,
			journalIndex: csdb.journal.length(),
		},
	)

	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (csdb *CommitStateDB) RevertToSnapshot(revID int) {
	if !csdb.ctx.IsCheckTx() {
		funcName := "RevertToSnapshot"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	// find the snapshot in the stack of valid snapshots
	idx := sort.Search(len(csdb.validRevisions), func(i int) bool {
		return csdb.validRevisions[i].id >= revID
	})

	if idx == len(csdb.validRevisions) || csdb.validRevisions[idx].id != revID {
		panic(fmt.Errorf("revision ID %v cannot be reverted", revID))
	}

	snapshot := csdb.validRevisions[idx].journalIndex

	// replay the journal to undo changes and remove invalidated snapshots
	csdb.journal.revert(csdb, snapshot)
	csdb.validRevisions = csdb.validRevisions[:idx]
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

// Database retrieves the low level database supporting the lower level trie
// ops. It is not used in Ethermint, so it returns nil.
func (csdb *CommitStateDB) Database() ethstate.Database {
	return csdb.db
}

// Empty returns whether the state object is either non-existent or empty
// according to the EIP161 specification (balance = nonce = code = 0).
func (csdb *CommitStateDB) Empty(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := "Empty"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	return so == nil || so.empty()
}

// Exist reports whether the given account address exists in the state. Notably,
// this also returns true for suicided accounts.
func (csdb *CommitStateDB) Exist(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := "Exist"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	return csdb.getStateObject(addr) != nil
}

// Error returns the first non-nil error the StateDB encountered.
func (csdb *CommitStateDB) Error() error {
	return csdb.dbErr
}

// Suicide marks the given account as suicided and clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (csdb *CommitStateDB) Suicide(addr ethcmn.Address) bool {
	if !csdb.ctx.IsCheckTx() {
		funcName := "Suicide"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so == nil {
		return false
	}

	csdb.journal.append(suicideChange{
		account:     &addr,
		prev:        so.suicided,
		prevBalance: sdk.NewDecFromBigIntWithPrec(so.Balance(), sdk.Precision), // int2dec
	})

	so.markSuicided()
	so.SetBalance(new(big.Int))

	return true
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

// ClearStateObjects clears cache of state objects to handle account changes outside of the EVM
func (csdb *CommitStateDB) ClearStateObjects() {
	csdb.stateObjects = make(map[ethcmn.Address]*stateObject)
	csdb.stateObjectsPending = make(map[ethcmn.Address]struct{})
	csdb.stateObjectsDirty = make(map[ethcmn.Address]struct{})
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
		funcName := "CreateAccount"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	newobj, prevobj := csdb.createObject(addr)
	if prevobj != nil {
		newobj.setBalance(sdk.DefaultBondDenom, sdk.NewDecFromBigIntWithPrec(prevobj.Balance(), sdk.Precision)) // int2dec
	}
}

// ForEachStorage iterates over each storage items, all invoke the provided
// callback on each key, value pair.
func (csdb *CommitStateDB) ForEachStorage(addr ethcmn.Address, cb func(key, value ethcmn.Hash) (stop bool)) error {
	if !csdb.ctx.IsCheckTx() {
		funcName := "ForEachStorage"
		trace.StartTxLog(funcName)
		defer trace.StopTxLog(funcName)
	}

	so := csdb.getStateObject(addr)
	if so == nil {
		return nil
	}

	return csdb.ForEachStorageMpt(so, cb)
}

// GetOrNewStateObject retrieves a state object or create a new state object if
// nil.
func (csdb *CommitStateDB) GetOrNewStateObject(addr ethcmn.Address) StateObject {
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

// SetError remembers the first non-nil error it is called with.
func (csdb *CommitStateDB) SetError(err error) {
	if err != nil {
		csdb.Logger().Debug("CommitStateDB", "error", err)
	}

	if csdb.dbErr == nil {
		csdb.dbErr = err
	}
}

// getStateObject attempts to retrieve a state object given by the address.
// Returns nil and sets an error if not found.
func (csdb *CommitStateDB) getStateObject(addr ethcmn.Address) (stateObject *stateObject) {
	if obj := csdb.getDeletedStateObject(addr); obj != nil && !obj.deleted {
		return obj
	}
	return nil
}

func (csdb *CommitStateDB) setStateObject(so *stateObject) {
	csdb.stateObjects[so.Address()] = so
}

// RawDump returns a raw state dump.
//
// TODO: Implement if we need it, especially for the RPC API.
func (csdb *CommitStateDB) RawDump() ethstate.Dump {
	return ethstate.Dump{}
}

type preimageEntry struct {
	// hash key of the preimage entry
	hash     ethcmn.Hash
	preimage []byte
}

func (csdb *CommitStateDB) SetLogSize(logSize uint) {
	csdb.logSize = logSize
}

func (csdb *CommitStateDB) GetLogSize() uint {
	return csdb.logSize
}

// SetContractDeploymentWhitelistMember sets the target address list into whitelist store
func (csdb *CommitStateDB) SetContractDeploymentWhitelist(addrList AddressList) {
	if csdb.ctx.GetWatcher().Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.ctx.GetWatcher().SaveContractDeploymentWhitelistItem(addrList[i])
		}
	}

	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	for i := 0; i < len(addrList); i++ {
		store.Set(GetContractDeploymentWhitelistMemberKey(addrList[i]), []byte(""))
	}
}

// DeleteContractDeploymentWhitelist deletes the target address list from whitelist store
func (csdb *CommitStateDB) DeleteContractDeploymentWhitelist(addrList AddressList) {
	if csdb.ctx.GetWatcher().Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.ctx.GetWatcher().DeleteContractDeploymentWhitelist(addrList[i])
		}
	}

	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	for i := 0; i < len(addrList); i++ {
		store.Delete(GetContractDeploymentWhitelistMemberKey(addrList[i]))
	}
}

// GetContractDeploymentWhitelist gets the whole contract deployment whitelist currently
func (csdb *CommitStateDB) GetContractDeploymentWhitelist() (whitelist AddressList) {
	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	iterator := sdk.KVStorePrefixIterator(store, KeyPrefixContractDeploymentWhitelist)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		whitelist = append(whitelist, splitApprovedDeployerAddress(iterator.Key()))
	}

	return
}

// IsDeployerInWhitelist checks whether the deployer is in the whitelist as a distributor
func (csdb *CommitStateDB) IsDeployerInWhitelist(deployerAddr sdk.AccAddress) bool {
	store := csdb.dbAdapter.NewStore(csdb.paramSpace.CustomKVStore(csdb.ctx), KeyPrefixContractDeploymentWhitelist)
	return store.Has(deployerAddr)
}

// SetContractBlockedList sets the target address list into blocked list store
func (csdb *CommitStateDB) SetContractBlockedList(addrList AddressList) {
	defer GetEvmParamsCache().SetNeedBlockedUpdate()
	if csdb.ctx.GetWatcher().Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.ctx.GetWatcher().SaveContractBlockedListItem(addrList[i])
		}
	}

	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	for i := 0; i < len(addrList); i++ {
		store.Set(GetContractBlockedListMemberKey(addrList[i]), []byte(""))
	}
}

// DeleteContractBlockedList deletes the target address list from blocked list store
func (csdb *CommitStateDB) DeleteContractBlockedList(addrList AddressList) {
	defer GetEvmParamsCache().SetNeedBlockedUpdate()
	if csdb.ctx.GetWatcher().Enabled() {
		for i := 0; i < len(addrList); i++ {
			csdb.ctx.GetWatcher().DeleteContractBlockedList(addrList[i])
		}
	}

	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	for i := 0; i < len(addrList); i++ {
		store.Delete(GetContractBlockedListMemberKey(addrList[i]))
	}
}

// GetContractBlockedList gets the whole contract blocked list currently
func (csdb *CommitStateDB) GetContractBlockedList() (blockedList AddressList) {
	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	iterator := sdk.KVStorePrefixIterator(store, KeyPrefixContractBlockedList)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if len(iterator.Value()) == 0 {
			blockedList = append(blockedList, splitBlockedContractAddress(iterator.Key()))
		}
	}
	return
}

// IsContractInBlockedList checks whether the contract address is in the blocked list
func (csdb *CommitStateDB) IsContractInBlockedList(contractAddr sdk.AccAddress) bool {
	bc := csdb.GetContractMethodBlockedByAddress(contractAddr)
	if bc == nil {
		//contractAddr is not blocked
		return false
	}
	// check contractAddr whether block full-method and special-method
	return bc.IsAllMethodBlocked()
}

// GetContractMethodBlockedByAddress gets contract methods blocked by address
func (csdb *CommitStateDB) GetContractMethodBlockedByAddress(contractAddr sdk.AccAddress) *BlockedContract {
	if csdb.ctx.UseParamCache() {
		tempEnableCache := true
		if GetEvmParamsCache().IsNeedBlockedUpdate() {
			bcl := csdb.GetContractMethodBlockedList()
			GetEvmParamsCache().UpdateBlockedContractMethod(bcl, csdb.ctx.IsCheckTx())
			// Note: when checktx GetEvmParamsCache().UpdateBlockedContractMethod will not be really update, so we must find GetBlockedContract from db.
			if csdb.ctx.IsCheckTx() {
				tempEnableCache = false
			}
		}
		if tempEnableCache {
			return GetEvmParamsCache().GetBlockedContractMethod(amino.BytesToStr(contractAddr))
		}
	}

	//use dbAdapter for watchdb or prefixdb
	store := csdb.dbAdapter.NewStore(csdb.paramSpace.CustomKVStore(csdb.ctx), KeyPrefixContractBlockedList)

	if ok := store.Has(contractAddr); !ok {
		// address is not exist
		return nil
	} else {
		value := store.Get(contractAddr)
		methods := ContractMethods{}
		var bc *BlockedContract
		if len(value) == 0 {
			//address is exist,but the blocked is old version.
			bc = NewBlockContract(contractAddr, methods)
		} else {
			// get block contract from cache without anmio
			if contractMethodBlockedCache != nil {
				if cm, ok := contractMethodBlockedCache.GetContractMethod(value); ok {
					return NewBlockContract(contractAddr, cm)
				}
			}
			//address is exist,but the blocked is new version.
			csdb.cdc.MustUnmarshalJSON(value, &methods)
			bc = NewBlockContract(contractAddr, methods)

			// write block contract into cache
			if contractMethodBlockedCache != nil {
				contractMethodBlockedCache.SetContractMethod(value, methods)
			}
		}
		return bc
	}
}

// InsertContractMethodBlockedList sets the list of contract method blocked into blocked list store
func (csdb *CommitStateDB) InsertContractMethodBlockedList(contractList BlockedContractList) sdk.Error {
	defer GetEvmParamsCache().SetNeedBlockedUpdate()
	if err := contractList.ValidateExtra(); err != nil {
		return err
	}
	for i := 0; i < len(contractList); i++ {
		bc := csdb.GetContractMethodBlockedByAddress(contractList[i].Address)
		if bc != nil {
			result, err := bc.BlockMethods.InsertContractMethods(contractList[i].BlockMethods)
			if err != nil {
				return err
			}
			bc.BlockMethods = result
		} else {
			bc = &contractList[i]
		}

		csdb.SetContractMethodBlocked(*bc)
	}
	return nil
}

// DeleteContractMethodBlockedList delete the list of contract method blocked  from blocked list store
func (csdb *CommitStateDB) DeleteContractMethodBlockedList(contractList BlockedContractList) sdk.Error {
	defer GetEvmParamsCache().SetNeedBlockedUpdate()
	for i := 0; i < len(contractList); i++ {
		bc := csdb.GetContractMethodBlockedByAddress(contractList[i].Address)
		if bc != nil {
			result, err := bc.BlockMethods.DeleteContractMethodMap(contractList[i].BlockMethods)
			if err != nil {
				return ErrBlockedContractMethodIsNotExist(contractList[i].Address, err)
			}
			bc.BlockMethods = result
			//if block contract method delete empty then remove contract from blocklist.
			if len(bc.BlockMethods) == 0 {
				addressList := AddressList{}
				addressList = append(addressList, bc.Address)
				//in watchdb contract blocked and contract method blocked use same prefix
				//so delete contract method blocked is can use function of delete contract blocked
				csdb.DeleteContractBlockedList(addressList)
			} else {
				csdb.SetContractMethodBlocked(*bc)
			}
		} else {
			return ErrBlockedContractMethodIsNotExist(contractList[i].Address, ErrorContractMethodBlockedIsNotExist)
		}
	}
	return nil
}

// GetContractMethodBlockedList get the list of contract method blocked from blocked list store
func (csdb *CommitStateDB) GetContractMethodBlockedList() (blockedContractList BlockedContractList) {

	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	iterator := sdk.KVStorePrefixIterator(store, KeyPrefixContractBlockedList)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		addr := sdk.AccAddress(splitBlockedContractAddress(iterator.Key()))
		value := iterator.Value()
		methods := ContractMethods{}
		if len(value) != 0 {
			csdb.cdc.MustUnmarshalJSON(value, &methods)
		}
		bc := NewBlockContract(addr, methods)
		blockedContractList = append(blockedContractList, *bc)
	}
	return
}

// IsContractMethodBlocked checks whether the contract method is blocked
func (csdb *CommitStateDB) IsContractMethodBlocked(contractAddr sdk.AccAddress, method string) bool {
	bc := csdb.GetContractMethodBlockedByAddress(contractAddr)
	if bc == nil {
		//contractAddr is not blocked
		return false
	}
	// it maybe happens,because ok_verifier verify called before getCode, for executing old logic follow code is return false
	if bc.IsAllMethodBlocked() {
		return false
	}
	// check contractAddr whether block full-method and special-method
	return bc.IsMethodBlocked(method)
}

// SetContractMethodBlocked sets contract method blocked into blocked list store
func (csdb *CommitStateDB) SetContractMethodBlocked(contract BlockedContract) {
	SortContractMethods(contract.BlockMethods)
	value := csdb.cdc.MustMarshalJSON(contract.BlockMethods)
	value = sdk.MustSortJSON(value)
	if csdb.ctx.GetWatcher().Enabled() {
		csdb.ctx.GetWatcher().SaveContractMethodBlockedListItem(contract.Address, value)
	}

	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	key := GetContractBlockedListMemberKey(contract.Address)
	store.Set(key, value)
}

func (csdb *CommitStateDB) GetAccount(addr ethcmn.Address) *ethermint.EthAccount {
	obj := csdb.getStateObject(addr)
	if obj == nil {
		return nil
	}
	return obj.account
}

func (csdb *CommitStateDB) UpdateContractBytecode(ctx sdk.Context, p ManageContractByteCodeProposal) sdk.Error {
	contract := ethcmn.BytesToAddress(p.Contract)
	substituteContract := ethcmn.BytesToAddress(p.SubstituteContract)

	revertContractByteCode := p.Contract.String() == p.SubstituteContract.String()

	preCode := csdb.GetCode(contract)
	contractAcc := csdb.GetAccount(contract)
	if contractAcc == nil {
		return ErrNotContracAddress(fmt.Errorf("%s", contract.String()))
	}
	preCodeHash := contractAcc.CodeHash

	var newCodeHash []byte
	if revertContractByteCode {
		newCodeHash = csdb.getInitContractCodeHash(p.Contract)
		if len(newCodeHash) == 0 || bytes.Equal(preCodeHash, newCodeHash) {
			return ErrContractCodeNotBeenUpdated(contract.String())
		}
	} else {
		newCodeHash = csdb.GetCodeHash(substituteContract).Bytes()
	}

	newCode := csdb.GetCodeByHash(ethcmn.BytesToHash(newCodeHash))
	// update code
	csdb.SetCode(contract, newCode)

	// store init code
	csdb.storeInitContractCodeHash(p.Contract, preCodeHash)

	// commit state db
	csdb.Commit(false)
	return csdb.afterUpdateContractByteCode(ctx, contract, substituteContract, preCodeHash, preCode, newCode)
}

var (
	EventTypeContractUpdateByProposal = "contract-update-by-proposal"
)

func (csdb *CommitStateDB) afterUpdateContractByteCode(ctx sdk.Context, contract, substituteContract ethcmn.Address, preCodeHash, preCode, newCode []byte) error {
	contractAfterUpdateCode := csdb.GetAccount(contract)
	if contractAfterUpdateCode == nil {
		return ErrNotContracAddress(fmt.Errorf("%s", contractAfterUpdateCode.String()))
	}

	// log
	ctx.Logger().Info("updateContractByteCode", "contract", contract, "preCodeHash", hex.EncodeToString(preCodeHash), "preCodeSize", len(preCode),
		"codeHashAfterUpdateCode", hex.EncodeToString(contractAfterUpdateCode.CodeHash), "codeSizeAfterUpdateCode", len(newCode))
	// emit event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		EventTypeContractUpdateByProposal,
		sdk.NewAttribute("contract", contract.String()),
		sdk.NewAttribute("preCodeHash", hex.EncodeToString(preCodeHash)),
		sdk.NewAttribute("preCodeSize", fmt.Sprintf("%d", len(preCode))),
		sdk.NewAttribute("SubstituteContract", substituteContract.String()),
		sdk.NewAttribute("codeHashAfterUpdateCode", hex.EncodeToString(contractAfterUpdateCode.CodeHash)),
		sdk.NewAttribute("codeSizeAfterUpdateCode", fmt.Sprintf("%d", len(newCode))),
	))
	// update watcher
	csdb.WithContext(ctx).IteratorCode(func(addr ethcmn.Address, c CacheCode) bool {
		ctx.GetWatcher().SaveContractCode(addr, c.Code, uint64(ctx.BlockHeight()))
		ctx.GetWatcher().SaveContractCodeByHash(c.CodeHash, c.Code)
		ctx.GetWatcher().SaveAccount(contractAfterUpdateCode)
		return true
	})
	return nil
}

func (csdb *CommitStateDB) storeInitContractCodeHash(addr sdk.AccAddress, codeHash []byte) {
	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	key := GetInitContractCodeHashKey(addr)
	if !store.Has(key) {
		store.Set(key, codeHash)
	}
}

func (csdb *CommitStateDB) getInitContractCodeHash(addr sdk.AccAddress) []byte {
	store := csdb.paramSpace.CustomKVStore(csdb.ctx)
	key := GetInitContractCodeHashKey(addr)
	return store.Get(key)
}
