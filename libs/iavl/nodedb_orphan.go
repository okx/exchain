package iavl


func (ndb *nodeDB) enqueueOrphanTask(version int64, rootHash []byte, newOrphans []*Node) {

	task := func() {
		ndb.mtx.Lock()
		defer ndb.mtx.Unlock()
		ndb.saveNewOrphans(version, newOrphans, false)
		ndb.oi.removeOldOrphans(version, rootHash)
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

