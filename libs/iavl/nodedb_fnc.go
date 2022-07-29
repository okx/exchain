package iavl

import (
	"sync"
)

type fastNodeChanges struct {
	additions map[string]*FastNode
	removals  map[string]interface{}
}

func newFastNodeChanges() *fastNodeChanges {
	return &fastNodeChanges{
		additions: make(map[string]*FastNode),
		removals:  make(map[string]interface{}),
	}
}

func (fnc *fastNodeChanges) get(key []byte) (*FastNode, bool) {
	if node, ok := fnc.additions[string(key)]; ok {
		return node, true
	}
	if _, ok := fnc.removals[string(key)]; ok {
		return nil, true
	}

	return nil, false
}

func (fnc *fastNodeChanges) add(key string, fastNode *FastNode) {
	fnc.additions[key] = fastNode
	delete(fnc.removals, key)
}

func (fnc *fastNodeChanges) remove(key string, value interface{}) {
	fnc.removals[key] = value
	delete(fnc.additions, key)
}

func (fnc *fastNodeChanges) reset() {
	for k := range fnc.additions {
		delete(fnc.additions, k)
	}
	for k := range fnc.removals {
		delete(fnc.removals, k)
	}
}

type fastNodeChangesWithVersion struct {
	mtx      sync.RWMutex
	versions []int64
	fncMap   map[int64]*fastNodeChanges
}

func newFastNodeChangesWithVersion() *fastNodeChangesWithVersion {
	return &fastNodeChangesWithVersion{
		fncMap: make(map[int64]*fastNodeChanges),
	}
}

func (fncv *fastNodeChangesWithVersion) add(version int64, fnc *fastNodeChanges) {
	fncv.mtx.Lock()
	defer fncv.mtx.Unlock()
	fncv.versions = append(fncv.versions, version)
	fncv.fncMap[version] = fnc
}

func (fncv *fastNodeChangesWithVersion) remove(version int64) {
	if len(fncv.versions) < 1 || version != fncv.versions[0] {
		return
	}
	fncv.mtx.Lock()
	defer fncv.mtx.Unlock()
	fncv.versions = fncv.versions[1:]
	delete(fncv.fncMap, version)
}

func (fncv *fastNodeChangesWithVersion) get(key []byte) (*FastNode, bool) {
	fncv.mtx.RLock()
	defer fncv.mtx.RUnlock()
	for i := len(fncv.versions) - 1; i >= 0; i-- {
		if fn, ok := fncv.fncMap[fncv.versions[i]].get(key); ok {
			return fn, ok
		}
	}
	return nil, false
}
