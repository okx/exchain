package iavl

import (
	"sync"
)

var onceTreeMap sync.Once
var treeMap *TreeMap

type TreeMap struct {
	mtx sync.RWMutex
	// used for checking whether a tree is saved or not
	mutableTreeSavedMap     map[string]*MutableTree
	totalPreCommitCacheSize int64
	lastUpdatedVersion      int64
	evm *MutableTree
}

func init() {
	onceTreeMap.Do(func() {
		treeMap = &TreeMap{
			mutableTreeSavedMap: make(map[string]*MutableTree),
		}
	})
}

func (tm *TreeMap) addNewTree(tree *MutableTree) {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()
	if _, ok := tm.mutableTreeSavedMap[tree.GetModuleName()]; !ok {
		tm.mutableTreeSavedMap[tree.GetModuleName()] = tree
		if tree.GetModuleName() == "evm" {
			tm.evm = tree
		}
		go tree.commitSchedule()
	}
}

func (tm *TreeMap) getTree(moduleName string) (tree *MutableTree, ok bool) {
	tm.mtx.RLock()
	defer tm.mtx.RUnlock()
	tree, ok = tm.mutableTreeSavedMap[moduleName]
	return
}

// updateMutableTreeMap marks into true when operation of save-version is done
func (tm *TreeMap) updateTotalPpnc(module string, version int64) {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()

	if version == tm.lastUpdatedVersion {
		return
	}
	var size int64 = 0
	for _, tree := range tm.mutableTreeSavedMap {
		size += int64(len(tree.ndb.prePersistNodeCache))
	}
	tm.totalPreCommitCacheSize = size
	tm.lastUpdatedVersion = version

	tm.evm.ndb.log(IavlInfo,"updateMutableTreeMap",
		"version", tm.lastUpdatedVersion-1,
		"module-by", module,
		"totalPreCommitCacheSize", tm.totalPreCommitCacheSize,
		)
}

// resetMap clear the TreeMap, only for test.
func (tm *TreeMap) resetMap() {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()
	tm.mutableTreeSavedMap = make(map[string]*MutableTree)
}
