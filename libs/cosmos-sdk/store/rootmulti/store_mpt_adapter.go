package rootmulti

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

const (
	AccStore    = "acc"
	EvmStore    = "evm"
	MptStore    = "mpt"       // new store for acc module, will use mpt instead of iavl as store engine
)

func evmAccStoreFilter(sName string, ver int64, forceFilter ...bool) bool {
	if (sName == AccStore || sName == EvmStore) && tmtypes.HigherThanMars(ver) {
		if len(forceFilter) > 0 && forceFilter[0] {
			return true
		}

		// if mpt.TrieDirtyDisabled == true, means is a full node, should still use acc and evm store to query history state, keep them!
		// else, no longer need them any more, filter them !!!
		return !mpt.TrieDirtyDisabled
	}
	return false
}

func newMptStoreFilter(sName string, ver int64) bool {
	if (sName == MptStore) && !tmtypes.HigherThanMars(ver) {
		return true
	}
	return false
}

func (rs *Store) commitInfoFilter(infos map[string]storeInfo, ver int64, storeName string) {
	evmConfigInfo := infos[storeName]
	if evmConfigInfo.Core.CommitID.Version == 0 {
		evmConfigInfo.Core.CommitID.Version = ver
		infos[storeName] = evmConfigInfo

		for key, param := range rs.storesParams {
			if key.Name() == storeName {
				param.initialVersion = uint64(ver)
				rs.storesParams[key] = param
			}
		}
	}
}
