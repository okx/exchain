package iavl

import (
	"container/list"
	"github.com/okex/exchain/libs/iavl/config"
	"github.com/tendermint/go-amino"
	"sync"
)

type NodeCache struct {
	nodeCache      map[string]*list.Element // Node cache.
	nodeCacheSize  int                      // Node cache size limit in elements.
	nodeCacheQueue *syncList                // LRU queue of cache elements. Used for deletion.
	nodeCacheMutex sync.RWMutex             // Mutex for node cache.
}

func newNodeCache(cacheSize int) *NodeCache {
	return &NodeCache{
		nodeCache:      makeNodeCacheMap(cacheSize, IavlCacheInitRatio),
		nodeCacheSize:  cacheSize,
		nodeCacheQueue: newSyncList(),
	}
}

// ===================================================
// ======= map[string]*list.Element implementation
// ===================================================

func (ndb *NodeCache) uncacheNode(hash []byte) {
	ndb.nodeCacheMutex.Lock()
	if elem, ok := ndb.nodeCache[string(hash)]; ok {
		ndb.nodeCacheQueue.Remove(elem)
		delete(ndb.nodeCache, string(hash))
	}
	ndb.nodeCacheMutex.Unlock()
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (ndb *NodeCache) cacheNode(node *Node) {
	ndb.nodeCacheMutex.Lock()
	elem := ndb.nodeCacheQueue.PushBack(node)
	ndb.nodeCache[string(node.hash)] = elem

	for ndb.nodeCacheQueue.Len() > config.DynamicConfig.GetIavlCacheSize() {
		oldest := ndb.nodeCacheQueue.Front()
		hash := ndb.nodeCacheQueue.Remove(oldest).(*Node).hash
		delete(ndb.nodeCache, amino.BytesToStr(hash))
	}
	ndb.nodeCacheMutex.Unlock()
}

func (ndb *NodeCache) cacheNodeByCheck(node *Node) {
	ndb.nodeCacheMutex.RLock()
	_, ok := ndb.nodeCache[string(node.hash)]
	ndb.nodeCacheMutex.RUnlock()
	if !ok {
		ndb.cacheNode(node)
	}
}

func (ndb *NodeCache) getNodeFromCache(hash []byte, promoteRecentNode bool) (n *Node) {
	// Check the cache.
	ndb.nodeCacheMutex.RLock()
	elem, ok := ndb.nodeCache[string(hash)]
	if ok {
		if promoteRecentNode {
			// Already exists. Move to back of nodeCacheQueue.
			ndb.nodeCacheQueue.MoveToBack(elem)
		}
		n = elem.Value.(*Node)
	}
	ndb.nodeCacheMutex.RUnlock()
	return
}

func (ndb *NodeCache) nodeCacheLen() int {
	return len(ndb.nodeCache)
}

// =========================================================
// ======= github.com/hashicorp/golang-lru implementation
// =========================================================


//func (ndb *nodeDB) cacheNode(node *Node) {
//	ndb.lruNodeCache.Add(string(node.hash), node)
//}
//
//func (ndb *nodeDB) cacheNodeByCheck(node *Node) {
//	if ndb.lruNodeCache.Contains(string(node.hash)) {
//		return
//	}
//	ndb.cacheNode(node)
//}
//
//
//func (ndb *nodeDB) getNodeFromCache(hash []byte, promoteRecentNode bool) (n *Node) {
//
//	var ok bool
//	var res interface{}
//	if promoteRecentNode {
//		res, ok = ndb.lruNodeCache.Get(string(hash))
//	} else {
//		res, ok = ndb.lruNodeCache.Peek(string(hash))
//	}
//
//	if ok {
//		n = res.(*Node)
//	}
//	return
//}
//
//
//func (ndb *nodeDB) uncacheNode(hash []byte) {
//	ndb.lruNodeCache.Remove(string(hash))
//}
//
//
//func (ndb *nodeDB) nodeCacheLen() int {
//	return ndb.lruNodeCache.Len()
//}

