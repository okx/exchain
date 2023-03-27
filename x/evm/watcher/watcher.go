package watcher

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	jsoniter "github.com/json-iterator/go"
	"github.com/okx/okbchain/app/rpc/namespaces/eth/state"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/crypto/tmhash"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	ctypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
	tmstate "github.com/okx/okbchain/libs/tendermint/state"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

const version = "v1"

var itjs = jsoniter.ConfigCompatibleWithStandardLibrary

type Watcher struct {
	store          *WatchStore
	height         uint64
	blockHash      common.Hash
	header         types.Header
	batch          []WatchMessage
	cumulativeGas  map[uint64]uint64
	gasUsed        uint64
	blockTxs       []common.Hash
	blockStdTxs    []common.Hash
	enable         bool
	firstUse       bool
	delayEraseKey  [][]byte
	eraseKeyFilter map[string][]byte
	log            log.Logger
	// for state delta transfering in network
	watchData     *WatchData
	jobChan       chan func()
	jobDone       *sync.WaitGroup
	evmTxIndex    uint64
	checkWd       bool
	filterMap     map[string]struct{}
	InfuraKeeper  InfuraKeeper
	delAccountMtx sync.Mutex
}

var (
	watcherEnable  = false
	watcherLruSize = 1000
	onceEnable     sync.Once
	onceLru        sync.Once
	onceCheckRpc   sync.Once
	isRpcNode      = false
)

func IsRpcNode() bool {
	onceCheckRpc.Do(func() {
		if viper.GetString("tx_index.indexer") == "kv" {
			isRpcNode = true
		}
	})
	return isRpcNode
}

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
	return &Watcher{store: InstanceOfWatchStore(),
		cumulativeGas:  make(map[uint64]uint64),
		enable:         IsWatcherEnabled(),
		firstUse:       true,
		delayEraseKey:  make([][]byte, 0),
		watchData:      &WatchData{},
		log:            logger,
		checkWd:        viper.GetBool(FlagCheckWd),
		filterMap:      make(map[string]struct{}),
		eraseKeyFilter: make(map[string][]byte),
	}
}

func (w *Watcher) IsFirstUse() bool {
	return w.firstUse
}

// SetFirstUse sets fistUse of Watcher only could use for ut
func (w *Watcher) SetFirstUse(v bool) {
	w.firstUse = v
}

func (w *Watcher) Used() {
	w.firstUse = false
}

func (w *Watcher) Enabled() bool {
	return w.enable
}

func (w *Watcher) Enable(enable bool) {
	w.enable = enable
}

func (w *Watcher) GetEvmTxIndex() uint64 {
	return w.evmTxIndex
}

func (w *Watcher) NewHeight(height uint64, blockHash common.Hash, header types.Header) {
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
	for k := range w.cumulativeGas {
		delete(w.cumulativeGas, k)
	}
	w.gasUsed = 0
	w.blockTxs = []common.Hash{}
	w.blockStdTxs = []common.Hash{}
}

func (w *Watcher) SaveTransactionReceipt(status uint32, msg *evmtypes.MsgEthereumTx, txHash common.Hash, txIndex uint64, data *evmtypes.ResultData, gasUsed uint64) {
	if !w.Enabled() {
		return
	}
	w.UpdateCumulativeGas(txIndex, gasUsed)
	tr := newTransactionReceipt(status, msg, txHash, w.blockHash, txIndex, w.height, data, w.cumulativeGas[txIndex], gasUsed)
	if w.InfuraKeeper != nil {
		w.InfuraKeeper.OnSaveTransactionReceipt(tr)
	}
	wMsg := NewMsgTransactionReceipt(tr, txHash)
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

func (w *Watcher) SaveAccount(account auth.Account) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgAccount(account)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) DeleteAccount(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	key1 := GetMsgAccountKey(addr.Bytes())
	key2 := append(prefixRpcDb, key1...)
	w.delAccountMtx.Lock()
	w.delayEraseKey = append(w.delayEraseKey, key1)
	w.delayEraseKey = append(w.delayEraseKey, key2)
	w.delAccountMtx.Unlock()
}

func (w *Watcher) DelayEraseKey() {
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

func (w *Watcher) ExecuteDelayEraseKey(delayEraseKey [][]byte) {
	if !w.Enabled() || len(delayEraseKey) <= 0 {
		return
	}
	for _, k := range delayEraseKey {
		w.eraseKeyFilter[bytes2Key(k)] = k
	}
	batch := w.store.db.NewBatch()
	defer batch.Close()
	for _, k := range w.eraseKeyFilter {
		batch.Delete(k)
	}
	batch.Write()
	for k := range w.eraseKeyFilter {
		delete(w.eraseKeyFilter, k)
	}
}

func (w *Watcher) SaveBlock(block evmtypes.Block, ethBlockHash common.Hash) {
	if !w.Enabled() {
		return
	}
	if w.InfuraKeeper != nil {
		w.InfuraKeeper.OnSaveBlock(block)
	}
	wMsg := NewMsgBlock(block, ethBlockHash)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}

	wInfo := NewMsgBlockInfo(w.height, ethBlockHash)
	if wInfo != nil {
		w.batch = append(w.batch, wInfo)
	}
	w.SaveLatestHeight(w.height)
}

func (w *Watcher) SaveBlockStdTxHash(blockHash common.Hash) {
	if !w.Enabled() || (len(w.blockStdTxs) == 0) {
		return
	}
	wMsg := NewMsgBlockStdTxHash(w.blockStdTxs, blockHash)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
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

func (w *Watcher) Commit() {
	if !w.Enabled() {
		return
	}
	//hold it in temp
	batch := w.batch
	// No need to write db when upload delta is enabled.
	if tmtypes.UploadDelta {
		return
	}
	w.dispatchJob(func() {
		w.commitBatch(batch)
	})
}

func (w *Watcher) CommitWatchData(data WatchData) {
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
	w.delayEraseKey = data.DelayEraseKey

	if w.checkWd {
		keys := make([][]byte, len(data.Batches))
		for i, _ := range data.Batches {
			keys[i] = data.Batches[i].Key
		}
		w.CheckWatchDB(keys, "consumer")
	}
}
func isDuplicated(key []byte, filterMap map[string]struct{}) bool {
	filterKey := bytes2Key(key)
	if _, exist := filterMap[filterKey]; exist {
		return true
	} else {
		filterMap[filterKey] = struct{}{}
		return false
	}
}
func (w *Watcher) commitBatch(batch []WatchMessage) {
	dbBatch := w.store.db.NewBatch()
	defer dbBatch.Close()
	for i := len(batch) - 1; i >= 0; i-- { //iterate batch from the end to start, to save the latest batch msgs
		//and to skip the duplicated batch msgs by key
		b := batch[i]
		key := b.GetKey()
		if isDuplicated(key, w.filterMap) {
			continue
		}
		value := []byte(b.GetValue())
		typeValue := b.GetType()
		if typeValue == TypeDelete {
			dbBatch.Delete(key)
		} else {
			dbBatch.Set(key, value)
			//need update params
			if typeValue == TypeEvmParams {
				msgParams := b.(*MsgParams)
				w.store.SetEvmParams(msgParams.Params)
			}
			if typeValue == TypeState {
				state.SetStateToLru(key, value)
			}
		}
	}
	dbBatch.Write()
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

func (w *Watcher) commitCenterBatch(batch []*Batch) {
	dbBatch := w.store.db.NewBatch()
	defer dbBatch.Close()
	for _, b := range batch {
		if b.TypeValue == TypeDelete {
			dbBatch.Delete(b.Key)
		} else {
			dbBatch.Set(b.Key, b.Value)
			if b.TypeValue == TypeState {
				state.SetStateToLru(b.Key, b.Value)
			}
		}
	}
	dbBatch.Write()
}

func (w *Watcher) delDirtyList(list [][]byte) {
	for _, key := range list {
		w.store.Delete(key)
	}
}

func (w *Watcher) commitBloomData(bloomData []*evmtypes.KV) {
	db := evmtypes.GetIndexer().GetDB()
	batch := db.NewBatch()
	defer batch.Close()
	for _, bd := range bloomData {
		batch.Set(bd.Key, bd.Value)
	}
	batch.Write()
}

func (w *Watcher) CreateWatchDataGenerator() func() ([]byte, error) {
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

func (w *Watcher) ApplyWatchData(watchData interface{}) {
	wd, ok := watchData.(WatchData)
	if !ok {
		panic("use watch data failed")
	}
	w.dispatchJob(func() { w.CommitWatchData(wd) })
}

func (w *Watcher) SetWatchDataManager() {
	go w.jobRoutine()
	tmstate.SetEvmWatchDataManager(w)
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
	return amino.BytesToStr(keyBytes)
}

func filterCopy(origin *WatchData) *WatchData {
	return &WatchData{
		Batches:       filterBatch(origin.Batches),
		DelayEraseKey: filterDelayEraseKey(origin.DelayEraseKey),
		BloomData:     filterBloomData(origin.BloomData),
		DirtyList:     filterDirtyList(origin.DirtyList),
	}
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

// ///////// job
func (w *Watcher) jobRoutine() {
	if !w.Enabled() {
		return
	}

	w.lazyInitialization()
	for job := range w.jobChan {
		job()
	}
	w.jobDone.Done()
}

func (w *Watcher) lazyInitialization() {
	// lazy initial:
	// now we will allocate chan memory
	// 5*3 means watcherCommitJob+DelayEraseKey+commitBatchJob(just in case)
	w.jobChan = make(chan func(), 5*3)
	w.jobDone = new(sync.WaitGroup)
	w.jobDone.Add(1)
}
func (w *Watcher) Stop() {
	if !w.Enabled() {
		return
	}
	close(w.jobChan)
	w.jobDone.Wait()
}
func (w *Watcher) dispatchJob(f func()) {
	// if jobRoutine were too slow to write data  to disk
	// we have to wait
	// why: something wrong happened: such as db panic(disk maybe is full)(it should be the only reason)
	//								  ApplyWatchData were executed every 4 seoncds(block schedual)
	w.jobChan <- f
}

func (w *Watcher) Height() uint64 {
	return w.height
}

func (w *Watcher) Collect(watchers ...sdk.IWatcher) {
	if !w.enable {
		return
	}
	for _, watcher := range watchers {
		batch := watcher.Destruct()
		w.batch = append(w.batch, batch...)
	}
}

func (w *Watcher) saveStdTxResponse(result *ctypes.ResultTx) {
	wMsg := NewStdTransactionResponse(result, w.header.Time, common.BytesToHash(result.Hash))
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}
