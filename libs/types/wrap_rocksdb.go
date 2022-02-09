//go:build rocksdb
// +build rocksdb

package types

import (
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tecbot/gorocksdb"
	db "github.com/tendermint/tm-db"
)

func init() {
	dbCreator := func(name string, dir string) (ethdb.KeyValueStore, error) {
		return NewWrapRocksDB(name, dir)
	}
	registerDBCreator(RocksDBBackend, dbCreator, false)
}

type WrapRocksDB struct {
	*db.RocksDB
}

func NewWrapRocksDB(name string, dir string) (*WrapRocksDB, error) {
	rdb, err := db.NewRocksDB(name, dir)

	return &WrapRocksDB{rdb}, err
}

func (db *WrapRocksDB) Put(key []byte, value []byte) error {
	return db.Set(key, value)
}

func (db *WrapRocksDB) NewBatch() ethdb.Batch {
	return NewWrapRocksDBBatch(db)
}

func (db *WrapRocksDB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	ro := gorocksdb.NewDefaultReadOptions()
	itr := db.DB().NewIterator(ro)

	st := append(prefix, start...)
	return NewWrapRocksDBIterator(itr, st, nil, false)
}

func (db *WrapRocksDB) Stat(property string) (string, error) {
	return db.DB().GetProperty(property), nil
}

func (db *WrapRocksDB) Compact(start []byte, limit []byte) error {
	db.DB().CompactRange(gorocksdb.Range{Start: start, Limit: limit})
	return nil
}
