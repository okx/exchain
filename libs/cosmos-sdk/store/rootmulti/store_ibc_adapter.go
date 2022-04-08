package rootmulti

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

var (
	ibcModules = map[string]struct{}{
		"ibc":            {},
		"mem_capability": {},
		"capability":     {},
		"transfer":       {},
		"erc20":          {},
	}
)

func isNeedFilterIbcModules(h int64, m string) bool {
	if !tmtypes.HigherThanVenus1(h) {
		_, exist := ibcModules[m]
		return exist
	}
	return false
}

func queryIbcProof(res *abci.ResponseQuery, info *commitInfo, storeName string) {
	// Restore origin path and append proof op.
	res.Proof.Ops = append(res.Proof.Ops, info.ProofOp(storeName))
}

func (s *Store) getFilterStores(h int64) map[types.StoreKey]types.CommitKVStore {
	if tmtypes.HigherThanVenus1(h) {
		return s.stores
	}
	m := make(map[types.StoreKey]types.CommitKVStore)
	for k, v := range s.stores {
		if _, exist := ibcModules[k.Name()]; exist {
			continue
		}
		m[k] = v
	}
	return m
}
