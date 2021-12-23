package state

import (
	"container/list"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
)

type deltaMap struct {
	mtx          sync.Mutex
	cacheMap     map[int64]*list.Element
	cacheList   *list.List
}

func newDataMap() *deltaMap {
	return &deltaMap {
		cacheMap: make(map[int64]*list.Element),
		cacheList: list.New(),
	}
}

type payload struct {
	h int64
	d *types.Deltas
}

func (m *deltaMap) insert(height int64, data *types.Deltas) {

	if data == nil {
		return
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	e := m.cacheList.PushBack(&payload{height, data})
	m.cacheMap[height] = e
}

func (m *deltaMap) fetch(height int64) *types.Deltas {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	popped := m.cacheMap[height]
	delete(m.cacheMap, height)
	if popped != nil {
		m.cacheList.Remove(popped)
		return popped.Value.(*payload).d
	}

	return nil
}

// remove all elements no higher than target
func (m *deltaMap) remove(target int64) (int, int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	num := 0
	for e := m.cacheList.Back(); e != nil; e = e.Prev() {
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
