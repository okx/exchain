package types

import (
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"
	"sync"
)

var _ AccStore = (*AccCommitStore)(nil)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type accValue struct {
	value   []byte
	deleted bool
	dirty   bool
}

// Store is a wrapper for a MemDB with Commiter implementation
type AccCommitStore struct {
	mtx    sync.Mutex
	trie   ethstate.Trie
	cache  map[string]*accValue
}

// Constructs new MemDB adapter
func NewAccCommitStore() *AccCommitStore {
	return &AccCommitStore{
		cache: make(map[string]*accValue),
	}
}

func (acs *AccCommitStore) CreateCacheStore() AccStore {
	return &AccCommitStore{
		cache: make(map[string]*accValue),
		trie:  acs.trie,
	}
}

func (acs *AccCommitStore) SetMptTrie(tr ethstate.Trie) {
	acs.mtx.Lock()
	defer acs.mtx.Unlock()

	acs.trie = tr
}

func (acs *AccCommitStore) Write() {
	acs.mtx.Lock()
	defer acs.mtx.Unlock()

	for key, dbValue := range acs.cache {
		if !dbValue.dirty {
			continue
		}

		if dbValue.deleted {
			// delete account
			if err := acs.trie.TryDelete([]byte(key)); err != nil {
				panic(err)
			}
		} else {
			if err := acs.trie.TryUpdate([]byte(key), dbValue.value); err != nil {
				panic(err)
			}
		}
	}

	// Clear the cache
	acs.cache = make(map[string]*accValue)
}

func (acs *AccCommitStore) IteratorCache(cb func(key string, value []byte, isDirty bool, isDelete bool)) {
	acs.mtx.Lock()
	defer acs.mtx.Unlock()

	if cb == nil || len(acs.cache) == 0 {
		return
	}

	for key, v := range acs.cache {
		cb(key, v.value, v.dirty, v.deleted)
	}
}

func (acs *AccCommitStore) Clean() {
	acs.mtx.Lock()
	defer acs.mtx.Unlock()

	// Clear the cache
	acs.cache = make(map[string]*accValue)
}

func (acs *AccCommitStore) Get(key string) []byte {
	acs.mtx.Lock()
	defer acs.mtx.Unlock()

	if cacheValue, ok := acs.cache[key]; ok {
		return cacheValue.value
	}

	enc, err := acs.trie.TryGet([]byte(key))
	if err != nil {
		return nil
	}
	if len(enc) == 0 {
		return nil
	}

	acs.setCacheValue(key, enc, false, false)
	return enc
}

func (acs *AccCommitStore) Has(key string) bool {
	return acs.Get(key) != nil
}

func (acs *AccCommitStore) Set(key string, value []byte) {
	acs.mtx.Lock()
	defer acs.mtx.Unlock()

	acs.setCacheValue(key, value, false, true)
}

func (acs *AccCommitStore) Delete(key string) {
	acs.mtx.Lock()
	defer acs.mtx.Unlock()

	acs.setCacheValue(key, nil, true, true)
}

func (acs *AccCommitStore) NewIterator(startKey []byte) *trie.Iterator {
	return trie.NewIterator(acs.trie.NodeIterator(startKey))
}

// Only entrypoint to mutate store.cache.
func (acs *AccCommitStore) setCacheValue(key string, value []byte, deleted bool, dirty bool) {
	acs.cache[key] = &accValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
}
