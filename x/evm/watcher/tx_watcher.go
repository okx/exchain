package watcher

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
)

// TxWatcher cache watch data when run tx
// Use Enabled() to check if watcher is enable when call methods of TxWatcher
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
	wMsg := NewMsgCode(addr, code, height)
	if wMsg != nil {
		w.staleBatch = append(w.staleBatch, wMsg)
	}
}

func (w *TxWatcher) SaveContractCodeByHash(hash []byte, code []byte) {
	wMsg := NewMsgCodeByHash(hash, code)
	w.staleBatch = append(w.staleBatch, wMsg)
}

func (w *TxWatcher) SaveAccount(account interface{}) {
	acc, ok := account.(auth.Account)
	if !ok {
		return
	}
	wMsg := NewMsgAccount(acc)
	w.staleBatch = append(w.staleBatch, wMsg)
}

func (w *TxWatcher) DeleteAccount(account interface{}) {
	acc, ok := account.(auth.Account)
	if !ok {
		return
	}
	wMsg := NewDelAccMsg(acc)
	w.batch = append(w.batch, wMsg)

}

func (w *TxWatcher) SaveState(addr common.Address, key, value []byte) {
	wMsg := NewMsgState(addr, key, value)
	w.staleBatch = append(w.staleBatch, wMsg)
}

func (w *TxWatcher) SaveContractBlockedListItem(addr interface{}) {
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractBlockedListItem(realAddr)
	w.batch = append(w.batch, wMsg)
}

func (w *TxWatcher) SaveContractMethodBlockedListItem(addr interface{}, methods []byte) {
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractMethodBlockedListItem(realAddr, methods)
	w.batch = append(w.batch, wMsg)
}

func (w *TxWatcher) SaveContractDeploymentWhitelistItem(addr interface{}) {
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgContractDeploymentWhitelistItem(realAddr)
	w.batch = append(w.batch, wMsg)
}

func (w *TxWatcher) DeleteContractBlockedList(addr interface{}) {
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgDelContractBlockedListItem(realAddr)
	w.batch = append(w.batch, wMsg)
}

func (w *TxWatcher) DeleteContractDeploymentWhitelist(addr interface{}) {
	realAddr, ok := addr.(sdk.AccAddress)
	if !ok {
		return
	}
	wMsg := NewMsgDelContractDeploymentWhitelistItem(realAddr)
	w.batch = append(w.batch, wMsg)
}

func (w *TxWatcher) Finalize() {
	if !w.Enabled() {
		return
	}
	w.batch = append(w.batch, w.staleBatch...)
	w.staleBatch = []WatchMessage{}
}

func (w *TxWatcher) Destruct() []WatchMessage {
	batch := w.batch
	w.staleBatch = []WatchMessage{}
	w.batch = []WatchMessage{}
	watcherPool.Put(w)
	return batch
}
