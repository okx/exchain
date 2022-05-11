package iavl

import (
	"github.com/tendermint/go-amino"
)

func (ndb *nodeDB) enqueueOrphanTask(version int64, rootHash []byte, newOrphans []*Node) {
	ndb.addOrphanItem(version, rootHash)

	go func(ndb *nodeDB, version int64, newOrphans []*Node) {
		ndb.mtx.Lock()
		ndb.saveNewOrphans(version, newOrphans, false)
		ndb.oi.removeOldOrphans(version)
		ndb.mtx.Unlock()
		ndb.oi.enqueueResult(version)
		ndb.uncacheNodeRontine(newOrphans)
	}(ndb, version, newOrphans)
}

func (ndb *nodeDB) addOrphanItem(version int64, rootHash []byte) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	ndb.oi.addOrphanItem(version, rootHash)
}

func (ndb *nodeDB) saveNewOrphans(version int64, orphans []*Node, lock bool) {
	if orphans == nil {
		return
	}

	version--
	ndb.log(IavlDebug, "saving orphan node to OrphanCache", "size", len(orphans))
	ndb.state.increasOrphanCount(len(orphans))

	if lock {
		ndb.mtx.Lock()
		defer ndb.mtx.Unlock()
	}

	ndb.oi.feedOrphansMap(version, orphans)
	for _, node := range orphans {
		ndb.oi.feedOrphanNodeCache(node)
		delete(ndb.prePersistNodeCache, amino.BytesToStr(node.hash))
		node.leftNode = nil
		node.rightNode = nil
	}
}

func (ndb *nodeDB) sanityCheckHandleOrphansResult(version int64) {
	ndb.oi.wait4Result(version)
}

func (ndb *nodeDB) findRootHash(version int64) (res []byte, found bool) {
	ndb.mtx.RLock()
	defer ndb.mtx.RUnlock()
	return ndb.oi.findRootHash(version)
}

func (ndb *nodeDB) orphanTask(version int64, orphans []*Node, rootHash []byte, persist bool) {
	ndb.addOrphanItem(version, rootHash)

	go func(ndb *nodeDB, version int64, orphans []*Node, persist bool) {
		if persist {
			ndb.mtx.Lock()
			ndb.oi.removeOldOrphans(version)
			ndb.mtx.Unlock()
		} else {
			ndb.mtx.Lock()
			ndb.saveNewOrphans(version, orphans, false)
			ndb.oi.removeOldOrphans(version)
			ndb.mtx.Unlock()
		}

		ndb.oi.enqueueResult(version)
		ndb.uncacheNodeRontine(orphans)
	}(ndb, version, orphans, persist)
}
