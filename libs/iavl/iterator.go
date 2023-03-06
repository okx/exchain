package iavl

// NOTE: This file favors int64 as opposed to int for size/counts.
// The Tree on the other hand favors int.  This is intentional.

import (
	"errors"

	dbm "github.com/okx/okbchain/libs/tm-db"
)

var errIteratorNilTreeGiven = errors.New("iterator must be created with an immutable tree but the tree was nil")

// Iterator is a dbm.Iterator for ImmutableTree
type Iterator struct {
	start, end []byte

	key, value []byte

	valid bool

	err error

	t *traversal
}

var _ dbm.Iterator = (*Iterator)(nil)

// Returns a new iterator over the immutable tree. If the tree is nil, the iterator will be invalid.
func NewIterator(start, end []byte, ascending bool, tree *ImmutableTree) dbm.Iterator {
	iter := &Iterator{
		start: start,
		end:   end,
	}

	if tree == nil {
		iter.err = errIteratorNilTreeGiven
	} else {
		iter.valid = true
		iter.t = tree.root.newTraversal(tree, start, end, ascending, false, false)
		// Move iterator before the first element
		iter.Next()
	}
	return iter
}

// Domain implements dbm.Iterator.
func (iter *Iterator) Domain() ([]byte, []byte) {
	return iter.start, iter.end
}

// Valid implements dbm.Iterator.
func (iter *Iterator) Valid() bool {
	return iter.valid
}

// Key implements dbm.Iterator
func (iter *Iterator) Key() []byte {
	return iter.key
}

// Value implements dbm.Iterator
func (iter *Iterator) Value() []byte {
	return iter.value
}

// Next implements dbm.Iterator
func (iter *Iterator) Next() {
	if iter.t == nil {
		return
	}

	node := iter.t.next()
	// TODO: double-check if this error is correctly handled.
	if node == nil {
		iter.t = nil
		iter.valid = false
		return
	}

	if node.height == 0 {
		iter.key, iter.value = node.key, node.value
		return
	}

	iter.Next()
}

// Close implements dbm.Iterator
func (iter *Iterator) Close() {
	iter.t = nil
	iter.valid = false
}

// Error implements dbm.Iterator
func (iter *Iterator) Error() error {
	return iter.err
}

// IsFast returnts true if iterator uses fast strategy
func (iter *Iterator) IsFast() bool {
	return false
}
