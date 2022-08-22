package iavl

import (
	"log"
	"sync"
)

type fastNodeChanges struct {
	additions map[string]*FastNode
	removals  map[string]interface{}
}

var fastNodeChangesPool = &sync.Pool{
	New: func() interface{} {
		return newFastNodeChanges()
	},
}

func getReusableFastNodeChanges() *fastNodeChanges {
	fnc := fastNodeChangesPool.Get().(*fastNodeChanges)
	fnc.reset()

	return fnc
}

func putReusableFastNodeChanges(fnc *fastNodeChanges) {
	fastNodeChangesPool.Put(fnc)
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

func (fnc *fastNodeChanges) checkRemovals(key string) bool {
	if _, ok := fnc.removals[key]; ok {
		return true
	}
	return false
}

func (fnc *fastNodeChanges) checkAdditions(key string) bool {
	if _, ok := fnc.additions[key]; ok {
		return true
	}

	return false
}

func (fnc *fastNodeChanges) reset() {
	for k := range fnc.additions {
		delete(fnc.additions, k)
	}
	for k := range fnc.removals {
		delete(fnc.removals, k)
	}
}

func (fnc *fastNodeChanges) clone() *fastNodeChanges {
	if fnc == nil {
		return nil
	}
	additions := make(map[string]*FastNode, len(fnc.additions))
	for k, v := range fnc.additions {
		additions[k] = v
	}
	removals := make(map[string]interface{}, len(fnc.removals))
	for k, v := range fnc.removals {
		removals[k] = v
	}
	return &fastNodeChanges{
		additions: additions,
		removals:  removals,
	}
}

func (fnc *fastNodeChanges) merge(src *fastNodeChanges) *fastNodeChanges {
	if fnc == nil {
		return src
	}
	if src == nil {
		return fnc
	}
	for k, v := range src.additions {
		if !fnc.checkAdditions(k) && !fnc.checkRemovals(k) {
			fnc.add(k, v)
		}
	}
	for k, v := range src.removals {
		if !fnc.checkAdditions(k) && !fnc.checkRemovals(k) {
			fnc.remove(k, v)
		}
	}
	return fnc
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
		if len(fncv.versions) > 0 && version != fncv.versions[0] {
			log.Printf("----giskook---- %v \n", fncv.versions)
		}
		return
	}
	fncv.mtx.Lock()
	defer fncv.mtx.Unlock()
	fncv.versions = fncv.versions[1:]

	putReusableFastNodeChanges(fncv.fncMap[version])
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

func (fncv *fastNodeChangesWithVersion) expand(changes *fastNodeChanges) *fastNodeChanges {
	fncv.mtx.RLock()
	defer fncv.mtx.RUnlock()
	ret := changes.clone()
	if ret == nil {
		ret = newFastNodeChanges()
	}
	for i := len(fncv.versions) - 1; i >= 0; i-- {
		for k, v := range fncv.fncMap[fncv.versions[i]].additions {
			if !ret.checkAdditions(k) && !ret.checkRemovals(k) {
				ret.add(k, v)
			}
		}
		for k, v := range fncv.fncMap[fncv.versions[i]].removals {
			if !ret.checkAdditions(k) && !ret.checkRemovals(k) {
				ret.remove(k, v)
			}
		}
	}

	return ret
}
