package watcher

func NewWatcher() IWatcher {
	var (
		ret IWatcher
	)

	wEnable := IsWatcherEnabled()
	asyncEnable := AsyncTxEnable
	baseW := newWatcher()
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
