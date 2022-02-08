//go:build !rocksdb
// +build !rocksdb

package main

func LtoR(_, _, _ string) {
	panic("Not supported rocksdb, must be built with rocksdb")
}

func RtoL(_, _, _ string) {
	panic("Not supported rocksdb, must be built with rocksdb")
}
