package iavl

import (
	"container/list"

	"github.com/okex/exchain/libs/iavl/config"
	"github.com/tendermint/go-amino"
)

func (ndb *nodeDB) uncacheNode(hash []byte) {
	ndb.nodeCache.RemoveCb(amino.BytesToStr(hash), func(key string, v interface{}, exists bool) bool {
		if !exists {
			return false
		}
		elem := v.(*list.Element)
		ndb.nodeCacheQueue.Remove(elem)
		return true
	})
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (ndb *nodeDB) cacheNode(node *Node) {
	_, count := ndb.nodeCacheQueue.PushBackCb(node, func(ele *list.Element) {
		ndb.nodeCache.Set(string(node.hash), ele)
	})

	iavlCacheMaxSize := config.DynamicConfig.GetIavlCacheSize()
	if count > iavlCacheMaxSize {
		ndb.nodeCacheQueue.RemoveFrontNCb(count-iavlCacheMaxSize, func(v interface{}) {
			ndb.nodeCache.Remove(amino.BytesToStr(v.(*Node).hash))
		})
	}
}

func (ndb *nodeDB) cacheNodeByCheck(node *Node) {
	if _, ok := ndb.nodeCache.Get(amino.BytesToStr(node.hash)); !ok {
		ndb.cacheNode(node)
	}
}

func (ndb *nodeDB) getNodeFromCache(hash []byte) (n *Node) {
	// Check the cache.
	if v, ok := ndb.nodeCache.Get(amino.BytesToStr(hash)); ok {
		elem := v.(*list.Element)
		// Already exists. Move to back of nodeCacheQueue.
		ndb.nodeCacheQueue.MoveToBack(elem)
		n = elem.Value.(*Node)
	}
	return
}

func (ndb *nodeDB) uncacheNodeRontine(n []*Node) {
	for _, node := range n {
		ndb.uncacheNode(node.hash)
	}
}
