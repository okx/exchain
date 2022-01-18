package watcher

import "sync"

type watchMap struct {
	mu sync.RWMutex
	wdBytes map[int64]*WatchData
}

func newWatchMap() *watchMap {
	return &watchMap{wdBytes: make(map[int64]*WatchData)}
}

func (m *watchMap) set(height int64, wd *WatchData) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wdBytes[height] = wd
}

func (m *watchMap) fetch(height int64) *WatchData {
	m.mu.Lock()
	defer m.mu.Unlock()
	wd := m.wdBytes[height]
	for k, _ := range m.wdBytes {
		if k <= height {
			delete(m.wdBytes, k)
		}
	}
	return wd
}