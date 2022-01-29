// +build rocksdb

package main

import (
	"log"

	"github.com/tecbot/gorocksdb"
	dbm "github.com/okex/exchain/libs/tm-db"
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
	log.Printf("convert %s(goleveldb => rocksdb) start...\n", name)

	ldb, err := dbm.NewGoLevelDB(name, fromDir)
	if err != nil {
		panic(err)
	}

	rdb, err := dbm.NewRocksDB(name, toDir)
	if err != nil {
		panic(err)
	}

	iter, err := ldb.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	for ; iter.Valid(); iter.Next() {
		rdb.Set(iter.Key(), iter.Value())
	}
	iter.Close()
	log.Printf("convert %s(goleveldb => rocksdb) end.\n", name)

	log.Printf("compact %s start...\n", name)
	rdb.DB().CompactRange(gorocksdb.Range{})
	log.Printf("compact %s end.\n", name)
}
