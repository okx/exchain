package watcher

import (
	"context"

	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	dbm "github.com/okex/exchain/libs/tm-db"
)

type readKVStore struct {
	sdk.KVStore
	prefix      []byte
	watchDBOnly bool
}

func WrapReadKVStore(ctx context.Context, store sdk.KVStore) sdk.KVStore {
	if !enableWatcher {
		return store
	}
	value := ctx.Value(QueryTypeKey)
	if v, ok := value.(string); ok {
		return &readKVStore{
			KVStore:     store,
			watchDBOnly: v == QueryWatchDBOnly,
		}
	}
	return store
}

func (r *readKVStore) Get(key []byte) []byte {
	value, _ := db.Get(key)
	if len(value) != 0 || r.watchDBOnly {
		return value
	}
	return r.KVStore.Get(key)
}

func (r *readKVStore) Has(key []byte) bool {
	has := dbadapter.Store{DB: db}.Has(key)
	if has || r.watchDBOnly {
		return has
	}
	return r.KVStore.Has(key)
}

func (r *readKVStore) Iterator(start, end []byte) dbm.Iterator {
	// TODO: make sure of consistence of watch db and blockchain
	return dbadapter.Store{DB: db}.Iterator(start, end)
}

func (r *readKVStore) ReverseIterator(start, end []byte) dbm.Iterator {
	return dbadapter.Store{DB: db}.ReverseIterator(start, end)
}
