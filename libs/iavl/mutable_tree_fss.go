package iavl

import (
	"fmt"
	dbm "github.com/okex/exchain/libs/tm-db"
	"runtime"
	"sort"
	"time"
)

type FastStorageSystem interface {
	FastGetFromChanges(key []byte) ([]byte, bool)
	AddUnsavedAddition(key []byte, node *FastNode)
	AddUnsavedRemoval(key []byte)
	GetUnsavedFastNodeAdditions() map[string]*FastNode
	GetUnsavedFastNodeRemovals() map[string]interface{}
	SetDelta(delta *TreeDelta, savedNodes map[string]*Node, orphans []*Node)
	Reset()

	FastGetFromNodeDB(fastCacheEnable bool, key []byte, version int64) (index int64, value []byte, succeed bool)
	SaveFastNodeVersion(batch dbm.Batch, changes *FastNodeChanges) error
	UpdateBranchFastNode()
	UpdateLatestVersion(version int64)
	Iterate(fastCacheEnable bool, fn func(key []byte, value []byte) bool) (bool, error)
	Iterator(fastCacheEnable bool, start, end []byte, ascending bool) dbm.Iterator
	DeleteFastNodesFromNodeDB(batch dbm.Batch, version int64) error
	AsyncPersistFastNodeStart(version int64) *FastNodeChanges
	AsyncPersistFastNodeFinished(event *commitEvent)

	EnableFastStorageAndCommitIfNotEnabled(tree *MutableTree) (bool, error)
	EnableFastStorageAndCommit(tree *MutableTree, batch dbm.Batch) error
}

type FssNone struct{}

func NewFssNone() *FssNone {
	return &FssNone{}
}
func (f *FssNone) FastGetFromChanges(key []byte) ([]byte, bool) { return nil, false }
func (f *FssNone) FastGetFromNodeDB(fastCacheEnable bool, key []byte, version int64) (index int64, value []byte, succeed bool) {
	return 0, nil, false
}
func (f *FssNone) AddUnsavedAddition(key []byte, node *FastNode)                           {}
func (f *FssNone) GetUnsavedFastNodeAdditions() map[string]*FastNode                       { return nil }
func (f *FssNone) AddUnsavedRemoval(key []byte)                                            {}
func (f *FssNone) GetUnsavedFastNodeRemovals() map[string]interface{}                      { return nil }
func (f *FssNone) SaveFastNodeVersion(batch dbm.Batch, changes *FastNodeChanges) error     { return nil }
func (f *FssNone) SetDelta(delta *TreeDelta, savedNodes map[string]*Node, orphans []*Node) {}
func (f *FssNone) UpdateBranchFastNode()                                                   {}
func (f *FssNone) Iterate(fastCacheEnable bool, fn func(key []byte, value []byte) bool) (bool, error) {
	return false, fmt.Errorf("fast cache is not enable")
}
func (f *FssNone) Iterator(fastCacheEnable bool, start, end []byte, ascending bool) dbm.Iterator {
	return nil
}
func (f *FssNone) Reset() {}
func (f *FssNone) EnableFastStorageAndCommitIfNotEnabled(tree *MutableTree) (bool, error) {
	return false, nil
}
func (f *FssNone) EnableFastStorageAndCommit(tree *MutableTree, batch dbm.Batch) error {
	return nil
}
func (f *FssNone) DeleteFastNodesFromNodeDB(batch dbm.Batch, version int64) error {
	return nil
}
func (f *FssNone) AsyncPersistFastNodeStart(version int64) *FastNodeChanges {
	return nil
}
func (f *FssNone) AsyncPersistFastNodeFinished(event *commitEvent) {}
func (f *FssNone) UpdateLatestVersion(version int64)               {}

type Fss struct {
	ndb *nodeDB

	unsavedFastNodeAdditions map[string]*FastNode   // FastNodes that have not yet been saved to disk
	unsavedFastNodeRemovals  map[string]interface{} // FastNodes that have not yet been removed from disk
}

func NewFss(ndb *nodeDB) *Fss {
	return &Fss{
		ndb:                      ndb,
		unsavedFastNodeAdditions: make(map[string]*FastNode),
		unsavedFastNodeRemovals:  make(map[string]interface{}),
	}
}

func (f *Fss) FastGetFromChanges(key []byte) ([]byte, bool) {
	if fastNode, ok := f.unsavedFastNodeAdditions[string(key)]; ok {
		return fastNode.value, ok
	}

	if _, ok := f.unsavedFastNodeRemovals[string(key)]; ok {
		// is deleted
		return nil, ok
	}
	return nil, false
}

func (f *Fss) FastGetFromNodeDB(fastCacheEnable bool, key []byte, version int64) (index int64, value []byte, succeed bool) {
	if fastCacheEnable {
		fastNode, _ := f.ndb.GetFastNode(key)
		if fastNode == nil && version == f.ndb.latestVersion {
			return -1, nil, true
		}

		if fastNode != nil && fastNode.versionLastUpdatedAt <= version {
			return fastNode.versionLastUpdatedAt, fastNode.value, true
		}
	}
	return 0, nil, false
}

// nolint: unused
func (f *Fss) GetUnsavedFastNodeAdditions() map[string]*FastNode {
	return f.unsavedFastNodeAdditions
}

// getUnsavedFastNodeRemovals returns unsaved FastNodes to remove
// nolint: unused
func (f *Fss) GetUnsavedFastNodeRemovals() map[string]interface{} {
	return f.unsavedFastNodeRemovals
}

func (f *Fss) AddUnsavedAddition(key []byte, node *FastNode) {
	delete(f.unsavedFastNodeRemovals, string(key))
	f.unsavedFastNodeAdditions[string(key)] = node
}

func (f *Fss) AddUnsavedRemoval(key []byte) {
	delete(f.unsavedFastNodeAdditions, string(key))
	f.unsavedFastNodeRemovals[string(key)] = true
}

func (f *Fss) SaveFastNodeVersion(batch dbm.Batch, changes *FastNodeChanges) error {
	if changes == nil {
		return nil
	}

	additions := changes.additions
	if additions == nil {
		additions = f.unsavedFastNodeAdditions
	}
	if err := f.saveFastNodeAdditions(batch, additions); err != nil {
		return err
	}

	removals := changes.removals
	if removals == nil {
		removals = f.unsavedFastNodeRemovals
	}
	if err := f.saveFastNodeRemovals(batch, removals); err != nil {
		return err
	}

	return f.ndb.setFastStorageVersionToBatch(batch, changes.version)
}

func (f *Fss) saveFastNodeAdditions(batch dbm.Batch, additions map[string]*FastNode) error {
	keysToSort := make([]string, 0, len(additions))
	for key := range additions {
		keysToSort = append(keysToSort, key)
	}
	sort.Strings(keysToSort)

	for _, key := range keysToSort {
		if err := f.ndb.SaveFastNode(additions[key], batch); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fss) saveFastNodeRemovals(batch dbm.Batch, removals map[string]interface{}) error {
	keysToSort := make([]string, 0, len(removals))
	for key := range removals {
		keysToSort = append(keysToSort, key)
	}
	sort.Strings(keysToSort)

	for _, key := range keysToSort {
		if err := f.ndb.DeleteFastNode([]byte(key), batch); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fss) SetDelta(delta *TreeDelta, savedNodes map[string]*Node, orphans []*Node) {
	if delta != nil {
		// fast node related
		for _, v := range savedNodes {
			if v.isLeaf() {
				f.unsavedFastNodeAdditions[string(v.key)] = NewFastNode(v.key, v.value, v.version)
			}
		}

		for _, v := range orphans {
			_, ok := f.unsavedFastNodeAdditions[string(v.key)]
			if v.isLeaf() && !ok {
				f.unsavedFastNodeRemovals[string(v.key)] = NewFastNode(v.key, v.value, v.version)
			}
		}
	}
}

func (f *Fss) UpdateBranchFastNode() {
	f.ndb.updateBranchForFastNode(f.unsavedFastNodeAdditions, f.unsavedFastNodeRemovals)
	f.unsavedFastNodeAdditions = make(map[string]*FastNode)
	f.unsavedFastNodeRemovals = make(map[string]interface{})
}

func (f *Fss) Iterate(fastCacheEnable bool, fn func(key []byte, value []byte) bool) (bool, error) {
	if fastCacheEnable {

		itr := NewUnsavedFastIterator(nil, nil, true, f.ndb, f.unsavedFastNodeAdditions, f.unsavedFastNodeRemovals)

		defer itr.Close()
		for ; itr.Valid(); itr.Next() {
			if fn(itr.Key(), itr.Value()) {
				return true, nil
			}
		}
		return false, nil
	}
	return false, fmt.Errorf("fast cache is not enable")
}

func (f *Fss) Iterator(fastCacheEnable bool, start, end []byte, ascending bool) dbm.Iterator {
	if fastCacheEnable {
		return NewUnsavedFastIterator(start, end, ascending, f.ndb, f.unsavedFastNodeAdditions, f.unsavedFastNodeRemovals)
	}
	return nil
}

func (f *Fss) Reset() {
	f.unsavedFastNodeAdditions = map[string]*FastNode{}
	f.unsavedFastNodeRemovals = map[string]interface{}{}
}

// EnableFastStorageAndCommitIfNotEnabled if nodeDB doesn't mark fast storage as enabled, enable it, and commit the update.
// Checks whether the fast cache on disk matches latest live state. If not, deletes all existing fast nodes and repopulates them
// from latest tree.
// nolint: unparam
func (f *Fss) EnableFastStorageAndCommitIfNotEnabled(tree *MutableTree) (bool, error) {
	shouldForceUpdate := tree.ndb.shouldForceFastStorageUpgrade()
	isFastStorageEnabled := tree.ndb.hasUpgradedToFastStorage()

	if !tree.IsUpgradeable() {
		return false, nil
	}

	if isFastStorageEnabled && shouldForceUpdate {
		// If there is a mismatch between which fast nodes are on disk and the live state due to temporary
		// downgrade and subsequent re-upgrade, we cannot know for sure which fast nodes have been removed while downgraded,
		// Therefore, there might exist stale fast nodes on disk. As a result, to avoid persisting the stale state, it might
		// be worth to delete the fast nodes from disk.
		batch := tree.NewBatch()
		fastItr := NewFastIterator(nil, nil, true, tree.ndb)
		defer fastItr.Close()
		for ; fastItr.Valid(); fastItr.Next() {
			if err := tree.ndb.DeleteFastNode(fastItr.Key(), batch); err != nil {
				return false, err
			}
		}

		if err := tree.ndb.Commit(batch); err != nil {
			return false, err
		}
	}

	// Force garbage collection before we proceed to enabling fast storage.
	runtime.GC()

	batch := tree.NewBatch()
	if err := f.EnableFastStorageAndCommit(tree, batch); err != nil {
		tree.ndb.storageVersion = defaultStorageVersionValue
		return false, err
	}
	return true, nil
}

func (f *Fss) EnableFastStorageAndCommit(tree *MutableTree, batch dbm.Batch) error {
	var err error

	// We start a new thread to keep on checking if we are above 4GB, and if so garbage collect.
	// This thread only lasts during the fast node migration.
	// This is done to keep RAM usage down.
	done := make(chan struct{})
	defer func() {
		done <- struct{}{}
		close(done)
	}()

	go func() {
		timer := time.NewTimer(time.Second)
		var m runtime.MemStats

		for {
			// Sample the current memory usage
			runtime.ReadMemStats(&m)

			if m.Alloc > 4*1024*1024*1024 {
				// If we are using more than 4GB of memory, we should trigger garbage collection
				// to free up some memory.
				runtime.GC()
			}

			select {
			case <-timer.C:
				timer.Reset(time.Second)
			case <-done:
				if !timer.Stop() {
					<-timer.C
				}
				return
			}
		}
	}()

	itr := NewIterator(nil, nil, true, tree.ImmutableTree)
	defer itr.Close()
	var upgradedNodes uint64
	const verboseGap = 50000
	for ; itr.Valid(); itr.Next() {
		if err = tree.ndb.SaveFastNodeNoCache(NewFastNode(itr.Key(), itr.Value(), tree.version), batch); err != nil {
			return err
		}
		upgradedNodes++
		if upgradedNodes%verboseGap == 0 {
			tree.log(IavlInfo, "Upgrading to fast IAVL...", "finished", upgradedNodes)
		}
	}

	if err = itr.Error(); err != nil {
		return err
	}

	if err = tree.ndb.setFastStorageVersionToBatch(batch, tree.ndb.getLatestVersion()); err != nil {
		return err
	}

	return tree.ndb.Commit(batch)
}

func (f *Fss) DeleteFastNodesFromNodeDB(batch dbm.Batch, version int64) error {
	return f.ndb.DeleteFastNodesFrom(batch, version)
}

func (f *Fss) AsyncPersistFastNodeStart(version int64) *FastNodeChanges {
	return f.ndb.asyncPersistFastNodeStart(version)
}
func (f *Fss) AsyncPersistFastNodeFinished(event *commitEvent) {
	f.ndb.asyncPersistFastNodeFinished(event)
}
func (f *Fss) UpdateLatestVersion(version int64) {
	f.ndb.updateLatestVersion4FastNode(version)
}
