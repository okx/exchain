package iavl

import (
	"bytes"
	"container/list"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/okex/exchain/libs/iavl/trace"

	dbm "github.com/tendermint/tm-db"
)

type heightOrphansItem struct {
	version  int64
	rootHash []byte
	orphans  []*Node
}

type tppItem struct {
	nodeMap  map[string]*Node
	listItem *list.Element
}

func (ndb *nodeDB) SaveOrphans(batch dbm.Batch, version int64, orphans []*Node) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if EnableAsyncCommit {
		ndb.log(IavlDebug, "saving orphan node(size:%d) to OrphanCache", len(orphans))
		version--
		atomic.AddInt64(&ndb.totalOrphanCount, int64(len(orphans)))
		orphansObj := ndb.heightOrphansMap[version]
		if orphansObj != nil {
			orphansObj.orphans = orphans
		}
		for _, node := range orphans {
			ndb.orphanNodeCache[string(node.hash)] = node
			ndb.uncacheNode(node.hash)
			delete(ndb.prePersistNodeCache, string(node.hash))
			node.leftNode = nil
			node.rightNode = nil
		}
	} else {
		toVersion := ndb.getPreviousVersion(version)
		for _, node := range orphans {
			ndb.log(IavlDebug, "SAVEORPHAN %v-%v %X", node.version, toVersion, node.hash)
			ndb.saveOrphan(batch, node.hash, node.version, toVersion)
		}
	}
}

func (ndb *nodeDB) setHeightOrphansItem(version int64, rootHash []byte) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
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
			delete(ndb.orphanNodeCache, string(node.hash))
		}
		delete(ndb.heightOrphansMap, oldHeightOrphanItem.version)
	}
}

func (ndb *nodeDB) dbGet(k []byte) ([]byte, error) {
	ts := time.Now()
	defer func() {
		ndb.addDBReadTime(time.Now().Sub(ts).Nanoseconds())
	}()
	ndb.addDBReadCount()
	return ndb.db.Get(k)
}

func (ndb *nodeDB) saveNodeToPrePersistCache(node *Node) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if node.hash == nil {
		panic("Expected to find node.hash, but none found.")
	}
	if node.persisted || node.prePersisted {
		panic("Shouldn't be calling save on an already persisted node.")
	}

	node.prePersisted = true
	ndb.prePersistNodeCache[string(node.hash)] = node
}

func (ndb *nodeDB) persistTpp(event *commitEvent, trc *trace.Tracer) {

	batch := event.batch
	tpp := event.tpp

	trc.Pin("batchSet")
	for _, node := range tpp {
		ndb.batchSet(node, batch)
	}
	atomic.AddInt64(&ndb.totalPersistedCount, int64(len(tpp)))
	ndb.addDBWriteCount(int64(len(tpp)))

	trc.Pin("batchCommit")
	if err := ndb.Commit(batch); err != nil {
		panic(err)
	}
	ndb.asyncPersistTppFinised(event, trc)
}

func (ndb *nodeDB) asyncPersistTppStart(version int64) map[string]*Node {
	ndb.log(IavlDebug, "moving prePersistCache(size:%d) to tempPrePersistCache", len(ndb.prePersistNodeCache))
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	tpp := ndb.prePersistNodeCache
	ndb.prePersistNodeCache = map[string]*Node{}

	lItem := ndb.tppVersionList.PushBack(version)
	ndb.tppMap[version] = &tppItem{
		nodeMap:  tpp,
		listItem: lItem,
	}

	for _, node := range tpp {
		if node.persisted || !node.prePersisted {
			panic("unexpected node state")
		}
		node.persisted = true
	}

	return tpp
}

func (ndb *nodeDB) asyncPersistTppFinised(event *commitEvent, trc *trace.Tracer) {

	version := event.version
	tpp := event.tpp
	iavlHeight := event.iavlHeight

	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	trc.Pin("cacheNode")
	for _, node := range tpp {
		if !node.persisted {
			panic("unexpected logic")
		}
		ndb.cacheNode(node)
	}

	nodeNum := ndb.getTppNodesNum()

	tItem := ndb.tppMap[version]
	if tItem != nil {
		ndb.tppVersionList.Remove(tItem.listItem)
	}
	delete(ndb.tppMap, version)

	ndb.log(IavlInfo, "CommitSchedule: Height<%d>, Tree<%s>, IavlHeight<%d>, NodeNum<%d>, %s",
		version, ndb.name, iavlHeight, nodeNum, trc.Format())
}

// SaveNode saves a node to disk.
func (ndb *nodeDB) batchSet(node *Node, batch dbm.Batch) {
	if node.hash == nil {
		panic("Expected to find node.hash, but none found.")
	}
	if !node.persisted {
		panic("Should set node.persisted to true before batchSet.")
	}

	if !node.prePersisted {
		panic("Should be calling save on an prePersisted node.")
	}

	// Save node bytes to db.
	var buf bytes.Buffer
	buf.Grow(node.aminoSize())

	if err := node.writeBytes(&buf); err != nil {
		panic(err)
	}

	nodeKey := ndb.nodeKey(node.hash)
	nodeValue := buf.Bytes()
	batch.Set(nodeKey, nodeValue)
	atomic.AddInt64(&ndb.totalPersistedSize, int64(len(nodeKey)+len(nodeValue)))
	ndb.log(IavlDebug, "BATCH SAVE %X %p", node.hash, node)
	//node.persisted = true // move to function MovePrePersistCacheToTempCache
}

func (ndb *nodeDB) addDBReadTime(ts int64) {
	atomic.AddInt64(&ndb.dbReadTime, ts)
}

func (ndb *nodeDB) addDBReadCount() {
	atomic.AddInt64(&ndb.dbReadCount, 1)
}

func (ndb *nodeDB) addDBWriteCount(count int64) {
	atomic.AddInt64(&ndb.dbWriteCount, count)
}

func (ndb *nodeDB) addNodeReadCount() {
	atomic.AddInt64(&ndb.nodeReadCount, 1)
}

func (ndb *nodeDB) resetDBReadTime() {
	atomic.StoreInt64(&ndb.dbReadTime, 0)
}

func (ndb *nodeDB) resetDBReadCount() {
	atomic.StoreInt64(&ndb.dbReadCount, 0)
}

func (ndb *nodeDB) resetDBWriteCount() {
	atomic.StoreInt64(&ndb.dbWriteCount, 0)
}

func (ndb *nodeDB) resetNodeReadCount() {
	atomic.StoreInt64(&ndb.nodeReadCount, 0)
}

func (ndb *nodeDB) getDBReadTime() int {
	return int(atomic.LoadInt64(&ndb.dbReadTime))
}

func (ndb *nodeDB) getDBReadCount() int {
	return int(atomic.LoadInt64(&ndb.dbReadCount))
}

func (ndb *nodeDB) getDBWriteCount() int {
	return int(atomic.LoadInt64(&ndb.dbWriteCount))
}

func (ndb *nodeDB) getNodeReadCount() int {
	return int(atomic.LoadInt64(&ndb.nodeReadCount))
}

func (ndb *nodeDB) resetCount() {
	ndb.resetDBReadTime()
	ndb.resetDBReadCount()
	ndb.resetDBWriteCount()
	ndb.resetNodeReadCount()
}

func (ndb *nodeDB) sprintCacheLog(version int64) string {
	if !EnableAsyncCommit {
		return ""
	}

	nodeReadCount := ndb.getNodeReadCount()
	cacheReadCount := ndb.getNodeReadCount() - ndb.getDBReadCount()
	printLog := fmt.Sprintf("Save Version<%d>: Tree<%s>", version, ndb.name)

	printLog += fmt.Sprintf(", TotalPreCommitCacheSize:%d", treeMap.totalPreCommitCacheSize)
	printLog += fmt.Sprintf(", nodeCCnt:%d", len(ndb.nodeCache))
	printLog += fmt.Sprintf(", orphanCCnt:%d", len(ndb.orphanNodeCache))
	printLog += fmt.Sprintf(", prePerCCnt:%d", len(ndb.prePersistNodeCache))
	printLog += fmt.Sprintf(", dbRCnt:%d", ndb.getDBReadCount())
	printLog += fmt.Sprintf(", dbWCnt:%d", ndb.getDBWriteCount())
	printLog += fmt.Sprintf(", nodeRCnt:%d", ndb.getNodeReadCount())

	if nodeReadCount > 0 {
		printLog += fmt.Sprintf(", CHit:%.2f", float64(cacheReadCount)/float64(nodeReadCount)*100)
	} else {
		printLog += ", CHit:0"
	}
	printLog += fmt.Sprintf(", TPersisCnt:%d", atomic.LoadInt64(&ndb.totalPersistedCount))
	printLog += fmt.Sprintf(", TPersisSize:%d", atomic.LoadInt64(&ndb.totalPersistedSize))
	printLog += fmt.Sprintf(", TDelCnt:%d", atomic.LoadInt64(&ndb.totalDeletedCount))
	printLog += fmt.Sprintf(", TOrphanCnt:%d", atomic.LoadInt64(&ndb.totalOrphanCount))

	return printLog
}

func (ndb *nodeDB) getTppNodesNum() int {
	var size = 0
	for _, mp := range ndb.tppMap {
		size += len(mp.nodeMap)
	}
	return size
}

func (ndb *nodeDB) NewBatch() dbm.Batch {
	return ndb.db.NewBatch()
}

func (ndb *nodeDB) updateBranch(node *Node, savedNodes map[string]*Node) []byte {
	if node.persisted || node.prePersisted {
		return node.hash
	}

	if node.leftNode != nil {
		node.leftHash = ndb.updateBranch(node.leftNode, savedNodes)
	}
	if node.rightNode != nil {
		node.rightHash = ndb.updateBranch(node.rightNode, savedNodes)
	}

	node._hash()
	ndb.saveNodeToPrePersistCache(node)

	node.leftNode = nil
	node.rightNode = nil

	// TODO: handle magic number
	savedNodes[hex.EncodeToString(node.hash)] = node

	return node.hash
}

func (ndb *nodeDB) getRootWithCache(version int64) ([]byte, error) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	orphansObj := ndb.heightOrphansMap[version]
	if orphansObj != nil {
		return orphansObj.rootHash, nil
	}
	return nil, fmt.Errorf("version %d is not in heightOrphansMap", version)
}

// Saves orphaned nodes to disk under a special prefix.
// version: the new version being saved.
// orphans: the orphan nodes created since version-1
func (ndb *nodeDB) saveCommitOrphans(batch dbm.Batch, version int64, orphans map[string]int64) {
	ndb.log(IavlDebug, "saving committed orphan node log to disk")
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	toVersion := ndb.getPreviousVersion(version)
	for hash, fromVersion := range orphans {
		ndb.log(IavlDebug, "SAVEORPHAN %v-%v %X", fromVersion, toVersion, hash)
		ndb.saveOrphan(batch, []byte(hash), fromVersion, toVersion)
	}
}

func (ndb *nodeDB) getNodeInTpp(hash []byte) (*Node, bool) {
	for v := ndb.tppVersionList.Back(); v != nil; v = v.Prev() {
		ver := v.Value.(int64)
		tppItem := ndb.tppMap[ver]

		if elem, ok := tppItem.nodeMap[string(hash)]; ok {
			return elem, ok
		}
	}
	return nil, false
}

func (ndb *nodeDB) getRootWithCacheAndDB(version int64) ([]byte, error) {
	if EnableAsyncCommit {
		root, err := ndb.getRootWithCache(version)
		if err == nil {
			return root, err
		}
	}
	return ndb.getRoot(version)
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
