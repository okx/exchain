package rootmulti

import (
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func queryIbcProof(res *abci.ResponseQuery, info *commitInfo, storeName string) {
	// Restore origin path and append proof op.
	res.Proof.Ops = append(res.Proof.Ops, info.ProofOp(storeName))
}

func (s *Store) getFilterStores(h int64) map[types.StoreKey]types.CommitKVStore {
	m := make(map[types.StoreKey]types.CommitKVStore)
	for k, v := range s.stores {
		if filter(k.Name(), h, v, s.pruneFilters) {
			continue
		}
		m[k] = v
	}
	return m
}

func (rs *Store) AppendCommitFilters(filters []storetypes.StoreFilter) {
	rs.commitFilters = append(rs.commitFilters, filters...)
}

func (rs *Store) AppendPruneFilters(filters []storetypes.StoreFilter) {
	rs.pruneFilters = append(rs.pruneFilters, filters...)
}

func (rs *Store) AppendVersionFilters(filters []storetypes.VersionFilter) {
	rs.versionFilters = append(rs.versionFilters, filters...)
}
