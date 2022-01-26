package db

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"
)

type badgerDBBatch struct {
	db *badger.DB
	wb *badger.WriteBatch

	// Calling db.Flush twice panics, so we must keep track of whether we've
	// flushed already on our own. If Write can receive from the firstFlush
	// channel, then it's the first and only Flush call we should do.
	//
	// Upstream bug report:
	// https://github.com/dgraph-io/badger/issues/1394
	firstFlush chan struct{}
}

var _ Batch = (*badgerDBBatch)(nil)

func newBadgerDBBatch(db *badger.DB) *badgerDBBatch {
	wb := &badgerDBBatch{
		db:         db,
		wb:         db.NewWriteBatch(),
		firstFlush: make(chan struct{}, 1),
	}
	wb.firstFlush <- struct{}{}
	return wb
}

func (b *badgerDBBatch) assertOpen() {
	if b.wb == nil {
		panic("batch has been written or closed")
	}
}

func (b *badgerDBBatch) Set(key, value []byte) {
	b.assertOpen()
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	b.wb.Set(key, value)
}

func (b *badgerDBBatch) Delete(key []byte) {
	b.assertOpen()
	key = nonNilBytes(key)
	b.wb.Delete(key)
}

func (b *badgerDBBatch) Write() error {
	b.assertOpen()
	select {
	case <-b.firstFlush:
		return b.wb.Flush()
	default:
		return fmt.Errorf("batch already flushed")
	}
}

func (b *badgerDBBatch) WriteSync() error {
	err := b.Write()
	if err != nil {
		return err
	}
	return b.db.Sync()
}

func (b *badgerDBBatch) Close() {
	b.assertOpen()
	select {
	case <-b.firstFlush: // a Flush after Cancel panics too
	default:
	}
	b.wb.Cancel()
}
