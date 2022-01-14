package watcher

import "github.com/okex/exchain/libs/tendermint/libs/log"

func NewWatcher(l log.Logger) IWatcher {
	var (
		ret IWatcher
	)
	if l==nil{
		l=log.NewNopLogger()
	}
	wEnable := IsWatcherEnabled()
	asyncEnable := AsyncTxEnable
	baseW := newWatcher(l)
	if !wEnable {
		return newDisableWatcher(baseW)
	}
	if asyncEnable {
		ret = newConcurrentWatcher(baseW)
	} else {
		ret = baseW
	}
	return ret
}
