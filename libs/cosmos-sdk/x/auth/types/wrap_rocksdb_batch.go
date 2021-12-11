// +build rocksdb

package types

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tecbot/gorocksdb"
)

var _ ethdb.Batch = (*WrapRocksDBBatch)(nil)

type WrapRocksDBBatch struct {
	db    *WrapRocksDB
	batch *gorocksdb.WriteBatch
	size int
}

func NewWrapRocksDBBatch(db *WrapRocksDB) *WrapRocksDBBatch {
	return &WrapRocksDBBatch{
		db:    db,
		batch: gorocksdb.NewWriteBatch(),
	}
}

func (wrsdbb *WrapRocksDBBatch) assertOpen() {
	if wrsdbb.batch == nil {
		panic("batch has been written or closed")
	}
}

func (wrsdbb *WrapRocksDBBatch) Put(key []byte, value []byte) error {
	wrsdbb.assertOpen()
	wrsdbb.batch.Put(key, value)
	wrsdbb.size += len(value)

	return nil
}

func (wrsdbb *WrapRocksDBBatch) Delete(key []byte) error {
	wrsdbb.assertOpen()
	wrsdbb.batch.Delete(key)
	wrsdbb.size += len(key)

	return nil
}

func (wrsdbb *WrapRocksDBBatch) ValueSize() int {
	return wrsdbb.size
}

func (wrsdbb *WrapRocksDBBatch) Write() error {
	wrsdbb.assertOpen()
	wo := gorocksdb.NewDefaultWriteOptions()
	err := wrsdbb.db.DB().Write(wo, wrsdbb.batch)
	if err != nil {
		return err
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	//wrsdbb.Close()
	return nil
}

func (wrsdbb *WrapRocksDBBatch) Reset() {
	wrsdbb.assertOpen()
	wrsdbb.batch.Clear()
	wrsdbb.size = 0
}

// Close implements Batch.
func (wrsdbb *WrapRocksDBBatch) Close() {
	if wrsdbb.batch != nil {
		wrsdbb.batch.Destroy()
		wrsdbb.batch = nil
	}
}

// Replay replays the batch contents.
func (wrsdbb *WrapRocksDBBatch) Replay(w ethdb.KeyValueWriter) error {
	rp := &replayer{writer: w}

	itr := wrsdbb.batch.NewIterator()
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


// replayer is a small wrapper to implement the correct replay methods.
type replayer struct {
	writer  ethdb.KeyValueWriter
	failure error
}

// Put inserts the given value into the key-value data store.
func (r *replayer) Put(key, value []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Put(key, value)
}

// Delete removes the key from the key-value data store.
func (r *replayer) Delete(key []byte) {
	// If the replay already failed, stop executing ops
	if r.failure != nil {
		return
	}
	r.failure = r.writer.Delete(key)
}