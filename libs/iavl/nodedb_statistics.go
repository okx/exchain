package iavl

import (
	"fmt"
	"sync/atomic"
)

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
	printLog += fmt.Sprintf(", nodeCCnt:%d", ndb.nodeCacheLen())
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
