package iavl

import (
	"bytes"
	"container/list"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/pkg/errors"
	dbm "github.com/tendermint/tm-db"
)

const (
	int64Size      = 8
	hashSize       = tmhash.Size
	genesisVersion = 1
)

var (
	// All node keys are prefixed with the byte 'n'. This ensures no collision is
	// possible with the other keys, and makes them easier to traverse. They are indexed by the node hash.
	nodeKeyFormat = NewKeyFormat('n', hashSize) // n<hash>

	// Orphans are keyed in the database by their expected lifetime.
	// The first number represents the *last* version at which the orphan needs
	// to exist, while the second number represents the *earliest* version at
	// which it is expected to exist - which starts out by being the version
	// of the node being orphaned.
	orphanKeyFormat = NewKeyFormat('o', int64Size, int64Size, hashSize) // o<last-version><first-version><hash>

	// Root nodes are indexed separately by their version
	rootKeyFormat = NewKeyFormat('r', int64Size) // r<version>

)

type nodeDB struct {
	mtx            sync.RWMutex     // Read/write lock.
	db             dbm.DB           // Persistent node storage.
	opts           Options          // Options to customize for pruning/writing
	versionReaders map[int64]uint32 // Number of active version readers

	latestVersion  int64
	nodeCache      map[string]*list.Element // Node cache.
	nodeCacheSize  int                      // Node cache size limit in elements.
	nodeCacheQueue *syncList                // LRU queue of cache elements. Used for deletion.

	orphanNodeCache         map[string]*Node
	heightOrphansCacheQueue *list.List
	heightOrphansCacheSize  int
	heightOrphansMap        map[int64]*heightOrphansItem

	prePersistNodeCache map[string]*Node
	tppMap              map[int64]*tppItem
	tppVersionList      *list.List

	dbReadTime    int64
	dbReadCount   int64
	nodeReadCount int64
	dbWriteCount  int64

	totalPersistedCount int64
	totalPersistedSize  int64
	totalDeletedCount   int64
	totalOrphanCount    int64

	name string
}

func newNodeDB(db dbm.DB, cacheSize int, opts *Options) *nodeDB {
	if opts == nil {
		o := DefaultOptions()
		opts = &o
	}
	return &nodeDB{
		db:                      db,
		opts:                    *opts,
		latestVersion:           0, // initially invalid
		nodeCache:               make(map[string]*list.Element),
		nodeCacheSize:           cacheSize,
		nodeCacheQueue:          newSyncList(),
		versionReaders:          make(map[int64]uint32, 8),
		orphanNodeCache:         make(map[string]*Node),
		heightOrphansCacheQueue: list.New(),
		heightOrphansCacheSize:  HeightOrphansCacheSize,
		heightOrphansMap:        make(map[int64]*heightOrphansItem),
		prePersistNodeCache:     make(map[string]*Node),
		tppMap:                  make(map[int64]*tppItem),
		tppVersionList:          list.New(),
		dbReadCount:             0,
		dbReadTime:              0,
		dbWriteCount:            0,
		name:                    ParseDBName(db),
	}
}

// GetNode gets a node from memory or disk. If it is an inner node, it does not
// load its children.
func (ndb *nodeDB) GetNode(hash []byte) *Node {

	res := func() *Node {

		ndb.mtx.RLock()
		defer ndb.mtx.RUnlock()
		ndb.addNodeReadCount()
		if len(hash) == 0 {
			panic("nodeDB.GetNode() requires hash")
		}
		if elem, ok := ndb.prePersistNodeCache[string(hash)]; ok {
			return elem
		}

		if elem, ok := ndb.getNodeInTpp(hash); ok { // GetNode from tpp
			return elem
		}
		// Check the cache.
		if elem, ok := ndb.nodeCache[string(hash)]; ok {
			// Already exists. Move to back of nodeCacheQueue.
			ndb.nodeCacheQueue.MoveToBack(elem)
			return elem.Value.(*Node)
		}
		if elem, ok := ndb.orphanNodeCache[string(hash)]; ok {
			return elem
		}

		return nil
	}()

	if res != nil {
		return res
	}

	// Doesn't exist, load.
	buf, err := ndb.dbGet(ndb.nodeKey(hash))
	if err != nil {
		panic(fmt.Sprintf("can't get node %X: %v", hash, err))
	}
	if buf == nil {
		panic(fmt.Sprintf("Value missing for hash %x corresponding to nodeKey %x", hash, ndb.nodeKey(hash)))
	}

	node, err := MakeNode(buf)
	if err != nil {
		panic(fmt.Sprintf("Error reading Node. bytes: %x, error: %v", buf, err))
	}

	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	node.hash = hash
	node.persisted = true
	ndb.cacheNodeByCheck(node)

	return node
}

func (ndb *nodeDB) getDbName() string {
	return ndb.name
}

// SaveNode saves a node to disk.
func (ndb *nodeDB) SaveNode(batch dbm.Batch, node *Node) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if node.hash == nil {
		panic("Expected to find node.hash, but none found.")
	}
	if node.persisted {
		panic("Shouldn't be calling save on an already persisted node.")
	}

	// Save node bytes to db.
	var buf bytes.Buffer
	buf.Grow(node.aminoSize())

	if err := node.writeBytes(&buf); err != nil {
		panic(err)
	}

	batch.Set(ndb.nodeKey(node.hash), buf.Bytes())
	ndb.log(IavlDebug, "BATCH SAVE %X %p", node.hash, node)
	node.persisted = true
	ndb.addDBWriteCount(1)
	ndb.cacheNode(node)
}

// Has checks if a hash exists in the database.
func (ndb *nodeDB) Has(hash []byte) (bool, error) {
	key := ndb.nodeKey(hash)

	if ldb, ok := ndb.db.(*dbm.GoLevelDB); ok {
		exists, err := ldb.DB().Has(key, nil)
		if err != nil {
			return false, err
		}
		return exists, nil
	}
	value, err := ndb.dbGet(key)
	if err != nil {
		return false, err
	}

	return value != nil, nil
}

// SaveBranch saves the given node and all of its descendants.
// NOTE: This function clears leftNode/rigthNode recursively and
// calls _hash() on the given node.
// TODO refactor, maybe use hashWithCount() but provide a callback.
func (ndb *nodeDB) SaveBranch(batch dbm.Batch, node *Node, savedNodes map[string]*Node) []byte {
	if node.persisted {
		return node.hash
	}

	if node.leftNode != nil {
		node.leftHash = ndb.SaveBranch(batch, node.leftNode, savedNodes)
	}
	if node.rightNode != nil {
		node.rightHash = ndb.SaveBranch(batch, node.rightNode, savedNodes)
	}

	node._hash()

	//resetBatch only working on generate a genesis block
	if node.version == genesisVersion {
		tmpBatch := ndb.NewBatch()
		ndb.SaveNode(tmpBatch, node)
		ndb.resetBatch(tmpBatch)
	} else {
		ndb.SaveNode(batch, node)
	}

	node.leftNode = nil
	node.rightNode = nil

	// TODO: handle magic number
	savedNodes[hex.EncodeToString(node.hash)] = node

	return node.hash
}

//resetBatch reset the db batch, keep low memory used
func (ndb *nodeDB) resetBatch(batch dbm.Batch) {
	var err error
	if ndb.opts.Sync {
		err = batch.WriteSync()
	} else {
		err = batch.Write()
	}
	if err != nil {
		panic(err)
	}
	batch.Close()
}

// DeleteVersion deletes a tree version from disk.
func (ndb *nodeDB) DeleteVersion(batch dbm.Batch, version int64, checkLatestVersion bool) error {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if ndb.versionReaders[version] > 0 {
		return errors.Errorf("unable to delete version %v, it has %v active readers", version, ndb.versionReaders[version])
	}

	ndb.deleteOrphans(batch, version)
	ndb.deleteRoot(batch, version, checkLatestVersion)
	return nil
}

// DeleteVersionsFrom permanently deletes all tree versions from the given version upwards.
func (ndb *nodeDB) DeleteVersionsFrom(batch dbm.Batch, version int64) error {
	latest := ndb.getLatestVersion()
	if latest < version {
		return nil
	}
	root, err := ndb.getRoot(latest)
	if err != nil {
		return err
	}
	if root == nil {
		return errors.Errorf("root for version %v not found", latest)
	}

	for v, r := range ndb.versionReaders {
		if v >= version && r != 0 {
			return errors.Errorf("unable to delete version %v with %v active readers", v, r)
		}
	}

	// First, delete all active nodes in the current (latest) version whose node version is after
	// the given version.
	err = ndb.deleteNodesFrom(batch, version, root)
	if err != nil {
		return err
	}

	// Next, delete orphans:
	// - Delete orphan entries *and referred nodes* with fromVersion >= version
	// - Delete orphan entries with toVersion >= version-1 (since orphans at latest are not orphans)
	ndb.traverseOrphans(func(key, hash []byte) {
		var fromVersion, toVersion int64
		orphanKeyFormat.Scan(key, &toVersion, &fromVersion)

		if fromVersion >= version {
			batch.Delete(key)
			batch.Delete(ndb.nodeKey(hash))
			ndb.uncacheNode(hash)
		} else if toVersion >= version-1 {
			batch.Delete(key)
		}
	})

	// Finally, delete the version root entries
	ndb.traverseRange(rootKeyFormat.Key(version), rootKeyFormat.Key(int64(math.MaxInt64)), func(k, v []byte) {
		batch.Delete(k)
	})

	return nil
}

// DeleteVersionsRange deletes versions from an interval (not inclusive).
func (ndb *nodeDB) DeleteVersionsRange(batch dbm.Batch, fromVersion, toVersion int64) error {
	if fromVersion >= toVersion {
		return errors.New("toVersion must be greater than fromVersion")
	}
	if toVersion == 0 {
		return errors.New("toVersion must be greater than 0")
	}

	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	latest := ndb.getLatestVersion()
	if latest < toVersion {
		return errors.Errorf("cannot delete latest saved version (%d)", latest)
	}

	predecessor := ndb.getPreviousVersion(fromVersion)

	for v, r := range ndb.versionReaders {
		if v < toVersion && v > predecessor && r != 0 {
			return errors.Errorf("unable to delete version %v with %v active readers", v, r)
		}
	}

	// If the predecessor is earlier than the beginning of the lifetime, we can delete the orphan.
	// Otherwise, we shorten its lifetime, by moving its endpoint to the predecessor version.
	for version := fromVersion; version < toVersion; version++ {
		ndb.traverseOrphansVersion(version, func(key, hash []byte) {
			var from, to int64
			orphanKeyFormat.Scan(key, &to, &from)
			batch.Delete(key)
			if from > predecessor {
				batch.Delete(ndb.nodeKey(hash))
				ndb.uncacheNode(hash)
			} else {
				ndb.saveOrphan(batch, hash, from, predecessor)
			}
		})
	}

	// Delete the version root entries
	ndb.traverseRange(rootKeyFormat.Key(fromVersion), rootKeyFormat.Key(toVersion), func(k, v []byte) {
		batch.Delete(k)
	})

	return nil
}

// deleteNodesFrom deletes the given node and any descendants that have versions after the given
// (inclusive). It is mainly used via LoadVersionForOverwriting, to delete the current version.
func (ndb *nodeDB) deleteNodesFrom(batch dbm.Batch, version int64, hash []byte) error {
	if len(hash) == 0 {
		return nil
	}

	node := ndb.GetNode(hash)
	if node.leftHash != nil {
		if err := ndb.deleteNodesFrom(batch, version, node.leftHash); err != nil {
			return err
		}
	}
	if node.rightHash != nil {
		if err := ndb.deleteNodesFrom(batch, version, node.rightHash); err != nil {
			return err
		}
	}

	if node.version >= version {
		batch.Delete(ndb.nodeKey(hash))
		ndb.uncacheNode(hash)
	}

	return nil
}

// Saves a single orphan to disk.
func (ndb *nodeDB) saveOrphan(batch dbm.Batch, hash []byte, fromVersion, toVersion int64) {
	if fromVersion > toVersion {
		panic(fmt.Sprintf("Orphan expires before it comes alive.  %d > %d", fromVersion, toVersion))
	}
	key := ndb.orphanKey(fromVersion, toVersion, hash)
	batch.Set(key, hash)
}

func (ndb *nodeDB) log(level int, format string, args ...interface{}) {
	iavlLog(ndb.name, level, format, args...)
}

// deleteOrphans deletes orphaned nodes from disk, and the associated orphan
// entries.
func (ndb *nodeDB) deleteOrphans(batch dbm.Batch, version int64) {
	// Will be zero if there is no previous version.
	predecessor := ndb.getPreviousVersion(version)

	// Traverse orphans with a lifetime ending at the version specified.
	// TODO optimize.
	ndb.traverseOrphansVersion(version, func(key, hash []byte) {
		var fromVersion, toVersion int64

		// See comment on `orphanKeyFmt`. Note that here, `version` and
		// `toVersion` are always equal.
		orphanKeyFormat.Scan(key, &toVersion, &fromVersion)

		// Delete orphan key and reverse-lookup key.
		batch.Delete(key)

		// If there is no predecessor, or the predecessor is earlier than the
		// beginning of the lifetime (ie: negative lifetime), or the lifetime
		// spans a single version and that version is the one being deleted, we
		// can delete the orphan.  Otherwise, we shorten its lifetime, by
		// moving its endpoint to the previous version.
		if predecessor < fromVersion || fromVersion == toVersion {
			ndb.log(IavlDebug, "DELETE predecessor:%v fromVersion:%v toVersion:%v %X", predecessor, fromVersion, toVersion, hash)
			batch.Delete(ndb.nodeKey(hash))
			ndb.uncacheNode(hash)
			ndb.totalDeletedCount++
		} else {
			ndb.log(IavlDebug, "MOVE predecessor:%v fromVersion:%v toVersion:%v %X", predecessor, fromVersion, toVersion, hash)
			ndb.saveOrphan(batch, hash, fromVersion, predecessor)
		}
	})
}

func (ndb *nodeDB) nodeKey(hash []byte) []byte {
	return nodeKeyFormat.KeyBytes(hash)
}

func (ndb *nodeDB) orphanKey(fromVersion, toVersion int64, hash []byte) []byte {
	return orphanKeyFormat.Key(toVersion, fromVersion, hash)
}

func (ndb *nodeDB) rootKey(version int64) []byte {
	return rootKeyFormat.Key(version)
}

func (ndb *nodeDB) getLatestVersion() int64 {
	if ndb.latestVersion == 0 {
		ndb.latestVersion = ndb.getPreviousVersion(1<<63 - 1)
	}
	return ndb.latestVersion
}

func (ndb *nodeDB) updateLatestVersion(version int64) {
	if ndb.latestVersion < version {
		ndb.latestVersion = version
	}
}

func (ndb *nodeDB) resetLatestVersion(version int64) {
	ndb.latestVersion = version
}

func (ndb *nodeDB) getPreviousVersion(version int64) int64 {
	itr, err := ndb.db.ReverseIterator(
		rootKeyFormat.Key(1),
		rootKeyFormat.Key(version),
	)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	pversion := int64(-1)
	for ; itr.Valid(); itr.Next() {
		k := itr.Key()
		rootKeyFormat.Scan(k, &pversion)
		return pversion
	}

	return 0
}

// deleteRoot deletes the root entry from disk, but not the node it points to.
func (ndb *nodeDB) deleteRoot(batch dbm.Batch, version int64, checkLatestVersion bool) {
	if checkLatestVersion && version == ndb.getLatestVersion() {
		panic("Tried to delete latest version")
	}
	batch.Delete(ndb.rootKey(version))
}

func (ndb *nodeDB) traverseOrphans(fn func(k, v []byte)) {
	ndb.traversePrefix(orphanKeyFormat.Key(), fn)
}

// Traverse orphans ending at a certain version.
func (ndb *nodeDB) traverseOrphansVersion(version int64, fn func(k, v []byte)) {
	ndb.traversePrefix(orphanKeyFormat.Key(version), fn)
}

// Traverse all keys.
func (ndb *nodeDB) traverse(fn func(key, value []byte)) {
	ndb.traverseRange(nil, nil, fn)
}

// Traverse all keys between a given range (excluding end).
func (ndb *nodeDB) traverseRange(start []byte, end []byte, fn func(k, v []byte)) {
	itr, err := ndb.db.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		fn(itr.Key(), itr.Value())
	}
}

// Traverse all keys with a certain prefix.
func (ndb *nodeDB) traversePrefix(prefix []byte, fn func(k, v []byte)) {
	itr, err := dbm.IteratePrefix(ndb.db, prefix)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		fn(itr.Key(), itr.Value())
	}
}

func (ndb *nodeDB) uncacheNode(hash []byte) {
	if elem, ok := ndb.nodeCache[string(hash)]; ok {
		ndb.nodeCacheQueue.Remove(elem)
		delete(ndb.nodeCache, string(hash))
	}
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (ndb *nodeDB) cacheNode(node *Node) {
	elem := ndb.nodeCacheQueue.PushBack(node)
	ndb.nodeCache[string(node.hash)] = elem

	if ndb.nodeCacheQueue.Len() > ndb.nodeCacheSize {
		oldest := ndb.nodeCacheQueue.Front()
		hash := ndb.nodeCacheQueue.Remove(oldest).(*Node).hash
		delete(ndb.nodeCache, string(hash))
	}
}

func (ndb *nodeDB) cacheNodeByCheck(node *Node) {
	if _, ok := ndb.nodeCache[string(node.hash)]; !ok {
		ndb.cacheNode(node)
	}
}

// Write to disk.
func (ndb *nodeDB) Commit(batch dbm.Batch) error {
	ndb.log(IavlDebug, "committing data to disk")
	var err error
	if ndb.opts.Sync {
		err = batch.WriteSync()
	} else {
		err = batch.Write()
	}
	if err != nil {
		return errors.Wrap(err, "failed to write batch")
	}

	batch.Close()

	return nil
}

func (ndb *nodeDB) getRoot(version int64) ([]byte, error) {
	return ndb.dbGet(ndb.rootKey(version))
}

func (ndb *nodeDB) getRoots() (map[int64][]byte, error) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	roots := map[int64][]byte{}

	ndb.traversePrefix(rootKeyFormat.Key(), func(k, v []byte) {
		var version int64
		rootKeyFormat.Scan(k, &version)
		roots[version] = v
	})
	return roots, nil
}

// SaveRoot creates an entry on disk for the given root, so that it can be
// loaded later.
func (ndb *nodeDB) SaveRoot(batch dbm.Batch, root *Node, version int64) error {
	if len(root.hash) == 0 {
		panic("SaveRoot: root hash should not be empty")
	}
	ndb.log(IavlDebug, "saving root hash(version %d) to disk", version)
	return ndb.saveRoot(batch, root.hash, version)
}

// SaveEmptyRoot creates an entry on disk for an empty root.
func (ndb *nodeDB) SaveEmptyRoot(batch dbm.Batch, version int64) error {
	return ndb.saveRoot(batch, []byte{}, version)
}

func (ndb *nodeDB) saveRoot(batch dbm.Batch, hash []byte, version int64) error {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if !EnableAsyncCommit {
		// We allow the initial version to be arbitrary
		latest := ndb.getLatestVersion()
		if !ignoreVersionCheck && latest > 0 && version != latest+1 {
			return fmt.Errorf("must save consecutive versions; expected %d, got %d", latest+1, version)
		}
	}

	batch.Set(ndb.rootKey(version), hash)
	ndb.updateLatestVersion(version)

	return nil
}

func (ndb *nodeDB) incrVersionReaders(version int64) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	ndb.versionReaders[version]++
}

func (ndb *nodeDB) decrVersionReaders(version int64) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	if ndb.versionReaders[version] > 0 {
		ndb.versionReaders[version]--
	}
}

// Utility and test functions

func (ndb *nodeDB) leafNodes() []*Node {
	leaves := []*Node{}

	ndb.traverseNodes(func(hash []byte, node *Node) {
		if node.isLeaf() {
			leaves = append(leaves, node)
		}
	})
	return leaves
}

func (ndb *nodeDB) nodes() []*Node {
	nodes := []*Node{}

	ndb.traverseNodes(func(hash []byte, node *Node) {
		nodes = append(nodes, node)
	})
	return nodes
}

func (ndb *nodeDB) orphans() [][]byte {
	orphans := [][]byte{}

	ndb.traverseOrphans(func(k, v []byte) {
		orphans = append(orphans, v)
	})
	return orphans
}

func (ndb *nodeDB) roots() map[int64][]byte {
	roots, _ := ndb.getRoots()
	return roots
}

// Not efficient.
// NOTE: DB cannot implement Size() because
// mutations are not always synchronous.
func (ndb *nodeDB) size() int {
	size := 0
	ndb.traverse(func(k, v []byte) {
		size++
	})
	return size
}

func (ndb *nodeDB) traverseNodes(fn func(hash []byte, node *Node)) {
	nodes := []*Node{}

	ndb.traversePrefix(nodeKeyFormat.Key(), func(key, value []byte) {
		node, err := MakeNode(value)
		if err != nil {
			panic(fmt.Sprintf("Couldn't decode node from database: %v", err))
		}
		nodeKeyFormat.Scan(key, &node.hash)
		nodes = append(nodes, node)
	})

	sort.Slice(nodes, func(i, j int) bool {
		return bytes.Compare(nodes[i].key, nodes[j].key) < 0
	})

	for _, n := range nodes {
		fn(n.hash, n)
	}
}

func (ndb *nodeDB) String() string {
	var str string
	index := 0

	ndb.traversePrefix(rootKeyFormat.Key(), func(key, value []byte) {
		str += fmt.Sprintf("%s: %x\n", string(key), value)
	})
	str += "\n"

	ndb.traverseOrphans(func(key, value []byte) {
		str += fmt.Sprintf("%s: %x\n", string(key), value)
	})
	str += "\n"

	ndb.traverseNodes(func(hash []byte, node *Node) {
		switch {
		case len(hash) == 0:
			str += "<nil>\n"
		case node == nil:
			str += fmt.Sprintf("%s%40x: <nil>\n", nodeKeyFormat.Prefix(), hash)
		case node.value == nil && node.height > 0:
			str += fmt.Sprintf("%s%40x: %s   %-16s h=%d version=%d\n",
				nodeKeyFormat.Prefix(), hash, node.key, "", node.height, node.version)
		default:
			str += fmt.Sprintf("%s%40x: %s = %-16s h=%d version=%d\n",
				nodeKeyFormat.Prefix(), hash, node.key, node.value, node.height, node.version)
		}
		index++
	})
	return "-" + "\n" + str + "-"
}
