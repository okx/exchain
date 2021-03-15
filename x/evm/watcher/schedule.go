package watcher

type EvmWatcher struct {
}

var gWatcher *EvmWatcher

func newEvmWatcher() *EvmWatcher {
	ret := &EvmWatcher{}
	return ret
}

func init() {
	gWatcher = newEvmWatcher()
}
