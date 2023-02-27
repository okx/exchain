package mpt

import (
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"sync"
)

var (
	gMptSnapshot        *snapshot.Tree = nil
	initMptSnapshotOnce sync.Once
)

func InstanceOfMptSnapshot() ethstate.Database {
	initMptSnapshotOnce.Do(func() {
		//		var recovery bool
		//
		//		version := ms.CurrentVersion()
		//		if layer := rawdb.ReadSnapshotRecoveryNumber(snap.GetDiskDB()); layer != nil && *layer > uint64(version) {
		//			ms.logger.Error("Enabling snapshot recovery", "chainhead", version, "diskbase", *layer)
		//			recovery = true
		//		}
		//		var err error
		//		ms.snaps, err = snapshot.New(snap.GetDiskDB(), ms.db.TrieDB(), 256, ms.originalRoot, false, true, recovery)
		//		if err != nil {
		//			return fmt.Errorf("open snapshot error %v", err)
		//		}
		//		return nil
		//	}
	})

	return gMptDatabase
}
