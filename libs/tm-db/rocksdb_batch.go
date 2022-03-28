//go:build rocksdb
// +build rocksdb

package db

import "github.com/tecbot/gorocksdb"

type RocksDBBatch struct {
	db    *RocksDB
	batch *gorocksdb.WriteBatch
}

var _ Batch = (*RocksDBBatch)(nil)

func NewRocksDBBatch(db *RocksDB) *RocksDBBatch {
	return &RocksDBBatch{
		db:    db,
		batch: gorocksdb.NewWriteBatch(),
	}
}

func (b *RocksDBBatch) assertOpen() {
	if b.batch == nil {
		panic("batch has been written or closed")
	}
}

// Set implements Batch.
func (b *RocksDBBatch) Set(key, value []byte) {
	b.assertOpen()
	b.batch.Put(key, value)
}

// Delete implements Batch.
func (b *RocksDBBatch) Delete(key []byte) {
	b.assertOpen()
	b.batch.Delete(key)
}

// Write implements Batch.
func (b *RocksDBBatch) Write() error {
	b.assertOpen()
	err := b.db.db.Write(b.db.wo, b.batch)
	if err != nil {
		return err
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	b.Close()
	return nil
}

// WriteSync implements Batch.
func (b *RocksDBBatch) WriteSync() error {
	b.assertOpen()
	err := b.db.db.Write(b.db.woSync, b.batch)
	if err != nil {
		return err
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	b.Close()
	return nil
}

// Close implements Batch.
func (b *RocksDBBatch) Close() {
	if b.batch != nil {
		b.batch.Destroy()
		b.batch = nil
	}
}

func (b *RocksDBBatch) Size() int {
	b.assertOpen()
	return b.batch.Count()
}

func (b *RocksDBBatch) Reset() {
	b.assertOpen()
	b.batch.Clear()
}

func (b *RocksDBBatch) NewIterator() *gorocksdb.WriteBatchIterator{
	b.assertOpen()
	return b.batch.NewIterator()
}

// WriteWithoutClose designed for ethdb.Batch: not close here for ethdb will use it again!!!
func (b *RocksDBBatch) WriteWithoutClose() error {
	b.assertOpen()
	err := b.db.db.Write(b.db.wo, b.batch)
	if err != nil {
		return err
	}
	// Never call b.Close() here!!!
	//b.Close()
	b.batch.Clear()
	return nil
}
