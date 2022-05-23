package cachekv

import (
	"bytes"
	"io"
	"sort"
	"sync"

	"github.com/okex/exchain/libs/iavl"

	"github.com/tendermint/go-amino"

	"github.com/okex/exchain/libs/cosmos-sdk/internal/conv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	kv "github.com/okex/exchain/libs/cosmos-sdk/types/kv"
	dbm "github.com/okex/exchain/libs/tm-db"
)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type cValue struct {
	value   []byte
	deleted bool
}

type PreChangesHandler func(keys []string, setOrDel []byte)

// Store wraps an in-memory cache around an underlying types.KVStore.
type Store struct {
	mtx           sync.Mutex
	dirtyCache    map[string]*cValue
	readListCache map[string][]byte
	deleted       map[string]struct{}
	unsortedCache map[string]struct{}
	sortedCache   *dbm.MemDB // always ascending sorted
	parent        types.KVStore

	preChangesHandler PreChangesHandler
}

var _ types.CacheKVStore = (*Store)(nil)

func NewStore(parent types.KVStore) *Store {
	return &Store{
		dirtyCache:    make(map[string]*cValue),
		readListCache: make(map[string][]byte),
		deleted:       make(map[string]struct{}),
		unsortedCache: make(map[string]struct{}),
		sortedCache:   dbm.NewMemDB(),
		parent:        parent,
	}
}

func NewStoreWithPreChangeHandler(parent types.KVStore, preChangesHandler PreChangesHandler) *Store {
	s := NewStore(parent)
	s.preChangesHandler = preChangesHandler
	return s
}

// Implements Store.
func (store *Store) GetStoreType() types.StoreType {
	return store.parent.GetStoreType()
}

// Implements types.KVStore.
func (store *Store) Get(key []byte) (value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	types.AssertValidKey(key)

	cacheValue, ok := store.dirtyCache[string(key)]
	if !ok {
		if c, ok := store.readListCache[string(key)]; ok {
			value = c
		} else {
			value = store.parent.Get(key)
			store.setCacheValue(key, value, false, false)
		}
	} else {
		value = cacheValue.value
	}

	return value
}

func (store *Store) IteratorCache(isdirty bool, cb func(key string, value []byte, isDirty bool, isDelete bool, sKey types.StoreKey) bool, sKey types.StoreKey) bool {
	if cb == nil {
		return true
	}
	store.mtx.Lock()
	defer store.mtx.Unlock()

	if isdirty {
		for key, v := range store.dirtyCache {
			if !cb(key, v.value, true, v.deleted, sKey) {
				return false
			}
		}
	} else {
		for key, v := range store.readListCache {
			if !cb(key, v, false, false, sKey) {
				return false
			}
		}
	}

	return true
}

// Implements types.KVStore.
func (store *Store) Set(key []byte, value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	types.AssertValidKey(key)
	types.AssertValidValue(value)

	store.setCacheValue(key, value, false, true)
}

// Implements types.KVStore.
func (store *Store) Has(key []byte) bool {
	value := store.Get(key)
	return value != nil
}

// Implements types.KVStore.
func (store *Store) Delete(key []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	types.AssertValidKey(key)

	store.setCacheValue(key, nil, true, true)
}

// Implements Cachetypes.KVStore.
func (store *Store) Write() {
	// if parent is cachekv.Store, we can write kv more efficiently
	if pStore, ok := store.parent.(*Store); ok {
		store.writeToCacheKv(pStore)
		return
	}

	store.mtx.Lock()
	defer store.mtx.Unlock()

	// We need a copy of all of the keys.
	// Not the best, but probably not a bottleneck depending.
	keys := make([]string, len(store.dirtyCache))
	index := 0
	for key, _ := range store.dirtyCache {
		keys[index] = key
		index++

	}

	sort.Strings(keys)

	store.preWrite(keys)

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for _, key := range keys {
		cacheValue := store.dirtyCache[key]
		switch {
		case cacheValue.deleted:
			store.parent.Delete([]byte(key))
		case cacheValue.value == nil:
			// Skip, it already doesn't exist in parent.
		default:
			store.parent.Set([]byte(key), cacheValue.value)
		}
	}

	// Clear the cache
	store.clearCache()
}

func (store *Store) preWrite(keys []string) {
	if store.preChangesHandler == nil || len(keys) < 4 {
		return
	}

	setOrDel := make([]byte, 0, len(keys))

	for _, key := range keys {
		cacheValue := store.dirtyCache[key]
		switch {
		case cacheValue.deleted:
			setOrDel = append(setOrDel, iavl.PreChangeOpDelete)
		case cacheValue.value == nil:
			// Skip, it already doesn't exist in parent.
			setOrDel = append(setOrDel, iavl.PreChangeNop)
		default:
			setOrDel = append(setOrDel, iavl.PreChangeOpSet)
		}
	}

	store.preChangesHandler(keys, setOrDel)
}

// writeToCacheKv will write cached kv to the parent Store, then clear the cache.
func (store *Store) writeToCacheKv(parent *Store) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for key, cacheValue := range store.dirtyCache {
		switch {
		case cacheValue.deleted:
			parent.Delete(amino.StrToBytes(key))
		case cacheValue.value == nil:
			// Skip, it already doesn't exist in parent.
		default:
			parent.Set(amino.StrToBytes(key), cacheValue.value)
		}
	}

	// Clear the cache
	store.clearCache()
}

func (store *Store) clearCache() {
	// https://github.com/golang/go/issues/20138
	for key := range store.dirtyCache {
		delete(store.dirtyCache, key)
	}

	for key := range store.readListCache {
		delete(store.readListCache, key)
	}
	
	for key := range store.deleted {
		delete(store.deleted, key)
	}
	
	for key := range store.unsortedCache {
		delete(store.unsortedCache, key)
	}
	store.sortedCache = dbm.NewMemDB()
}

//----------------------------------------
// To cache-wrap this Store further.

// Implements CacheWrapper.
func (store *Store) CacheWrap() types.CacheWrap {
	return NewStore(store)
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (store *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return NewStore(tracekv.NewStore(store, w, tc))
}

//----------------------------------------
// Iteration

// Implements types.KVStore.
func (store *Store) Iterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, true)
}

// Implements types.KVStore.
func (store *Store) ReverseIterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, false)
}

func (store *Store) iterator(start, end []byte, ascending bool) types.Iterator {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	var parent, cache types.Iterator

	if ascending {
		parent = store.parent.Iterator(start, end)
	} else {
		parent = store.parent.ReverseIterator(start, end)
	}

	store.dirtyItems(start, end)
	cache = newMemIterator(start, end, store.sortedCache, store.deleted, ascending)

	return newCacheMergeIterator(parent, cache, ascending)
}

func findStartIndex(strL []string, startQ string) int {
	// Modified binary search to find the very first element in >=startQ.
	if len(strL) == 0 {
		return -1
	}

	var left, right, mid int
	right = len(strL) - 1
	for left <= right {
		mid = (left + right) >> 1
		midStr := strL[mid]
		if midStr == startQ {
			// Handle condition where there might be multiple values equal to startQ.
			// We are looking for the very first value < midStL, that i+1 will be the first
			// element >= midStr.
			for i := mid - 1; i >= 0; i-- {
				if strL[i] != midStr {
					return i + 1
				}
			}
			return 0
		}
		if midStr < startQ {
			left = mid + 1
		} else { // midStrL > startQ
			right = mid - 1
		}
	}
	if left >= 0 && left < len(strL) && strL[left] >= startQ {
		return left
	}
	return -1
}

func findEndIndex(strL []string, endQ string) int {
	if len(strL) == 0 {
		return -1
	}

	// Modified binary search to find the very first element <endQ.
	var left, right, mid int
	right = len(strL) - 1
	for left <= right {
		mid = (left + right) >> 1
		midStr := strL[mid]
		if midStr == endQ {
			// Handle condition where there might be multiple values equal to startQ.
			// We are looking for the very first value < midStL, that i+1 will be the first
			// element >= midStr.
			for i := mid - 1; i >= 0; i-- {
				if strL[i] < midStr {
					return i + 1
				}
			}
			return 0
		}
		if midStr < endQ {
			left = mid + 1
		} else { // midStrL > startQ
			right = mid - 1
		}
	}

	// Binary search failed, now let's find a value less than endQ.
	for i := right; i >= 0; i-- {
		if strL[i] < endQ {
			return i
		}
	}

	return -1
}

type sortState int

const (
	stateUnsorted sortState = iota
	stateAlreadySorted
)

// Constructs a slice of dirtyCache items, to use w/ memIterator.
func (store *Store) dirtyItems(start, end []byte) {
	startStr, endStr := conv.UnsafeBytesToStr(start), conv.UnsafeBytesToStr(end)
	if startStr > endStr {
		// Nothing to do here.
		return
	}

	n := len(store.unsortedCache)
	unsorted := make([]*kv.Pair, 0)
	// If the unsortedCache is too big, its costs too much to determine
	// whats in the subset we are concerned about.
	// If you are interleaving iterator calls with writes, this can easily become an
	// O(N^2) overhead.
	// Even without that, too many range checks eventually becomes more expensive
	// than just not having the cache.
	if n < 1024 {
		for key := range store.unsortedCache {
			if dbm.IsKeyInDomain(conv.UnsafeStrToBytes(key), start, end) {
				cacheValue := store.dirtyCache[key]
				unsorted = append(unsorted, &kv.Pair{Key: []byte(key), Value: cacheValue.value})
			}
		}
		store.clearUnsortedCacheSubset(unsorted, stateUnsorted)
		return
	}

	// Otherwise it is large so perform a modified binary search to find
	// the target ranges for the keys that we should be looking for.
	strL := make([]string, 0, n)
	for key := range store.unsortedCache {
		strL = append(strL, key)
	}
	sort.Strings(strL)

	// Now find the values within the domain
	//  [start, end)
	startIndex := findStartIndex(strL, startStr)
	endIndex := findEndIndex(strL, endStr)

	if endIndex < 0 {
		endIndex = len(strL) - 1
	}
	if startIndex < 0 {
		startIndex = 0
	}

	kvL := make([]*kv.Pair, 0)
	for i := startIndex; i <= endIndex; i++ {
		key := strL[i]
		cacheValue := store.dirtyCache[key]
		kvL = append(kvL, &kv.Pair{Key: []byte(key), Value: cacheValue.value})
	}

	// kvL was already sorted so pass it in as is.
	store.clearUnsortedCacheSubset(kvL, stateAlreadySorted)
}

func (store *Store) clearUnsortedCacheSubset(unsorted []*kv.Pair, sortState sortState) {
	n := len(store.unsortedCache)
	if len(unsorted) == n { // This pattern allows the Go compiler to emit the map clearing idiom for the entire map.
		for key := range store.unsortedCache {
			delete(store.unsortedCache, key)
		}
	} else { // Otherwise, normally delete the unsorted keys from the map.
		for _, kv := range unsorted {
			delete(store.unsortedCache, conv.UnsafeBytesToStr(kv.Key))
		}
	}

	if sortState == stateUnsorted {
		sort.Slice(unsorted, func(i, j int) bool {
			return bytes.Compare(unsorted[i].Key, unsorted[j].Key) < 0
		})
	}

	for _, item := range unsorted {
		if item.Value == nil {
			// deleted element, tracked by store.deleted
			// setting arbitrary value
			store.sortedCache.Set(item.Key, []byte{})
			continue
		}
		err := store.sortedCache.Set(item.Key, item.Value)
		if err != nil {
			panic(err)
		}
	}
}

//----------------------------------------
// etc

// Only entrypoint to mutate store.cache.
func (store *Store) setCacheValue(key, value []byte, deleted bool, dirty bool) {
	keyStr := conv.UnsafeBytesToStr(key)
	if !dirty {
		store.readListCache[keyStr] = value
		return
	}

	if deleted {
		store.deleted[keyStr] = struct{}{}
	} else {
		delete(store.deleted, keyStr)
	}

	store.dirtyCache[keyStr] = &cValue{
		value:   value,
		deleted: deleted,
	}
	if dirty {
		store.unsortedCache[keyStr] = struct{}{}
	}
}
