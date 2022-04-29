package iavl

import (
	"fmt"
	"github.com/tendermint/go-amino"
)

//type OrphanInfo struct {
//	mtx            sync.RWMutex     // Read/write lock.
//	orphanNodeCache         map[string]*Node
//	heightOrphansCacheQueue *list.List
//	heightOrphansCacheSize  int
//	heightOrphansMap        map[int64]*heightOrphansItem
//}
//
//func newOrphanInfo() *OrphanInfo {
//
//	oi := &OrphanInfo{
//		orphanNodeCache:         make(map[string]*Node),
//		heightOrphansCacheQueue: list.New(),
//		heightOrphansCacheSize:  HeightOrphansCacheSize,
//		heightOrphansMap:        make(map[int64]*heightOrphansItem),
//	}
//
//	return oi
//}


func (ndb *nodeDB) handleOrphansRoutine() {
	for task := range ndb.taskChan {
		task()
	}
}

func (ndb *nodeDB) handleOrphans(version int64, rootHash []byte, newOrphans []*Node) {
	ndb.log(IavlInfo, "issue handleOrphans", "version", version)
	defer ndb.log(IavlInfo, "issue handleOrphans return", "version", version)

	task := func() {
		ndb.log(IavlInfo, "handleOrphans", "version", version)

		ndb.mtx.Lock()
		defer ndb.mtx.Unlock()

		defer ndb.log(IavlInfo, "handleOrphans locked", "version", version)

		ndb.saveOrphansAsync(version, newOrphans, false)
		ndb.setHeightOrphansItem(version, rootHash)
		ndb.taskResultChan <- version
	}

	ndb.taskChan <- task
}


func (ndb *nodeDB) handleOrphans2(version int64) {

	task := func() {
		ndb.taskResultChan <- version
	}

	ndb.taskChan <- task
}


func (ndb *nodeDB) getHandleOrphansResult(version int64) {
	if version <= 1 {
		return
	}
	version--
	ndb.log(IavlInfo, "getHandleOrphansResult enter ", "version", version)
	defer ndb.log(IavlInfo, "getHandleOrphansResult exit ", "version", version)

	for versionCompleted := range ndb.taskResultChan {
		if version == versionCompleted {
			break
		}
	}
}

func (ndb *nodeDB) setHeightOrphansItem(version int64, rootHash []byte) {
	if rootHash == nil {
		rootHash = []byte{}
	}
	orphanObj := &heightOrphansItem{
		version:  version,
		rootHash: rootHash,
	}
	ndb.heightOrphansCacheQueue.PushBack(orphanObj)
	ndb.heightOrphansMap[version] = orphanObj

	for ndb.heightOrphansCacheQueue.Len() > ndb.heightOrphansCacheSize {
		orphans := ndb.heightOrphansCacheQueue.Front()
		oldHeightOrphanItem := ndb.heightOrphansCacheQueue.Remove(orphans).(*heightOrphansItem)
		for _, node := range oldHeightOrphanItem.orphans {
			delete(ndb.orphanNodeCache, amino.BytesToStr(node.hash))
		}
		delete(ndb.heightOrphansMap, oldHeightOrphanItem.version)
	}
}

func (ndb *nodeDB) getRootWithCache(version int64) ([]byte, error) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	orphansObj, ok := ndb.heightOrphansMap[version]
	if ok {
		return orphansObj.rootHash, nil
	}
	return nil, fmt.Errorf("version %d is not in heightOrphansMap", version)
}


func (ndb *nodeDB) inVersionCacheMap(version int64) ([]byte, bool) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	item := ndb.heightOrphansMap[version]
	if item != nil {
		return item.rootHash, true
	}
	return nil, false
}

func (ndb *nodeDB) containedByOrphansMap(version int64) bool {
	ndb.mtx.RLock()
	defer ndb.mtx.RUnlock()
	_, ok := ndb.heightOrphansMap[version]
	return ok
}

func (ndb *nodeDB) feedOrphansMap(version int64, orphans []*Node) {
	v, ok := ndb.heightOrphansMap[version]
	if !ok {
		return
	}
	v.orphans = orphans
}


func (ndb *nodeDB) getNodeFromOrphanCache(hash []byte) *Node {
	elem, ok := ndb.orphanNodeCache[string(hash)]
	if ok {
		return elem
	}
	return nil
}

func (ndb *nodeDB) feedOrphanNodeCache(node *Node) {
	ndb.orphanNodeCache[string(node.hash)] = node
}

func (ndb *nodeDB) orphanNodeCacheLen() int {
	return len(ndb.orphanNodeCache)
}