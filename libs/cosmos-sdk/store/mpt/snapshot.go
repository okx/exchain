package mpt

func (ms *MptStore) openSnapshot() error {
	//	if ms == nil || ms.db == nil || ms.trie == nil {
	//		return fmt.Errorf("mpt store is nil or mpt trie is nil")
	//	}
	//	// If the chain was rewound past the snapshot persistent layer (causing
	//	// a recovery block number to be persisted to disk), check if we're still
	//	// in recovery mode and in that case, don't invalidate the snapshot on a
	//	// head mismatch.
	//	var recovery bool
	//
	//	version := ms.CurrentVersion()
	//	if layer := rawdb.ReadSnapshotRecoveryNumber(snap.GetDiskDB()); layer != nil && *layer > uint64(version) {
	//		ms.logger.Error("Enabling snapshot recovery", "chainhead", version, "diskbase", *layer)
	//		recovery = true
	//	}
	//	var err error
	//	ms.snaps, err = snapshot.New(snap.GetDiskDB(), ms.db.TrieDB(), 256, ms.originalRoot, false, true, recovery)
	//	if err != nil {
	//		return fmt.Errorf("open snapshot error %v", err)
	//	}
	return nil
}
