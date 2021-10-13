// +build !rocksdb

package main

import (
	dbm "github.com/tendermint/tm-db"
)

func compactRocksDB(db dbm.DB) {
	panic("unsupported rocksdb")
}
