//go:build rocksdb
// +build rocksdb

package innertx

import (
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	dbm "github.com/okex/exchain/libs/tm-db"
)

func newRocksDBCreator() ethvm.DBCreator {
	return func(name string, dir string) (ethvm.OKDB, error) {
		return dbm.NewRocksDB(name, dir)
	}
}
