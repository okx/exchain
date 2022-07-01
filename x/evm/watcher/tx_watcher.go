package watcher

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
)

type TxWatcher struct {
	enable     bool
	staleBatch []WatchMessage
	batch      []WatchMessage
}

var watcherPool = sync.Pool{
	New: func() interface{} {
		return &TxWatcher{
			enable: IsWatcherEnabled(),
		}
	},
}

func NewTxWatcher() *TxWatcher {
	return watcherPool.Get().(*TxWatcher)
}

func (w *TxWatcher) Enabled() bool {
	return w.enable
}

func (w *TxWatcher) SaveContractCode(addr common.Address, code []byte, height uint64) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCode(addr, code, height)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *TxWatcher) SaveContractCodeByHash(hash []byte, code []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgCodeByHash(hash, code)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *TxWatcher) SaveAccount(account interface{}, isDirectly bool) {
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

func (w *TxWatcher) DeleteAccount(account interface{}, isDirectly bool) {
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

func (w *TxWatcher) SaveState(addr common.Address, key, value []byte) {
	if !w.Enabled() {
		return
	}
	wMsg := NewMsgState(addr, key, value)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *TxWatcher) SaveContractBlockedListItem(addr interface{}) {
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

func (w *TxWatcher) SaveContractMethodBlockedListItem(addr interface{}, methods []byte) {
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

func (w *TxWatcher) SaveContractDeploymentWhitelistItem(addr interface{}) {
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

func (w *TxWatcher) DeleteContractBlockedList(addr interface{}) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractBlockedListItem(realAddr)
	if wMsg != nil {
		//key := wMsg.GetKey()
		//w.store.Delete(key)
		//w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *TxWatcher) DeleteContractDeploymentWhitelist(addr interface{}) {
	if !w.Enabled() {
		return
	}
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractDeploymentWhitelistItem(realAddr)
	if wMsg != nil {
		//key := wMsg.GetKey()
		//w.store.Delete(key)
		//w.watchData.DirtyList = append(w.watchData.DirtyList, key)
	}
}

func (w *TxWatcher) Finalize() {
	if !w.Enabled() {
		return
	}
	w.batch = append(w.batch, w.staleBatch...)
	w.staleBatch = []WatchMessage{}
}

func (w *TxWatcher) Destruct() []WatchMessage {
	if !w.Enabled() {
		return nil
	}
	batch := w.batch
	w.batch = []WatchMessage{}
	watcherPool.Put(w)
	return batch
}
