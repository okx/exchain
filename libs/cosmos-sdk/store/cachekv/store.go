package cachekv

import (
	"bytes"
	"container/list"
	"io"
	"reflect"
	"sort"
	"sync"
	"unsafe"

	"github.com/tendermint/go-amino"

	tmkv "github.com/okex/exchain/libs/tendermint/libs/kv"
	dbm "github.com/okex/exchain/libs/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type cValue struct {
	value   []byte
	deleted bool
	dirty   bool
}

// Store wraps an in-memory cache around an underlying types.KVStore.
type Store struct {
	mtx           sync.Mutex
	dirty         map[string]cValue
	unsortedCache map[string]struct{}
	sortedCache   *list.List // always ascending sorted
	parent        types.KVStore
	ReadList      map[string][]byte
}

var _ types.CacheKVStore = (*Store)(nil)

func NewStore(parent types.KVStore) *Store {
	return &Store{
		dirty:         make(map[string]cValue),
		ReadList:      make(map[string][]byte),
		unsortedCache: make(map[string]struct{}),
		sortedCache:   list.New(),
		parent:        parent,
	}
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

	sKey := string(key)
	cacheValue, ok := store.dirty[sKey]
	if !ok {
		if c, ok := store.ReadList[sKey]; ok {
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

func (store *Store) IteratorCache(cb func(key, value []byte, isDirty bool, isDelete bool, sKey types.StoreKey) bool, sKey types.StoreKey) bool {
	if cb == nil {
		return true
	}
	store.mtx.Lock()
	defer store.mtx.Unlock()

	for key, v := range store.dirty {
		if !cb([]byte(key), v.value, v.dirty, v.deleted, sKey) {
			return false
		}
	}
	return true
}

func (store *Store) GetRWSet(rSet map[string][]byte, wSet map[string][]byte) {
	for k, v := range store.ReadList {
		rSet[k] = v
	}
	for k, v := range store.dirty {
		wSet[k] = v.value
	}
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
	keys := make([]string, len(store.dirty))
	index := 0
	for key, _ := range store.dirty {
		keys[index] = key
		index++

	}

	sort.Strings(keys)

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for _, key := range keys {
		cacheValue := store.dirty[key]
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

// writeToCacheKv will write cached kv to the parent Store, then clear the cache.
func (store *Store) writeToCacheKv(parent *Store) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for key, cacheValue := range store.dirty {
		if !cacheValue.dirty {
			continue
		}
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
	for key := range store.dirty {
		delete(store.dirty, key)
	}

	for Key := range store.ReadList {
		delete(store.ReadList, Key)
	}
	for key := range store.unsortedCache {
		delete(store.unsortedCache, key)
	}

	store.sortedCache.Init()
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
	cache = newMemIterator(start, end, store.sortedCache, ascending)

	return newCacheMergeIterator(parent, cache, ascending)
}

// strToByte is meant to make a zero allocation conversion
// from string -> []byte to speed up operations, it is not meant
// to be used generally, but for a specific pattern to check for available
// keys within a domain.
func strToByte(s string) []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = len(s)
	hdr.Len = len(s)
	hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	return b
}

// byteSliceToStr is meant to make a zero allocation conversion
// from []byte -> string to speed up operations, it is not meant
// to be used generally, but for a specific pattern to delete keys
// from a map.
func byteSliceToStr(b []byte) string {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(hdr))
}

// Constructs a slice of dirty items, to use w/ memIterator.
func (store *Store) dirtyItems(start, end []byte) {
	unsorted := make([]*tmkv.Pair, 0)

	n := len(store.unsortedCache)
	for key := range store.unsortedCache {
		if dbm.IsKeyInDomain(strToByte(key), start, end) {
			cacheValue := store.dirty[key]
			unsorted = append(unsorted, &tmkv.Pair{Key: []byte(key), Value: cacheValue.value})
		}
	}

	if len(unsorted) == n { // This pattern allows the Go compiler to emit the map clearing idiom for the entire map.
		for key := range store.unsortedCache {
			delete(store.unsortedCache, key)
		}
	} else { // Otherwise, normally delete the unsorted keys from the map.
		for _, kv := range unsorted {
			delete(store.unsortedCache, byteSliceToStr(kv.Key))
		}
	}

	sort.Slice(unsorted, func(i, j int) bool {
		return bytes.Compare(unsorted[i].Key, unsorted[j].Key) < 0
	})

	for e := store.sortedCache.Front(); e != nil && len(unsorted) != 0; {
		uitem := unsorted[0]
		sitem := e.Value.(*tmkv.Pair)
		comp := bytes.Compare(uitem.Key, sitem.Key)
		switch comp {
		case -1:
			unsorted = unsorted[1:]
			store.sortedCache.InsertBefore(uitem, e)
		case 1:
			e = e.Next()
		case 0:
			unsorted = unsorted[1:]
			e.Value = uitem
			e = e.Next()
		}
	}

	for _, kvp := range unsorted {
		store.sortedCache.PushBack(kvp)
	}

}

//----------------------------------------
// etc

// Only entrypoint to mutate store.cache.
func (store *Store) setCacheValue(key, value []byte, deleted bool, dirty bool) {
	keyStr := string(key)
	if !dirty {
		store.ReadList[keyStr] = value
		return
	}

	store.dirty[keyStr] = cValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
	if dirty {
		store.unsortedCache[keyStr] = struct{}{}
	}
}
