package iavl

import (
	"container/list"
	"fmt"
	"github.com/tendermint/go-amino"
)

type OrphanInfo struct {
	orphanNodeCache         map[string]*Node
	orphanItemCacheQueue *list.List
	orphanItemCacheSize  int
	orphanItemMap        map[int64]*orphanItem

	orphanTaskChan   chan func()
	resultChan       chan int64
}

type orphanItem struct {
	version  int64
	rootHash []byte
	orphans  []*Node
}

func newOrphanInfo() *OrphanInfo {

	oi := &OrphanInfo{
		orphanNodeCache:         make(map[string]*Node),
		orphanItemCacheQueue: list.New(),
		orphanItemCacheSize:  HeightOrphansCacheSize,
		orphanItemMap:        make(map[int64]*orphanItem),
		orphanTaskChan:          make(chan func(), 1),
		resultChan:              make(chan int64, 1),
	}

	oi.enqueueResult(0)
	go oi.handleOrphansRoutine()
	return oi
}

func (oi *OrphanInfo) enqueueResult(res int64) {
	oi.resultChan <- res
}

func (oi *OrphanInfo) enqueueTask(t func()) {
	oi.orphanTaskChan <- t
}

func (oi *OrphanInfo) handleOrphansRoutine() {
	for task := range oi.orphanTaskChan {
		task()
	}
}

func (oi *OrphanInfo) wait4Result(version int64) {

	version--
	for versionCompleted := range oi.resultChan {
		if versionCompleted == version {
			break
		} else if versionCompleted == 0 {
			break
		}
	}
}

func (oi *OrphanInfo) addOrphanItem(version int64, rootHash []byte) {
	if rootHash == nil {
		rootHash = []byte{}
	}
	orphanObj := &orphanItem{
		version:  version,
		rootHash: rootHash,
	}
	oi.orphanItemCacheQueue.PushBack(orphanObj)
	_, ok := oi.orphanItemMap[version]
	if ok {
		panic(fmt.Sprintf("unexpected orphanItemMap, version: %d", version))
	}
	oi.orphanItemMap[version] = orphanObj
}


func (oi *OrphanInfo) removeOldOrphans() {
	for oi.orphanItemCacheQueue.Len() > oi.orphanItemCacheSize {
		orphans := oi.orphanItemCacheQueue.Front()
		oldHeightOrphanItem := oi.orphanItemCacheQueue.Remove(orphans).(*orphanItem)
		for _, node := range oldHeightOrphanItem.orphans {
			delete(oi.orphanNodeCache, amino.BytesToStr(node.hash))
		}
		delete(oi.orphanItemMap, oldHeightOrphanItem.version)
	}
}


func (oi *OrphanInfo) feedOrphansMap(version int64, orphans []*Node) {
	v, ok := oi.orphanItemMap[version]
	if !ok {
		return
	}
	v.orphans = orphans
}

func (oi *OrphanInfo) feedOrphanNodeCache(node *Node) {
	oi.orphanNodeCache[string(node.hash)] = node
}


func (oi *OrphanInfo) getNodeFromOrphanCache(hash []byte) *Node {
	elem, ok := oi.orphanNodeCache[string(hash)]
	if ok {
		return elem
	}
	return nil
}


func (oi *OrphanInfo) orphanNodeCacheLen() int {
	return len(oi.orphanNodeCache)
}

func (oi *OrphanInfo) findRootHash(version int64) (res []byte, found bool) {
	v, ok := oi.orphanItemMap[version]
	if ok {
		res = v.rootHash
		found = true
	}
	return
}