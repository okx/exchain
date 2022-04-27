package iavl

import (
	cmap "github.com/orcaman/concurrent-map"
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
