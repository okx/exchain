package iavl

func (ndb *nodeDB) addOrphanItem(version int64, rootHash []byte) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	ndb.oi.addOrphanItem(version, rootHash)
}

func (ndb *nodeDB) enqueueOrphanTask(version int64, rootHash []byte, newOrphans []*Node) {

	ndb.addOrphanItem(version, rootHash)

	task := func() {
		ndb.mtx.Lock()
		defer ndb.mtx.Unlock()
		ndb.saveNewOrphans(version, newOrphans, false)
		ndb.oi.removeOldOrphans()
		ndb.oi.resultChan <- version
	}

	ndb.oi.orphanTaskChan <- task
}

func (ndb *nodeDB) sanityCheckHandleOrphansResult(version int64) {
	ndb.oi.wait4Result(version)
}

func (ndb *nodeDB) findRootHash(version int64) (res []byte, found bool) {
	ndb.mtx.RLock()
	defer ndb.mtx.RUnlock()
	return ndb.oi.findRootHash(version)
}

