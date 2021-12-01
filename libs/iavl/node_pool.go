package iavl

import (
	"sync"
)

var nodePool = sync.Pool{}

func GetNodeFromPool() *Node {
	np := nodePool.Get()
	if np == nil {
		return nil
	}
	n, ok := np.(*Node)
	if !ok {
		return nil
	}
	return n
}

func SetNodeToPool(node *Node) {
	nodePool.Put(node)
}

func (node *Node) Reset(key, value, hash, leftHash, rightHash []byte, version, size int64, leftNode, rightNode *Node, height int8, persisted, prePersisted bool) {
	node.key = key
	node.value = value
	node.hash = hash
	node.leftHash = leftHash
	node.rightHash = rightHash

	node.version = version
	node.size = size
	node.leftNode = leftNode
	node.rightNode = rightNode
	node.height = height
	node.persisted = persisted
	node.prePersisted = prePersisted
}
