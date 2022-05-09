package rootmulti

import (
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

const (
	AccStore    = "acc"
	EvmStore    = "evm"
	MptStore    = "mpt"       // new store for acc module, will use mpt instead of iavl as store engine
	NewEvmStore = "evmlegacy" //new store for evm module, evm store will del after migration. the chainconfig, whiteList and block list info will store in evmlegacy store
)

func evmAccStoreFilter(sName string, ver int64) bool {
	if (sName == AccStore || sName == EvmStore) && tmtypes.HigherThanMars(ver) {
		return true
	}
	return false
}

func newEvmStoreFilter(sName string, ver int64) bool {
	if sName == NewEvmStore && !tmtypes.HigherThanMars(ver) {
		return true
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
