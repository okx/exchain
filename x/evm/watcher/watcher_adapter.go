package watcher

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/tendermint/abci/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

var (
	_ IWatcher = (*Watcher)(nil)
	_ IWatcher = (*concurrentWatcher)(nil)
	_ IWatcher = (*disableWatcher)(nil)
)

type IWatcher interface {
	IsFirstUse() bool
	Used()
	Enabled() bool
	Enable(sw bool)
	NewHeight(height uint64, blockHash common.Hash, header types.Header)
	SaveEthereumTx(msg evmtypes.MsgEthereumTx, txHash common.Hash, index uint64)
	SaveContractCode(addr common.Address, code []byte)
	SaveContractCodeByHash(hash []byte, code []byte)
	SaveTransactionReceipt(status uint32, msg evmtypes.MsgEthereumTx, txHash common.Hash, txIndex uint64, data *evmtypes.ResultData, gasUsed uint64)
	UpdateCumulativeGas(txIndex, gasUsed uint64)
	UpdateBlockTxs(txHash common.Hash)
	SaveAccount(account auth.Account, isDirectly bool)
	DeleteAccount(addr sdk.AccAddress)
	AddDirtyAccount(addr *sdk.AccAddress)
	ExecuteDelayEraseKey()
	SaveState(addr common.Address, key, value []byte)
	SaveBlock(bloom ethtypes.Bloom)
	SaveLatestHeight(height uint64)
	SaveParams(params evmtypes.Params)
	SaveContractBlockedListItem(addr sdk.AccAddress)
	SaveContractMethodBlockedListItem(addr sdk.AccAddress, methods []byte)
	SaveContractDeploymentWhitelistItem(addr sdk.AccAddress)
	DeleteContractBlockedList(addr sdk.AccAddress)
	DeleteContractDeploymentWhitelist(addr sdk.AccAddress)
	Finalize()
	CommitStateToRpcDb(addr common.Address, key, value []byte)
	CommitAccountToRpcDb(account auth.Account)
	CommitCodeHashToDb(hash []byte, code []byte)
	Reset()
	Commit()
	CommitWatchData()
	GetWatchData() ([]byte, error)
	UseWatchData(wdByte []byte)

	GetBloomDataPoint() *[]*evmtypes.KV
	Init() error
}

type concurrentWatcher struct {
	mtx *sync.RWMutex
	w   *Watcher
}

func newConcurrentWatcher(w *Watcher) *concurrentWatcher {
	ret := &concurrentWatcher{mtx: &sync.RWMutex{}, w: w}

	return ret
}

func (c *concurrentWatcher) IsFirstUse() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.w.IsFirstUse()
}

func (c *concurrentWatcher) Used() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.Used()
}

func (c *concurrentWatcher) Enabled() bool {
	return true
}

func (c *concurrentWatcher) Enable(sw bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.Enable(sw)
}

func (c *concurrentWatcher) NewHeight(height uint64, blockHash common.Hash, header types.Header) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.NewHeight(height, blockHash, header)
}

func (c *concurrentWatcher) SaveEthereumTx(msg evmtypes.MsgEthereumTx, txHash common.Hash, index uint64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveEthereumTx(msg, txHash, index)
}

func (c *concurrentWatcher) SaveContractCode(addr common.Address, code []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveContractCode(addr, code)
}

func (c *concurrentWatcher) SaveContractCodeByHash(hash []byte, code []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveContractCodeByHash(hash, code)
}

func (c *concurrentWatcher) SaveTransactionReceipt(status uint32, msg evmtypes.MsgEthereumTx, txHash common.Hash, txIndex uint64, data *evmtypes.ResultData, gasUsed uint64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveTransactionReceipt(status, msg, txHash, txIndex, data, gasUsed)
}

func (c *concurrentWatcher) UpdateCumulativeGas(txIndex, gasUsed uint64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.UpdateCumulativeGas(txIndex, gasUsed)
}

func (c *concurrentWatcher) UpdateBlockTxs(txHash common.Hash) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.UpdateBlockTxs(txHash)
}

func (c *concurrentWatcher) SaveAccount(account auth.Account, isDirectly bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveAccount(account, isDirectly)
}

func (c *concurrentWatcher) DeleteAccount(addr sdk.AccAddress) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.DeleteAccount(addr)
}

func (c *concurrentWatcher) AddDirtyAccount(addr *sdk.AccAddress) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.AddDirtyAccount(addr)
}

func (c *concurrentWatcher) ExecuteDelayEraseKey() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.ExecuteDelayEraseKey()
}

func (c *concurrentWatcher) SaveState(addr common.Address, key, value []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveState(addr, key, value)
}

func (c *concurrentWatcher) SaveBlock(bloom ethtypes.Bloom) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveBlock(bloom)
}

func (c *concurrentWatcher) SaveLatestHeight(height uint64) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveLatestHeight(height)
}

func (c *concurrentWatcher) SaveParams(params evmtypes.Params) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveParams(params)
}

func (c *concurrentWatcher) SaveContractBlockedListItem(addr sdk.AccAddress) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveContractBlockedListItem(addr)
}

func (c *concurrentWatcher) SaveContractMethodBlockedListItem(addr sdk.AccAddress, methods []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveContractMethodBlockedListItem(addr, methods)
}

func (c *concurrentWatcher) SaveContractDeploymentWhitelistItem(addr sdk.AccAddress) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.SaveContractDeploymentWhitelistItem(addr)
}

func (c *concurrentWatcher) DeleteContractBlockedList(addr sdk.AccAddress) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.DeleteContractBlockedList(addr)
}

func (c *concurrentWatcher) DeleteContractDeploymentWhitelist(addr sdk.AccAddress) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.DeleteContractDeploymentWhitelist(addr)
}

func (c *concurrentWatcher) Finalize() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.Finalize()
}

func (c *concurrentWatcher) CommitStateToRpcDb(addr common.Address, key, value []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.CommitStateToRpcDb(addr, key, value)
}

func (c *concurrentWatcher) CommitAccountToRpcDb(account auth.Account) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.CommitAccountToRpcDb(account)
}

func (c *concurrentWatcher) CommitCodeHashToDb(hash []byte, code []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.CommitCodeHashToDb(hash, code)
}

func (c *concurrentWatcher) Reset() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.Reset()
}

func (c *concurrentWatcher) Commit() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.Commit()
}

func (c *concurrentWatcher) CommitWatchData() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.CommitWatchData()
}

func (c *concurrentWatcher) GetWatchData() ([]byte, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.w.GetWatchData()
}

func (c *concurrentWatcher) UseWatchData(wdByte []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.w.UseWatchData(wdByte)
}

func (c *concurrentWatcher) GetBloomDataPoint() *[]*evmtypes.KV {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.w.GetBloomDataPoint()
}

func (c *concurrentWatcher) Init() error {
	return c.w.Init()
}

//==============
type disableWatcher struct {
	w *Watcher
}

func newDisableWatcher(w *Watcher) *disableWatcher {
	ret := &disableWatcher{w: w}
	return ret
}

func (d *disableWatcher) IsFirstUse() bool {
	return d.w.IsFirstUse()
}

func (d *disableWatcher) Used() {
	d.w.Used()
}

func (d *disableWatcher) Enabled() bool {
	return false
}

func (d *disableWatcher) Enable(sw bool) {
	panic("not supported")
}

func (d *disableWatcher) NewHeight(height uint64, blockHash common.Hash, header types.Header) {}

func (d *disableWatcher) SaveEthereumTx(msg evmtypes.MsgEthereumTx, txHash common.Hash, index uint64) {
}

func (d *disableWatcher) SaveContractCode(addr common.Address, code []byte) {}

func (d *disableWatcher) SaveContractCodeByHash(hash []byte, code []byte) {}

func (d *disableWatcher) SaveTransactionReceipt(status uint32, msg evmtypes.MsgEthereumTx, txHash common.Hash, txIndex uint64, data *evmtypes.ResultData, gasUsed uint64) {
}

func (d *disableWatcher) UpdateCumulativeGas(txIndex, gasUsed uint64) {}

func (d *disableWatcher) UpdateBlockTxs(txHash common.Hash) {}

func (d *disableWatcher) SaveAccount(account auth.Account, isDirectly bool) {}

func (d *disableWatcher) DeleteAccount(addr sdk.AccAddress) {}

func (d *disableWatcher) AddDirtyAccount(addr *sdk.AccAddress) {}

func (d *disableWatcher) ExecuteDelayEraseKey() {}

func (d *disableWatcher) SaveState(addr common.Address, key, value []byte) {}

func (d *disableWatcher) SaveBlock(bloom ethtypes.Bloom) {}

func (d *disableWatcher) SaveLatestHeight(height uint64) {}

func (d *disableWatcher) SaveParams(params evmtypes.Params) {}

func (d *disableWatcher) SaveContractBlockedListItem(addr sdk.AccAddress) {}

func (d *disableWatcher) SaveContractMethodBlockedListItem(addr sdk.AccAddress, methods []byte) {}

func (d *disableWatcher) SaveContractDeploymentWhitelistItem(addr sdk.AccAddress) {}

func (d *disableWatcher) DeleteContractBlockedList(addr sdk.AccAddress) {}

func (d *disableWatcher) DeleteContractDeploymentWhitelist(addr sdk.AccAddress) {}

func (d *disableWatcher) Finalize() {}

func (d *disableWatcher) CommitStateToRpcDb(addr common.Address, key, value []byte) {}

func (d *disableWatcher) CommitAccountToRpcDb(account auth.Account) {}

func (d *disableWatcher) CommitCodeHashToDb(hash []byte, code []byte) {}

func (d *disableWatcher) Reset() {}

func (d *disableWatcher) Commit() {}

func (d *disableWatcher) CommitWatchData() {}

func (d *disableWatcher) GetWatchData() ([]byte, error) { return d.w.GetWatchData() }

func (d *disableWatcher) UseWatchData(wdByte []byte) { d.w.UseWatchData(wdByte) }

func (d *disableWatcher) GetBloomDataPoint() *[]*evmtypes.KV { return d.w.GetBloomDataPoint() }

func (d *disableWatcher) Init() error { return d.w.Init() }
