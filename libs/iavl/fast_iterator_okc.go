package iavl

import dbm "github.com/okx/okbchain/libs/tm-db"

type FastIteratorWithCache struct {
	*UnsavedFastIterator
}

var _ dbm.Iterator = (*FastIteratorWithCache)(nil)

func NewFastIteratorWithCache(start, end []byte, ascending bool, ndb *nodeDB) *FastIteratorWithCache {
	iter := &FastIteratorWithCache{
		UnsavedFastIterator: &UnsavedFastIterator{},
	}

	if ndb == nil {
		iter.UnsavedFastIterator = newUnsavedFastIterator(start, end, ascending, ndb, nil, nil)
		return iter
	}
	var fnc *fastNodeChanges

	if ndb.tpfv != nil {
		fnc = ndb.tpfv.expand(ndb.prePersistFastNode)
	} else {
		fnc = ndb.prePersistFastNode.clone()
	}

	iter.UnsavedFastIterator = newUnsavedFastIterator(start, end, ascending, ndb, fnc.additions, fnc.removals)

	return iter
}
