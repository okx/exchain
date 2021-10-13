// +build rocksdb

package main

import (
	"github.com/tecbot/gorocksdb"
	dbm "github.com/tendermint/tm-db"
)

func compactRocksDB(db dbm.DB) {
	if rdb, ok := db.(*dbm.RocksDB); ok {
		for i := 0; i < 5; i++ {
			rdb.DB().CompactRange(gorocksdb.Range{})
		}
	}
}
