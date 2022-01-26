//go:build !rocksdb
// +build !rocksdb

package main

func ldbToRdb(_, _, _ string) {
	panic("Not supported rocksdb, must be built with rocksdb")
}

func rdbtoBdb(_, _, _ string) {
	panic("Not supported rocksdb, must be built with rocksdb")
}
