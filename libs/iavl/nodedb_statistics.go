package iavl

import (
	"fmt"
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
}

func (ndb *RuntimeState) addDBReadTime(ts int64) {
	atomic.AddInt64(&ndb.dbReadTime, ts)
}

func (ndb *RuntimeState) addDBReadCount() {
	atomic.AddInt64(&ndb.dbReadCount, 1)
}

func (ndb *RuntimeState) addDBWriteCount(count int64) {
	atomic.AddInt64(&ndb.dbWriteCount, count)
}

func (ndb *RuntimeState) addNodeReadCount() {
	atomic.AddInt64(&ndb.nodeReadCount, 1)
}

func (ndb *RuntimeState) resetDBReadTime() {
	atomic.StoreInt64(&ndb.dbReadTime, 0)
}

func (ndb *RuntimeState) resetDBReadCount() {
	atomic.StoreInt64(&ndb.dbReadCount, 0)
}

func (ndb *RuntimeState) resetDBWriteCount() {
	atomic.StoreInt64(&ndb.dbWriteCount, 0)
}

func (ndb *RuntimeState) resetNodeReadCount() {
	atomic.StoreInt64(&ndb.nodeReadCount, 0)
}

func (ndb *RuntimeState) getDBReadTime() int {
	return int(atomic.LoadInt64(&ndb.dbReadTime))
}

func (ndb *RuntimeState) getDBReadCount() int {
	return int(atomic.LoadInt64(&ndb.dbReadCount))
}

func (ndb *RuntimeState) getDBWriteCount() int {
	return int(atomic.LoadInt64(&ndb.dbWriteCount))
}

func (ndb *RuntimeState) getNodeReadCount() int {
	return int(atomic.LoadInt64(&ndb.nodeReadCount))
}

func (ndb *RuntimeState) resetCount() {
	ndb.resetDBReadTime()
	ndb.resetDBReadCount()
	ndb.resetDBWriteCount()
	ndb.resetNodeReadCount()
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

//================================
func (ndb *nodeDB) sprintCacheLog(version int64) string {
	if !EnableAsyncCommit {
		return ""
	}

	nodeReadCount := ndb.state.getNodeReadCount()
	cacheReadCount := ndb.state.getNodeReadCount() - ndb.state.getDBReadCount()
	printLog := fmt.Sprintf("Save Version<%d>: Tree<%s>", version, ndb.name)

	printLog += fmt.Sprintf(", TotalPreCommitCacheSize:%d", treeMap.totalPreCommitCacheSize)
	printLog += fmt.Sprintf(", nodeCCnt:%d", ndb.nc.nodeCacheLen())
	printLog += fmt.Sprintf(", orphanCCnt:%d", ndb.oi.orphanNodeCacheLen())
	printLog += fmt.Sprintf(", prePerCCnt:%d", len(ndb.prePersistNodeCache))
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

	return printLog
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