package iavl

import dbm "github.com/okex/exchain/libs/tm-db"

type FastIteratorOKC struct {
	*UnsavedFastIterator
}

var _ dbm.Iterator = (*FastIterator)(nil)

func NewFastIteratorOKC(start, end []byte, ascending bool, ndb *nodeDB) *FastIteratorOKC {
	iter := &FastIteratorOKC{
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
		fnc = ndb.prePersistFastNode
	}

	iter.UnsavedFastIterator = newUnsavedFastIterator(start, end, ascending, ndb, fnc.additions, fnc.removals)

	return iter
}
