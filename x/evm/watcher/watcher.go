package watcher

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/okex/exchain/app/rpc/namespaces/eth/state"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/tendermint/abci/types"
	tmstate "github.com/okex/exchain/libs/tendermint/state"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
)

var itjs = jsoniter.ConfigCompatibleWithStandardLibrary

type Watcher struct {
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
	watchData *WatchData

	jobChan chan func()
}

var (
	watcherEnable  = false
	watcherLruSize = 1000
	checkWd        = false
	onceEnable     sync.Once
	onceLru        sync.Once
)

func IsWatcherEnabled() bool {
	onceEnable.Do(func() {
		watcherEnable = viper.GetBool(FlagFastQuery)
	})
	return watcherEnable
}

func GetWatchLruSize() int {
	onceLru.Do(func() {
		watcherLruSize = viper.GetInt(FlagFastQueryLru)
	})
	return watcherLruSize
}

func NewWatcher(logger log.Logger) *Watcher {
	watcher := &Watcher{store: InstanceOfWatchStore(), cumulativeGas: make(map[uint64]uint64), sw: IsWatcherEnabled(), firstUse: true, delayEraseKey: make([][]byte, 0), watchData: &WatchData{}, log: logger}
	checkWd = viper.GetBool(FlagCheckWd)
	return watcher
}

func (w *Watcher) IsFirstUse() bool {
	return w.firstUse
}

func (w *Watcher) Used() {
	w.firstUse = false
}

func (w *Watcher) Enabled() bool {
	return w.sw
}

func (w *Watcher) Enable(sw bool) {
	w.sw = sw
}

func (w *Watcher) NewHeight(height uint64, blockHash common.Hash, header types.Header) {
	if !w.Enabled() {
		return
	}
	w.batch = []WatchMessage{} // reset batch
	w.header = header
	w.height = height
	w.blockHash = blockHash
	w.cumulativeGas = make(map[uint64]uint64)
	w.gasUsed = 0
	w.blockTxs = []common.Hash{}

	// ResetTransferWatchData
	w.watchData = &WatchData{}
}

func (w *Watcher) SaveEthereumTx(msg evmtypes.MsgEthereumTx, txHash common.Hash, index uint64) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgEthTx(&msg, txHash, w.blockHash, w.height, index)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
	w.UpdateBlockTxs(txHash)
}

func (w *Watcher) SaveContractCode(addr common.Address, code []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCode(addr, code, w.height)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *Watcher) SaveContractCodeByHash(hash []byte, code []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCodeByHash(hash, code)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *Watcher) SaveTransactionReceipt(status uint32, msg evmtypes.MsgEthereumTx, txHash common.Hash, txIndex uint64, data *evmtypes.ResultData, gasUsed uint64) {
	if !w.Enabled() {
		return
	}
	w.UpdateCumulativeGas(txIndex, gasUsed)
	wMsg := NewMsgTransactionReceipt(status, &msg, txHash, w.blockHash, txIndex, w.height, data, w.cumulativeGas[txIndex], gasUsed)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) UpdateCumulativeGas(txIndex, gasUsed uint64) {
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

func (w *Watcher) UpdateBlockTxs(txHash common.Hash) {
	if !w.Enabled() {
		return
	}
	w.blockTxs = append(w.blockTxs, txHash)
}

func (w *Watcher) SaveAccount(account auth.Account, isDirectly bool) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgAccount(account)
	if wMsg != nil {
		if isDirectly {
			w.batch = append(w.batch, wMsg)
		} else {
			w.staleBatch = append(w.staleBatch, wMsg)
		}

	}
}

func (w *Watcher) AddDelAccMsg(account auth.Account, isDirectly bool) {
	if !w.Enabled() {
		return
	}
	wMsg := NewDelAccMsg(account)
	if wMsg != nil {
		if isDirectly {
			w.batch = append(w.batch, wMsg)
		} else {
			w.staleBatch = append(w.staleBatch, wMsg)
		}

	}
}

func (w *Watcher) DeleteAccount(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	w.store.Delete(GetMsgAccountKey(addr.Bytes()))
	key := append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...)
	w.delayEraseKey = append(w.delayEraseKey, key)
}

func (w *Watcher) AddDirtyAccount(addr *sdk.AccAddress) {
	w.watchData.DirtyAccount = append(w.watchData.DirtyAccount, addr)
}

func (w *Watcher) ExecuteDelayEraseKey() {
	if !w.Enabled() {
		return
	}
	if len(w.delayEraseKey) <= 0 {
		return
	}
	for _, k := range w.delayEraseKey {
		w.store.Delete(k)
	}
	w.delayEraseKey = make([][]byte, 0)
}

func (w *Watcher) SaveState(addr common.Address, key, value []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgState(addr, key, value)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *Watcher) SaveBlock(bloom ethtypes.Bloom) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgBlock(w.height, bloom, w.blockHash, w.header, uint64(0xffffffff), big.NewInt(int64(w.gasUsed)), w.blockTxs)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}

	wInfo := NewMsgBlockInfo(w.height, w.blockHash)
	if wInfo != nil {
		w.batch = append(w.batch, wInfo)
	}
	w.SaveLatestHeight(w.height)
}

func (w *Watcher) SaveLatestHeight(height uint64) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgLatestHeight(height)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveParams(params evmtypes.Params) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgParams(params)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveContractBlockedListItem(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgContractBlockedListItem(addr)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveContractMethodBlockedListItem(addr sdk.AccAddress, methods []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgContractMethodBlockedListItem(addr, methods)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveContractDeploymentWhitelistItem(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgContractDeploymentWhitelistItem(addr)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) DeleteContractBlockedList(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgContractBlockedListItem(addr)
	if wMsg != nil {
		key := wMsg.GetKey()
		w.store.Delete(key)
		w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *Watcher) DeleteContractDeploymentWhitelist(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgContractDeploymentWhitelistItem(addr)
	if wMsg != nil {
		key := wMsg.GetKey()
		w.store.Delete(key)
		w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *Watcher) Finalize() {
	if !w.Enabled() {
		return
	}
	w.batch = append(w.batch, w.staleBatch...)
	w.Reset()
}

func (w *Watcher) CommitStateToRpcDb(addr common.Address, key, value []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgState(addr, key, value)
	if wMsg != nil {
		w.store.Set(append(prefixRpcDb, wMsg.GetKey()...), []byte(wMsg.GetValue()))
	}
}

func (w *Watcher) CommitAccountToRpcDb(account auth.Account) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgAccount(account)
	if wMsg != nil {
		key := append(prefixRpcDb, wMsg.GetKey()...)
		w.store.Set(key, []byte(wMsg.GetValue()))
	}
}

func (w *Watcher) CommitCodeHashToDb(hash []byte, code []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCodeByHash(hash, code)
	if wMsg != nil {
		w.store.Set(wMsg.GetKey(), []byte(wMsg.GetValue()))
	}
}

func (w *Watcher) Reset() {
	if !w.Enabled() {
		return
	}
	w.staleBatch = []WatchMessage{}
}

func (w *Watcher) Commit() {
	if !w.Enabled() {
		return
	}
	//hold it in temp
	batch := w.batch
	w.dispatchJob(func() { w.commitBatch(batch) })

	// we dont do deduplicatie here,we do it in `commit routine`
	// get centerBatch for sending to DataCenter
	ddsBatch := make([]*Batch, len(batch))
	for i, b := range batch {
		ddsBatch[i] = &Batch{b.GetKey(), []byte(b.GetValue()), b.GetType()}
	}
	w.watchData.Batches = ddsBatch
}

func (w *Watcher) CommitWatchData(data WatchData) {
	if data.Size() == 0 {
		return
	}
	if data.Batches != nil {
		w.commitCenterBatch(data.Batches)
	}
	if data.DirtyAccount != nil {
		w.delDirtyAccount(data.DirtyAccount)
	}
	if data.DirtyList != nil {
		w.delDirtyList(data.DirtyList)
	}
	if data.BloomData != nil {
		w.commitBloomData(data.BloomData)
	}

	if checkWd {
		keys := make([][]byte, len(data.Batches))
		for i, _ := range data.Batches {
			keys[i] = data.Batches[i].Key
		}
		w.CheckWatchDB(keys, "consumer")
	}
}

func (w *Watcher) commitBatch(batch []WatchMessage) {
	filterMap := make(map[string]WatchMessage)
	for _, b := range batch {
		filterMap[bytes2Key(b.GetKey())] = b
	}

	for _, b := range filterMap {
		key := b.GetKey()
		value := []byte(b.GetValue())
		typeValue := b.GetType()
		if typeValue == TypeDelete {
			w.store.Delete(key)
		} else {
			w.store.Set(key, value)
			//need update params
			if bytes.Compare(key, prefixParams) == 0 {
				if msgParams, ok := b.(*MsgParams); ok {
					w.store.SetEvmParams(msgParams.Params)
				}
			}
			if typeValue == TypeState {
				state.SetStateToLru(common.BytesToHash(key), value)
			}
		}
	}

	if checkWd {
		keys := make([][]byte, len(batch))
		for i, _ := range batch {
			keys[i] = batch[i].GetKey()
		}
		w.CheckWatchDB(keys, "producer")
	}
}

func (w *Watcher) commitCenterBatch(batch []*Batch) {
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

func (w *Watcher) delDirtyAccount(accounts []*sdk.AccAddress) {
	for _, account := range accounts {
		w.store.Delete(GetMsgAccountKey(account.Bytes()))
	}
}

func (w *Watcher) delDirtyList(list [][]byte) {
	for _, key := range list {
		w.store.Delete(key)
	}
}

func (w *Watcher) commitBloomData(bloomData []*evmtypes.KV) {
	db := evmtypes.GetIndexer().GetDB()
	for _, bd := range bloomData {
		db.Set(bd.Key, bd.Value)
	}
}

func (w *Watcher) GetWatchDataFunc() func() ([]byte, error) {
	value := w.watchData
	value.DelayEraseKey = w.delayEraseKey

	return func() ([]byte, error) {
		filterWatcher := filterCopy(value)
		valueByte, err := filterWatcher.MarshalToAmino(nil)
		if err != nil {
			return nil, err
		}
		return valueByte, nil
	}
}

func (w *Watcher) UnmarshalWatchData(wdByte []byte) (interface{}, error) {
	if len(wdByte) == 0 {
		return nil, fmt.Errorf("failed unmarshal watch data: empty data")
	}
	wd := WatchData{}
	if err := wd.UnmarshalFromAmino(nil, wdByte); err != nil {
		return nil, err
	}
	return wd, nil
}

func (w *Watcher) UseWatchData(watchData interface{}) {
	wd, ok := watchData.(WatchData)
	if !ok {
		panic("use watch data failed")
	}
	w.delayEraseKey = wd.DelayEraseKey

	w.dispatchJob(func() { w.CommitWatchData(wd) })
}

func (w *Watcher) SetWatchDataFunc() {
	go w.jobRoutine()
	tmstate.SetWatchDataFunc(w.GetWatchDataFunc, w.UnmarshalWatchData, w.UseWatchData)
}

func (w *Watcher) GetBloomDataPoint() *[]*evmtypes.KV {
	return &w.watchData.BloomData
}

func (w *Watcher) CheckWatchDB(keys [][]byte, mode string) {
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

func bytes2Key(keyBytes []byte) string {
	return string(keyBytes)
}

func key2Bytes(key string) []byte {
	return []byte(key)
}

func filterCopy(origin *WatchData) *WatchData {
	return &WatchData{
		DirtyAccount:  filterAccount(origin.DirtyAccount),
		Batches:       filterBatch(origin.Batches),
		DelayEraseKey: filterDelayEraseKey(origin.DelayEraseKey),
		BloomData:     filterBloomData(origin.BloomData),
		DirtyList:     filterDirtyList(origin.DirtyList),
	}
}

func filterAccount(accounts []*sdk.AccAddress) []*sdk.AccAddress {
	if len(accounts) == 0 {
		return nil
	}

	filterAccountMap := make(map[string]*sdk.AccAddress)
	for _, account := range accounts {
		filterAccountMap[bytes2Key(account.Bytes())] = account
	}

	ret := make([]*sdk.AccAddress, len(filterAccountMap))
	i := 0
	for _, acc := range filterAccountMap {
		ret[i] = acc
		i++
	}

	return ret
}

func filterBatch(datas []*Batch) []*Batch {
	if len(datas) == 0 {
		return nil
	}

	filterBatch := make(map[string]*Batch)
	for _, b := range datas {
		filterBatch[bytes2Key(b.Key)] = b
	}

	ret := make([]*Batch, len(filterBatch))
	i := 0
	for _, b := range filterBatch {
		ret[i] = b
		i++
	}

	return ret
}

func filterDelayEraseKey(datas [][]byte) [][]byte {
	if len(datas) == 0 {
		return nil
	}

	filterDelayEraseKey := make(map[string][]byte, 0)
	for _, b := range datas {
		filterDelayEraseKey[bytes2Key(b)] = b
	}

	ret := make([][]byte, len(filterDelayEraseKey))
	i := 0
	for _, k := range filterDelayEraseKey {
		ret[i] = k
		i++
	}

	return ret
}
func filterBloomData(datas []*evmtypes.KV) []*evmtypes.KV {
	if len(datas) == 0 {
		return nil
	}

	filterBloomData := make(map[string]*evmtypes.KV, 0)
	for _, k := range datas {
		filterBloomData[bytes2Key(k.Key)] = k
	}

	ret := make([]*evmtypes.KV, len(filterBloomData))
	i := 0
	for _, k := range filterBloomData {
		ret[i] = k
		i++
	}

	return ret
}

func filterDirtyList(datas [][]byte) [][]byte {
	if len(datas) == 0 {
		return nil
	}

	filterDirtyList := make(map[string][]byte, 0)
	for _, k := range datas {
		filterDirtyList[bytes2Key(k)] = k
	}

	ret := make([][]byte, len(filterDirtyList))
	i := 0
	for _, k := range filterDirtyList {
		ret[i] = k
		i++
	}

	return ret
}

/////////// job
func (w *Watcher) jobRoutine() {
	if !w.Enabled() {
		return
	}

	w.lazyInitialization()

	for job := range w.jobChan {
		job()
	}
}

func (w *Watcher) lazyInitialization() {
	// lazy initial:
	// now we will allocate chan memory
	// 5*2 means watcherCommitJob+commitBatchJob(just in case)
	w.jobChan = make(chan func(), 5*2)
}

func (w *Watcher) dispatchJob(f func()) {
	// if jobRoutine were too slow to write data  to disk
	// we have to wait
	// why: something wrong happened: such as db panic(disk maybe is full)(it should be the only reason)
	//								  UseWatchData were executed every 4 seoncds(block schedual)
	w.jobChan <- f
}
