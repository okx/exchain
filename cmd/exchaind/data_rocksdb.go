//go:build rocksdb
// +build rocksdb

package main

import (
	"github.com/okex/exchain/libs/tm-db/tikv"
	"log"

	"github.com/cosmos/gorocksdb"
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

func R2TiKV(name, fromDir string) {
	log.Printf("convert %s(rocksdb => tikv) start...\n", name)

	rdb, err := dbm.NewRocksDB(name, fromDir)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	iter, err := rdb.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	tidb, err := tikv.NewTiKV("", "127.0.0.1:2379")
	counter := 0
	const commitGap = 50000

	keys := make([][]byte, 0)
	values := make([][]byte, 0)
	for ; iter.Valid(); iter.Next() {
		if counter%commitGap == 0 {
			log.Printf("convert %v ...\n", counter)
		}
		k := iter.Key()
		keys = append(keys, k)
		v := iter.Value()
		values = append(values, v)
		if len(k) > 100 {
			err = tidb.(*tikv.TiKV).BatchSet(keys, values)
			if err != nil {
				panic(err)
			}
			keys = keys[:0]
			values = values[:0]
		}

		//tidb.Set(iter.Key(), iter.Value())
		counter++
	}
	err = tidb.(*tikv.TiKV).BatchSet(keys, values)
	if err != nil {
		panic(err)
	}
	log.Printf("convert %v done \n", counter)
	iter.Close()

	log.Printf("convert %s(rocksdb => tikv) end.\n", name)
}
