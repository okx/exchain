package watcher

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"

	"github.com/okex/exchain/app/rpc/namespaces/eth/state"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/abci/types"
	tmstate "github.com/tendermint/tendermint/state"
	tmtypes "github.com/tendermint/tendermint/types"
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
	// dirtyAccount, centerBatch send to DataCenter
	dirtyAccount []*sdk.AccAddress
	centerBatch  []*Batch
	needDelAcc   bool
}

var (
	watcherEnable  = false
	watcherLruSize = 1000
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

func NewWatcher() *Watcher {
	watcher := &Watcher{store: InstanceOfWatchStore(), sw: IsWatcherEnabled(), firstUse: true, delayEraseKey: make([][]byte, 0)}
	watcher.SetGetBatchFunc()
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
	w.batch = []WatchMessage{}
	w.header = header
	w.height = height
	w.blockHash = blockHash
	w.cumulativeGas = make(map[uint64]uint64)
	w.gasUsed = 0
	w.blockTxs = []common.Hash{}
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

func (w *Watcher) DeleteAccount(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	w.store.Delete(GetMsgAccountKey(addr.Bytes()))
	key := append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...)
	w.delayEraseKey = append(w.delayEraseKey, key)
	w.needDelAcc = false
}

func (w *Watcher) AddDirtyAccount(addr *sdk.AccAddress) {
	if w.dirtyAccount == nil {
		w.dirtyAccount = []*sdk.AccAddress{}
	}
	w.dirtyAccount = append(w.dirtyAccount, addr)
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
		w.store.Delete(wMsg.GetKey())
	}
}

func (w *Watcher) DeleteContractDeploymentWhitelist(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgContractDeploymentWhitelistItem(addr)
	if wMsg != nil {
		w.store.Delete(wMsg.GetKey())
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
	go w.commitBatch(w.batch)

	// if get data from DataCenter
	if w.centerBatch != nil {
		go w.commitCenterBatch(w.centerBatch)
		w.centerBatch = nil
	}
	if w.needDelAcc && w.dirtyAccount != nil {
		go w.delDirtyAccount(w.dirtyAccount)
		w.dirtyAccount = nil
		w.needDelAcc = false
	}

	// get centerBatch for sending to DataCenter
	centerBatch := make([]*Batch, len(batch))
	for i, b := range batch {
		centerBatch[i] = &Batch{b.GetKey(), []byte(b.GetValue()), b.GetType()}
	}
	w.centerBatch = centerBatch
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

func (w *Watcher) SetGetBatchFunc() {
	gb := func(height int64) bool {
		msg := tmstate.DataCenterMsg{Height: height}
		msgBody, err := tmtypes.Json.Marshal(&msg)
		if  err != nil {
			return false
		}
		response, err := http.Post(viper.GetString(tmtypes.DataCenterUrl)+"loadBatch", "application/json", bytes.NewBuffer([]byte(msgBody)))
		if err != nil {
			return false
		}
		defer response.Body.Close()
		rlt, _ := ioutil.ReadAll(response.Body)
		data := WatchData{}
		if err = itjs.Unmarshal(rlt, &data); err != nil {
			return false
		}
		if data.Account == nil || data.Batches == nil {
			return false
		}
		w.centerBatch = data.Batches
		w.dirtyAccount = data.Account
		w.needDelAcc = true
		return true
	}

	tmstate.GetBatch = gb
}

// sendToDatacenter send bcBlockResponseMessage to DataCenter
func (w *Watcher) SendToDatacenter(height int64) {
	value := WatchData{w.dirtyAccount, w.centerBatch}
	valueByte, err := itjs.Marshal(&value)
	if err != nil {
		return
	}
	msg := tmstate.DataCenterMsg{Height: height, Value: valueByte}
	msgBody, err := tmtypes.Json.Marshal(&msg)
	response, err := http.Post(viper.GetString(tmtypes.DataCenterUrl)+"saveBatch", "application/json", bytes.NewBuffer(msgBody))
	if err != nil {
		return
	}
	defer response.Body.Close()
}
