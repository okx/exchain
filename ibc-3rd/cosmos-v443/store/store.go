package store

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/store/cache"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/store/rootmulti"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewCommitKVStoreCacheManager() types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
}
