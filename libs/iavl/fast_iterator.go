package iavl

import (
	"bytes"
	"errors"
	"sort"

	dbm "github.com/okex/exchain/libs/tm-db"
)

var errFastIteratorNilNdbGiven = errors.New("fast iterator must be created with a nodedb but it was nil")

const latestNodeDBVersion = 1<<63 - 1

// FastIterator is a dbm.Iterator for ImmutableTree
// it iterates over the latest state via fast nodes,
// taking advantage of keys being located in sequence in the underlying database.
type FastIterator struct {
	start, end []byte

	valid bool

	ascending bool

	err error

	ndb *nodeDB

	nextKey []byte

	nextVal []byte

	nextUnsavedNodeIdx int

	unsavedNodeDBFastNodesToSort []string
	unsavedNodeDBFastNodesMap    map[string]int64

	fastIterator dbm.Iterator
}

var _ dbm.Iterator = (*FastIterator)(nil)

func NewFastIterator(start, end []byte, ascending bool, ndb *nodeDB) *FastIterator {
	iter := &FastIterator{
		start:                     start,
		end:                       end,
		err:                       nil,
		ascending:                 ascending,
		ndb:                       ndb,
		fastIterator:              nil,
		unsavedNodeDBFastNodesMap: make(map[string]int64),
	}
	if iter.ndb == nil {
		iter.err = errFastIteratorNilNdbGiven
		iter.valid = false
		return iter
	}

	// We need to ensure that we iterate over saved and unsaved state in order.
	// The strategy is to sort unsaved nodes, the fast node on disk are already sorted.
	// Then, we keep a pointer to both the unsaved and saved nodes, and iterate over them in order efficiently.

	// perPersistFastNode
	for _, fastNode := range ndb.prePersistFastNode.additions {
		if start != nil && bytes.Compare(fastNode.key, start) < 0 {
			continue
		}

		if end != nil && bytes.Compare(fastNode.key, end) >= 0 {
			continue
		}

		iter.unsavedNodeDBFastNodesToSort = append(iter.unsavedNodeDBFastNodesToSort, string(fastNode.key))
		iter.unsavedNodeDBFastNodesMap[string(fastNode.key)] = latestNodeDBVersion
	}
	// tpfv
	for i := len(ndb.tpfv.versions) - 1; i >= 0; i-- {
		for _, fn := range ndb.tpfv.fncMap[ndb.tpfv.versions[i]].additions {
			if start != nil && bytes.Compare(fn.key, start) < 0 {
				continue
			}

			if end != nil && bytes.Compare(fn.key, end) >= 0 {
				continue
			}

			if _, ok := iter.unsavedNodeDBFastNodesMap[string(fn.key)]; ok {
				continue
			}
			if _, ok := ndb.prePersistFastNode.removals[string(fn.key)]; ok {
				continue
			}

			if ndb.tpfv.checkRemovals(string(fn.key)) {
				continue
			}

			iter.unsavedNodeDBFastNodesToSort = append(iter.unsavedNodeDBFastNodesToSort, string(fn.key))
			iter.unsavedNodeDBFastNodesMap[string(fn.key)] = ndb.tpfv.versions[i]
		}
	}

	sort.Slice(iter.unsavedNodeDBFastNodesToSort, func(i, j int) bool {
		if ascending {
			return iter.unsavedNodeDBFastNodesToSort[i] < iter.unsavedNodeDBFastNodesToSort[j]
		}
		return iter.unsavedNodeDBFastNodesToSort[i] > iter.unsavedNodeDBFastNodesToSort[j]
	})

	// Move iterator before the first element
	iter.Next()
	return iter
}

// Domain implements dbm.Iterator.
// Maps the underlying nodedb iterator domain, to the 'logical' keys involved.
func (iter *FastIterator) Domain() ([]byte, []byte) {
	if iter.fastIterator == nil {
		return iter.start, iter.end
	}

	start, end := iter.fastIterator.Domain()

	if start != nil {
		start = start[1:]
		if len(start) == 0 {
			start = nil
		}
	}

	if end != nil {
		end = end[1:]
		if len(end) == 0 {
			end = nil
		}
	}

	return start, end
}

// Valid implements dbm.Iterator.
func (iter *FastIterator) Valid() bool {
	if iter.start != nil && iter.end != nil {
		if bytes.Compare(iter.end, iter.start) != 1 {
			return false
		}
	}

	return (iter.fastIterator != nil && iter.fastIterator.Valid()) || iter.nextUnsavedNodeIdx < len(iter.unsavedNodeDBFastNodesToSort) || (iter.nextKey != nil && iter.nextVal != nil)
	//return iter.fastIterator != nil && iter.fastIterator.Valid() && iter.valid
}

// Key implements dbm.Iterator
func (iter *FastIterator) Key() []byte {
	return iter.nextKey
}

// Value implements dbm.Iterator
func (iter *FastIterator) Value() []byte {
	return iter.nextVal
}

// Next implements dbm.Iterator
func (iter *FastIterator) Next() {
	if iter.ndb == nil {
		iter.err = errFastIteratorNilNdbGiven
		iter.valid = false
		return
	}

	if iter.fastIterator == nil {
		iter.fastIterator, iter.err = iter.ndb.getFastIterator(iter.start, iter.end, iter.ascending)
		iter.valid = true
	}
	if iter.fastIterator.Valid() && iter.nextUnsavedNodeIdx < len(iter.unsavedNodeDBFastNodesToSort) {
		diskKeyStr := string(iter.fastIterator.Key()[1:])

		if iter.ndb.prePersistFastNode.checkRemovals(diskKeyStr) {
			// If next fast node from disk is to be removed, skip it.
			iter.fastIterator.Next()
			iter.Next()
			return
		}
		if !iter.ndb.prePersistFastNode.checkAdditions(diskKeyStr) &&
			iter.ndb.tpfv.checkRemovals(diskKeyStr) {
			// If next fast node from disk is to be removed, skip it.
			iter.fastIterator.Next()
			iter.Next()
			return
		}

		nextUnsavedKey := iter.unsavedNodeDBFastNodesToSort[iter.nextUnsavedNodeIdx]
		nextUnsavedNodeVersion := iter.unsavedNodeDBFastNodesMap[nextUnsavedKey]
		var nextUnsavedNode *FastNode
		if nextUnsavedNodeVersion == latestNodeDBVersion {
			nextUnsavedNode = iter.ndb.prePersistFastNode.additions[nextUnsavedKey]
		} else {
			nextUnsavedNode = iter.ndb.tpfv.fncMap[nextUnsavedNodeVersion].additions[nextUnsavedKey]
		}

		var isUnsavedNext bool
		if iter.ascending {
			isUnsavedNext = diskKeyStr >= nextUnsavedKey
		} else {
			isUnsavedNext = diskKeyStr <= nextUnsavedKey
		}

		if isUnsavedNext {
			// Unsaved node is next
			if diskKeyStr == nextUnsavedKey {
				// Unsaved update prevails over saved copy so we skip the copy from disk
				iter.fastIterator.Next()
			}

			iter.nextKey = nextUnsavedNode.key
			iter.nextVal = nextUnsavedNode.value

			iter.nextUnsavedNodeIdx++
			return
		}
		// Disk node is next
		iter.loadFromDisk()

		iter.fastIterator.Next()
		return
	}

	// if only nodes on disk are left, we return them
	if iter.fastIterator.Valid() {
		nextKey := string(iter.fastIterator.Key()[1:])
		if iter.ndb.prePersistFastNode.checkRemovals(nextKey) {
			// If next fast node from disk is to be removed, skip it.
			iter.fastIterator.Next()
			iter.Next()
			return
		}
		if !iter.ndb.prePersistFastNode.checkAdditions(nextKey) &&
			iter.ndb.tpfv.checkRemovals(nextKey) {
			// If next fast node from disk is to be removed, skip it.
			iter.fastIterator.Next()
			iter.Next()
			return
		}
		// Disk node is next
		iter.loadFromDisk()

		iter.fastIterator.Next()
		return
	}

	// if only unsaved nodes are left, we can just iterate
	if iter.nextUnsavedNodeIdx < len(iter.unsavedNodeDBFastNodesToSort) {
		nextUnsavedKey := iter.unsavedNodeDBFastNodesToSort[iter.nextUnsavedNodeIdx]
		nextUnsavedNodeVersion := iter.unsavedNodeDBFastNodesMap[nextUnsavedKey]
		var nextUnsavedNode *FastNode
		if nextUnsavedNodeVersion == latestNodeDBVersion {
			nextUnsavedNode = iter.ndb.prePersistFastNode.additions[nextUnsavedKey]
		} else {
			nextUnsavedNode = iter.ndb.tpfv.fncMap[nextUnsavedNodeVersion].additions[nextUnsavedKey]
		}

		iter.nextKey = nextUnsavedNode.key
		iter.nextVal = nextUnsavedNode.value

		iter.nextUnsavedNodeIdx++
		return
	}

	iter.nextKey = nil
	iter.nextVal = nil
}

// Close implements dbm.Iterator
func (iter *FastIterator) Close() {
	if iter.fastIterator != nil {
		iter.fastIterator.Close()
	}
	iter.valid = false
	iter.fastIterator = nil
}

// Error implements dbm.Iterator
func (iter *FastIterator) Error() error {
	return iter.err
}

func (iter *FastIterator) loadFromDisk() {
	var nextFastNode *FastNode
	if iter.err == nil {
		iter.err = iter.fastIterator.Error()
	}
	iter.valid = iter.valid && iter.fastIterator.Valid()
	if iter.valid {
		nextFastNode, iter.err = DeserializeFastNode(iter.fastIterator.Key()[1:], iter.fastIterator.Value())
		iter.valid = iter.err == nil
		iter.nextKey = nextFastNode.key
		iter.nextVal = nextFastNode.value
	}
}
