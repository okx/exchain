package mpt

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"sync/atomic"
	"time"
)

type StatKeyValueStore struct {
	ethdb.KeyValueStore
	stat *RuntimeState
}

func NewStatKeyValueStore(db ethdb.KeyValueStore, stat *RuntimeState) ethdb.KeyValueStore {
	if stat == nil {
		stat = &RuntimeState{}
	}
	return &StatKeyValueStore{db, stat}
}

func (s *StatKeyValueStore) Get(key []byte) ([]byte, error) {
	s.stat.addDBReadCount()
	ts := time.Now()
	defer func() {
		s.stat.addDBReadTime(time.Now().Sub(ts).Nanoseconds())
	}()
	return s.KeyValueStore.Get(key)
}

func (s *StatKeyValueStore) Put(key []byte, value []byte) error {
	s.stat.addDBWriteCount(1)
	s.stat.increasePersistedCount(1)
	s.stat.increasePersistedSize(len(key) + len(value))
	return s.KeyValueStore.Put(key, value)
}

func (s *StatKeyValueStore) Delete(key []byte) error {
	s.stat.addDBWriteCount(1)
	s.stat.increaseDeletedCount()
	s.stat.increasePersistedSize(len(key))
	return s.KeyValueStore.Delete(key)
}

func (s *StatKeyValueStore) NewBatch() ethdb.Batch {
	nb := &Batch{
		s.KeyValueStore.NewBatch(),
		s.stat,
	}
	return nb
}

type Batch struct {
	ethdb.Batch
	stat *RuntimeState
}

func (s *Batch) Put(key []byte, value []byte) error {
	s.stat.addDBWriteCount(1)
	s.stat.increasePersistedCount(1)
	s.stat.increasePersistedSize(len(key) + len(value))
	return s.Batch.Put(key, value)
}

func (s *Batch) Delete(key []byte) error {
	s.stat.addDBWriteCount(1)
	s.stat.increasePersistedSize(len(key))
	s.stat.increaseDeletedCount()
	return s.Batch.Delete(key)
}

type RuntimeState struct {
	dbReadTime    int64
	dbReadCount   int64
	nodeReadCount int64
	dbWriteCount  int64

	totalPersistedCount int64
	totalPersistedSize  int64
	totalDeletedCount   int64
}

func NewRuntimeState() *RuntimeState {
	return &RuntimeState{}
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

func (s *RuntimeState) increaseDeletedCount() {
	atomic.AddInt64(&s.totalDeletedCount, 1)
}

func (s *RuntimeState) getPersistedSize() int {
	return int(atomic.LoadInt64(&s.totalPersistedSize))
}

func (s *RuntimeState) getPersistedCount() int {
	return int(atomic.LoadInt64(&s.totalPersistedCount))
}

func (s *RuntimeState) getDeletedCount() int {
	return int(atomic.LoadInt64(&s.totalDeletedCount))
}

//================================
func (ms *MptStore) sprintDebugLog(version int64) {
	if ms.logger == nil {
		return
	}
	nodeReadCount := ms.GetNodeReadCount()
	cacheReadCount := ms.GetCacheReadCount()
	nodeFromDBCount := 0
	if nodeReadCount > cacheReadCount {
		nodeFromDBCount = nodeReadCount - cacheReadCount
	}
	header := fmt.Sprintf("Save mpt Version<%d> , ", version)

	printLog := fmt.Sprintf("getNodeFrom<dbRCnt:%d, nodeCache=%d, nodeRCnt:%d>, DBCount<dbRCnt:%d, dbWCnt:%d>",
		nodeFromDBCount, cacheReadCount, ms.GetNodeReadCount(), ms.GetDBReadCount(), ms.GetDBWriteCount())

	if nodeReadCount > 0 {
		printLog += fmt.Sprintf(", NodeCHit:%.2f", float64(cacheReadCount)/float64(nodeReadCount)*100)
	} else {
		printLog += ", NodeCHit:0"
	}
	printLog += fmt.Sprintf(", TPersisCnt:%d", gStatic.getPersistedCount())
	printLog += fmt.Sprintf(", TPersisSize:%d", gStatic.getPersistedSize())
	printLog += fmt.Sprintf(", TDelCnt:%d", gStatic.getDeletedCount())

	ms.logger.Debug(header + printLog)
}
