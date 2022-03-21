package rootmulti

import (
	types2 "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func queryIbcProof(res *abci.ResponseQuery, info *commitInfo, storeName string) {
	// Restore origin path and append proof op.
	res.Proof.Ops = append(res.Proof.Ops, info.ProofOp(storeName))
}

func (s *Store) getStores(h int64) map[types.StoreKey]types.CommitKVStore {
	return s.stores
}

func (s *Store) getFilterStores(h int64) map[types.StoreKey]types.CommitKVStore {
	f := s.pruneHeightFilterPipeline(h)
	// TODO FILTER:
	m := make(map[types.StoreKey]types.CommitKVStore)
	for k, v := range s.stores {
		if f(k.Name()) {
			continue
		}
		m[k] = v
	}
	return m
}

func (rs *Store) SetCommitHeightFilterPipeline(f types2.HeightFilterPipeline) {
	rs.commitHeightFilterPipeline = types2.LinkPipeline(f, rs.commitHeightFilterPipeline)
}
func (rs *Store) SetPruneHeightFilterPipeline(f types2.HeightFilterPipeline) {
	rs.pruneHeightFilterPipeline = types2.LinkPipeline(f, rs.pruneHeightFilterPipeline)
}
