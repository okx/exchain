package watcher

import "sync"

type watchMap struct {
	mu sync.RWMutex
	wdBytes map[int64]WatchData
}

func newWatchMap() *watchMap {
	return &watchMap{wdBytes: make(map[int64]WatchData)}
}

func (m *watchMap) set(height int64, wd WatchData) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wdBytes[height] = wd
}

func (m *watchMap) get(height int64) WatchData {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.wdBytes[height]
}

func (m *watchMap) del(height int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.wdBytes, height)
}

func (m *watchMap) fetch(height int64) WatchData {
	m.mu.Lock()
	defer m.mu.Unlock()
	wd := m.wdBytes[height]
	delete(m.wdBytes, height)
	return wd
}