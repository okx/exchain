package db

import "github.com/cockroachdb/pebble"

type pebbleBatch struct {
	batch *pebble.Batch
}

var _ Batch = (*pebbleBatch)(nil)

func (pb *pebbleBatch) assertOpen() {
	if pb.batch == nil {
		panic("batch has been written or closed")
	}
}

// Set implements Batch.
func (pb *pebbleBatch) Set(key, value []byte) {
	pb.assertOpen()
	_ = pb.batch.Set(key, value, nil)
}

// Delete implements Batch.
func (pb *pebbleBatch) Delete(key []byte) {
	pb.assertOpen()
	_ = pb.batch.Delete(key, nil)
}

// Write implements Batch.
func (pb *pebbleBatch) Write() error {
	return pb.write(false)
}

// WriteSync implements Batch.
func (pb *pebbleBatch) WriteSync() error {
	return pb.write(true)
}

func (b *pebbleBatch) write(sync bool) error {
	b.assertOpen()
	err := b.batch.Commit(&pebble.WriteOptions{Sync: sync})
	if err != nil {
		return err
	}
	// Make sure batch cannot be used afterwards. Callers should still call Close(), for errors.
	b.Close()
	return nil
}

// Close implements Batch.
func (pb *pebbleBatch) Close() {
	if pb.batch != nil {
		pb.batch.Reset()
		pb.batch = nil
	}
}
