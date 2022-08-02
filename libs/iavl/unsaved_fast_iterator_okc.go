package iavl

import dbm "github.com/okex/exchain/libs/tm-db"

type UnsavedFastIteratorWithCache struct {
	*UnsavedFastIterator
}

var _ dbm.Iterator = (*UnsavedFastIteratorWithCache)(nil)

func NewUnsavedFastIteratorWithCache(start, end []byte, ascending bool, ndb *nodeDB, unsavedFastNodeAdditions map[string]*FastNode, unsavedFastNodeRemovals map[string]interface{}) *UnsavedFastIteratorWithCache {
	iter := &UnsavedFastIteratorWithCache{
		UnsavedFastIterator: &UnsavedFastIterator{},
	}

	if ndb == nil || unsavedFastNodeAdditions == nil || unsavedFastNodeRemovals == nil {
		iter.UnsavedFastIterator = newUnsavedFastIterator(start, end, ascending, ndb, unsavedFastNodeAdditions, unsavedFastNodeRemovals)
		return iter
	}

	fnc := &fastNodeChanges{
		additions: unsavedFastNodeAdditions,
		removals:  unsavedFastNodeRemovals,
	}
	fnc = fnc.clone()
	fnc.merge(ndb.prePersistFastNode)

	if ndb.tpfv != nil {
		fnc = ndb.tpfv.expand(fnc)
	}

	iter.UnsavedFastIterator = newUnsavedFastIterator(start, end, ascending, ndb, fnc.additions, fnc.removals)

	return iter
}
