//go:build rocksdb
// +build rocksdb

package types

import (
	"github.com/ethereum/go-ethereum/ethdb"
	tmdb "github.com/okx/okbchain/libs/tm-db"
)

var _ ethdb.Batch = (*WrapRocksDBBatch)(nil)

type record struct {
	key   []byte
	value []byte
}

type WrapRocksDBBatch struct {
	db      *tmdb.RocksDB
	records []record
	size    int
}

func NewWrapRocksDBBatch(db *tmdb.RocksDB) *WrapRocksDBBatch {
	return &WrapRocksDBBatch{db: db}
}

func (wrb *WrapRocksDBBatch) Put(key []byte, value []byte) error {
	wrb.records = append(wrb.records, record{
		key:   key,
		value: nonNilBytes(value),
	})
	wrb.size += len(value)
	return nil
}

func (wrb *WrapRocksDBBatch) Delete(key []byte) error {
	wrb.records = append(wrb.records, record{
		key: key,
	})
	wrb.size += len(key)
	return nil
}

func (wrb *WrapRocksDBBatch) ValueSize() int {
	return wrb.size
}

func (wrb *WrapRocksDBBatch) Write() error {
	batch := tmdb.NewRocksDBBatch(wrb.db)
	defer batch.Close()
	for _, rcd := range wrb.records {
		if rcd.value != nil {
			batch.Set(rcd.key, rcd.value)
		} else {
			batch.Delete(rcd.key)
		}
	}

	return batch.Write()
}

// Replay replays the batch contents.
func (wrb *WrapRocksDBBatch) Replay(w ethdb.KeyValueWriter) error {
	var err error
	for _, rcd := range wrb.records {
		if rcd.value != nil {
			err = w.Put(rcd.key, rcd.value)
		} else {
			err = w.Delete(rcd.key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (wrb *WrapRocksDBBatch) Reset() {
	wrb.records = wrb.records[:0]
}

func nonNilBytes(bz []byte) []byte {
	if bz == nil {
		return []byte{}
	}
	return bz
}
