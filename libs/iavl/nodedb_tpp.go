package iavl

import (
	"container/list"
	"sync"
)

type tempPrePersistNodes struct {
	mtx            sync.RWMutex
	tppMap         map[int64]*tppItem
	tppVersionList *list.List
}

func newTempPrePersistNodes() *tempPrePersistNodes {
	tpp := &tempPrePersistNodes{
		tppMap:              make(map[int64]*tppItem),
		tppVersionList:      list.New(),
	}
	return tpp
}

func (tpp *tempPrePersistNodes) getNode(hash []byte) (*Node, bool) {
	tpp.mtx.RLock()
	defer tpp.mtx.RUnlock()
	for v := tpp.tppVersionList.Back(); v != nil; v = v.Prev() {
		ver := v.Value.(int64)
		tppItem := tpp.tppMap[ver]

		if elem, ok := tppItem.nodeMap[string(hash)]; ok {
			return elem, ok
		}
	}
	return nil, false
}

func (tpp *tempPrePersistNodes) pushToTpp(version int64, tppMap map[string]*Node) {
	tpp.mtx.Lock()
	lItem := tpp.tppVersionList.PushBack(version)
	tpp.tppMap[version] = &tppItem{
		nodeMap:  tppMap,
		listItem: lItem,
	}
	tpp.mtx.Unlock()
}

func (tpp *tempPrePersistNodes) removeFromTpp(version int64) {
	tpp.mtx.Lock()
	tItem := tpp.tppMap[version]
	if tItem != nil {
		tpp.tppVersionList.Remove(tItem.listItem)
	}
	delete(tpp.tppMap, version)
	tpp.mtx.Unlock()
}

func (tpp *tempPrePersistNodes) getTppNodesNum() int {
	var size = 0
	tpp.mtx.RLock()
	for _, mp := range tpp.tppMap {
		size += len(mp.nodeMap)
	}
	tpp.mtx.RUnlock()
	return size
}

