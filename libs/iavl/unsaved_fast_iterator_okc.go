package iavl

import dbm "github.com/okx/okbchain/libs/tm-db"

type UnsavedFastIteratorWithCache struct {
	*UnsavedFastIterator
}

var _ dbm.Iterator = (*UnsavedFastIteratorWithCache)(nil)

func NewUnsavedFastIteratorWithCache(start, end []byte, ascending bool, ndb *nodeDB, fncIn *fastNodeChanges) *UnsavedFastIteratorWithCache {
	iter := &UnsavedFastIteratorWithCache{
		UnsavedFastIterator: &UnsavedFastIterator{},
	}

	fnc := fncIn.clone()
	if ndb == nil || fnc.additions == nil || fnc.removals == nil {
		iter.UnsavedFastIterator = newUnsavedFastIterator(start, end, ascending, ndb, fnc.additions, fnc.removals)
		return iter
	}

	fnc.mergePrev(ndb.prePersistFastNode)

	if ndb.tpfv != nil {
		fnc = ndb.tpfv.expand(fnc)
	}

	iter.UnsavedFastIterator = newUnsavedFastIterator(start, end, ascending, ndb, fnc.additions, fnc.removals)

	return iter
}
