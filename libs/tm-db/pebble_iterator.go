package db

import (
	"bytes"

	"github.com/cockroachdb/pebble"
)

type pebbleIterator struct {
	source    *pebble.Iterator
	start     []byte
	end       []byte
	isReverse bool
	isInvalid bool
	moved     bool
}

func newPebbleIterator(source *pebble.Iterator, start, end []byte, isReverse bool) *pebbleIterator {
	if isReverse {
		if end == nil {
			source.Last()
		} else {
			valid := source.SeekGE(end)
			if valid {
				eoakey := source.Key() // end or after key
				if bytes.Compare(end, eoakey) <= 0 {
					source.Prev()
				}
			} else {
				source.Last()
			}
		}
	} else {
		if start == nil {
			source.First()
		} else {
			source.SeekGE(start)
		}
	}
	return &pebbleIterator{
		source:    source,
		start:     start,
		end:       end,
		isReverse: isReverse,
		isInvalid: false,
		moved:     true,
	}
}

func (itr pebbleIterator) Domain() (start []byte, end []byte) {
	return itr.start, itr.end
}

func (itr pebbleIterator) Valid() bool {
	if itr.isInvalid {
		return false
	}
	// Panic on DB error.  No way to recover.
	itr.assertNoError()
	if !itr.source.Valid() {
		itr.isInvalid = true
		return false
	}

	// If key is end or past it, invalid.
	var start = itr.start
	var end = itr.end
	var key = itr.source.Key()

	if itr.isReverse {
		if start != nil && bytes.Compare(key, start) < 0 {
			itr.isInvalid = true
			return false
		}
	} else {
		if end != nil && bytes.Compare(end, key) <= 0 {
			itr.isInvalid = true
			return false
		}
	}

	// Valid
	return true
}

func (itr pebbleIterator) Next() {
	if itr.moved {
		itr.moved = false
		return
	}
	if itr.isReverse {
		itr.source.Prev()
	} else {
		itr.source.Next()
	}
}

func (itr pebbleIterator) Key() (key []byte) {
	itr.assertNoError()
	itr.assertIsValid()
	return cp(itr.source.Key())
}

func (itr pebbleIterator) Value() (value []byte) {
	itr.assertNoError()
	itr.assertIsValid()
	return cp(itr.source.Value())
}

func (itr pebbleIterator) Error() error {
	return itr.source.Error()
}

func (itr pebbleIterator) Close() {
	itr.source.Close()
}

func (itr pebbleIterator) assertNoError() {
	err := itr.source.Error()
	if err != nil {
		panic(err)
	}
}

func (itr pebbleIterator) assertIsValid() {
	if !itr.Valid() {
		panic("goLevelDBIterator is invalid")
	}
}
