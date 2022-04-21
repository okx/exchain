package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	dbm "github.com/tendermint/tm-db"
	"reflect"
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
	return sa.parent.Iterator(start, end)
}

// Iterator over a domain of keys in descending order. End is exclusive.
// Start must be less than end, or the Iterator is invalid.
// Iterator must be closed by caller.
func (sa StoreAdapter) ReverseIterator(start, end []byte) dbm.Iterator {
	return sa.parent.ReverseIterator(start, end)
}
