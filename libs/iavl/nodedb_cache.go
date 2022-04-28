package iavl

import (
	cmap "github.com/orcaman/concurrent-map"
	"github.com/tendermint/go-amino"
	"github.com/okex/exchain/libs/iavl/config"

)

func (ndb *nodeDB) uncacheNodeRontine(n []*Node) {
	for _, node := range n {
		ndb.uncacheNode(node.hash)
	}
}

func (ndb *nodeDB) initPreWriteCache() {
	if ndb.preWriteNodeCache == nil {
		ndb.preWriteNodeCache = cmap.New()
	}
}

func (ndb *nodeDB) cacheNodeToPreWriteCache(n *Node) {
	ndb.preWriteNodeCache.Set(string(n.hash), n)
}

func (ndb *nodeDB) finishPreWriteCache() {
	ndb.preWriteNodeCache.IterCb(func(key string, v interface{}) {
		ndb.cacheNode(v.(*Node))
	})
	ndb.preWriteNodeCache = nil
}


// ===================================================
// ======= map[string]*list.Element implementation
// ===================================================

func (ndb *nodeDB) uncacheNode(hash []byte) {
	ndb.nodeCacheMutex.Lock()
	if elem, ok := ndb.nodeCache[string(hash)]; ok {
		ndb.nodeCacheQueue.Remove(elem)
		delete(ndb.nodeCache, string(hash))
	}
	ndb.nodeCacheMutex.Unlock()
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (ndb *nodeDB) cacheNode(node *Node) {
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

func (ndb *nodeDB) cacheNodeByCheck(node *Node) {
	ndb.nodeCacheMutex.RLock()
	_, ok := ndb.nodeCache[string(node.hash)]
	ndb.nodeCacheMutex.RUnlock()
	if !ok {
		ndb.cacheNode(node)
	}
}


func (ndb *nodeDB) getNodeFromCache(hash []byte, promoteRecentNode bool) (n *Node) {
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

func (ndb *nodeDB) nodeCacheLen() int {
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

