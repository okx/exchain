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
}

func newPebbleIterator(iter *pebble.Iterator, start, end []byte, reverse bool) Iterator {
	if reverse {
		if end == nil {
			iter.Last()
		} else {
			iter.SeekLT(end)
		}
	} else {
		if start == nil {
			iter.First()
		} else {
			iter.SeekGE(start)
		}
	}
	return &pebbleIterator{
		source:    iter,
		start:     start,
		end:       end,
		isReverse: reverse,
	}
}

func (itr *pebbleIterator) Domain() ([]byte, []byte) {
	return itr.start, itr.end
}

// Valid implements Iterator.
func (itr *pebbleIterator) Valid() bool {

	// Once invalid, forever invalid.
	if itr.isInvalid {
		return false
	}

	// Panic on DB error.  No way to recover.
	itr.assertNoError()

	// If source is invalid, invalid.
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

// Key returns a copy of the current key.
func (itr *pebbleIterator) Key() []byte {
	itr.assertNoError()
	itr.assertIsValid()
	return cp(itr.source.Key())
}

// Value returns a copy of the current value.
func (itr *pebbleIterator) Value() []byte {
	itr.assertNoError()
	itr.assertIsValid()
	return cp(itr.source.Value())
}

// Next implements Iterator.
func (itr *pebbleIterator) Next() {
	itr.assertNoError()
	itr.assertIsValid()
	if itr.isReverse {
		itr.source.Prev()
	} else {
		itr.source.Next()
	}
}

// Error implements Iterator.
func (itr *pebbleIterator) Error() error {
	return itr.source.Error()
}

// Close implements Iterator.
func (itr *pebbleIterator) Close() {
	_ = itr.source.Close()
}

// Seek seek
/*
TODO：
	关于seek, seek方法将iterator定位到大于等于key的一个位置，如果存在这样一个key则返true，否则返回false。
	在leveldb中，seek返回false时，可以继续调用Prev()方法向前遍历。
	pebble中seek方法分为SeekGE(key)和SeekLT(key)，且当这两个方法返回false时，无法再进一步调用Next()或者Prev()方法向后或者向前遍历。
	因此在封装 *pebbleIt.Next() 方法时作了特殊处理。
*/

func (itr *pebbleIterator) assertNoError() {
	err := itr.source.Error()
	if err != nil {
		panic(err)
	}
}

func (itr *pebbleIterator) assertIsValid() {
	if !itr.Valid() {
		panic("pebbleIterator is invalid")
	}
}
