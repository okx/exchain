//go:build go1.18

package cachekv

// https://github.com/golang/go/issues/53157

//go:noinline
func (store *Store) clearCache() {
	// https://github.com/golang/go/issues/20138
	for key := range store.dirty {
		delete(store.dirty, key)
	}

	for Key := range store.readList {
		delete(store.readList, Key)
	}
	for key := range store.unsortedCache {
		delete(store.unsortedCache, key)
	}
	store.disableCacheReadList = false
	store.sortedCache.Init()
}
