package iavl

import (
	"sync"
	"testing"
)

type treeWrapper struct {
	tree *MutableTree
	mtx  sync.RWMutex
}

func newTreeWrapper(t *testing.T) *treeWrapper {
	ret := &treeWrapper{
		tree: nil,
	}
	ret.tree = newTestTree(t, false, 10000, "test")
	return ret
}

func TestPruneWithReadWrite(t *testing.T) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		treeMap.resetMap()
	}()
	treeWp := newTreeWrapper(t)
	// 读线程
}
