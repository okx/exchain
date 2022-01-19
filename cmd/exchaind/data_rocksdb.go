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

	pairCounter := 0

	smaleCounter := 0

	midCounter := 0

	largeCounter := 0

	smale := 128

	mid := 1024

	large := 16384

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

	keySize := iter.Key().Size()

	valueSize := iter.Value().Size()

	if valueSize < mid {
		smaleCounter++
	}

	if valueSize > mid && valueSize < large {
		midCounter++
	}

	if valueSize > large {
		largeCounter++
	}

	for ; iter.Valid(); iter.Next() {
		bdb.Set(iter.Key(), iter.Value())
	}
	pairCounter++

	iter.Close()
	log.Printf("convert %s(rocksdb => badgerdb) end.\n", name)
	log.Printf("pairs %s", pairCounter)
	log.Printf("smale  %s", smaleCounter)
	log.Printf("mid %s", midCounter)
	log.Printf("large %s", largeCounter)

	//log.Printf("compact %s start...\n", name)
	////bdb.DB()(gorocksdb.Range{})
	//log.Printf("compact %s end.\n", name)
}
