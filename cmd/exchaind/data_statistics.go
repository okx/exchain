//go:build rocksdb
// +build rocksdb

package main

import (
	"log"

	dbm "github.com/tendermint/tm-db"
)

func Statistic(name, fromDir string) {

	log.Printf("statistics started")

	pairCounter := 0

	smallCounter := 0

	midCounter := 0

	largeCounter := 0

	small := 128

	mid := 1024

	large := 16384

	rdb, err := dbm.NewRocksDB(name, fromDir)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	iter, err := rdb.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	for ; iter.Valid(); iter.Next() {
		valueSize := len(iter.Value())

		pairCounter++

		if valueSize < small {
			smallCounter++
		}

		if valueSize > mid && valueSize < large {
			midCounter++
		}

		if valueSize > large {
			largeCounter++
		}
	}

	iter.Close()

	log.Printf("pairs count: %s", pairCounter)
	log.Printf("value small count: %s", smallCounter)
	log.Printf("value mid count: %s", midCounter)
	log.Printf("value large count: %s", largeCounter)
}
