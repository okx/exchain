package types

import (
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"sync"
)

var _ types.AccCacheStore = (*CacheStore)(nil)

// Store wraps an in-memory cache around an underlying types.KVStore.
type CacheStore struct {
	mtx           sync.Mutex
	cache         map[string]*accValue
	parent        *AccRootKVStore
}

func NewCacheStore(parent *AccRootKVStore) *CacheStore {
	return &CacheStore{
		cache: make(map[string]*accValue),
		parent: parent,
	}
}

// Implements types.KVStore.
func (store *CacheStore) Get(addr types.AccAddress) (value interface{}) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	cacheValue, ok := store.cache[addr.String()]
	if ok {
		return cacheValue.value
	} else {
		val := store.parent.Get(addr)
		if val != nil {
			store.setCacheValue(addr.String(), val, false, false)
		}

		return val
	}
}

// Implements types.KVStore.
func (store *CacheStore) Set(addr types.AccAddress, value interface{}) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.setCacheValue(addr.String(), value.(exported.Account), false, true)
}

// Implements types.KVStore.
func (store *CacheStore) Has(addr types.AccAddress) bool {
	value := store.Get(addr)
	return value != nil
}

// Implements types.KVStore.
func (store *CacheStore) Delete(addr types.AccAddress) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.setCacheValue(addr.String(), nil, true, true)
}

// Implements Cachetypes.KVStore.
func (store *CacheStore) Write() {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	for key, dbValue := range store.cache {
		if !dbValue.dirty {
			continue
		}

		addr, _ := types.AccAddressFromBech32(key)
		if dbValue.deleted {
			store.parent.Delete(addr)
		} else {
			if dbValue.value == nil {
				continue
			}

			store.parent.Set(addr, dbValue.value)
		}
	}

	// Clear the cache
	store.cache = make(map[string]*accValue)
}

// Only entrypoint to mutate store.cache.
func (store *CacheStore) setCacheValue(key string, value exported.Account, deleted bool, dirty bool) {
	store.cache[key] = &accValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}

func (store *CacheStore) NewIterator(startKey []byte) *trie.Iterator {
	return store.parent.NewIterator(startKey)
}
