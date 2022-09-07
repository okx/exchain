package iavl

import (
	"sync"

	"github.com/tendermint/go-amino"
)

type fastNodeChanges struct {
	additions map[string]*FastNode
	removals  map[string]interface{}
	mtx       sync.RWMutex
}

func newFastNodeChanges() *fastNodeChanges {
	return &fastNodeChanges{
		additions: make(map[string]*FastNode),
		removals:  make(map[string]interface{}),
	}
}

func (fnc *fastNodeChanges) get(key []byte) (*FastNode, bool) {
	fnc.mtx.RLock()
	defer fnc.mtx.RUnlock()
	if node, ok := fnc.additions[string(key)]; ok {
		return node, true
	}
	if _, ok := fnc.removals[string(key)]; ok {
		return nil, true
	}

	return nil, false
}

func (fnc *fastNodeChanges) add(key []byte, fastNode *FastNode) {
	fnc.mtx.Lock()
	fnc.additions[string(key)] = fastNode
	delete(fnc.removals, string(key))
	fnc.mtx.Unlock()
}

func (fnc *fastNodeChanges) addAdditions(key []byte, fastNode *FastNode) {
	fnc.mtx.Lock()
	fnc.additions[string(key)] = fastNode
	fnc.mtx.Unlock()
}

func (fnc *fastNodeChanges) remove(key []byte, value interface{}) {
	fnc.mtx.Lock()
	fnc.removals[string(key)] = value
	delete(fnc.additions, string(key))
	fnc.mtx.Unlock()
}

func (fnc *fastNodeChanges) addRemovals(key []byte) {
	fnc.mtx.Lock()
	fnc.removals[string(key)] = true
	fnc.mtx.Unlock()
}

func (fnc *fastNodeChanges) checkRemovals(key []byte) bool {
	fnc.mtx.RLock()
	defer fnc.mtx.RUnlock()

	if _, ok := fnc.removals[string(key)]; ok {
		return true
	}
	return false
}

func (fnc *fastNodeChanges) checkAdditions(key []byte) bool {
	fnc.mtx.RLock()
	defer fnc.mtx.RUnlock()
	if _, ok := fnc.additions[string(key)]; ok {
		return true
	}

	return false
}

func (fnc *fastNodeChanges) getAdditions() map[string]*FastNode {
	return fnc.additions
}

func (fnc *fastNodeChanges) getRemovals() map[string]interface{} {
	return fnc.removals
}

func (fnc *fastNodeChanges) clone() *fastNodeChanges {
	if fnc == nil {
		return nil
	}
	fnc.mtx.RLock()
	defer fnc.mtx.RUnlock()
	var additions map[string]*FastNode
	if fnc.additions != nil {
		additions = make(map[string]*FastNode, len(fnc.additions))
		for k, v := range fnc.additions {
			additions[k] = v
		}
	}

	var removals map[string]interface{}
	if fnc.removals != nil {
		removals = make(map[string]interface{}, len(fnc.removals))
		for k, v := range fnc.removals {
			removals[k] = v
		}
	}
	return &fastNodeChanges{
		additions: additions,
		removals:  removals,
	}
}

// mergePrev merge previous to fnc, prev is old than fnc
func (fnc *fastNodeChanges) mergePrev(prev *fastNodeChanges) *fastNodeChanges {
	if fnc == nil {
		return prev
	}
	if prev == nil {
		return fnc
	}

	for k, v := range prev.additions {
		if !fnc.checkAdditions(amino.StrToBytes(k)) && !fnc.checkRemovals(amino.StrToBytes(k)) {
			fnc.add(amino.StrToBytes(k), v)
		}
	}
	for k, v := range prev.removals {
		if !fnc.checkAdditions(amino.StrToBytes(k)) && !fnc.checkRemovals(amino.StrToBytes(k)) {
			fnc.remove(amino.StrToBytes(k), v)
		}
	}
	return fnc
}

// mergeLater merge later to fnc, later is new than fnc
func (fnc *fastNodeChanges) mergeLater(later *fastNodeChanges) {
	for k, v := range later.additions {
		fnc.add(amino.StrToBytes(k), v)
	}
	for k, v := range later.removals {
		fnc.remove(amino.StrToBytes(k), v)
	}
}

func (fnc *fastNodeChanges) reset() {
	fnc.mtx.Lock()
	for k := range fnc.additions {
		delete(fnc.additions, k)
	}
	for k := range fnc.removals {
		delete(fnc.removals, k)
	}
	fnc.mtx.Unlock()
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

func (fncv *fastNodeChangesWithVersion) expand(changes *fastNodeChanges) *fastNodeChanges {
	fncv.mtx.RLock()
	defer fncv.mtx.RUnlock()
	ret := changes.clone()
	if ret == nil {
		ret = newFastNodeChanges()
	}
	for i := len(fncv.versions) - 1; i >= 0; i-- {
		for k, v := range fncv.fncMap[fncv.versions[i]].additions {
			if !ret.checkAdditions(amino.StrToBytes(k)) && !ret.checkRemovals(amino.StrToBytes(k)) {
				ret.add(amino.StrToBytes(k), v)
			}
		}
		for k, v := range fncv.fncMap[fncv.versions[i]].removals {
			if !ret.checkAdditions(amino.StrToBytes(k)) && !ret.checkRemovals(amino.StrToBytes(k)) {
				ret.remove(amino.StrToBytes(k), v)
			}
		}
	}

	return ret
}
