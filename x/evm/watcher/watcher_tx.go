package watcher

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/app/rpc/namespaces/eth/state"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmstate "github.com/okex/exchain/libs/tendermint/state"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
)

type WatcherTx struct {
	store         *WatchStore
	height        uint64
	blockHash     common.Hash
	header        types.Header
	batch         []WatchMessage
	staleBatch    []WatchMessage
	cumulativeGas map[uint64]uint64
	gasUsed       uint64
	blockTxs      []common.Hash
	sw            bool
	firstUse      bool
	delayEraseKey [][]byte
	log           log.Logger
	// for state delta transfering in network
	watchData  *WatchData
	jobChan    chan func()
	evmTxIndex uint64
	checkWd    bool
	filterMap  map[string]WatchMessage
}

var (
	watcherEnableTx  = false
	watcherLruSizeTx = 1000
	onceEnableTx     sync.Once
	onceLruTx        sync.Once
)

func IsWatcherEnabledTx() bool {
	onceEnable.Do(func() {
		watcherEnable = viper.GetBool(FlagFastQuery)
	})
	return watcherEnable
}

func GetWatchLruSizeTx() int {
	onceLru.Do(func() {
		watcherLruSize = viper.GetInt(FlagFastQueryLru)
	})
	return watcherLruSize
}

func NewWatcherTx(logger log.Logger) *WatcherTx {
	return &WatcherTx{store: InstanceOfWatchStore(),
		cumulativeGas: make(map[uint64]uint64),
		sw:            IsWatcherEnabled(),
		firstUse:      true,
		delayEraseKey: make([][]byte, 0),
		watchData:     &WatchData{},
		log:           logger,
		checkWd:       viper.GetBool(FlagCheckWd),
		filterMap:     make(map[string]WatchMessage)}
}

func (w *WatcherTx) IsFirstUse() bool {
	return w.firstUse
}

// SetFirstUse sets fistUse of Watcher only could use for ut
func (w *WatcherTx) SetFirstUse(v bool) {
	w.firstUse = v
}

func (w *WatcherTx) Used() {
	w.firstUse = false
}

func (w *WatcherTx) Enabled() bool {
	return w.sw
}

func (w *WatcherTx) Enable(sw bool) {
	w.sw = sw
}

func (w *WatcherTx) GetEvmTxIndex() uint64 {
	return w.evmTxIndex
}

func (w *WatcherTx) NewHeight(height uint64, blockHash common.Hash, header types.Header) {
	if !w.Enabled() {
		return
	}
	w.header = header
	w.height = height
	w.blockHash = blockHash
	w.batch = []WatchMessage{} // reset batch
	// ResetTransferWatchData
	w.watchData = &WatchData{}
	w.evmTxIndex = 0
}

func (w *WatcherTx) clean() {
	for k := range w.cumulativeGas {
		delete(w.cumulativeGas, k)
	}
	w.gasUsed = 0
	w.blockTxs = []common.Hash{}
}

func (w *WatcherTx) SaveContractCode(addr common.Address, code []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCode(addr, code, w.height)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *WatcherTx) SaveContractCodeByHash(hash []byte, code []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCodeByHash(hash, code)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *WatcherTx) SaveTransactionReceipt(status uint32, msg interface{}, txHash common.Hash, txIndex uint64, data interface{}, gasUsed uint64) {
	if !w.Enabled() {
		return
	}
	realMsg, ok := msg.(*evmtypes.MsgEthereumTx)
	if !ok {
		return
	}
	realData, ok := data.(*evmtypes.ResultData)
	if !ok {
		return
	}
	w.UpdateCumulativeGas(txIndex, gasUsed)
	tr := newTransactionReceipt(status, realMsg, txHash, w.blockHash, txIndex, w.height, realData, w.cumulativeGas[txIndex], gasUsed)

	wMsg := NewMsgTransactionReceipt(tr, txHash)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *WatcherTx) UpdateCumulativeGas(txIndex, gasUsed uint64) {
	if !w.Enabled() {
		return
	}
	if len(w.cumulativeGas) == 0 {
		w.cumulativeGas[txIndex] = gasUsed
	} else {
		w.cumulativeGas[txIndex] = w.cumulativeGas[txIndex-1] + gasUsed
	}
	w.gasUsed += gasUsed
}

func (w *WatcherTx) SaveAccount(account interface{}, isDirectly bool) {
	if !w.Enabled() {
		return
	}
	acc, ok := account.(auth.Account)
	if !ok {
		return
	}
	wMsg := NewMsgAccount(acc)
	if wMsg != nil {
		if isDirectly {
			w.batch = append(w.batch, wMsg)
		} else {
			w.staleBatch = append(w.staleBatch, wMsg)
		}

	}
}

func (w *WatcherTx) AddDelAccMsg(account interface{}, isDirectly bool) {
	if !w.Enabled() {
		return
	}
	acc, ok := account.(auth.Account)
	if !ok {
		return
	}
	wMsg := NewDelAccMsg(acc)
	if wMsg != nil {
		if isDirectly {
			w.batch = append(w.batch, wMsg)
		} else {
			w.staleBatch = append(w.staleBatch, wMsg)
		}
	}
}

func (w *WatcherTx) DeleteAccount(addr interface{}) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	key1 := GetMsgAccountKey(realAddr.Bytes())
	key2 := append(prefixRpcDb, key1...)
	w.delayEraseKey = append(w.delayEraseKey, key1)
	w.delayEraseKey = append(w.delayEraseKey, key2)
}

func (w *WatcherTx) DelayEraseKey() {
	if !w.Enabled() {
		return
	}
	//hold it in temp
	delayEraseKey := w.delayEraseKey
	w.delayEraseKey = make([][]byte, 0)
	w.dispatchJob(func() {
		w.ExecuteDelayEraseKey(delayEraseKey)
	})
}

func (w *WatcherTx) ExecuteDelayEraseKey(delayEraseKey [][]byte) {
	if !w.Enabled() {
		return
	}
	if len(delayEraseKey) <= 0 {
		return
	}
	for _, k := range delayEraseKey {
		w.store.Delete(k)
	}
}

func (w *WatcherTx) SaveState(addr common.Address, key, value []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgState(addr, key, value)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *WatcherTx) SaveBlock(bloom ethtypes.Bloom) {
	if !w.Enabled() {
		return
	}
	block := newBlock(w.height, bloom, w.blockHash, w.header, uint64(0xffffffff), big.NewInt(int64(w.gasUsed)), w.blockTxs)
	wMsg := NewMsgBlock(block)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}

	wInfo := NewMsgBlockInfo(w.height, w.blockHash)
	if wInfo != nil {
		w.batch = append(w.batch, wInfo)
	}
	w.SaveLatestHeight(w.height)
}

func (w *WatcherTx) SaveLatestHeight(height uint64) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgLatestHeight(height)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *WatcherTx) SaveParams(params interface{}) {
	if !w.Enabled() {
		return
	}
	realParam, ok := params.(evmtypes.Params)
	if !ok {
		return
	}
	wMsg := NewMsgParams(realParam)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *WatcherTx) SaveContractBlockedListItem(addr interface{}) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractBlockedListItem(realAddr)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *WatcherTx) SaveContractMethodBlockedListItem(addr interface{}, methods []byte) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractMethodBlockedListItem(realAddr, methods)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *WatcherTx) SaveContractDeploymentWhitelistItem(addr interface{}) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractDeploymentWhitelistItem(realAddr)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *WatcherTx) DeleteContractBlockedList(addr interface{}) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractBlockedListItem(realAddr)
	if wMsg != nil {
		key := wMsg.GetKey()
		w.store.Delete(key)
		w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *WatcherTx) DeleteContractDeploymentWhitelist(addr interface{}) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractDeploymentWhitelistItem(realAddr)
	if wMsg != nil {
		key := wMsg.GetKey()
		w.store.Delete(key)
		w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *WatcherTx) Finalize() {
	if !w.Enabled() {
		return
	}
	w.batch = append(w.batch, w.staleBatch...)
	w.Reset()
}

func (w *WatcherTx) CommitStateToRpcDb(addr common.Address, key, value []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgState(addr, key, value)
	if wMsg != nil {
		w.store.Set(append(prefixRpcDb, wMsg.GetKey()...), []byte(wMsg.GetValue()))
	}
}

func (w *WatcherTx) CommitAccountToRpcDb(account auth.Account) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgAccount(account)
	if wMsg != nil {
		key := append(prefixRpcDb, wMsg.GetKey()...)
		w.store.Set(key, []byte(wMsg.GetValue()))
	}
}

func (w *WatcherTx) CommitCodeHashToDb(hash []byte, code []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCodeByHash(hash, code)
	if wMsg != nil {
		w.store.Set(wMsg.GetKey(), []byte(wMsg.GetValue()))
	}
}

func (w *WatcherTx) Reset() {
	if !w.Enabled() {
		return
	}
	w.staleBatch = []WatchMessage{}
}

func (w *WatcherTx) Commit() {
	if !w.Enabled() {
		return
	}
	//hold it in temp
	batch := w.batch
	w.clean()
	w.dispatchJob(func() {
		w.commitBatch(batch)
	})
}

func (w *WatcherTx) CommitWatchData(d interface{}) {
	data, ok := d.(WatchData)
	if !ok {
		return
	}
	if data.Size() == 0 {
		return
	}
	if data.Batches != nil {
		w.commitCenterBatch(data.Batches)
	}
	if data.DirtyList != nil {
		w.delDirtyList(data.DirtyList)
	}
	if data.BloomData != nil {
		w.commitBloomData(data.BloomData)
	}

	if w.checkWd {
		keys := make([][]byte, len(data.Batches))
		for i, _ := range data.Batches {
			keys[i] = data.Batches[i].Key
		}
		w.CheckWatchDB(keys, "consumer")
	}
}

func (w *WatcherTx) commitBatch(batch []WatchMessage) {
	for _, b := range batch {
		w.filterMap[bytes2Key(b.GetKey())] = b
	}

	for _, b := range w.filterMap {
		key := b.GetKey()
		value := []byte(b.GetValue())
		typeValue := b.GetType()
		if typeValue == TypeDelete {
			w.store.Delete(key)
		} else {
			w.store.Set(key, value)
			//need update params
			if typeValue == TypeEvmParams {
				msgParams := b.(*MsgParams)
				w.store.SetEvmParams(msgParams.Params)
			}
			if typeValue == TypeState {
				state.SetStateToLru(common.BytesToHash(key), value)
			}
		}
	}

	for k := range w.filterMap {
		delete(w.filterMap, k)
	}

	if w.checkWd {
		keys := make([][]byte, len(batch))
		for i, _ := range batch {
			keys[i] = batch[i].GetKey()
		}
		w.CheckWatchDB(keys, "producer")
	}
}

func (w *WatcherTx) commitCenterBatch(batch []*Batch) {
	for _, b := range batch {
		if b.TypeValue == TypeDelete {
			w.store.Delete(b.Key)
		} else {
			w.store.Set(b.Key, b.Value)
			if b.TypeValue == TypeState {
				state.SetStateToLru(common.BytesToHash(b.Key), b.Value)
			}
		}
	}
}

func (w *WatcherTx) delDirtyAccount(accounts []*sdk.AccAddress) {
	for _, account := range accounts {
		w.store.Delete(GetMsgAccountKey(account.Bytes()))
	}
}

func (w *WatcherTx) delDirtyList(list [][]byte) {
	for _, key := range list {
		w.store.Delete(key)
	}
}

func (w *WatcherTx) commitBloomData(bloomData []*evmtypes.KV) {
	db := evmtypes.GetIndexer().GetDB()
	for _, bd := range bloomData {
		db.Set(bd.Key, bd.Value)
	}
}

func (w *WatcherTx) GetWatchDataFunc() func() ([]byte, error) {
	value := w.watchData
	value.DelayEraseKey = w.delayEraseKey

	// hold it in temp
	batch := w.batch
	return func() ([]byte, error) {
		ddsBatch := make([]*Batch, len(batch))
		for i, b := range batch {
			ddsBatch[i] = &Batch{b.GetKey(), []byte(b.GetValue()), b.GetType()}
		}
		value.Batches = ddsBatch

		filterWatcher := filterCopy(value)
		valueByte, err := filterWatcher.MarshalToAmino(nil)
		if err != nil {
			return nil, err
		}
		return valueByte, nil
	}
}

func (w *WatcherTx) UnmarshalWatchData(wdByte []byte) (interface{}, error) {
	if len(wdByte) == 0 {
		return nil, fmt.Errorf("failed unmarshal watch data: empty data")
	}
	wd := WatchData{}
	if err := wd.UnmarshalFromAmino(nil, wdByte); err != nil {
		return nil, err
	}
	return wd, nil
}

func (w *WatcherTx) UseWatchData(watchData interface{}) {
	wd, ok := watchData.(WatchData)
	if !ok {
		panic("use watch data failed")
	}
	w.dispatchJob(func() { w.CommitWatchData(wd) })
}

func (w *WatcherTx) SetWatchDataFunc() {
	go w.jobRoutine()
	tmstate.SetWatchDataFunc(w.GetWatchDataFunc, w.UnmarshalWatchData, w.UseWatchData)
}

func (w *WatcherTx) GetBloomDataPoint() *[]*evmtypes.KV {
	return &w.watchData.BloomData
}

func (w *WatcherTx) CheckWatchDB(keys [][]byte, mode string) {
	output := make(map[string]string, len(keys))
	kvHash := tmhash.New()
	for _, key := range keys {
		value, err := w.store.Get(key)
		if err != nil {
			continue
		}
		kvHash.Write(key)
		kvHash.Write(value)
		output[hex.EncodeToString(key)] = string(value)
	}

	w.log.Info("watchDB delta", "mode", mode, "height", w.height, "hash", hex.EncodeToString(kvHash.Sum(nil)), "kv", output)
}

/////////// job
func (w *WatcherTx) jobRoutine() {
	if !w.Enabled() {
		return
	}

	w.lazyInitialization()
	for job := range w.jobChan {
		job()
	}
}

func (w *WatcherTx) lazyInitialization() {
	// lazy initial:
	// now we will allocate chan memory
	// 5*3 means watcherCommitJob+DelayEraseKey+commitBatchJob(just in case)
	w.jobChan = make(chan func(), 5*3)
}

func (w *WatcherTx) dispatchJob(f func()) {
	// if jobRoutine were too slow to write data  to disk
	// we have to wait
	// why: something wrong happened: such as db panic(disk maybe is full)(it should be the only reason)
	//								  UseWatchData were executed every 4 seoncds(block schedual)
	w.jobChan <- f
}

func (w *WatcherTx) Height() uint64 {
	return w.height
}
