// +build rocksdb

package types

import (
	"bytes"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tecbot/gorocksdb"
)


var _ ethdb.Iterator = (*WrapRocksDBIterator)(nil)

type WrapRocksDBIterator struct {
	source     *gorocksdb.Iterator
	start, end []byte
	isReverse  bool
	isInvalid  bool
}

func NewWrapRocksDBIterator(source *gorocksdb.Iterator, start, end []byte, isReverse bool) *WrapRocksDBIterator {
	if isReverse {
		if end == nil {
			source.SeekToLast()
		} else {
			source.Seek(end)
			if source.Valid() {
				eoakey := moveSliceToBytes(source.Key()) // end or after key
				if bytes.Compare(end, eoakey) <= 0 {
					source.Prev()
				}
			} else {
				source.SeekToLast()
			}
		}
	} else {
		if start == nil {
			source.SeekToFirst()
		} else {
			source.Seek(start)
		}
	}
	return &WrapRocksDBIterator{
		source:    source,
		start:     start,
		end:       end,
		isReverse: isReverse,
		isInvalid: false,
	}
}

// Valid implements Iterator.
func (witr *WrapRocksDBIterator) Valid() bool {

	// Once invalid, forever invalid.
	if witr.isInvalid {
		return false
	}

	// Panic on DB error.  No way to recover.
	witr.assertNoError()

	// If source is invalid, invalid.
	if !witr.source.Valid() {
		witr.isInvalid = true
		return false
	}

	// If key is end or past it, invalid.
	var start = witr.start
	var end = witr.end
	var key = moveSliceToBytes(witr.source.Key())
	if witr.isReverse {
		if start != nil && bytes.Compare(key, start) < 0 {
			witr.isInvalid = true
			return false
		}
	} else {
		if end != nil && bytes.Compare(end, key) <= 0 {
			witr.isInvalid = true
			return false
		}
	}

	// It's valid.
	return true
}

func (witr *WrapRocksDBIterator) Next() bool {
	witr.assertNoError()
	witr.assertIsValid()
	if witr.isReverse {
		witr.source.Prev()
	} else {
		witr.source.Next()
	}

	return true
}

func (witr *WrapRocksDBIterator) Error() error {
	return witr.source.Err()
}

func (witr *WrapRocksDBIterator) Key() []byte {
	witr.assertNoError()
	witr.assertIsValid()
	return moveSliceToBytes(witr.source.Key())
}

func (witr *WrapRocksDBIterator) Value() []byte {
	witr.assertNoError()
	witr.assertIsValid()
	return moveSliceToBytes(witr.source.Value())
}

func (witr *WrapRocksDBIterator) Release() {
	witr.assertNoError()
	witr.assertIsValid()

	witr = NewWrapRocksDBIterator(witr.source, witr.start, witr.end, witr.isReverse)
}

// Close implements Iterator.
func (witr *WrapRocksDBIterator) Close() {
	witr.source.Close()
}

func (witr *WrapRocksDBIterator) assertNoError() {
	if err := witr.source.Err(); err != nil {
		panic(err)
	}
}

func (witr *WrapRocksDBIterator) assertIsValid() {
	if !witr.Valid() {
		panic("rocksDBIterator is invalid")
	}
}

// moveSliceToBytes will free the slice and copy out a go []byte
// This function can be applied on *Slice returned from Key() and Value()
// of an Iterator, because they are marked as freed.
func moveSliceToBytes(s *gorocksdb.Slice) []byte {
	defer s.Free()
	if !s.Exists() {
		return nil
	}
	v := make([]byte, len(s.Data()))
	copy(v, s.Data())
	return v
}
