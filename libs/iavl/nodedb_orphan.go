package iavl

import (
	"github.com/tendermint/go-amino"
	"sync/atomic"
)

func (ndb *nodeDB) enqueueOrphanTask(version int64, rootHash []byte, newOrphans []*Node) {

	ndb.addOrphanItem(version, rootHash)

	task := func() {
		ndb.mtx.Lock()
		defer ndb.mtx.Unlock()
		ndb.saveNewOrphans(version, newOrphans, false)
		ndb.oi.removeOldOrphans(version)
		ndb.oi.enqueueResult(version)
	}

	ndb.oi.enqueueTask(task)
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
	atomic.AddInt64(&ndb.totalOrphanCount, int64(len(orphans)))

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
	ndb.uncacheNodeRontine(orphans)
}

func (ndb *nodeDB) sanityCheckHandleOrphansResult(version int64) {
	ndb.oi.wait4Result(version)
}

func (ndb *nodeDB) findRootHash(version int64) (res []byte, found bool) {
	ndb.mtx.RLock()
	defer ndb.mtx.RUnlock()
	return ndb.oi.findRootHash(version)
}

