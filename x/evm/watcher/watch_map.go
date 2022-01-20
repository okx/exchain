package watcher

import (
	"container/list"
	"sync"
)

type watchMap struct {
	mtx          sync.Mutex
	cacheMap     map[int64]*list.Element
	cacheList   *list.List
	mrh          int64
}

func newDataMap() *watchMap {
	return &watchMap {
		cacheMap: make(map[int64]*list.Element),
		cacheList: list.New(),
	}
}

type payload struct {
	h int64
	d *WatchData
}

func (m *watchMap) insert(height int64, data *WatchData, mrh int64) {

	if data == nil {
		return
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	e := m.cacheList.PushBack(&payload{height, data})
	m.cacheMap[height] = e
	m.mrh = mrh
}

func (m *watchMap) fetch(height int64) (*WatchData, int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	popped := m.cacheMap[height]
	delete(m.cacheMap, height)
	if popped != nil {
		m.cacheList.Remove(popped)
		return popped.Value.(*payload).d, m.mrh
	}

	return nil, m.mrh
}

// remove all elements no higher than target
func (m *watchMap) remove(target int64) (int, int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	num := 0
	for {
		e := m.cacheList.Front()
		if e == nil {
			break
		}
		h := e.Value.(*payload).h
		if h > target {
			break
		}
		m.cacheList.Remove(e)
		delete(m.cacheMap, h)
		num++
	}

	return num, len(m.cacheMap)
}

func (m *watchMap) info() (int, int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return len(m.cacheMap), m.cacheList.Len()
}


