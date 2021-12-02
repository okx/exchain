package iavl

import (
	"sync"
)

var EnableNodePool = true
var nodePool = sync.Pool{}
var GetNodeFromPoolCounter = 0
var SetNodeFromPoolCounter = 0

func GetNodeFromPool() *Node {
	if EnableNodePool {
		np := nodePool.Get()
		if np == nil {
			return nil
		}
		n, ok := np.(*Node)
		if !ok {
			return nil
		}
		GetNodeFromPoolCounter++
		return n
	}
	return nil
}

func SetNodeToPool(node *Node) {
	if EnableNodePool {
		SetNodeFromPoolCounter++
		nodePool.Put(node)
	}
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
