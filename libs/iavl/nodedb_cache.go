package iavl

import (
	"container/list"

	cmap "github.com/orcaman/concurrent-map"

	"github.com/okex/exchain/libs/iavl/config"
	"github.com/tendermint/go-amino"
)

func (ndb *nodeDB) uncacheNode(hash []byte) {
	if v, ok := ndb.nodeCache.Get(amino.BytesToStr(hash)); ok {
		ndb.nodeCache.Remove(amino.BytesToStr(hash))
		elem := v.(*list.Element)
		ndb.nodeCacheQueue.Remove(elem)
	}
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (ndb *nodeDB) cacheNode(node *Node) {
	elem, count := ndb.nodeCacheQueue.PushBack(node)
	ndb.nodeCache.Set(string(node.hash), elem)

	if count > config.DynamicConfig.GetIavlCacheSize() {
		needRemove := count - config.DynamicConfig.GetIavlCacheSize()

		for i := 0; i < needRemove; i++ {
			oldest := ndb.nodeCacheQueue.Front()
			ndb.nodeCache.Remove(amino.BytesToStr(oldest.Value.(*Node).hash))
			_ = ndb.nodeCacheQueue.Remove(oldest)
		}
	}
}

func (ndb *nodeDB) cacheNodeByCheck(node *Node) {
	if _, ok := ndb.nodeCache.Get(amino.BytesToStr(node.hash)); !ok {
		ndb.cacheNode(node)
	}
}

func (ndb *nodeDB) getNodeFromCache(hash []byte, promoteRecentNode bool) (n *Node) {
	// Check the cache.
	if v, ok := ndb.nodeCache.Get(amino.BytesToStr(hash)); ok {
		elem := v.(*list.Element)
		if promoteRecentNode {
			// Already exists. Move to back of nodeCacheQueue.
			ndb.nodeCacheQueue.MoveToBack(elem)
		}
		n = elem.Value.(*Node)
	}
	return
}

func (ndb *nodeDB) uncacheNodeRontine(n []*Node) {
	for _, node := range n {
		ndb.uncacheNode(node.hash)
	}
}

func (ndb *nodeDB) InitPreWriteCache() {
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
