// +build rocksdb

package main

import (
	"github.com/tecbot/gorocksdb"
	dbm "github.com/tendermint/tm-db"
)

func init() {
	dbCompactor := func(db dbm.DB) {
		if rdb, ok := db.(*dbm.RocksDB); ok {
			rdb.DB().CompactRange(gorocksdb.Range{})
		}
	}

	registerDBCompactor(dbm.RocksDBBackend, dbCompactor)
}
