package mpt

import (
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

var _ types.Iterator = (*mptIterator)(nil)

type mptIterator struct {
	// Domain
	start, end []byte

	key   []byte // The current key (mutable)
	value []byte // The current value (mutable)

	// Underlying store
	iterator *trie.Iterator

	valid bool
}

func newMptIterator(t ethstate.Trie, start, end []byte) *mptIterator {
	iter := &mptIterator{
		iterator: trie.NewIterator(t.NodeIterator(start)),

		start: types.Cp(start),
		end:   types.Cp(end),
		valid: true,
	}
	return iter
}

func (it *mptIterator) Domain() (start []byte, end []byte) {
	return it.start, it.end
}

func (it *mptIterator) Valid() bool {
	// return it.invalid
	return it.valid
}

func (it *mptIterator) Next() {
	if !it.iterator.Next() {
		it.valid = false
	}
}

func (it *mptIterator) Key() (key []byte) {
	return it.iterator.Key
}

func (it *mptIterator) Value() (value []byte) {
	return it.iterator.Value
}

func (it *mptIterator) Error() error {
	return it.iterator.Err
}

func (it *mptIterator) Close() {
	return
}
