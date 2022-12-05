//go:build rocksdb
// +build rocksdb

package db

// #include "rocksdb/c.h"
// #include "tmrocksdb.h"
import "C"

import (
	"unsafe"

	"github.com/cosmos/gorocksdb"
)

func boolToChar(b bool) C.uchar {
	if b {
		return 1
	}
	return 0
}

type privateOptions struct {
	c *C.rocksdb_options_t
}

func enableUnorderedWrite(opts *gorocksdb.Options, enable bool) {
	myOpts := (*privateOptions)(unsafe.Pointer(opts))
	C.rocksdb_options_set_unordered_write(myOpts.c, boolToChar(enable))
}
