package iavl

import "sync"

type SyncMap struct {
	mp map[int64]bool
	lock sync.RWMutex
}

func NewSyncMap() *SyncMap {
	return &SyncMap{
		mp: make(map[int64]bool),
		lock: sync.RWMutex{},
	}
}

func (sm *SyncMap) Get(key int64) bool {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	return sm.mp[key]
}

func (sm *SyncMap) Set(key int64, value bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.mp[key] = value
}

func (sm *SyncMap) Has(key int64) bool {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	_, ok := sm.mp[key]
	return ok
}

func (sm *SyncMap) Delete(key int64) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.DeleteWithoutLock(key)
}

func (sm *SyncMap) DeleteWithoutLock(key int64) {
	delete(sm.mp, key)
}

func (sm *SyncMap) Len() int {
	return len(sm.mp)
}

func (sm *SyncMap) Range(f func(key int64, value bool)bool) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	for k, v := range sm.mp {
		ok := f(k, v)
		if !ok {
			break
		}
	}
}

func (sm *SyncMap) Clone() map[int64]bool{
	mp := make(map[int64]bool, sm.Len())
	sm.Range(func(key int64, value bool) bool {
		mp[key] = value
		return true
	})
	return mp
}

