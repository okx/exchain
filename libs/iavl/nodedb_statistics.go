package iavl

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	"sync/atomic"
)

type RuntimeState struct {
	dbReadTime    int64
	dbReadCount   int64
	nodeReadCount int64
	dbWriteCount  int64

	totalPersistedCount int64
	totalPersistedSize  int64
	totalDeletedCount   int64
	totalOrphanCount    int64

	fromPpnc         int64
	fromTpp          int64
	fromNodeCache    int64
	fromOrphanCache  int64
	fromDisk         int64
}

type retrieveType int
const (
	unknown retrieveType = iota
	fromPpnc
	fromTpp
	fromNodeCache
	fromOrphanCache
	fromDisk
)

func newRuntimeState() *RuntimeState {
	r := &RuntimeState{}
	return r
}

func (s *RuntimeState) onLoadNode(from retrieveType) {

	switch from {
	case fromPpnc:
		atomic.AddInt64(&s.fromPpnc, 1)
	case fromTpp:
		atomic.AddInt64(&s.fromTpp, 1)
	case fromNodeCache:
		atomic.AddInt64(&s.fromNodeCache, 1)
	case fromOrphanCache:
		atomic.AddInt64(&s.fromOrphanCache, 1)
	case fromDisk:
		atomic.AddInt64(&s.fromDisk, 1)
	}
}

func (s *RuntimeState) addDBReadTime(ts int64) {
	atomic.AddInt64(&s.dbReadTime, ts)
}

func (s *RuntimeState) addDBReadCount() {
	atomic.AddInt64(&s.dbReadCount, 1)
}

func (s *RuntimeState) addDBWriteCount(count int64) {
	atomic.AddInt64(&s.dbWriteCount, count)
}

func (s *RuntimeState) addNodeReadCount() {
	atomic.AddInt64(&s.nodeReadCount, 1)
}

func (s *RuntimeState) resetDBReadTime() {
	atomic.StoreInt64(&s.dbReadTime, 0)
}

func (s *RuntimeState) resetDBReadCount() {
	atomic.StoreInt64(&s.dbReadCount, 0)
}

func (s *RuntimeState) resetDBWriteCount() {
	atomic.StoreInt64(&s.dbWriteCount, 0)
}

func (s *RuntimeState) resetNodeReadCount() {
	atomic.StoreInt64(&s.nodeReadCount, 0)
}

func (s *RuntimeState) getDBReadTime() int {
	return int(atomic.LoadInt64(&s.dbReadTime))
}

func (s *RuntimeState) getDBReadCount() int {
	return int(atomic.LoadInt64(&s.dbReadCount))
}

func (s *RuntimeState) getDBWriteCount() int {
	return int(atomic.LoadInt64(&s.dbWriteCount))
}

func (s *RuntimeState) getNodeReadCount() int {
	return int(atomic.LoadInt64(&s.nodeReadCount))
}

func (s *RuntimeState) resetCount() {
	s.resetDBReadTime()
	s.resetDBReadCount()
	s.resetDBWriteCount()
	s.resetNodeReadCount()
}

func (s *RuntimeState) increasePersistedSize(num int) {
	atomic.AddInt64(&s.totalPersistedSize, int64(num))
}
func (s *RuntimeState) increasePersistedCount(num int) {
	atomic.AddInt64(&s.totalPersistedCount, int64(num))
}
func (s *RuntimeState) increasOrphanCount(num int) {
	atomic.AddInt64(&s.totalOrphanCount, int64(num))
}
func (s *RuntimeState) increaseDeletedCount() {
	s.totalDeletedCount++
}

func inOutputModules(name string) bool {
	v, ok := OutputModules[name]
	return ok && v != 0
}

//================================
func (ndb *nodeDB) sprintCacheLog(version int64) (printLog string) {
	if !EnableAsyncCommit {
		return
	}

	if !inOutputModules(ndb.name) {
		return
	}

	nodeReadCount := ndb.state.getNodeReadCount()
	cacheReadCount := ndb.state.getNodeReadCount() - ndb.state.getDBReadCount()
	header := fmt.Sprintf("Save Version<%d>: Tree<%s>, ", version, ndb.name)

	printLog = fmt.Sprintf("getNodeFrom<ppnc=%d, tpp=%d, nodeCache=%d, orphanCache=%d, disk=%d>",
		ndb.state.fromPpnc,
		ndb.state.fromTpp,
		ndb.state.fromNodeCache,
		ndb.state.fromOrphanCache,
		ndb.state.fromDisk)
	printLog += fmt.Sprintf(", ppncCache:%d", len(ndb.prePersistNodeCache))
	printLog += fmt.Sprintf(", nodeCache:%d", ndb.nc.nodeCacheLen())
	printLog += fmt.Sprintf(", orphanCache:%d", ndb.oi.orphanNodeCacheLen())
	printLog += fmt.Sprintf(", totalPpnc:%d", treeMap.totalPpncSize)
	printLog += fmt.Sprintf(", evmPpnc:%d", treeMap.evmPpncSize)
	printLog += fmt.Sprintf(", accPpnc:%d", treeMap.accPpncSize)
	printLog += fmt.Sprintf(", dbRCnt:%d", ndb.state.getDBReadCount())
	printLog += fmt.Sprintf(", dbWCnt:%d", ndb.state.getDBWriteCount())
	printLog += fmt.Sprintf(", nodeRCnt:%d", ndb.state.getNodeReadCount())

	if nodeReadCount > 0 {
		printLog += fmt.Sprintf(", CHit:%.2f", float64(cacheReadCount)/float64(nodeReadCount)*100)
	} else {
		printLog += ", CHit:0"
	}
	printLog += fmt.Sprintf(", TPersisCnt:%d", atomic.LoadInt64(&ndb.state.totalPersistedCount))
	printLog += fmt.Sprintf(", TPersisSize:%d", atomic.LoadInt64(&ndb.state.totalPersistedSize))
	printLog += fmt.Sprintf(", TDelCnt:%d", atomic.LoadInt64(&ndb.state.totalDeletedCount))
	printLog += fmt.Sprintf(", TOrphanCnt:%d", atomic.LoadInt64(&ndb.state.totalOrphanCount))

	if ndb.name == "evm" {
		trace.GetElapsedInfo().AddInfo(trace.IavlRuntime, printLog)
	}

	return header+printLog
}


func (ndb *nodeDB) getDBReadTime() int {
	return ndb.state.getDBReadTime()
}

func (ndb *nodeDB) getDBReadCount() int {
	return ndb.state.getDBReadCount()
}

func (ndb *nodeDB) getDBWriteCount() int {
	return ndb.state.getDBWriteCount()
}

func (ndb *nodeDB) getNodeReadCount() int {
	return ndb.state.getNodeReadCount()
}

func (ndb *nodeDB) resetCount() {
	ndb.state.resetCount()
}

func (ndb *nodeDB) addDBReadTime(ts int64) {
	ndb.state.addDBReadTime(ts)
}

func (ndb *nodeDB) addDBReadCount() {
	ndb.state.addDBReadCount()
}

func (ndb *nodeDB) addDBWriteCount(count int64) {
	ndb.state.addDBWriteCount(count)
}

func (ndb *nodeDB) addNodeReadCount() {
	ndb.state.addNodeReadCount()
}