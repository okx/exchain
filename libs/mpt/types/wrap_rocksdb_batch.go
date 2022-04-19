//go:build rocksdb
// +build rocksdb

package types

import (
	"container/list"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/ethdb"
	tmdb "github.com/okex/exchain/libs/tm-db"
	"github.com/tecbot/gorocksdb"
)

type BatchCache struct {
	batchList  *list.List
	batchCache map[int64]*list.Element
	maxSize    int

	lock sync.Mutex
}

func NewBatchCache(maxSize int) *BatchCache {
	return &BatchCache{
		batchList:  list.New(),
		batchCache: make(map[int64]*list.Element, maxSize),
		maxSize:    maxSize,
	}
}

func (bc *BatchCache) PushBack(batch *WrapRocksDBBatch) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	ele := bc.batchList.PushBack(batch)
	bc.batchCache[batch.GetID()] = ele
}

func (bc *BatchCache) TryPopFront() *WrapRocksDBBatch {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	if bc.batchList.Len() > bc.maxSize {
		deathEle := bc.batchList.Front()
		bc.batchList.Remove(deathEle)

		deathBatch := deathEle.Value.(*WrapRocksDBBatch)
		delete(bc.batchCache, deathBatch.GetID())

		return deathBatch
	}

	return nil
}

func (bc *BatchCache) MoveToBack(id int64) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	if ele, ok := bc.batchCache[id]; ok {
		bc.batchList.MoveToBack(ele)
	}
}

var (
	gBatchCache *BatchCache
	batchIdSeed int64

	initRocksdbBatchOnce sync.Once
)

func InstanceBatchCache() *BatchCache {
	initRocksdbBatchOnce.Do(func() {
		gBatchCache = NewBatchCache(int(TrieRocksdbBatchSize))
	})

	return gBatchCache
}

var _ ethdb.Batch = (*WrapRocksDBBatch)(nil)

type WrapRocksDBBatch struct {
	*tmdb.RocksDBBatch
	id int64
}

func NewWrapRocksDBBatch(db *tmdb.RocksDB) *WrapRocksDBBatch {
	sed := atomic.LoadInt64(&batchIdSeed)
	batch := &WrapRocksDBBatch{tmdb.NewRocksDBBatch(db), sed}
	atomic.AddInt64(&batchIdSeed, 1)

	batchCache := InstanceBatchCache()
	batchCache.PushBack(batch)
	if deathBatch := batchCache.TryPopFront(); deathBatch != nil {
		deathBatch.Close()
	}

	return batch
}

func (wrsdbb *WrapRocksDBBatch) Put(key []byte, value []byte) error {
	InstanceBatchCache().MoveToBack(wrsdbb.GetID())

	wrsdbb.Set(key, value)
	return nil
}

func (wrsdbb *WrapRocksDBBatch) Delete(key []byte) error {
	InstanceBatchCache().MoveToBack(wrsdbb.GetID())

	wrsdbb.RocksDBBatch.Delete(key)
	return nil
}

func (wrsdbb *WrapRocksDBBatch) ValueSize() int {
	return wrsdbb.Size()
}

func (wrsdbb *WrapRocksDBBatch) Write() error {
	InstanceBatchCache().MoveToBack(wrsdbb.GetID())

	return wrsdbb.RocksDBBatch.WriteWithoutClose()
}

// Replay replays the batch contents.
func (wrsdbb *WrapRocksDBBatch) Replay(w ethdb.KeyValueWriter) error {
	InstanceBatchCache().MoveToBack(wrsdbb.GetID())

	rp := &replayer{writer: w}

	itr := wrsdbb.NewIterator()
	for itr.Next() {
		rcd := itr.Record()

		switch rcd.Type {
		case gorocksdb.WriteBatchValueRecord:
			rp.Put(rcd.Key, rcd.Value)
		case gorocksdb.WriteBatchDeletionRecord:
			rp.Delete(rcd.Key)
		}
	}

	return nil
}

func (wrsdbb *WrapRocksDBBatch) GetID() int64 {
	return wrsdbb.id
}
