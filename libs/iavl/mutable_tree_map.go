package iavl

import "sync"

var onceTreeMap sync.Once
var treeMap *TreeMap
type TreeMap struct {
	mtx sync.Mutex
	// used for checking whether a tree is saved or not
	mutableTreeList         []*MutableTree
	totalPreCommitCacheSize int64
	mutableTreeSavedMap     map[string]bool
}

func init() {
	onceTreeMap.Do(func() {
		treeMap = &TreeMap{
			mutableTreeSavedMap: make(map[string]bool),
		}
	})
}

func (tm *TreeMap) addNewTree(tree *MutableTree) {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()
	if _, ok := tm.mutableTreeSavedMap[tree.GetModuleName()]; !ok {
		tm.mutableTreeList = append(tm.mutableTreeList, tree)
		tm.mutableTreeSavedMap[tree.GetModuleName()] = false
		go tree.commitSchedule()
	}
}

// updateMutableTreeMap marks into true when operation of save-version is done
func (tm *TreeMap) updateMutableTreeMap(module string) {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()
	if _, ok := tm.mutableTreeSavedMap[module]; !ok {
		return
	}
	tm.mutableTreeSavedMap[module] = true
	if tm.isMutableTreeSavedMapAllReady() {
		tm.updateTotalPreCommitCacheSize()
	}
}

// isMutableTreeSavedMapAllReady check if all trees are saved or not
func (tm *TreeMap) isMutableTreeSavedMapAllReady() bool {
	for _, isReady := range tm.mutableTreeSavedMap {
		if !isReady {
			return false
		}
	}
	for key := range tm.mutableTreeSavedMap {
		tm.mutableTreeSavedMap[key] = false
	}
	return true
}

// updateTotalPreCommitCacheSize counts the number of prePersis node
func (tm *TreeMap) updateTotalPreCommitCacheSize() {
	var size int64 = 0
	for _, tree := range tm.mutableTreeList {
		size += int64(len(tree.ndb.prePersistNodeCache))
	}
	tm.totalPreCommitCacheSize = size
}

// resetMap clear the TreeMap, only for test.
func (tm *TreeMap) resetMap() {
	tm.mtx.Lock()
	defer tm.mtx.Unlock()
	tm.mutableTreeSavedMap = make(map[string]bool)
	tm.mutableTreeList = []*MutableTree{}
}