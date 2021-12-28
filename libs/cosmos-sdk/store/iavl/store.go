package iavl

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/okex/exchain/libs/iavl"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	tmkv "github.com/okex/exchain/libs/tendermint/libs/kv"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

var (
	FlagIavlCacheSize = "iavl-cache-size"

	IavlCacheSize = 1000000
)

var (
	_ types.KVStore       = (*Store)(nil)
	_ types.CommitStore   = (*Store)(nil)
	_ types.CommitKVStore = (*Store)(nil)
	_ types.Queryable     = (*Store)(nil)
)

// Store Implements types.KVStore and CommitKVStore.
type Store struct {
	tree             Tree
	flatKVDB         dbm.DB
	cache            map[string][]byte
	flatKVReadTime   int64
	flatKVReadCount  int64
	flatKVWriteCount int64
}

func (st *Store) GetFlatKVReadTime() int {
	return int(atomic.LoadInt64(&st.flatKVReadTime))
}

func (st *Store) addFlatKVReadTime(ts int64) {
	atomic.AddInt64(&st.flatKVReadTime, ts)
}

func (st *Store) resetFlatKVReadTime() {
	atomic.StoreInt64(&st.flatKVReadTime, 0)
}

func (st *Store) GetFlatKVReadCount() int {
	return int(atomic.LoadInt64(&st.flatKVReadCount))
}

func (st *Store) addFlatKVReadCount() {
	atomic.AddInt64(&st.flatKVReadCount, 1)
}

func (st *Store) resetFlatKVReadCount() {
	atomic.StoreInt64(&st.flatKVReadCount, 0)
}

func (st *Store) GetFlatKVWriteCount() int {
	return int(atomic.LoadInt64(&st.flatKVWriteCount))
}

func (st *Store) addFlatKVWriteCount() {
	atomic.AddInt64(&st.flatKVWriteCount, 1)
}

func (st *Store) resetFlatKVWriteCount() {
	atomic.StoreInt64(&st.flatKVWriteCount, 0)
}

func (st *Store) getCache(key string) (value []byte, ok bool) {
	value, ok = st.cache[key]
	return
}

func (st *Store) addCache(key string, value []byte) {
	st.cache[key] = value
}

func (st *Store) deleteCache(key string) {
	delete(st.cache, key)
}

func (st *Store) StopStore() {
	tr := st.tree.(*iavl.MutableTree)
	tr.StopTree()
}

func (st *Store) GetHeights() map[int64][]byte {
	return st.tree.GetPersistedRoots()
}

// LoadStore returns an IAVL Store as a CommitKVStore. Internally, it will load the
// store's version (id) from the provided DB. An error is returned if the version
// fails to load.
func LoadStore(db dbm.DB, flatKVDB dbm.DB, id types.CommitID, lazyLoading bool, startVersion int64) (types.CommitKVStore, error) {
	return LoadStoreWithInitialVersion(db, flatKVDB, id, lazyLoading, uint64(startVersion))
}

// LoadStore returns an IAVL Store as a CommitKVStore setting its initialVersion
// to the one given. Internally, it will load the store's version (id) from the
// provided DB. An error is returned if the version fails to load.
func LoadStoreWithInitialVersion(db dbm.DB, flatKVDB dbm.DB, id types.CommitID, lazyLoading bool, initialVersion uint64) (types.CommitKVStore, error) {
	tree, err := iavl.NewMutableTreeWithOpts(db, IavlCacheSize, &iavl.Options{InitialVersion: initialVersion})
	if err != nil {
		return nil, err
	}

	if lazyLoading {
		_, err = tree.LazyLoadVersion(id.Version)
	} else {
		_, err = tree.LoadVersion(id.Version)
	}

	if err != nil {
		return nil, err
	}

	return &Store{
		tree:             tree,
		flatKVDB:         flatKVDB,
		cache:            make(map[string][]byte),
		flatKVReadTime:   0,
		flatKVReadCount:  0,
		flatKVWriteCount: 0,
	}, nil
}

func GetCommitVersion(db dbm.DB) (int64, error) {
	tree, err := iavl.NewMutableTreeWithOpts(db, IavlCacheSize, &iavl.Options{InitialVersion: 0})
	if err != nil {
		return 0, err
	}
	return tree.GetCommitVersion(), nil
}

// UnsafeNewStore returns a reference to a new IAVL Store with a given mutable
// IAVL tree reference. It should only be used for testing purposes.
//
// CONTRACT: The IAVL tree should be fully loaded.
// CONTRACT: PruningOptions passed in as argument must be the same as pruning options
// passed into iavl.MutableTree
func UnsafeNewStore(tree *iavl.MutableTree) *Store {
	return &Store{
		tree: tree,
	}
}

// GetImmutable returns a reference to a new store backed by an immutable IAVL
// tree at a specific version (height) without any pruning options. This should
// be used for querying and iteration only. If the version does not exist or has
// been pruned, an empty immutable IAVL tree will be used.
// Any mutable operations executed will result in a panic.
func (st *Store) GetImmutable(version int64) (*Store, error) {
	var iTree *iavl.ImmutableTree
	var err error
	if !abci.GetDisableABCIQueryMutex() {
		if !st.VersionExists(version) {
			return &Store{tree: &immutableTree{&iavl.ImmutableTree{}}}, nil
		}

		iTree, err = st.tree.GetImmutable(version)
		if err != nil {
			return nil, err
		}
	} else {
		iTree, err = st.tree.GetImmutable(version)
		if err != nil {
			return &Store{tree: &immutableTree{&iavl.ImmutableTree{}}}, nil
		}
	}
	return &Store{
		tree: &immutableTree{iTree},
	}, nil
}

// Commit commits the current store state and returns a CommitID with the new
// version and hash.
func (st *Store) Commit(inDelta *iavl.TreeDelta, deltas []byte) (types.CommitID, iavl.TreeDelta, []byte) {
	flag := false
	if (tmtypes.EnableApplyP2PDelta() || tmtypes.EnableDownloadDelta()) && len(deltas) != 0 {
		flag = true
		st.tree.SetDelta(inDelta)
	}
	hash, version, delta, err := st.tree.SaveVersion(flag)
	if err != nil {
		panic(err)
	}

	// commit to flat kv db
	batch := st.flatKVDB.NewBatch()
	defer batch.Close()
	for key, value := range st.cache {
		batch.Set([]byte(key), value)
	}
	batch.Write()
	st.addFlatKVWriteCount()
	// clear cache
	st.cache = make(map[string][]byte)

	return types.CommitID{
		Version: version,
		Hash:    hash,
	}, delta, nil
}

// Implements Committer.
func (st *Store) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: st.tree.Version(),
		Hash:    st.tree.Hash(),
	}
}

// SetPruning panics as pruning options should be provided at initialization
// since IAVl accepts pruning options directly.
func (st *Store) SetPruning(_ types.PruningOptions) {
	panic("cannot set pruning options on an initialized IAVL store")
}

// VersionExists returns whether or not a given version is stored.
func (st *Store) VersionExists(version int64) bool {
	return st.tree.VersionExists(version)
}

// Implements Store.
func (st *Store) GetStoreType() types.StoreType {
	return types.StoreTypeIAVL
}

// Implements Store.
func (st *Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(st)
}

// CacheWrapWithTrace implements the Store interface.
func (st *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(st, w, tc))
}

// Implements types.KVStore.
func (st *Store) Set(key, value []byte) {
	types.AssertValidValue(value)
	st.tree.Set(key, value)
	strKey := string(key)
	st.addCache(strKey, value)
}

// Implements types.KVStore.
func (st *Store) Get(key []byte) []byte {
	strKey := string(key)
	if cacheVal, ok := st.getCache(strKey); ok {
		return cacheVal
	}
	ts := time.Now()
	value, err := st.flatKVDB.Get(key)
	st.addFlatKVReadTime(time.Now().Sub(ts).Nanoseconds())
	st.addFlatKVReadCount()
	if err == nil && len(value) != 0 {
		return value
	}
	_, value = st.tree.Get(key)
	if value != nil {
		st.addCache(strKey, value)
	}

	return value
}

// Implements types.KVStore.
func (st *Store) Has(key []byte) (exists bool) {
	strKey := string(key)
	if _, ok := st.getCache(strKey); ok {
		return true
	}
	st.addFlatKVReadCount()
	if ok, err := st.flatKVDB.Has(key); err == nil && ok {
		return true
	}

	return st.tree.Has(key)
}

// Implements types.KVStore.
func (st *Store) Delete(key []byte) {
	st.tree.Remove(key)
	st.flatKVDB.Delete(key)
	st.addFlatKVWriteCount()
	st.deleteCache(string(key))
}

// DeleteVersions deletes a series of versions from the MutableTree. An error
// is returned if any single version is invalid or the delete fails. All writes
// happen in a single batch with a single commit.
func (st *Store) DeleteVersions(versions ...int64) error {
	return st.tree.DeleteVersions(versions...)
}

// Implements types.KVStore.
func (st *Store) Iterator(start, end []byte) types.Iterator {
	var iTree *iavl.ImmutableTree

	switch tree := st.tree.(type) {
	case *immutableTree:
		iTree = tree.ImmutableTree
	case *iavl.MutableTree:
		iTree = tree.ImmutableTree
	}

	return newIAVLIterator(iTree, start, end, true)
}

// Implements types.KVStore.
func (st *Store) ReverseIterator(start, end []byte) types.Iterator {
	var iTree *iavl.ImmutableTree

	switch tree := st.tree.(type) {
	case *immutableTree:
		iTree = tree.ImmutableTree
	case *iavl.MutableTree:
		iTree = tree.ImmutableTree
	}

	return newIAVLIterator(iTree, start, end, false)
}

// Handle gatest the latest height, if height is 0
func getHeight(tree Tree, req abci.RequestQuery) int64 {
	height := req.Height
	if height == 0 {
		latest := tree.Version()
		if tree.VersionExists(latest - 1) {
			height = latest - 1
		} else {
			height = latest
		}
	}
	return height
}

// Query implements ABCI interface, allows queries
//
// by default we will return from (latest height -1),
// as we will have merkle proofs immediately (header height = data height + 1)
// If latest-1 is not present, use latest (which must be present)
// if you care to have the latest data to see a tx results, you must
// explicitly set the height you want to see
func (st *Store) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrTxDecode, "query cannot be zero length"))
	}

	tree := st.tree

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	res.Height = getHeight(tree, req)

	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes

		res.Key = key
		if !st.VersionExists(res.Height) {
			res.Log = iavl.ErrVersionDoesNotExist.Error()
			break
		}

		if req.Prove {
			value, proof, err := tree.GetVersionedWithProof(key, res.Height)
			if err != nil {
				res.Log = err.Error()
				break
			}
			if proof == nil {
				// Proof == nil implies that the store is empty.
				if value != nil {
					panic("unexpected value for an empty proof")
				}
			}
			if value != nil {
				// value was found
				res.Value = value
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewValueOp(key, proof).ProofOp()}}
			} else {
				// value wasn't found
				res.Value = nil
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{iavl.NewAbsenceOp(key, proof).ProofOp()}}
			}
		} else {
			_, res.Value = tree.GetVersioned(key, res.Height)
		}

	case "/subspace":
		var KVs []types.KVPair

		subspace := req.Data
		res.Key = subspace

		iterator := types.KVStorePrefixIterator(st, subspace)
		for ; iterator.Valid(); iterator.Next() {
			KVs = append(KVs, types.KVPair{Key: iterator.Key(), Value: iterator.Value()})
		}

		iterator.Close()
		res.Value = cdc.MustMarshalBinaryLengthPrefixed(KVs)

	default:
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unexpected query path: %v", req.Path))
	}

	return res
}

func (st *Store) GetDBReadTime() int {
	return st.tree.GetDBReadTime()
}

func (st *Store) GetDBWriteCount() int {
	return st.tree.GetDBWriteCount()
}

func (st *Store) GetDBReadCount() int {
	return st.tree.GetDBReadCount()
}

func (st *Store) GetNodeReadCount() int {
	return st.tree.GetNodeReadCount()
}

func (st *Store) ResetCount() {
	st.tree.ResetCount()
	st.resetFlatKVReadTime()
	st.resetFlatKVReadCount()
	st.resetFlatKVWriteCount()
}

//----------------------------------------

// Implements types.Iterator.
type iavlIterator struct {
	// Domain
	start, end []byte

	key   []byte // The current key (mutable)
	value []byte // The current value (mutable)

	// Underlying store
	tree *iavl.ImmutableTree

	// Channel to push iteration values.
	iterCh chan tmkv.Pair

	// Close this to release goroutine.
	quitCh chan struct{}

	// Close this to signal that state is initialized.
	initCh chan struct{}

	mtx sync.Mutex

	ascending bool // Iteration order

	invalid bool // True once, true forever (mutable)
}

var _ types.Iterator = (*iavlIterator)(nil)

// newIAVLIterator will create a new iavlIterator.
// CONTRACT: Caller must release the iavlIterator, as each one creates a new
// goroutine.
func newIAVLIterator(tree *iavl.ImmutableTree, start, end []byte, ascending bool) *iavlIterator {
	iter := &iavlIterator{
		tree:      tree,
		start:     types.Cp(start),
		end:       types.Cp(end),
		ascending: ascending,
		iterCh:    make(chan tmkv.Pair), // Set capacity > 0?
		quitCh:    make(chan struct{}),
		initCh:    make(chan struct{}),
	}
	go iter.iterateRoutine()
	go iter.initRoutine()
	return iter
}

// Run this to funnel items from the tree to iterCh.
func (iter *iavlIterator) iterateRoutine() {
	iter.tree.IterateRange(
		iter.start, iter.end, iter.ascending,
		func(key, value []byte) bool {
			select {
			case <-iter.quitCh:
				return true // done with iteration.
			case iter.iterCh <- tmkv.Pair{Key: key, Value: value}:
				return false // yay.
			}
		},
	)
	close(iter.iterCh) // done.
}

// Run this to fetch the first item.
func (iter *iavlIterator) initRoutine() {
	iter.receiveNext()
	close(iter.initCh)
}

// Implements types.Iterator.
func (iter *iavlIterator) Domain() (start, end []byte) {
	return iter.start, iter.end
}

// Implements types.Iterator.
func (iter *iavlIterator) Valid() bool {
	iter.waitInit()
	iter.mtx.Lock()

	validity := !iter.invalid
	iter.mtx.Unlock()
	return validity
}

// Implements types.Iterator.
func (iter *iavlIterator) Next() {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	iter.receiveNext()
	iter.mtx.Unlock()
}

// Implements types.Iterator.
func (iter *iavlIterator) Key() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	key := iter.key
	iter.mtx.Unlock()
	return key
}

// Implements types.Iterator.
func (iter *iavlIterator) Value() []byte {
	iter.waitInit()
	iter.mtx.Lock()
	iter.assertIsValid(true)

	val := iter.value
	iter.mtx.Unlock()
	return val
}

// Close closes the IAVL iterator by closing the quit channel and waiting for
// the iterCh to finish/close.
func (iter *iavlIterator) Close() {
	close(iter.quitCh)
	// wait iterCh to close
	for range iter.iterCh {
	}
}

// Error performs a no-op.
func (iter *iavlIterator) Error() error {
	return nil
}

//----------------------------------------

func (iter *iavlIterator) setNext(key, value []byte) {
	iter.assertIsValid(false)

	iter.key = key
	iter.value = value
}

func (iter *iavlIterator) setInvalid() {
	iter.assertIsValid(false)

	iter.invalid = true
}

func (iter *iavlIterator) waitInit() {
	<-iter.initCh
}

func (iter *iavlIterator) receiveNext() {
	kvPair, ok := <-iter.iterCh
	if ok {
		iter.setNext(kvPair.Key, kvPair.Value)
	} else {
		iter.setInvalid()
	}
}

// assertIsValid panics if the iterator is invalid. If unlockMutex is true,
// it also unlocks the mutex before panicing, to prevent deadlocks in code that
// recovers from panics
func (iter *iavlIterator) assertIsValid(unlockMutex bool) {
	if iter.invalid {
		if unlockMutex {
			iter.mtx.Unlock()
		}
		panic("invalid iterator")
	}
}

// SetInitialVersion sets the initial version of the IAVL tree. It is used when
// starting a new chain at an arbitrary height.
func (st *Store) SetInitialVersion(version int64) {
	st.tree.SetInitialVersion(uint64(version))
}

// Exports the IAVL store at the given version, returning an iavl.Exporter for the tree.
func (st *Store) Export(version int64) (*iavl.Exporter, error) {
	istore, err := st.GetImmutable(version)
	if err != nil {
		return nil, fmt.Errorf("iavl export failed for version %v: %w", version, err)
	}
	tree, ok := istore.tree.(*immutableTree)
	if !ok || tree == nil {
		return nil, fmt.Errorf("iavl export failed: unable to fetch tree for version %v", version)
	}
	return tree.Export(), nil
}

// Import imports an IAVL tree at the given version, returning an iavl.Importer for importing.
func (st *Store) Import(version int64) (*iavl.Importer, error) {
	tree, ok := st.tree.(*iavl.MutableTree)
	if !ok {
		return nil, errors.New("iavl import failed: unable to find mutable tree")
	}
	return tree.Import(version)
}
