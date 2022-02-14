// +build rocksdb

package types

import (
	"github.com/ethereum/go-ethereum/ethdb"
	tmdb "github.com/okex/exchain/libs/tm-db"
)

var _ ethdb.Iterator = (*WrapRocksDBIterator)(nil)

type WrapRocksDBIterator struct {
	*tmdb.RocksDBIterator
}

func NewWrapRocksDBIterator(db *tmdb.RocksDB, start, end []byte) *WrapRocksDBIterator {
	itr, _ := db.Iterator(start, end)
	return &WrapRocksDBIterator{itr.(*tmdb.RocksDBIterator)}
}

func (wrsdi *WrapRocksDBIterator) Next() bool {
	if wrsdi.Error() !=nil ||  !wrsdi.Valid() {
		return false
	}

	wrsdi.RocksDBIterator.Next()

	return true
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (wrsdi *WrapRocksDBIterator) Release() {
	wrsdi.RocksDBIterator = wrsdi.RocksDBIterator.Release()
}
