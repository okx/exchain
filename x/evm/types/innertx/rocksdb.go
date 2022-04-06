//go:build !rocksdb
// +build !rocksdb

package innertx

import (
	"errors"
	"fmt"

	ethvm "github.com/ethereum/go-ethereum/core/vm"
)

func newRocksDBCreator() ethvm.DBCreator {
	panic(errors.New(fmt.Sprintf("Rocks DB has not build")))
}
