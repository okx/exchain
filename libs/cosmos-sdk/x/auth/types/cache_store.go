package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"sync"
)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type cValue struct {
	value   exported.Account
	deleted bool
	dirty   bool
}

// Store wraps an in-memory cache around an underlying types.KVStore.
type CacheStore struct {
	mtx           sync.Mutex
	cache         map[string]*cValue
}

func NewCacheStore() *CacheStore {
	return &CacheStore{
		cache: make(map[string]*cValue),
	}
}

// Implements types.KVStore.
func (store *CacheStore) Get(key string) (value exported.Account) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	cacheValue, ok := store.cache[key]
	if ok  && !cacheValue.deleted{
		value = cacheValue.value
	}

	return value
}


// Implements types.KVStore.
func (store *CacheStore) Set(key string, value exported.Account) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.setCacheValue(key, value, false, true)
}

// Implements types.KVStore.
func (store *CacheStore) Has(key string) bool {
	value := store.Get(key)
	return value != nil
}

// Implements types.KVStore.
func (store *CacheStore) Delete(key string) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.setCacheValue(key, nil, true, true)
}

func (store *CacheStore) Clean() {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	// Clear the cache
	store.cache = make(map[string]*cValue)
}

func (store *CacheStore) IteratorCache(cb func(key string, value exported.Account, isDirty bool, isDelete bool)) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	if cb == nil || len(store.cache) == 0 {
		return
	}

	for key, v := range store.cache {
		cb(key, v.value, v.dirty, v.deleted)
	}
}

// Only entrypoint to mutate store.cache.
func (store *CacheStore) setCacheValue(key string, value exported.Account, deleted bool, dirty bool) {
	store.cache[key] = &cValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}
