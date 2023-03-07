package iavl

import (
	"fmt"
	"strings"

	dbm "github.com/okx/okbchain/libs/tm-db"
)

// ImmutableTree contains the immutable tree at a given version. It is typically created by calling
// MutableTree.GetImmutable(), in which case the returned tree is safe for concurrent access as
// long as the version is not deleted via DeleteVersion() or the tree's pruning settings.
//
// Returned key/value byte slices must not be modified, since they may point to data located inside
// IAVL which would also be modified.
type ImmutableTree struct {
	root           *Node
	ndb            *nodeDB
	version        int64
	upgradeVersion int64
}

// NewImmutableTree creates both in-memory and persistent instances
func NewImmutableTree(db dbm.DB, cacheSize int) *ImmutableTree {
	if db == nil {
		// In-memory Tree.
		return &ImmutableTree{}
	}
	return &ImmutableTree{
		// NodeDB-backed Tree.
		ndb: newNodeDB(db, cacheSize, nil),
	}
}

// NewImmutableTreeWithOpts creates an ImmutableTree with the given options.
func NewImmutableTreeWithOpts(db dbm.DB, cacheSize int, opts *Options) *ImmutableTree {
	return &ImmutableTree{
		// NodeDB-backed Tree.
		ndb: newNodeDB(db, cacheSize, opts),
	}
}

// String returns a string representation of Tree.
func (t *ImmutableTree) String() string {
	leaves := []string{}
	t.Iterate(func(key []byte, val []byte) (stop bool) {
		leaves = append(leaves, fmt.Sprintf("%x: %x", key, val))
		return false
	})
	return "Tree{" + strings.Join(leaves, ", ") + "}"
}

// RenderShape provides a nested tree shape, ident is prepended in each level
// Returns an array of strings, one per line, to join with "\n" or display otherwise
func (t *ImmutableTree) RenderShape(indent string, encoder NodeEncoder) []string {
	if encoder == nil {
		encoder = defaultNodeEncoder
	}
	return t.renderNode(t.root, indent, 0, encoder)
}

// NodeEncoder will take an id (hash, or key for leaf nodes), the depth of the node,
// and whether or not this is a leaf node.
// It returns the string we wish to print, for iaviwer
type NodeEncoder func(id []byte, depth int, isLeaf bool) string

// defaultNodeEncoder can encode any node unless the client overrides it
func defaultNodeEncoder(id []byte, depth int, isLeaf bool) string {
	prefix := "- "
	if isLeaf {
		prefix = "* "
	}
	if len(id) == 0 {
		return fmt.Sprintf("%s<nil>", prefix)
	}
	return fmt.Sprintf("%s%X", prefix, id)
}

func (t *ImmutableTree) renderNode(node *Node, indent string, depth int, encoder func([]byte, int, bool) string) []string {
	prefix := strings.Repeat(indent, depth)
	// handle nil
	if node == nil {
		return []string{fmt.Sprintf("%s<nil>", prefix)}
	}
	// handle leaf
	if node.isLeaf() {
		here := fmt.Sprintf("%s%s", prefix, encoder(node.key, depth, true))
		return []string{here}
	}

	// recurse on inner node
	here := fmt.Sprintf("%s%s", prefix, encoder(node.hash, depth, false))
	left := t.renderNode(node.getLeftNode(t), indent, depth+1, encoder)
	right := t.renderNode(node.getRightNode(t), indent, depth+1, encoder)
	result := append(left, here)
	result = append(result, right...)
	return result
}

// Size returns the number of leaf nodes in the tree.
func (t *ImmutableTree) Size() int64 {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

// Version returns the version of the tree.
func (t *ImmutableTree) Version() int64 {
	return t.version
}

// Height returns the height of the tree.
func (t *ImmutableTree) Height() int8 {
	if t.root == nil {
		return 0
	}
	return t.root.height
}

// Has returns whether or not a key exists.
func (t *ImmutableTree) Has(key []byte) bool {
	if t.root == nil {
		return false
	}
	return t.root.has(t, key)
}

// Hash returns the root hash.
func (t *ImmutableTree) Hash() []byte {
	if t.root == nil {
		return nil
	}
	hash, _ := t.root.hashWithCount()
	return hash
}

// hashWithCount returns the root hash and hash count.
func (t *ImmutableTree) hashWithCount() ([]byte, int64) {
	if t.root == nil {
		return nil, 0
	}
	return t.root.hashWithCount()
}

// Export returns an iterator that exports tree nodes as ExportNodes. These nodes can be
// imported with MutableTree.Import() to recreate an identical tree.
func (t *ImmutableTree) Export() *Exporter {
	return newExporter(t)
}

// GetWithIndex returns the index and value of the specified key if it exists, or nil and the next index
// otherwise. The returned value must not be modified, since it may point to data stored within
// IAVL.
func (t *ImmutableTree) GetWithIndex(key []byte) (index int64, value []byte) {
	if t.root == nil {
		return 0, nil
	}
	return t.root.get(t, key)
}

// Get returns the value of the specified key if it exists, or nil.
// The returned value must not be modified, since it may point to data stored within IAVL.
// Get potentially employs a more performant strategy than GetWithIndex for retrieving the value.
func (t *ImmutableTree) Get(key []byte) []byte {
	if t.root == nil {
		return nil
	}

	if GetEnableFastStorage() {
		// attempt to get a FastNode directly from db/cache.
		// if call fails, fall back to the original IAVL logic in place.
		fastNode, err := t.ndb.GetFastNode(key)
		if err != nil {
			_, result := t.root.get(t, key)
			return result
		}

		if fastNode == nil {
			// If the tree is of the latest version and fast node is not in the tree
			// then the regular node is not in the tree either because fast node
			// represents live state.
			if t.version == t.ndb.getLatestMemoryVersion() {
				return nil
			}

			_, result := t.root.get(t, key)
			return result
		}

		if fastNode.versionLastUpdatedAt <= t.version {
			return fastNode.value
		}
	}
	// Otherwise the cached node was updated later than the current tree. In this case,
	// we need to use the regular stategy for reading from the current tree to avoid staleness.
	_, result := t.root.get(t, key)
	return result
}

// GetByIndex gets the key and value at the specified index.
func (t *ImmutableTree) GetByIndex(index int64) (key []byte, value []byte) {
	if t.root == nil {
		return nil, nil
	}
	return t.root.getByIndex(t, index)
}

// Iterate iterates over all keys of the tree. The keys and values must not be modified,
// since they may point to data stored within IAVL.Returns true if stopped by callback, false otherwise
func (t *ImmutableTree) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	if t.root == nil {
		return false
	}
	itr := t.Iterator(nil, nil, true)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		if fn(itr.Key(), itr.Value()) {
			return true
		}
	}
	return false
}

// Iterator returns an iterator over the immutable tree.
func (t *ImmutableTree) Iterator(start, end []byte, ascending bool) dbm.Iterator {
	if t.IsFastCacheEnabled() {
		return NewFastIteratorWithCache(start, end, ascending, t.ndb)
	}
	return NewIterator(start, end, ascending, t)
}

// IterateRange makes a callback for all nodes with key between start and end non-inclusive.
// If either are nil, then it is open on that side (nil, nil is the same as Iterate). The keys and
// values must not be modified, since they may point to data stored within IAVL.
func (t *ImmutableTree) IterateRange(start, end []byte, ascending bool, fn func(key []byte, value []byte) bool) (stopped bool) {
	if t.root == nil {
		return false
	}
	return t.root.traverseInRange(t, start, end, ascending, false, 0, false, func(node *Node, _ uint8) bool {
		if node.height == 0 {
			return fn(node.key, node.value)
		}
		return false
	})
}

// IterateRangeInclusive makes a callback for all nodes with key between start and end inclusive.
// If either are nil, then it is open on that side (nil, nil is the same as Iterate). The keys and
// values must not be modified, since they may point to data stored within IAVL.
func (t *ImmutableTree) IterateRangeInclusive(start, end []byte, ascending bool, fn func(key, value []byte, version int64) bool) (stopped bool) {
	if t.root == nil {
		return false
	}
	return t.root.traverseInRange(t, start, end, ascending, true, 0, false, func(node *Node, _ uint8) bool {
		if node.height == 0 {
			return fn(node.key, node.value, node.version)
		}
		return false
	})
}

// IsFastCacheEnabled returns true if fast cache is enabled, false otherwise.
// For fast cache to be enabled, the following 2 conditions must be met:
// 1. The tree is of the latest version.
// 2. The underlying storage has been upgraded to fast cache
func (t *ImmutableTree) IsFastCacheEnabled() bool {
	return GetEnableFastStorage() && t.isLatestTreeVersion() && t.ndb.hasUpgradedToFastStorage()
}

func (t *ImmutableTree) isLatestTreeVersion() bool {
	return t.version == t.ndb.getLatestMemoryVersion()
}

// Clone creates a clone of the tree.
// Used internally by MutableTree.
func (t *ImmutableTree) clone() *ImmutableTree {
	return &ImmutableTree{
		root:           t.root,
		ndb:            t.ndb,
		version:        t.version,
		upgradeVersion: -1,
	}
}

// nodeSize is like Size, but includes inner nodes too.
func (t *ImmutableTree) nodeSize() int {
	size := 0
	t.root.traverse(t, true, func(n *Node) bool {
		size++
		return false
	})
	return size
}

func (t *ImmutableTree) SetUpgradeVersion(version int64) {
	t.upgradeVersion = version
}

func (t *ImmutableTree) GetUpgradeVersion() int64 {
	return t.upgradeVersion
}

// Only used for debug!
func (t *ImmutableTree) DebugGetNode(nodeHash []byte) *Node {
	if string(t.Hash()) == string(nodeHash) {
		return t.root
	}
	return t.ndb.GetNode(nodeHash)
}

// Only used for debug!
func (t *ImmutableTree) DebugFssVersion() ([]byte, error) {
	return t.ndb.db.Get(metadataKeyFormat.Key([]byte(storageVersionKey)))
}
