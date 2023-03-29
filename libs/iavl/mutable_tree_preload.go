package iavl

import (
	"bytes"
	"fmt"
	"github.com/tendermint/go-amino"
	"runtime"
	"sync"
)

const (
	PreloadConcurrencyThreshold = 4

	PreChangeOpSet    byte = 1
	PreChangeOpDelete byte = 0
	PreChangeNop      byte = 0xFF
)

type preWriteJob struct {
	key      []byte
	setOrDel byte
}

func (tree *MutableTree) PreChanges(keys []string, setOrDel []byte) {
	if tree.root == nil {
		return
	}

	maxNums := runtime.NumCPU()
	keyCount := len(keys)
	if maxNums > keyCount {
		maxNums = keyCount
	}
	if maxNums < PreloadConcurrencyThreshold {
		return
	}

	tree.ndb.initPreWriteCache()

	txJobChan := make(chan preWriteJob, keyCount)
	var wg sync.WaitGroup
	wg.Add(keyCount)

	for index := 0; index < maxNums; index++ {
		go func(ch chan preWriteJob, wg *sync.WaitGroup) {
			for j := range ch {
				tree.preChangeWithOutCache(tree.root, j.key, j.setOrDel)
				wg.Done()
			}
		}(txJobChan, &wg)
	}

	for i, key := range keys {
		setOrDelFlag := setOrDel[i]
		if setOrDelFlag != PreChangeNop {
			txJobChan <- preWriteJob{amino.StrToBytes(key), setOrDel[i]}
		}
	}
	close(txJobChan)
	wg.Wait()

	tree.ndb.finishPreWriteCache()
}

func (tree *MutableTree) preChangeWithOutCache(node *Node, key []byte, setOrDel byte) (find bool) {
	if node.isLeaf() {
		if bytes.Equal(node.key, key) {
			return true
		}
		return
	} else {
		var isSet = setOrDel == PreChangeOpSet
		if bytes.Compare(key, node.key) < 0 {
			node.leftNode = tree.preGetLeftNode(node)
			if find = tree.preChangeWithOutCache(node.leftNode, key, setOrDel); (!find && isSet) || (find && !isSet) {
				tree.preGetRightNode(node)
			}
		} else {
			node.rightNode = tree.preGetRightNode(node)
			if find = tree.preChangeWithOutCache(node.rightNode, key, setOrDel); (!find && isSet) || (find && !isSet) {
				tree.preGetLeftNode(node)
			}
		}
		return
	}
}

func (tree *MutableTree) preGetNode(hash []byte) (n *Node) {
	var fromDisk bool
	n, fromDisk = tree.ImmutableTree.ndb.GetNodeWithoutUpdateCache(hash)
	if fromDisk {
		tree.ndb.cacheNodeToPreWriteCache(n)
	}
	return
}

func (tree *MutableTree) preGetLeftNode(node *Node) (n *Node) {
	if node.leftNode != nil {
		return node.leftNode
	}
	return tree.preGetNode(node.leftHash)
}

func (tree *MutableTree) preGetRightNode(node *Node) (n *Node) {
	if node.rightNode != nil {
		return node.rightNode
	}
	return tree.preGetNode(node.rightHash)
}

func (tree *MutableTree) makeOrphansSliceReady() []*Node {
	maxOrphansNum := int(tree.Height()) + 3
	if cap(tree.readableOrphansSlice) < maxOrphansNum {
		tree.readableOrphansSlice = make([]*Node, 0, maxOrphansNum)
	} else {
		tree.readableOrphansSlice = tree.readableOrphansSlice[:0]
	}
	return tree.readableOrphansSlice
}

func (tree *MutableTree) setWithOrphansSlice(key []byte, value []byte, orphans *[]*Node) (updated bool) {
	if value == nil {
		panic(fmt.Sprintf("Attempt to store nil value at key '%s'", key))
	}

	if tree.ImmutableTree.root == nil {
		tree.addUnsavedAddition(key, value, tree.version+1)
		tree.ImmutableTree.root = NewNode(key, value, tree.version+1)
		return updated
	}

	tree.ImmutableTree.root, updated = tree.recursiveSet(tree.ImmutableTree.root, key, value, orphans)
	return updated
}

func (tree *MutableTree) removeWithOrphansSlice(key []byte, orphaned *[]*Node) (value []byte, removed bool) {
	if tree.root == nil {
		return nil, false
	}
	newRootHash, newRoot, _, value := tree.recursiveRemove(tree.root, key, orphaned)
	if len(*orphaned) == 0 {
		return nil, false
	}
	tree.addUnsavedRemoval(key)

	if newRoot == nil && newRootHash != nil {
		tree.root = tree.ndb.GetNode(newRootHash)
	} else {
		tree.root = newRoot
	}
	return value, true
}
