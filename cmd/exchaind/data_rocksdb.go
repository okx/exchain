//go:build rocksdb
// +build rocksdb

package main

import (
	"log"

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

func LtoR(name, fromDir, toDir string) {
	log.Printf("convert %s(rocksdb => badgerdb) start...\n", name)

	rdb, err := dbm.NewRocksDB(name, fromDir)
	if err != nil {
		panic(err)
	}

	bdb, err := dbm.NewBadgerDB(name, toDir)
	if err != nil {
		panic(err)
	}

	iter, err := rdb.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	for ; iter.Valid(); iter.Next() {
		bdb.Set(iter.Key(), iter.Value())
	}
	iter.Close()
	log.Printf("convert %s(rocksdb => badgerdb) end.\n", name)

	//log.Printf("compact %s start...\n", name)
	////bdb.DB()(gorocksdb.Range{})
	//log.Printf("compact %s end.\n", name)
}
