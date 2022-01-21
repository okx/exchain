package watcher

import (
	"encoding/hex"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"math/big"
	"strings"
	"sync"

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
	store      *WatchStore
	height     uint64
	blockHash  common.Hash
	header     types.Header
	batch      []WatchMessage
	staleBatch []WatchMessage

	cumulativeGas map[uint64]uint64
	gasUsed       uint64
	blockTxs      []common.Hash
	sw            bool
	firstUse      bool
	delayEraseKey [][]byte
	log           log.Logger
	// for state delta transfering in network
	watchData *WatchData

	regionKeySet map[cacheNameSpace]map[string]int
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
	watcher := &Watcher{store: InstanceOfWatchStore(), cumulativeGas: make(map[uint64]uint64), sw: IsWatcherEnabled(), firstUse: true, delayEraseKey: make([][]byte, 0), watchData: &WatchData{}, log: logger, regionKeySet: make(map[cacheNameSpace]map[string]int)}
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
	w.regionKeySet = make(map[cacheNameSpace]map[string]int)

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
	if wMsg == nil {
		return
	}
	regionId := regionAccountDirectly
	if !isDirectly {
		regionId = regionAccountIndirectly
	}
	key := wMsg.GetKey()
	uniqueKey := buildKey(key)

	w.appendOrSwap(regionId, uniqueKey, func() int {
		index := -1
		if isDirectly {
			w.batch = append(w.batch, wMsg)
			index = len(w.batch) - 1
		} else {
			w.staleBatch = append(w.staleBatch, wMsg)
			index = len(w.staleBatch) - 1
		}
		return index
	}, func(index int) {
		if isDirectly {
			w.batch[index] = wMsg
		} else {
			w.staleBatch[index] = wMsg
		}
	})
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
	storeKey := GetMsgAccountKey(addr.Bytes())
	key := append(prefixRpcDb, storeKey...)
	uniqueKey := buildKey(storeKey)

	w.store.Delete(storeKey)
	w.appendOrSwap(regionDelayEraseKey, uniqueKey, func() int {
		w.delayEraseKey = append(w.delayEraseKey, key)
		return len(w.delayEraseKey) - 1
	}, func(index int) {})
}

func (w *Watcher) AddDirtyAccount(addr *sdk.AccAddress) {
	key := buildKey(GetMsgAccountKey(addr.Bytes()))
	w.appendOrSwap(regionDirtyAccount, key, func() int {
		w.watchData.DirtyAccount = append(w.watchData.DirtyAccount, addr)
		return len(w.watchData.DirtyAccount) - 1
	}, func(index int) {})
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

	w.flushUselessCacheRegion(regionAccountIndirectly)

	w.staleBatch = []WatchMessage{}
}

func (w *Watcher) flushUselessCacheRegion(regionIds ...cacheNameSpace) {
	for _, id := range regionIds {
		delete(w.regionKeySet, id)
	}
}

func (w *Watcher) Commit() {
	if !w.Enabled() {
		return
	}
	//hold it in temp
	batch := w.batch
	go w.commitBatch(batch)

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
	w.delayEraseKey = data.DelayEraseKey

	if checkWd {
		keys := make([][]byte, len(data.Batches))
		for i, _ := range data.Batches {
			keys[i] = data.Batches[i].Key
		}
		w.CheckWatchDB(keys, "consumer")
	}
}

func (w *Watcher) commitBatch(batch []WatchMessage) {
	for _, b := range batch {
		key := b.GetKey()
		value := []byte(b.GetValue())
		typeValue := b.GetType()
		if typeValue == TypeDelete {
			w.store.Delete(key)
		} else {
			w.store.Set(key, value)
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
		addr := *account
		w.store.Delete(GetMsgAccountKey(addr.Bytes()))
		key := append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...)
		w.delayEraseKey = append(w.delayEraseKey, key)
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
		valueByte, err := itjs.Marshal(value)
		if err != nil {
			return nil, err
		}
		return valueByte, nil
	}
}

func (w *Watcher) UseWatchData(wdByte []byte) {
	wd := WatchData{}
	if len(wdByte) > 0 {
		if err := itjs.Unmarshal(wdByte, &wd); err != nil {
			return
		}
	}

	go w.CommitWatchData(wd)
}

func (w *Watcher) SetWatchDataFunc() {
	tmstate.SetWatchDataFunc(w.GetWatchDataFunc, w.UseWatchData)
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

// f: if data is not exists, replace: means data already exists
// if panic: `index out of range` or  `map concurrent read write`: race condition
func (w *Watcher) appendOrSwap(regionId cacheNameSpace, key string, f func() int, replace func(index int)) {
	region := w.regionKeySet[regionId]
	if region == nil {
		region = make(map[string]int)
		w.regionKeySet[regionId] = region
	}
	originIndex, exist := region[key]
	if !exist {
		region[key] = f()
	} else {
		replace(originIndex)
	}
}

// helper[0]: prefix
// helper[1]...: data
func buildKey(helper ...interface{}) string {
	str := "%v-"
	if len(helper) > 1 {
		str += strings.Repeat("%v-", len(helper)-1)
	}
	return fmt.Sprintf(str, helper...)
}
