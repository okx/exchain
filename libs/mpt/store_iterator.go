package mpt

import (
	"sync"

	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	tmkv "github.com/okex/exchain/libs/tendermint/libs/kv"
)

var _ types.Iterator = (*mptIterator)(nil)

type mptIterator struct {
	// Domain
	start, end []byte

	key   []byte // The current key (mutable)
	value []byte // The current value (mutable)

	// Underlying store
	iterator *trie.Iterator

	// Channel to push iteration values.
	iterCh chan tmkv.Pair

	// Close this to release goroutine.
	quitCh chan struct{}

	// Close this to signal that state is initialized.
	initCh chan struct{}

	mtx sync.Mutex

	ascending bool // Iteration order

	invalid bool // True once, true forever (mutable)
}

func newMptIterator(t ethstate.Trie, start, end []byte) *mptIterator {
	iter := &mptIterator{
		iterator: trie.NewIterator(t.NodeIterator(start)),

		start:     types.Cp(start),
		end:       types.Cp(end),
		ascending: true,
		iterCh:    make(chan tmkv.Pair), // Set capacity > 0?
		quitCh:    make(chan struct{}),
		initCh:    make(chan struct{}),
	}
	return iter
}

func (it *mptIterator) Domain() (start []byte, end []byte) {
	return it.start, it.end
}

func (it *mptIterator) Valid() bool {
	// return it.invalid
	return false
}

func (it *mptIterator) Next() {
	it.iterator.Next()
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
