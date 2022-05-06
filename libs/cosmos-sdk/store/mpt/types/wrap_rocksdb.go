//go:build rocksdb
// +build rocksdb

package types

import (
	"github.com/ethereum/go-ethereum/ethdb"
	tmdb "github.com/okex/exchain/libs/tm-db"
	"github.com/pkg/errors"
	"github.com/tecbot/gorocksdb"
)

//------------------------------------------
//	Register go-ethereum gorocksdb
//------------------------------------------
func init() {
	dbCreator := func(name string, dir string) (ethdb.KeyValueStore, error) {
		return NewWrapRocksDB(name, dir)
	}
	registerDBCreator(RocksDBBackend, dbCreator, false)
}

type WrapRocksDB struct {
	*tmdb.RocksDB
}

func NewWrapRocksDB(name string, dir string) (*WrapRocksDB, error) {
	rdb, err := tmdb.NewRocksDB(name, dir)
	return &WrapRocksDB{rdb}, err
}

func (db *WrapRocksDB) Put(key []byte, value []byte) error {
	return db.Set(key, value)
}

func (db *WrapRocksDB) NewBatch() ethdb.Batch {
	return NewWrapRocksDBBatch(db.RocksDB)
}

func (db *WrapRocksDB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	limit := bytesPrefix(prefix)
	return NewWrapRocksDBIterator(db.RocksDB, append(prefix, start...), limit)
}

func (db *WrapRocksDB) Stat(property string) (string, error) {
	stats := db.RocksDB.Stats()
	if pro, ok := stats[property]; ok {
		return pro, nil
	}

	return "", errors.New("property not exist")
}

func (db *WrapRocksDB) Compact(start []byte, limit []byte) error {
	db.DB().CompactRange(gorocksdb.Range{Start: start, Limit: limit})
	return nil
}

// BytesPrefix returns key range that satisfy the given prefix.
// This only applicable for the standard 'bytes comparer'.
func bytesPrefix(prefix []byte) []byte {
	var limit []byte
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c < 0xff {
			limit = make([]byte, i+1)
			copy(limit, prefix)
			limit[i] = c + 1
			break
		}
	}
	return limit
}
