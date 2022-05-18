package iavl

import (
	"fmt"

	"github.com/tendermint/go-amino"
)

type OrphanInfo struct {
	ndb             *nodeDB
	orphanNodeCache map[string]*Node
	orphanItemMap   map[int64]*orphanItem
	itemSize        int

	orphanTaskChan chan func()
	resultChan     chan int64
}

type orphanItem struct {
	rootHash []byte
	orphans  []*Node
}

func newOrphanInfo(ndb *nodeDB) *OrphanInfo {
	oi := &OrphanInfo{
		ndb:             ndb,
		orphanNodeCache: make(map[string]*Node),
		orphanItemMap:   make(map[int64]*orphanItem),
		itemSize:        HeightOrphansCacheSize,
		orphanTaskChan:  make(chan func(), 1),
		resultChan:      make(chan int64, 1),
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
		rootHash: rootHash,
	}
	_, ok := oi.orphanItemMap[version]
	if ok {
		panic(fmt.Sprintf("unexpected orphanItemMap, version: %d", version))
	}
	oi.orphanItemMap[version] = orphanObj
}

func (oi *OrphanInfo) removeOldOrphans(version int64) {
	expiredVersion := version - int64(oi.itemSize)
	expiredItem, ok := oi.orphanItemMap[expiredVersion]
	if !ok {
		return
	}
	for _, node := range expiredItem.orphans {
		delete(oi.orphanNodeCache, amino.BytesToStr(node.hash))
	}
	delete(oi.orphanItemMap, expiredVersion)
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
