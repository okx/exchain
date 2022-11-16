package iavl

import (
	"container/list"
	"sync"

	"github.com/tendermint/go-amino"
)

type FastNodeCache struct {
	items      map[string]*list.Element // items.
	cacheSize  int                      // cache size limit in elements.
	cacheQueue *syncList                // LRU queue of cache elements. Used for deletion.
	cacheMutex sync.RWMutex             // Mutex for node cache.
}

func newFastNodeCache(dbName string, cacheSize int) *FastNodeCache {
	if dbName == "evm" {
		return &FastNodeCache{
			items:      makeFastNodeCacheMap(cacheSize, 1),
			cacheSize:  cacheSize,
			cacheQueue: newSyncList(),
		}
	} else {
		return &FastNodeCache{
			items:      make(map[string]*list.Element),
			cacheSize:  cacheSize,
			cacheQueue: newSyncList(),
		}
	}
}

func makeFastNodeCacheMap(cacheSize int, initRatio float64) map[string]*list.Element {
	if initRatio <= 0 {
		return make(map[string]*list.Element)
	}
	if initRatio >= 1 {
		return make(map[string]*list.Element, cacheSize)
	}
	cacheSize = int(float64(cacheSize) * initRatio)
	return make(map[string]*list.Element, cacheSize)
}

// ===================================================
// ======= map[string]*list.Element implementation
// ===================================================

func (fnc *FastNodeCache) uncache(key []byte) {
	fnc.cacheMutex.Lock()
	if elem, ok := fnc.items[string(key)]; ok {
		fnc.cacheQueue.Remove(elem)
		delete(fnc.items, string(key))
	}
	fnc.cacheMutex.Unlock()
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (fnc *FastNodeCache) cache(node *FastNode) {
	fnc.cacheMutex.Lock()

	if elem, ok := fnc.items[string(node.key)]; ok {
		fnc.cacheQueue.MoveToBack(elem)
		elem.Value = node
	} else {
		elem := fnc.cacheQueue.PushBack(node)
		fnc.items[string(node.key)] = elem

		for fnc.cacheQueue.Len() > GetFastNodeCacheSize() {
			oldest := fnc.cacheQueue.Front()
			key := fnc.cacheQueue.Remove(oldest).(*FastNode).key
			delete(fnc.items, amino.BytesToStr(key))
		}
	}

	fnc.cacheMutex.Unlock()
}

func (fnc *FastNodeCache) get(key []byte, promoteRecentNode bool) (n *FastNode) {
	// Check the cache.
	fnc.cacheMutex.RLock()
	elem, ok := fnc.items[string(key)]
	if ok {
		if promoteRecentNode {
			// Already exists. Move to back of cacheQueue.
			fnc.cacheQueue.MoveToBack(elem)
		}
		n = elem.Value.(*FastNode)
	}
	fnc.cacheMutex.RUnlock()
	return
}

func (fnc *FastNodeCache) cacheLen() int {
	return len(fnc.items)
}
