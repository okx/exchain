package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	dbm "github.com/tendermint/tm-db"
)

type StoreAdapter struct {
	parent prefix.Store
}

func NewStoreAdapter(parent prefix.Store) StoreAdapter {
	return StoreAdapter{parent: parent}
}

func (sa StoreAdapter) Get(key []byte) []byte {
	return sa.parent.Get(key)
}

func (sa StoreAdapter) Set(key, value []byte) {
	sa.parent.Set(key, value)
}
func (sa StoreAdapter) Delete(key []byte) {
	sa.parent.Delete(key)
}

// Iterator over a domain of keys in ascending order. End is exclusive.
// Start must be less than end, or the Iterator is invalid.
// Iterator must be closed by caller.
// To iterate over entire domain, use store.Iterator(nil, nil)
func (sa StoreAdapter) Iterator(start, end []byte) dbm.Iterator {
	iter := sa.parent.Iterator(start, end)
	adapter := newIteratorAdapter(iter)
	return adapter
}

// Iterator over a domain of keys in descending order. End is exclusive.
// Start must be less than end, or the Iterator is invalid.
// Iterator must be closed by caller.
func (sa StoreAdapter) ReverseIterator(start, end []byte) dbm.Iterator {
	iter := sa.parent.ReverseIterator(start, end)
	adapter := newIteratorAdapter(iter)
	return adapter
}

type iteratorAdapter struct {
	parent types.Iterator
}

func newIteratorAdapter(iter types.Iterator) *iteratorAdapter {
	return &iteratorAdapter{parent: iter}
}

// Domain returns the start (inclusive) and end (exclusive) limits of the iterator.
// CONTRACT: start, end readonly []byte
func (iter *iteratorAdapter) Domain() (start []byte, end []byte) {
	return iter.parent.Domain()
}

// Valid returns whether the current iterator is valid. Once invalid, the Iterator remains
// invalid forever.
func (iter *iteratorAdapter) Valid() bool {
	return iter.parent.Valid()
}

// Next moves the iterator to the next key in the database, as defined by order of iteration.
// If Valid returns false, this method will panic.
func (iter *iteratorAdapter) Next() {
	iter.parent.Next()
}

// Key returns the key at the current position. Panics if the iterator is invalid.
// CONTRACT: key readonly []byte
func (iter *iteratorAdapter) Key() (key []byte) {
	return iter.parent.Key()
}

// Value returns the value at the current position. Panics if the iterator is invalid.
// CONTRACT: value readonly []byte
func (iter *iteratorAdapter) Value() (value []byte) {
	return iter.parent.Value()
}

// Error returns the last error encountered by the iterator, if any.
func (iter *iteratorAdapter) Error() error {
	return iter.parent.Error()
}

// Close closes the iterator, relasing any allocated resources.
func (iter *iteratorAdapter) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("couldn't close db iteratorAdapter : %v", r)
			return
		}
	}()
	iter.parent.Close()
	return
}
