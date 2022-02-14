// +build rocksdb

package types

import (
	"github.com/ethereum/go-ethereum/ethdb"
	tmdb "github.com/okex/exchain/libs/tm-db"
	"github.com/tecbot/gorocksdb"
)

var _ ethdb.Batch = (*WrapRocksDBBatch)(nil)

type WrapRocksDBBatch struct {
	*tmdb.RocksDBBatch
}

func (wrsdbb *WrapRocksDBBatch) Put(key []byte, value []byte) error {
	wrsdbb.Set(key, value)
	return nil
}

func (wrsdbb *WrapRocksDBBatch) Delete(key []byte) error {
	wrsdbb.RocksDBBatch.Delete(key)
	return nil
}

func (wrsdbb *WrapRocksDBBatch) ValueSize() int {
	return wrsdbb.Size()
}

func NewWrapRocksDBBatch(db *tmdb.RocksDB) *WrapRocksDBBatch {
	return &WrapRocksDBBatch{tmdb.NewRocksDBBatch(db)}
}

// Replay replays the batch contents.
func (wrsdbb *WrapRocksDBBatch) Replay(w ethdb.KeyValueWriter) error {
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
