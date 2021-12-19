package types

import (
	"github.com/ethereum/go-ethereum/trie"
	"sync"
)

var _ AccStore = (*AccCacheCommitStore)(nil)

// Store wraps an in-memory cache around an underlying types.KVStore.
type AccCacheCommitStore struct {
	mtx    sync.Mutex
	cache  map[string]*accValue
	parent AccStore
}

func NewAccCacheCommitStore(parent AccStore) *AccCacheCommitStore {
	return &AccCacheCommitStore{
		cache: make(map[string]*accValue),
		parent: parent,
	}
}

func (store *AccCacheCommitStore) CreateCacheStore() AccStore {
	return NewAccCacheCommitStore(store)
}

// Implements types.KVStore.
func (store *AccCacheCommitStore) Get(key string) (value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	cacheValue, ok := store.cache[key]
	if ok {
		return cacheValue.value
	} else {
		val := store.parent.Get(key)
		if val != nil {
			store.setCacheValue(key, val, false, false)
		}

		return val
	}
}

// Implements types.KVStore.
func (store *AccCacheCommitStore) Set(key string, value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.setCacheValue(key, value, false, true)
}

// Implements types.KVStore.
func (store *AccCacheCommitStore) Has(key string) bool {
	value := store.Get(key)
	return value != nil
}

// Implements types.KVStore.
func (store *AccCacheCommitStore) Delete(key string) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.setCacheValue(key, nil, true, true)
}

// Implements Cachetypes.KVStore.
func (store *AccCacheCommitStore) Write() {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	for key, dbValue := range store.cache {
		if !dbValue.dirty {
			continue
		}

		if dbValue.deleted {
			store.parent.Delete(key)
		} else {
			if dbValue.value == nil {
				continue
			}

			store.parent.Set(key, dbValue.value)
		}
	}

	// Clear the cache
	store.cache = make(map[string]*accValue)
}

func (store *AccCacheCommitStore) Clean() {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	// Clear the cache
	store.cache = make(map[string]*accValue)
}

// Only entrypoint to mutate store.cache.
func (store *AccCacheCommitStore) setCacheValue(key string, value []byte, deleted bool, dirty bool) {
	store.cache[key] = &accValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}

func (store *AccCacheCommitStore) NewIterator(startKey []byte) *trie.Iterator {
	return store.parent.NewIterator(startKey)
}
