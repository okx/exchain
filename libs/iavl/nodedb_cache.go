package iavl

import (
	lru "github.com/hashicorp/golang-lru"
)

type NodeCache struct {
	nodeCache *lru.Cache
	//nodeCache      map[string]*list.Element // Node cache.
	//nodeCacheSize  int                      // Node cache size limit in elements.
	//nodeCacheQueue *syncList                // LRU queue of cache elements. Used for deletion.
	//nodeCacheMutex sync.RWMutex             // Mutex for node cache.

	// for test
	count int
}

func newNodeCache(dbName string, cacheSize int) *NodeCache {
	//if dbName == "evm" {
	//	return &NodeCache{
	//		nodeCache:      makeNodeCacheMap(cacheSize, IavlCacheInitRatio),
	//		nodeCacheSize:  cacheSize,
	//		nodeCacheQueue: newSyncList(),
	//	}
	//} else {
	//	return &NodeCache{
	//		nodeCache:      make(map[string]*list.Element),
	//		nodeCacheSize:  cacheSize,
	//		nodeCacheQueue: newSyncList(),
	//	}
	//}
	cache, err := lru.New(cacheSize)
	if err != nil {
		panic(err)
	}
	return &NodeCache{
		nodeCache: cache,
	}
}

//func makeNodeCacheMap(cacheSize int, initRatio float64) map[string]*list.Element {
//	if initRatio <= 0 {
//		return make(map[string]*list.Element)
//	}
//	if initRatio >= 1 {
//		return make(map[string]*list.Element, cacheSize)
//	}
//	cacheSize = int(float64(cacheSize) * initRatio)
//	return make(map[string]*list.Element, cacheSize)
//}
//
// ===================================================
// ======= map[string]*list.Element implementation
// ===================================================

func (ndb *NodeCache) uncache(hash []byte) {
	//ndb.nodeCacheMutex.Lock()
	//if elem, ok := ndb.nodeCache[string(hash)]; ok {
	//	ndb.count--
	//	ndb.nodeCacheQueue.Remove(elem)
	//	delete(ndb.nodeCache, string(hash))
	//}
	//ndb.nodeCacheMutex.Unlock()
	ndb.nodeCache.Remove(string(hash))
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (ndb *NodeCache) cache(node *Node) {
	ndb.nodeCache.Add(string(node.hash), node)
	//ndb.nodeCacheMutex.Lock()
	//ndb.count++
	//if ele, ok := ndb.nodeCache[string(node.hash)]; ok {
	//	ndb.nodeCacheQueue.MoveToBack(ele)
	//} else {
	//	elem := ndb.nodeCacheQueue.PushBack(node)
	//	ndb.nodeCache[string(node.hash)] = elem
	//
	//	for ndb.nodeCacheQueue.Len() > config.DynamicConfig.GetIavlCacheSize() {
	//		oldest := ndb.nodeCacheQueue.Front()
	//		hash := ndb.nodeCacheQueue.Remove(oldest).(*Node).hash
	//		delete(ndb.nodeCache, amino.BytesToStr(hash))
	//	}
	//}
	//ndb.nodeCacheMutex.Unlock()
}

func (ndb *NodeCache) cacheWithKey(key string, node *Node) {
	ndb.nodeCache.Add(key, node)
	//ndb.nodeCacheMutex.Lock()
	//ndb.count++
	//elem := ndb.nodeCacheQueue.PushBack(node)
	//ndb.nodeCache[key] = elem
	//
	//for ndb.nodeCacheQueue.Len() > config.DynamicConfig.GetIavlCacheSize() {
	//	fmt.Println("")
	//	oldest := ndb.nodeCacheQueue.Front()
	//	hash := ndb.nodeCacheQueue.Remove(oldest).(*Node).hash
	//	delete(ndb.nodeCache, amino.BytesToStr(hash))
	//}
	//ndb.nodeCacheMutex.Unlock()
}

func (ndb *NodeCache) get(hash []byte, promoteRecentNode bool) (n *Node) {
	//// Check the cache.
	//ndb.nodeCacheMutex.RLock()
	//elem, ok := ndb.nodeCache[string(hash)]
	//if ok {
	//	if promoteRecentNode {
	//		// Already exists. Move to back of nodeCacheQueue.
	//		ndb.nodeCacheQueue.MoveToBack(elem)
	//	}
	//	n = elem.Value.(*Node)
	//}
	//ndb.nodeCacheMutex.RUnlock()

	var ok bool
	var res interface{}
	if promoteRecentNode {
		res, ok = ndb.nodeCache.Get(string(hash))
	} else {
		res, ok = ndb.nodeCache.Peek(string(hash))
	}

	if ok {
		n = res.(*Node)
	}
	return
	return
}

func (ndb *NodeCache) nodeCacheLen() int {
	return ndb.nodeCache.Len()
}

// =========================================================
// ======= github.com/hashicorp/golang-lru implementation
// =========================================================

//func (ndb *nodeDB) cache(node *Node) {
//	ndb.nc.nodeCache.Add(string(node.hash), node)
//}
//
//func (ndb *nodeDB) cacheWithKey(key string, node *Node) {
//	ndb.nc.nodeCache.Add(key, node)
//}
//
//func (ndb *nodeDB) cacheNodeByCheck(node *Node) {
//	if ndb.nc.nodeCache.Contains(string(node.hash)) {
//		return
//	}
//	ndb.cacheNode(node)
//}
//
//func (ndb *nodeDB) getNodeFromCache(hash []byte, promoteRecentNode bool) (n *Node) {
//
//	var ok bool
//	var res interface{}
//	if promoteRecentNode {
//		res, ok = ndb.nc.nodeCache.Get(string(hash))
//	} else {
//		res, ok = ndb.nc.nodeCache.Peek(string(hash))
//	}
//
//	if ok {
//		n = res.(*Node)
//	}
//	return
//}
//
//func (ndb *nodeDB) uncacheNode(hash []byte) {
//	ndb.nc.nodeCache.Remove(string(hash))
//}
//
//func (ndb *nodeDB) nodeCacheLen() int {
//	return ndb.nc.nodeCache.Len()
//}
