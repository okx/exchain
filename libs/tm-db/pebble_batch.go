package db

import "github.com/cockroachdb/pebble"

type pebbleBatch struct {
	b *pebble.Batch
}

func (p pebbleBatch) Set(key, value []byte) {
	p.b.Set(key, value, nil)
}

func (p pebbleBatch) Delete(key []byte) {
	p.b.Delete(key, nil)
}

func (p pebbleBatch) Write() error {
	return p.b.Commit(pebble.NoSync)
}

func (p pebbleBatch) WriteSync() error {
	return p.b.Commit(pebble.Sync)
}

func (p pebbleBatch) Close() {
	p.b.Close()
}
