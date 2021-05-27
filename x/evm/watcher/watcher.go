package watcher

import (
	"fmt"
	"github.com/okex/exchain/x/stream/distrlock"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"os"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	streamTypes "github.com/okex/exchain/x/stream/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/abci/types"
)

type Watcher struct {
	store           *WatchStore
	height          uint64
	blockHash       common.Hash
	header          types.Header
	batch           []WatchMessage
	staleBatch      []WatchMessage
	cumulativeGas   map[uint64]uint64
	gasUsed         uint64
	blockTxs        []common.Hash
	sw              bool
	firstUse        bool
	enableScheduler bool
	scheduler       streamTypes.IDistributeStateService
}

const (
	lockerID              = "evm_lock_id"
	distributeLock        = "evm_watcher_lock"
	distributeLockTimeout = 1000
	latestHeightKey       = "latest_Height_key"
)

func IsWatcherEnabled() bool {
	return viper.GetBool(FlagFastQuery)
}

func NewWatcher() *Watcher {
	var scheduler streamTypes.IDistributeStateService

	if IsWatcherEnabled() && viper.GetString(FlagWatcherDBType) != DBTypeLevel {
		logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
		var err error
		scheduler, err = distrlock.NewRedisDistributeStateService(viper.GetString(FlagWatcherDisLockUrl), viper.GetString(FlagWatcherDisLockUrlPassword), logger, lockerID)
		if err != nil {
			panic(err)
		}
	}

	return &Watcher{store: InstanceOfWatchStore(), sw: IsWatcherEnabled(), firstUse: true, enableScheduler: viper.GetString(FlagWatcherDBType) != DBTypeLevel, scheduler: scheduler}
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
	wMsg := NewMsgCodeByHash(hash, code, w.height)
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
	w.store.Delete(key)
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
	w.store.Delete(evmtypes.GetContractBlockedListMemberKey(addr))
}

func (w *Watcher) DeleteContractDeploymentWhitelist(addr sdk.AccAddress) {
	if !w.Enabled() {
		return
	}
	w.store.Delete(evmtypes.GetContractDeploymentWhitelistMemberKey(addr))
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
	wMsg := NewMsgCodeByHash(hash, code, w.height)
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
	if !w.sw {
		return
	}

	if w.enableScheduler {
		w.CommitScheduler()
	} else {
		//hold it in temp
		batch := w.batch
		go func() {
			for _, b := range batch {
				w.store.Set(b.GetKey(), []byte(b.GetValue()))
			}
		}()
	}
}

func (w *Watcher) CommitScheduler() {
	if !w.sw {
		return
	}
	if !w.enableScheduler {
		return
	}

	//hold it in temp
	batch := w.batch
	// auto garbage collection
	w.batch = nil
	start := time.Now()
	locked, err := w.scheduler.FetchDistLock(distributeLock, lockerID, distributeLockTimeout)
	if !locked || err != nil {
		fmt.Println(time.Since(start))
		return
	}
	latestHeight, err := w.scheduler.GetDistState(latestHeightKey)
	if err != nil {
		w.scheduler.ReleaseDistLock(distributeLock, lockerID)
		fmt.Println(time.Since(start))
		return
	}

	// maybe first get latestHeightKey
	if len(latestHeight) == 0 {
		latestHeight = "0"
	}
	latestHeightNum, _ := strconv.Atoi(latestHeight)
	if uint64(latestHeightNum) < w.height {
		w.scheduler.SetDistState(latestHeightKey, strconv.FormatUint(w.height, 10))
		w.scheduler.ReleaseDistLock(distributeLock, lockerID)
		// set data
		fmt.Println("hbase to write")
		fmt.Println(time.Since(start))
		go func() {
			for _, b := range batch {
				w.store.Set(b.GetKey(), []byte(b.GetValue()))
			}
		}()
	} else {
		w.scheduler.ReleaseDistLock(distributeLock, lockerID)
		fmt.Println(time.Since(start))
	}
}
