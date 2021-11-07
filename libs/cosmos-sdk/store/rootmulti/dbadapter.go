package rootmulti

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

var commithash = []byte("FAKE_HASH")

//----------------------------------------
// commitDBStoreWrapper should only be used for simulation/debugging,
// as it doesn't compute any commit hash, and it cannot load older state.

// Wrapper type for dbm.Db with implementation of KVStore
type commitDBStoreAdapter struct {
	dbadapter.Store
}

func (cdsa commitDBStoreAdapter) Commit() types.CommitID {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}
}

func (cdsa commitDBStoreAdapter) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: -1,
		Hash:    commithash,
	}
}

func (cdsa commitDBStoreAdapter) SetPruning(_ types.PruningOptions) {}

func (cdsa commitDBStoreAdapter) GetDBWriteCount() int { return 0 }

func (cdsa commitDBStoreAdapter) GetDBReadCount() int { return 0 }
func (cdsa commitDBStoreAdapter) GetNodeReadCount() int { return 0 }

func (cdsa commitDBStoreAdapter) ResetCount() {}