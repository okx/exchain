package iavl

import (
	"bytes"
	"container/list"
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"
	dbm "github.com/tendermint/tm-db"
)

func SetIgnoreVersionCheck(check bool) {
	ignoreVersionCheck = check
}

var (
	ignoreVersionCheck = false
)

// MutableTree is a persistent tree which keeps track of versions. It is not safe for concurrent
// use, and should be guarded by a Mutex or RWLock as appropriate. An immutable tree at a given
// version can be returned via GetImmutable, which is safe for concurrent access.
//
// Given and returned key/value byte slices must not be modified, since they may point to data
// located inside IAVL which would also be modified.
//
// The inner ImmutableTree should not be used directly by callers.
type MutableTree struct {
	*ImmutableTree                   // The current, working tree.
	lastSaved       *ImmutableTree   // The most recently saved tree.
	orphans         []*Node          // Nodes removed by changes to working tree.Will refresh after each block
	commitOrphans   map[string]int64 // Nodes removed by changes to working tree.Will refresh after each commit.
	versions        *SyncMap         // The previous, saved versions of the tree.
	removedVersions sync.Map         // The removed versions of the tree.
	ndb             *nodeDB

	committedHeightQueue *list.List
	committedHeightMap   map[int64]bool
	historyStateNum      int

	commitCh          chan commitEvent
	lastPersistHeight int64
}

// NewMutableTree returns a new tree with the specified cache size and datastore.
func NewMutableTree(db dbm.DB, cacheSize int) (*MutableTree, error) {
	return NewMutableTreeWithOpts(db, cacheSize, nil)
}

// NewMutableTreeWithOpts returns a new tree with the specified options.
func NewMutableTreeWithOpts(db dbm.DB, cacheSize int, opts *Options) (*MutableTree, error) {
	ndb := newNodeDB(db, cacheSize, opts)
	head := &ImmutableTree{ndb: ndb}
	var initVersion int64
	if opts != nil {
		initVersion = int64(opts.InitialVersion)
	} else {
		initVersion = 0
	}
	tree := &MutableTree{
		ImmutableTree: head,
		lastSaved:     head.clone(),
		orphans:       []*Node{},
		commitOrphans: map[string]int64{},
		versions:      NewSyncMap(),
		ndb:           ndb,

		committedHeightMap:   map[int64]bool{},
		committedHeightQueue: list.New(),
		historyStateNum:      MaxCommittedHeightNum,

		commitCh:          make(chan commitEvent),
		lastPersistHeight: initVersion,
	}

	if tree.historyStateNum < minHistoryStateNum {
		tree.historyStateNum = minHistoryStateNum
	}

	if EnableAsyncCommit {
		treeMap.addNewTree(tree)
	}

	return tree, nil
}

// IsEmpty returns whether or not the tree has any keys. Only trees that are
// not empty can be saved.
func (tree *MutableTree) IsEmpty() bool {
	return tree.ImmutableTree.Size() == 0
}

// VersionExists returns whether or not a version exists.
func (tree *MutableTree) VersionExists(version int64) bool {
	tree.ndb.mtx.Lock()
	defer tree.ndb.mtx.Unlock()
	if tree.ndb.heightOrphansMap[version] != nil {
		return true
	}
	return tree.versions.Get(version)
}

// AvailableVersions returns all available versions in ascending order
func (tree *MutableTree) AvailableVersions() []int {
	res := make([]int, 0, tree.versions.Len())
	tree.versions.Range(func(key int64, value bool) bool {
		if value {
			res = append(res, int(key))
		}
		return true
	})
	sort.Ints(res)
	return res
}

// Hash returns the hash of the latest saved version of the tree, as returned
// by SaveVersion. If no versions have been saved, Hash returns nil.
func (tree *MutableTree) Hash() []byte {
	if tree.version > 0 {
		return tree.lastSaved.Hash()
	}
	return nil
}

// WorkingHash returns the hash of the current working tree.
func (tree *MutableTree) WorkingHash() []byte {
	return tree.ImmutableTree.Hash()
}

// String returns a string representation of the tree.
func (tree *MutableTree) String() string {
	return tree.ndb.String()
}

// Set/Remove will orphan at most tree.Height nodes,
// balancing the tree after a Set/Remove will orphan at most 3 nodes.
func (tree *MutableTree) prepareOrphansSlice() []*Node {
	return make([]*Node, 0, tree.Height()+3)
}

// Set sets a key in the working tree. Nil values are invalid. The given key/value byte slices must
// not be modified after this call, since they point to slices stored within IAVL.
func (tree *MutableTree) Set(key, value []byte) bool {
	orphaned, updated := tree.set(key, value)
	tree.addOrphans(orphaned)
	return updated
}

// Import returns an importer for tree nodes previously exported by ImmutableTree.Export(),
// producing an identical IAVL tree. The caller must call Close() on the importer when done.
//
// version should correspond to the version that was initially exported. It must be greater than
// or equal to the highest ExportNode version number given.
//
// Import can only be called on an empty tree. It is the callers responsibility that no other
// modifications are made to the tree while importing.
func (tree *MutableTree) Import(version int64) (*Importer, error) {
	return newImporter(tree, version)
}

func (tree *MutableTree) set(key []byte, value []byte) (orphans []*Node, updated bool) {
	if value == nil {
		panic(fmt.Sprintf("Attempt to store nil value at key '%s'", key))
	}

	if tree.ImmutableTree.root == nil {
		tree.ImmutableTree.root = NewNode(key, value, tree.version+1)
		return nil, updated
	}

	orphans = tree.prepareOrphansSlice()
	tree.ImmutableTree.root, updated = tree.recursiveSet(tree.ImmutableTree.root, key, value, &orphans)
	return orphans, updated
}

func (tree *MutableTree) recursiveSet(node *Node, key []byte, value []byte, orphans *[]*Node) (
	newSelf *Node, updated bool,
) {
	version := tree.version + 1

	if node.isLeaf() {
		switch bytes.Compare(key, node.key) {
		case -1:
			return &Node{
				key:       node.key,
				height:    1,
				size:      2,
				leftNode:  NewNode(key, value, version),
				rightNode: node,
				version:   version,
			}, false
		case 1:
			return &Node{
				key:       key,
				height:    1,
				size:      2,
				leftNode:  node,
				rightNode: NewNode(key, value, version),
				version:   version,
			}, false
		default:
			*orphans = append(*orphans, node)
			return NewNode(key, value, version), true
		}
	} else {
		*orphans = append(*orphans, node)
		node = node.clone(version)

		if bytes.Compare(key, node.key) < 0 {
			node.leftNode, updated = tree.recursiveSet(node.getLeftNode(tree.ImmutableTree), key, value, orphans)
			node.leftHash = nil // leftHash is yet unknown
		} else {
			node.rightNode, updated = tree.recursiveSet(node.getRightNode(tree.ImmutableTree), key, value, orphans)
			node.rightHash = nil // rightHash is yet unknown
		}

		if updated {
			return node, updated
		}
		node.calcHeightAndSize(tree.ImmutableTree)
		newNode := tree.balance(node, orphans)
		return newNode, updated
	}
}

// Remove removes a key from the working tree. The given key byte slice should not be modified
// after this call, since it may point to data stored inside IAVL.
func (tree *MutableTree) Remove(key []byte) ([]byte, bool) {
	val, orphaned, removed := tree.remove(key)
	tree.addOrphans(orphaned)
	return val, removed
}

// remove tries to remove a key from the tree and if removed, returns its
// value, nodes orphaned and 'true'.
func (tree *MutableTree) remove(key []byte) (value []byte, orphaned []*Node, removed bool) {
	if tree.root == nil {
		return nil, nil, false
	}
	orphaned = tree.prepareOrphansSlice()
	newRootHash, newRoot, _, value := tree.recursiveRemove(tree.root, key, &orphaned)
	if len(orphaned) == 0 {
		return nil, nil, false
	}

	if newRoot == nil && newRootHash != nil {
		tree.root = tree.ndb.GetNode(newRootHash)
	} else {
		tree.root = newRoot
	}
	return value, orphaned, true
}

// removes the node corresponding to the passed key and balances the tree.
// It returns:
// - the hash of the new node (or nil if the node is the one removed)
// - the node that replaces the orig. node after remove
// - new leftmost leaf key for tree after successfully removing 'key' if changed.
// - the removed value
// - the orphaned nodes.
func (tree *MutableTree) recursiveRemove(node *Node, key []byte, orphans *[]*Node) (newHash []byte, newSelf *Node, newKey []byte, newValue []byte) {
	version := tree.version + 1

	if node.isLeaf() {
		if bytes.Equal(key, node.key) {
			*orphans = append(*orphans, node)
			return nil, nil, nil, node.value
		}
		return node.hash, node, nil, nil
	}

	// node.key < key; we go to the left to find the key:
	if bytes.Compare(key, node.key) < 0 {
		newLeftHash, newLeftNode, newKey, value := tree.recursiveRemove(node.getLeftNode(tree.ImmutableTree), key, orphans) //nolint:govet

		if len(*orphans) == 0 {
			return node.hash, node, nil, value
		}
		*orphans = append(*orphans, node)
		if newLeftHash == nil && newLeftNode == nil { // left node held value, was removed
			return node.rightHash, node.rightNode, node.key, value
		}

		newNode := node.clone(version)
		newNode.leftHash, newNode.leftNode = newLeftHash, newLeftNode
		newNode.calcHeightAndSize(tree.ImmutableTree)
		newNode = tree.balance(newNode, orphans)
		return newNode.hash, newNode, newKey, value
	}
	// node.key >= key; either found or look to the right:
	newRightHash, newRightNode, newKey, value := tree.recursiveRemove(node.getRightNode(tree.ImmutableTree), key, orphans)

	if len(*orphans) == 0 {
		return node.hash, node, nil, value
	}
	*orphans = append(*orphans, node)
	if newRightHash == nil && newRightNode == nil { // right node held value, was removed
		return node.leftHash, node.leftNode, nil, value
	}

	newNode := node.clone(version)
	newNode.rightHash, newNode.rightNode = newRightHash, newRightNode
	if newKey != nil {
		newNode.key = newKey
	}
	newNode.calcHeightAndSize(tree.ImmutableTree)
	newNode = tree.balance(newNode, orphans)
	return newNode.hash, newNode, nil, value
}

// Load the latest versioned tree from disk.
func (tree *MutableTree) Load() (int64, error) {
	return tree.LoadVersion(int64(0))
}

// LazyLoadVersion attempts to lazy load only the specified target version
// without loading previous roots/versions. Lazy loading should be used in cases
// where only reads are expected. Any writes to a lazy loaded tree may result in
// unexpected behavior. If the targetVersion is non-positive, the latest version
// will be loaded by default. If the latest version is non-positive, this method
// performs a no-op. Otherwise, if the root does not exist, an error will be
// returned.
func (tree *MutableTree) LazyLoadVersion(targetVersion int64) (int64, error) {
	latestVersion := tree.ndb.getLatestVersion()
	if latestVersion < targetVersion {
		return latestVersion, fmt.Errorf("wanted to load target %d but only found up to %d", targetVersion, latestVersion)
	}

	// no versions have been saved if the latest version is non-positive
	if latestVersion <= 0 {
		return 0, nil
	}

	// default to the latest version if the targeted version is non-positive
	if targetVersion <= 0 {
		targetVersion = latestVersion
	}

	rootHash, err := tree.ndb.getRoot(targetVersion)
	if err != nil {
		return 0, err
	}
	if rootHash == nil {
		return latestVersion, ErrVersionDoesNotExist
	}

	tree.versions.Set(targetVersion, true)

	iTree := &ImmutableTree{
		ndb:     tree.ndb,
		version: targetVersion,
		root:    tree.ndb.GetNode(rootHash),
	}

	tree.orphans = []*Node{}
	tree.commitOrphans = map[string]int64{}
	tree.ImmutableTree = iTree
	tree.lastSaved = iTree.clone()

	return targetVersion, nil
}

// Returns the version number of the latest version found
func (tree *MutableTree) LoadVersion(targetVersion int64) (int64, error) {
	roots, err := tree.ndb.getRoots()
	if err != nil {
		return 0, err
	}

	if len(roots) == 0 {
		return 0, nil
	}

	firstVersion := int64(0)
	latestVersion := int64(0)

	var latestRoot []byte
	for version, r := range roots {
		tree.versions.Set(version, true)
		if version > latestVersion && (targetVersion == 0 || version <= targetVersion) {
			latestVersion = version
			latestRoot = r
		}
		if firstVersion == 0 || version < firstVersion {
			firstVersion = version
		}
	}

	if !(targetVersion == 0 || latestVersion == targetVersion) {
		return latestVersion, fmt.Errorf("wanted to load target %v but only found up to %v",
			targetVersion, latestVersion)
	}

	if firstVersion > 0 && firstVersion < int64(tree.ndb.opts.InitialVersion) {
		return latestVersion, fmt.Errorf("initial version set to %v, but found earlier version %v",
			tree.ndb.opts.InitialVersion, firstVersion)
	}

	t := &ImmutableTree{
		ndb:     tree.ndb,
		version: latestVersion,
	}

	if len(latestRoot) != 0 {
		t.root = tree.ndb.GetNode(latestRoot)
	}

	tree.orphans = []*Node{}
	tree.commitOrphans = map[string]int64{}
	tree.ImmutableTree = t
	tree.lastSaved = t.clone()
	tree.lastPersistHeight = latestVersion

	return latestVersion, nil
}

// LoadVersionForOverwriting attempts to load a tree at a previously committed
// version, or the latest version below it. Any versions greater than targetVersion will be deleted.
func (tree *MutableTree) LoadVersionForOverwriting(targetVersion int64) (int64, error) {
	latestVersion, err := tree.LoadVersion(targetVersion)
	if err != nil {
		return latestVersion, err
	}

	batch := tree.NewBatch()
	if err = tree.ndb.DeleteVersionsFrom(batch, targetVersion+1); err != nil {
		return latestVersion, err
	}

	if err = tree.ndb.Commit(batch); err != nil {
		return latestVersion, err
	}

	tree.ndb.resetLatestVersion(latestVersion)

	tree.versions.Range(func(key int64, value bool) bool {
		if key > targetVersion {
			tree.versions.DeleteWithoutLock(key)
		}
		return true
	})
	return latestVersion, nil
}

// GetImmutable loads an ImmutableTree at a given version for querying. The returned tree is
// safe for concurrent access, provided the version is not deleted, e.g. via `DeleteVersion()`.
func (tree *MutableTree) GetImmutable(version int64) (*ImmutableTree, error) {
	rootHash, err := tree.ndb.getRootWithCacheAndDB(version)
	if err != nil {
		return nil, err
	}
	if rootHash == nil {
		return nil, ErrVersionDoesNotExist
	} else if len(rootHash) == 0 {
		return &ImmutableTree{
			ndb:     tree.ndb,
			version: version,
		}, nil
	}
	return &ImmutableTree{
		root:    tree.ndb.GetNode(rootHash),
		ndb:     tree.ndb,
		version: version,
	}, nil
}

// Rollback resets the working tree to the latest saved version, discarding
// any unsaved modifications.
func (tree *MutableTree) Rollback() {
	if tree.version > 0 {
		tree.ImmutableTree = tree.lastSaved.clone()
	} else {
		tree.ImmutableTree = &ImmutableTree{ndb: tree.ndb, version: 0}
	}
	tree.orphans = []*Node{}
	tree.commitOrphans = map[string]int64{}
}

// GetVersioned gets the value at the specified key and version. The returned value must not be
// modified, since it may point to data stored within IAVL.
func (tree *MutableTree) GetVersioned(key []byte, version int64) (
	index int64, value []byte,
) {
	tree.log(IavlErr, "GetVersioned KEY:%s, VERSION:%d\n", fmt.Sprintf("%x", key), version)
	if tree.versions.Get(version) {
		tree.log(IavlErr, "GetVersioned KEY:%s, VERSION:%d exist\n", fmt.Sprintf("%x", key), version)
		t, err := tree.GetImmutable(version)
		if err != nil {
			tree.log(IavlErr, "GetVersioned KEY:%s, VERSION:%d, ERROR:%s\n", fmt.Sprintf("%x", key), version, err.Error())
			return -1, nil
		}

		index, value = t.Get(key)
		tree.log(IavlErr,"GetVersioned KEY:%s, INDEX:%d, VALUE:%s, VERSION:%d\n", fmt.Sprintf("%x", key), index, fmt.Sprintf("%x", value), version)
		return
	}
	return -1, nil
}

// SaveVersion saves a new tree version to disk, based on the current state of
// the tree. Returns the hash and new version number.
func (tree *MutableTree) SaveVersion() ([]byte, int64, error) {
	version := tree.version + 1
	if version == 1 && tree.ndb.opts.InitialVersion > 0 {
		version = int64(tree.ndb.opts.InitialVersion) + 1
	}
	if !ignoreVersionCheck && tree.versions.Get(version) {
		// If the version already exists, return an error as we're attempting to overwrite.
		// However, the same hash means idempotent (i.e. no-op).
		existingHash, err := tree.ndb.getRoot(version)
		if err != nil {
			return nil, version, err
		}

		var newHash = tree.WorkingHash()

		if bytes.Equal(existingHash, newHash) {
			tree.version = version
			tree.ImmutableTree = tree.ImmutableTree.clone()
			tree.lastSaved = tree.ImmutableTree.clone()
			tree.orphans = []*Node{}
			tree.commitOrphans = map[string]int64{}
			return existingHash, version, nil
		}

		return nil, version, fmt.Errorf("version %d was already saved to different hash %X (existing hash %X)", version, newHash, existingHash)
	}

	if EnableAsyncCommit {
		return tree.SaveVersionAsync(version)
	}
	return tree.SaveVersionSync(version)
}

func (tree *MutableTree) SaveVersionSync(version int64) ([]byte, int64, error) {
	batch := tree.NewBatch()
	if tree.root == nil {
		// There can still be orphans, for example if the root is the node being
		// removed.
		tree.log(IavlDebug, "SAVE EMPTY TREE %v", version)
		tree.ndb.SaveOrphans(batch, version, tree.orphans)
		if err := tree.ndb.SaveEmptyRoot(batch, version); err != nil {
			return nil, 0, err
		}
	} else {
		tree.log(IavlDebug, "SAVE TREE %v", version)
		tree.ndb.SaveBranch(batch, tree.root)
		tree.ndb.SaveOrphans(batch, version, tree.orphans)
		if err := tree.ndb.SaveRoot(batch, tree.root, version); err != nil {
			return nil, 0, err
		}
	}

	if err := tree.ndb.Commit(batch); err != nil {
		return nil, version, err
	}

	tree.version = version
	tree.versions.Set(version, true)

	// set new working tree
	tree.ImmutableTree = tree.ImmutableTree.clone()
	tree.lastSaved = tree.ImmutableTree.clone()
	tree.orphans = []*Node{}

	tree.ndb.log(IavlDebug, tree.ndb.sprintCacheLog(version))
	return tree.Hash(), version, nil
}

func (tree *MutableTree) deleteVersion(batch dbm.Batch, version int64, versions map[int64]bool) error {
	if version == 0 {
		return errors.New("version must be greater than 0")
	}
	if version == tree.version {
		return errors.Errorf("cannot delete latest saved version (%d)", version)
	}
	if _, ok := versions[version]; !ok {
		var logStr string
		for v, isTrue := range versions {
			logStr += fmt.Sprintf("%d:%t, ", v, isTrue)
		}
		tree.log(IavlDebug, logStr)
		return errors.Wrap(ErrVersionDoesNotExist, fmt.Sprintf("%d", version))
	}

	if err := tree.ndb.DeleteVersion(batch, version, true); err != nil {
		return err
	}

	return nil
}

// SetInitialVersion sets the initial version of the tree, replacing Options.InitialVersion.
// It is only used during the initial SaveVersion() call for a tree with no other versions,
// and is otherwise ignored.
func (tree *MutableTree) SetInitialVersion(version uint64) {
	tree.ndb.opts.InitialVersion = version
}

// DeleteVersions deletes a series of versions from the MutableTree.
// Deprecated: please use DeleteVersionsRange instead.
func (tree *MutableTree) DeleteVersions(versions ...int64) error {
	tree.log(IavlDebug, "DELETING VERSIONS: %v", versions)

	if len(versions) == 0 {
		return nil
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	// Find ordered data and delete by interval
	intervals := map[int64]int64{}
	var fromVersion int64
	for _, version := range versions {
		if version-fromVersion != intervals[fromVersion] {
			fromVersion = version
		}
		intervals[fromVersion]++
	}

	for fromVersion, sortedBatchSize := range intervals {
		if err := tree.DeleteVersionsRange(fromVersion, fromVersion+sortedBatchSize); err != nil {
			return err
		}
	}

	return nil
}

// DeleteVersionsRange removes versions from an interval from the MutableTree (not inclusive).
// An error is returned if any single version has active readers.
// All writes happen in a single batch with a single commit.
func (tree *MutableTree) DeleteVersionsRange(fromVersion, toVersion int64) error {
	batch := tree.NewBatch()
	if err := tree.ndb.DeleteVersionsRange(batch, fromVersion, toVersion); err != nil {
		return err
	}

	if err := tree.ndb.Commit(batch); err != nil {
		return err
	}

	for version := fromVersion; version < toVersion; version++ {
		tree.versions.Delete(version)
	}
	return nil
}

// DeleteVersion deletes a tree version from disk. The version can then no
// longer be accessed.
func (tree *MutableTree) DeleteVersion(version int64) error {
	tree.log(IavlDebug, "DELETE VERSION: %d", version)
	batch := tree.NewBatch()
	if err := tree.deleteVersion(batch, version, tree.versions.Clone()); err != nil {
		return err
	}

	if err := tree.ndb.Commit(batch); err != nil {
		return err
	}

	tree.versions.Delete(version)
	return nil
}

// Rotate right and return the new node and orphan.
func (tree *MutableTree) rotateRight(node *Node) (*Node, *Node) {
	version := tree.version + 1

	// TODO: optimize balance & rotate.
	node = node.clone(version)
	orphaned := node.getLeftNode(tree.ImmutableTree)
	newNode := orphaned.clone(version)

	newNoderHash, newNoderCached := newNode.rightHash, newNode.rightNode
	newNode.rightHash, newNode.rightNode = node.hash, node
	node.leftHash, node.leftNode = newNoderHash, newNoderCached

	node.calcHeightAndSize(tree.ImmutableTree)
	newNode.calcHeightAndSize(tree.ImmutableTree)

	return newNode, orphaned
}

// Rotate left and return the new node and orphan.
func (tree *MutableTree) rotateLeft(node *Node) (*Node, *Node) {
	version := tree.version + 1

	// TODO: optimize balance & rotate.
	node = node.clone(version)
	orphaned := node.getRightNode(tree.ImmutableTree)
	newNode := orphaned.clone(version)

	newNodelHash, newNodelCached := newNode.leftHash, newNode.leftNode
	newNode.leftHash, newNode.leftNode = node.hash, node
	node.rightHash, node.rightNode = newNodelHash, newNodelCached

	node.calcHeightAndSize(tree.ImmutableTree)
	newNode.calcHeightAndSize(tree.ImmutableTree)

	return newNode, orphaned
}

// NOTE: assumes that node can be modified
// TODO: optimize balance & rotate
func (tree *MutableTree) balance(node *Node, orphans *[]*Node) (newSelf *Node) {
	if node.persisted || node.prePersisted {
		panic("Unexpected balance() call on persisted node")
	}
	balance := node.calcBalance(tree.ImmutableTree)

	if balance > 1 {
		if node.getLeftNode(tree.ImmutableTree).calcBalance(tree.ImmutableTree) >= 0 {
			// Left Left Case
			newNode, orphaned := tree.rotateRight(node)
			*orphans = append(*orphans, orphaned)
			return newNode
		}
		// Left Right Case
		var leftOrphaned *Node

		left := node.getLeftNode(tree.ImmutableTree)
		node.leftHash = nil
		node.leftNode, leftOrphaned = tree.rotateLeft(left)
		newNode, rightOrphaned := tree.rotateRight(node)
		*orphans = append(*orphans, left, leftOrphaned, rightOrphaned)
		return newNode
	}
	if balance < -1 {
		if node.getRightNode(tree.ImmutableTree).calcBalance(tree.ImmutableTree) <= 0 {
			// Right Right Case
			newNode, orphaned := tree.rotateLeft(node)
			*orphans = append(*orphans, orphaned)
			return newNode
		}
		// Right Left Case
		var rightOrphaned *Node

		right := node.getRightNode(tree.ImmutableTree)
		node.rightHash = nil
		node.rightNode, rightOrphaned = tree.rotateRight(right)
		newNode, leftOrphaned := tree.rotateLeft(node)

		*orphans = append(*orphans, right, leftOrphaned, rightOrphaned)
		return newNode
	}
	// Nothing changed
	return node
}

func (tree *MutableTree) addOrphans(orphans []*Node) {
	if EnableAsyncCommit {
		tree.addOrphansOptimized(orphans)
	} else {
		for _, node := range orphans {
			if !node.persisted {
				// We don't need to orphan nodes that were never persisted.
				continue
			}
			if len(node.hash) == 0 {
				panic("Expected to find node hash, but was empty")
			}
			tree.orphans = append(tree.orphans, node)
		}
	}
}
