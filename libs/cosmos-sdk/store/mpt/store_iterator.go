package mpt

import (
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
)

var _ types.Iterator = (*mptIterator)(nil)

type mptIterator struct {
	// Domain
	start, end []byte

	// Underlying store
	iterator *trie.Iterator
	trie     ethstate.Trie

	valid bool
}

func newMptIterator(t ethstate.Trie, start, end []byte) types.Iterator {
	return newWrapIterator(t, start, end)
}

func newOriginIterator(t ethstate.Trie, start, end []byte) *mptIterator {
	iter := &mptIterator{
		iterator: trie.NewIterator(t.NodeIterator(start)),
		trie:     t,
		start:    types.Cp(start),
		end:      nil, // enforce end is nil, because trie iterator origin key is out of order
		valid:    true,
	}
	iter.Next()
	return iter
}

func (it *mptIterator) Domain() (start []byte, end []byte) {
	return it.start, it.end
}

func (it *mptIterator) Valid() bool {
	return it.valid
}

func (it *mptIterator) Next() {
	if !it.iterator.Next() || it.iterator.Key == nil {
		it.valid = false
	}
}

func (it *mptIterator) Key() []byte {
	key := it.iterator.Key
	originKey := it.trie.GetKey(key)
	return originKey
}

func (it *mptIterator) Value() []byte {
	value := it.iterator.Value
	return value
}

func (it *mptIterator) Error() error {
	return it.iterator.Err
}

func (it *mptIterator) Close() {
	return
}
