package store

import (
	dbm "github.com/okx/okbchain/libs/tm-db"

	"github.com/okx/okbchain/libs/cosmos-sdk/store/cache"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/rootmulti"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db)
}

func NewCommitKVStoreCacheManager() types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
}
