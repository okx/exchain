package watcher

import (
	"encoding/hex"
	jsoniter "github.com/json-iterator/go"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"math/big"
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
	store         *WatchStore
	height        uint64
	blockHash     common.Hash
	header        types.Header
	batch         []WatchMessage
	staleBatch    []WatchMessage
	cumulativeGas map[uint64]uint64
	gasUsed       uint64
	blockTxs      []common.Hash
	//sw            bool
	firstUse      bool
	delayEraseKey [][]byte
	log           log.Logger
	// for state delta transfering in network
	watchData *WatchData
}

var (
	watcherEnable  = false
	AsyncTxEnable=false
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

func newWatcher(log log.Logger) *Watcher {
	watcher := &Watcher{store: InstanceOfWatchStore(), cumulativeGas: make(map[uint64]uint64),  firstUse: true, delayEraseKey: make([][]byte, 0), watchData: &WatchData{},log: log}
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
	return true
}

func (w *Watcher) Enable(sw bool) {
	// do nothing for now
	// we will enable watcher  dynamically
}

func (w *Watcher) NewHeight(height uint64, blockHash common.Hash, header types.Header) {
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
	wMsg := NewMsgEthTx(&msg, txHash, w.blockHash, w.height, index)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
	w.UpdateBlockTxs(txHash)
}

func (w *Watcher) SaveContractCode(addr common.Address, code []byte) {
	wMsg := NewMsgCode(addr, code, w.height)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *Watcher) SaveContractCodeByHash(hash []byte, code []byte) {
	wMsg := NewMsgCodeByHash(hash, code)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *Watcher) SaveTransactionReceipt(status uint32, msg evmtypes.MsgEthereumTx, txHash common.Hash, txIndex uint64, data *evmtypes.ResultData, gasUsed uint64) {
	w.UpdateCumulativeGas(txIndex, gasUsed)
	wMsg := NewMsgTransactionReceipt(status, &msg, txHash, w.blockHash, txIndex, w.height, data, w.cumulativeGas[txIndex], gasUsed)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) UpdateCumulativeGas(txIndex, gasUsed uint64) {
	if len(w.cumulativeGas) == 0 {
		w.cumulativeGas[txIndex] = gasUsed
	} else {
		w.cumulativeGas[txIndex] = w.cumulativeGas[txIndex-1] + gasUsed
	}
	w.gasUsed += gasUsed
}

func (w *Watcher) UpdateBlockTxs(txHash common.Hash) {
	w.blockTxs = append(w.blockTxs, txHash)
}

func (w *Watcher) SaveAccount(account auth.Account, isDirectly bool) {
	wMsg := NewMsgAccount(account)
	if wMsg != nil {
		if isDirectly {
			w.batch = append(w.batch, wMsg)
		} else {
			w.staleBatch = append(w.staleBatch, wMsg)
		}

	}
}

func (w *Watcher) DeleteAccount(addr sdk.AccAddress) {
	w.store.Delete(GetMsgAccountKey(addr.Bytes()))
	key := append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...)
	w.delayEraseKey = append(w.delayEraseKey, key)
}

func (w *Watcher) AddDirtyAccount(addr *sdk.AccAddress) {
	w.watchData.DirtyAccount = append(w.watchData.DirtyAccount, addr)
}

func (w *Watcher) ExecuteDelayEraseKey() {
	if len(w.delayEraseKey) <= 0 {
		return
	}
	for _, k := range w.delayEraseKey {
		w.store.Delete(k)
	}
	w.delayEraseKey = make([][]byte, 0)
}

func (w *Watcher) SaveState(addr common.Address, key, value []byte) {
	wMsg := NewMsgState(addr, key, value)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *Watcher) SaveBlock(bloom ethtypes.Bloom) {
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
	wMsg := NewMsgLatestHeight(height)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveParams(params evmtypes.Params) {
	wMsg := NewMsgParams(params)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveContractBlockedListItem(addr sdk.AccAddress) {
	wMsg := NewMsgContractBlockedListItem(addr)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveContractMethodBlockedListItem(addr sdk.AccAddress, methods []byte) {
	wMsg := NewMsgContractMethodBlockedListItem(addr, methods)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveContractDeploymentWhitelistItem(addr sdk.AccAddress) {
	wMsg := NewMsgContractDeploymentWhitelistItem(addr)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) DeleteContractBlockedList(addr sdk.AccAddress) {
	wMsg := NewMsgContractBlockedListItem(addr)
	if wMsg != nil {
		key := wMsg.GetKey()
		w.store.Delete(key)
		w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *Watcher) DeleteContractDeploymentWhitelist(addr sdk.AccAddress) {
	wMsg := NewMsgContractDeploymentWhitelistItem(addr)
	if wMsg != nil {
		key := wMsg.GetKey()
		w.store.Delete(key)
		w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *Watcher) Finalize() {
	w.batch = append(w.batch, w.staleBatch...)
	w.Reset()
}

func (w *Watcher) CommitStateToRpcDb(addr common.Address, key, value []byte) {
	wMsg := NewMsgState(addr, key, value)
	if wMsg != nil {
		w.store.Set(append(prefixRpcDb, wMsg.GetKey()...), []byte(wMsg.GetValue()))
	}
}

func (w *Watcher) CommitAccountToRpcDb(account auth.Account) {
	wMsg := NewMsgAccount(account)
	if wMsg != nil {
		key := append(prefixRpcDb, wMsg.GetKey()...)
		w.store.Set(key, []byte(wMsg.GetValue()))
	}
}

func (w *Watcher) CommitCodeHashToDb(hash []byte, code []byte) {
	wMsg := NewMsgCodeByHash(hash, code)
	if wMsg != nil {
		w.store.Set(wMsg.GetKey(), []byte(wMsg.GetValue()))
	}
}

func (w *Watcher) Reset() {
	w.staleBatch = []WatchMessage{}
}

func (w *Watcher) Commit() {
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
		w.store.Set(key, value)
		if typeValue == TypeState {
			state.SetStateToLru(common.BytesToHash(key), value)
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
		w.store.Set(b.Key, b.Value)
		if b.TypeValue == TypeState {
			state.SetStateToLru(common.BytesToHash(b.Key), b.Value)
		}
	}
}

func (w *Watcher) delDirtyAccount(accounts []*sdk.AccAddress) {
	for _, account := range accounts {
		w.DeleteAccount(*account)
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

func (w *Watcher) GetWatchData() ([]byte, error) {
	value := w.watchData
	value.DelayEraseKey = w.delayEraseKey
	valueByte, err := itjs.Marshal(value)
	if err != nil {
		return nil, err
	}
	return valueByte, nil
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
	tmstate.SetWatchDataFunc(w.GetWatchData, w.UseWatchData)
}

func (w *Watcher) GetBloomDataPoint() *[]*evmtypes.KV {
	return &w.watchData.BloomData
}


func(w *Watcher)Init()error{
	w.SetWatchDataFunc()
	return nil
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
