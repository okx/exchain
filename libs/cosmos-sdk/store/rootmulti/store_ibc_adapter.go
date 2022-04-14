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
	f := s.pruneHeightFilterPipeline(h)
	m := make(map[types.StoreKey]types.CommitKVStore)
	for k, v := range s.stores {
		if f(k.Name()) {
			continue
		}
		m[k] = v
	}
	return m
}

func (rs *Store) SetCommitHeightFilterPipeline(f storetypes.HeightFilterPipeline) {
	rs.commitHeightFilterPipeline = storetypes.LinkPipeline(f, rs.commitHeightFilterPipeline)
}
func (rs *Store) SetPruneHeightFilterPipeline(f storetypes.HeightFilterPipeline) {
	rs.pruneHeightFilterPipeline = storetypes.LinkPipeline(f, rs.pruneHeightFilterPipeline)
}
func (rs *Store) SetVersionFilterPipeline(f storetypes.VersionFilterPipeline) {
	rs.versionPipeline = storetypes.LinkPipeline2(f, rs.versionPipeline)
}
