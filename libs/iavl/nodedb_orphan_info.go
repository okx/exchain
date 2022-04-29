package iavl

import (
	"container/list"
	"github.com/tendermint/go-amino"
)

type OrphanInfo struct {
	orphanNodeCache         map[string]*Node
	heightOrphansCacheQueue *list.List
	heightOrphansCacheSize  int
	heightOrphansMap        map[int64]*heightOrphansItem

	orphanTaskChan   chan func()
	orphanResultChan chan int64
}

func newOrphanInfo() *OrphanInfo {

	oi := &OrphanInfo{
		orphanNodeCache:         make(map[string]*Node),
		heightOrphansCacheQueue: list.New(),
		heightOrphansCacheSize:  HeightOrphansCacheSize,
		heightOrphansMap:        make(map[int64]*heightOrphansItem),
		orphanTaskChan:          make(chan func(), 1),
		orphanResultChan:        make(chan int64, 1),
	}

	go oi.handleOrphansRoutine()
	oi.orphanResultChan <- 0
	return oi
}


func (oi *OrphanInfo) handleOrphansRoutine() {
	for task := range oi.orphanTaskChan {
		task()
	}
}

func (oi *OrphanInfo) wait4Result(version int64) {

	version--
	for versionCompleted := range oi.orphanResultChan {
		if versionCompleted == version {
			break
		} else if versionCompleted == 0 {
			break
		}
	}
}

func (oi *OrphanInfo) removeOldOrphans(version int64, rootHash []byte) {
	if rootHash == nil {
		rootHash = []byte{}
	}
	orphanObj := &heightOrphansItem{
		version:  version,
		rootHash: rootHash,
	}
	oi.heightOrphansCacheQueue.PushBack(orphanObj)
	oi.heightOrphansMap[version] = orphanObj

	for oi.heightOrphansCacheQueue.Len() > oi.heightOrphansCacheSize {
		orphans := oi.heightOrphansCacheQueue.Front()
		oldHeightOrphanItem := oi.heightOrphansCacheQueue.Remove(orphans).(*heightOrphansItem)
		for _, node := range oldHeightOrphanItem.orphans {
			delete(oi.orphanNodeCache, amino.BytesToStr(node.hash))
		}
		delete(oi.heightOrphansMap, oldHeightOrphanItem.version)
	}
}


func (oi *OrphanInfo) feedOrphansMap(version int64, orphans []*Node) {
	v, ok := oi.heightOrphansMap[version]
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
	v, ok := oi.heightOrphansMap[version]
	if ok {
		res = v.rootHash
		found = true
	}
	return
}