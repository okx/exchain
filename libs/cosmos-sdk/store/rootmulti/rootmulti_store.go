package rootmulti

import (
	"fmt"

	sdkmaps "github.com/okex/exchain/libs/cosmos-sdk/store/internal/maps"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mem"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"

	"io"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/store/flatkv"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/spf13/viper"

	jsoniter "github.com/json-iterator/go"

	iavltree "github.com/okex/exchain/libs/iavl"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	//"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/pkg/errors"

	"github.com/okex/exchain/libs/cosmos-sdk/store/cachemulti"
	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/transient"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

var itjs = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	latestVersionKey      = "s/latest"
	pruneHeightsKey       = "s/pruneheights"
	versionsKey           = "s/versions"
	commitInfoKeyFmt      = "s/%d" // s/<version>
	maxPruneHeightsLength = 100
)

// Store is composed of many CommitStores. Name contrasts with
// cacheMultiStore which is for cache-wrapping other MultiStores. It implements
// the CommitMultiStore interface.
type Store struct {
	db             dbm.DB
	flatKVDB       dbm.DB
	lastCommitInfo commitInfo
	pruningOpts    types.PruningOptions
	storesParams   map[types.StoreKey]storeParams
	stores         map[types.StoreKey]types.CommitKVStore
	keysByName     map[string]types.StoreKey
	lazyLoading    bool
	pruneHeights   []int64
	versions       []int64

	traceWriter  io.Writer
	traceContext types.TraceContext

	interBlockCache types.MultiStorePersistentCache

	logger tmlog.Logger

	commitHeightFilterPipeline func(h int64) func(str string) bool
	pruneHeightFilterPipeline  func(h int64) func(str string) bool
	upgradeVersion             int64
}

var (
	_ types.CommitMultiStore = (*Store)(nil)
	_ types.Queryable        = (*Store)(nil)
	_ types.CommitMultiStore = (*Store)(nil)
)

// NewStore returns a reference to a new Store object with the provided DB. The
// store will be created with a PruneNothing pruning strategy by default. After
// a store is created, KVStores must be mounted and finally LoadLatestVersion or
// LoadVersion must be called.
func NewStore(db dbm.DB) *Store {
	var flatKVDB dbm.DB
	if viper.GetBool(flatkv.FlagEnable) {
		flatKVDB = newFlatKVDB()
	}
	ret := &Store{
		db:                         db,
		flatKVDB:                   flatKVDB,
		pruningOpts:                types.PruneNothing,
		storesParams:               make(map[types.StoreKey]storeParams),
		stores:                     make(map[types.StoreKey]types.CommitKVStore),
		keysByName:                 make(map[string]types.StoreKey),
		pruneHeights:               make([]int64, 0),
		versions:                   make([]int64, 0),
		commitHeightFilterPipeline: types.DefaultAcceptAll,
		pruneHeightFilterPipeline:  types.DefaultAcceptAll,
		upgradeVersion:             -1,
	}

	return ret
}

func newFlatKVDB() dbm.DB {
	rootDir := viper.GetString("home")
	dataDir := filepath.Join(rootDir, "data")
	var err error
	flatKVDB, err := sdk.NewLevelDB("flat", dataDir)
	if err != nil {
		panic(err)
	}
	return flatKVDB
}

// SetPruning sets the pruning strategy on the root store and all the sub-stores.
// Note, calling SetPruning on the root store prior to LoadVersion or
// LoadLatestVersion performs a no-op as the stores aren't mounted yet.
func (rs *Store) SetPruning(pruningOpts types.PruningOptions) {
	rs.pruningOpts = pruningOpts
}

// SetLazyLoading sets if the iavl store should be loaded lazily or not
func (rs *Store) SetLazyLoading(lazyLoading bool) {
	rs.lazyLoading = lazyLoading
}

// Implements Store.
func (rs *Store) GetStoreType() types.StoreType {
	return types.StoreTypeMulti
}

func (rs *Store) GetStores() map[types.StoreKey]types.CommitKVStore {
	return rs.stores
}

func (rs *Store) GetVersions() []int64 {
	return rs.versions
}

func (rs *Store) GetPruningHeights() []int64 {
	return rs.pruneHeights
}

// Implements CommitMultiStore.
func (rs *Store) MountStoreWithDB(key types.StoreKey, typ types.StoreType, db dbm.DB) {
	if key == nil {
		panic("MountIAVLStore() key cannot be nil")
	}
	if _, ok := rs.storesParams[key]; ok {
		panic(fmt.Sprintf("Store duplicate store key %v", key))
	}
	if _, ok := rs.keysByName[key.Name()]; ok {
		panic(fmt.Sprintf("Store duplicate store key name %v", key))
	}
	rs.storesParams[key] = storeParams{
		key: key,
		typ: typ,
		db:  db,
	}
	rs.keysByName[key.Name()] = key
}

// GetCommitStore returns a mounted CommitStore for a given StoreKey. If the
// store is wrapped in an inter-block cache, it will be unwrapped before returning.
func (rs *Store) GetCommitStore(key types.StoreKey) types.CommitStore {
	return rs.GetCommitKVStore(key)
}

// GetCommitKVStore returns a mounted CommitKVStore for a given StoreKey. If the
// store is wrapped in an inter-block cache, it will be unwrapped before returning.
func (rs *Store) GetCommitKVStore(key types.StoreKey) types.CommitKVStore {
	// If the Store has an inter-block cache, first attempt to lookup and unwrap
	// the underlying CommitKVStore by StoreKey. If it does not exist, fallback to
	// the main mapping of CommitKVStores.
	if rs.interBlockCache != nil {
		if store := rs.interBlockCache.Unwrap(key); store != nil {
			return store
		}
	}

	return rs.stores[key]
}

// LoadLatestVersionAndUpgrade implements CommitMultiStore
func (rs *Store) LoadLatestVersionAndUpgrade(upgrades *types.StoreUpgrades) error {
	ver := getLatestVersion(rs.db)
	return rs.loadVersion(ver, upgrades)
}

// LoadVersionAndUpgrade allows us to rename substores while loading an older version
func (rs *Store) LoadVersionAndUpgrade(ver int64, upgrades *types.StoreUpgrades) error {
	return rs.loadVersion(ver, upgrades)
}

// LoadLatestVersion implements CommitMultiStore.
func (rs *Store) LoadLatestVersion() error {
	ver := getLatestVersion(rs.db)
	return rs.loadVersion(ver, nil)
}

func (rs *Store) GetLatestVersion() int64 {
	return getLatestVersion(rs.db)
}

// LoadVersion implements CommitMultiStore.
func (rs *Store) LoadVersion(ver int64) error {
	return rs.loadVersion(ver, nil)
}

func (rs *Store) GetCommitVersion() (int64, error) {
	var minVersion int64 = 1<<63 - 1
	for _, storeParams := range rs.storesParams {
		if storeParams.typ != types.StoreTypeIAVL {
			continue
		}
		commitVersion, err := rs.getCommitVersionFromParams(storeParams)
		if err != nil {
			return 0, err
		}
		if commitVersion < minVersion {
			minVersion = commitVersion
		}
	}
	return minVersion, nil
}

func (rs *Store) loadVersion(ver int64, upgrades *types.StoreUpgrades) error {
	infos := make(map[string]storeInfo)
	var cInfo commitInfo
	cInfo.Version = tmtypes.GetStartBlockHeight()

	// load old data if we are not version 0
	if ver != 0 {
		var err error
		cInfo, err = getCommitInfo(rs.db, ver)
		if err != nil {
			return err
		}

		// convert StoreInfos slice to map
		for _, storeInfo := range cInfo.StoreInfos {
			infos[storeInfo.Name] = storeInfo
		}
	}

	roots := make(map[int64][]byte)
	// load each Store (note this doesn't panic on unmounted keys now)
	var newStores = make(map[types.StoreKey]types.CommitKVStore)
	for key, storeParams := range rs.storesParams {
		commitID := rs.getCommitID(infos, key.Name())

		// If it has been added, set the initial version
		if upgrades.IsAdded(key.Name()) {
			storeParams.initialVersion = uint64(ver) + 1
		}

		// Load it
		store, err := rs.loadCommitStoreFromParams(key, commitID, storeParams)
		if err != nil {
			return fmt.Errorf("failed to load Store: %v", err)
		}
		newStores[key] = store

		if storeParams.typ == types.StoreTypeIAVL {
			if len(roots) == 0 {
				iStore := store.(*iavl.Store)
				roots = iStore.GetHeights()
			}
		}

		// If it was deleted, remove all data
		if upgrades.IsDeleted(key.Name()) {
			if err := deleteKVStore(store.(types.KVStore)); err != nil {
				return fmt.Errorf("failed to delete store %s: %v", key.Name(), err)
			}
		} else if oldName := upgrades.RenamedFrom(key.Name()); oldName != "" {
			// handle renames specially
			// make an unregistered key to satify loadCommitStore params
			oldKey := types.NewKVStoreKey(oldName)
			oldParams := storeParams
			oldParams.key = oldKey

			// load from the old name
			oldStore, err := rs.loadCommitStoreFromParams(oldKey, rs.getCommitID(infos, oldName), oldParams)
			if err != nil {
				return fmt.Errorf("failed to load old Store '%s': %v", oldName, err)
			}

			// move all data
			if err := moveKVStoreData(oldStore.(types.KVStore), store.(types.KVStore)); err != nil {
				return fmt.Errorf("failed to move store %s -> %s: %v", oldName, key.Name(), err)
			}
		}
	}

	rs.lastCommitInfo = cInfo
	rs.stores = newStores

	err := rs.checkAndResetPruningHeights(roots)
	if err != nil {
		return err
	}

	vs, err := getVersions(rs.db)
	if err != nil {
		return err
	}
	if len(vs) > 0 {
		rs.versions = vs
	}
	if rs.logger != nil {
		rs.logger.Info("loadVersion info", "pruned heights length", len(rs.pruneHeights), "versions", len(rs.versions))
	}
	if len(rs.pruneHeights) > maxPruneHeightsLength {
		return fmt.Errorf("Pruned heights length <%d> exceeds <%d>, "+
			"need to prune them with command "+
			"<exchaind data prune-compact all --home your_exchaind_home_directory> before running exchaind",
			len(rs.pruneHeights), maxPruneHeightsLength)
	}
	return nil
}

func (rs *Store) checkAndResetPruningHeights(roots map[int64][]byte) error {
	ph, err := getPruningHeights(rs.db, false)
	if err != nil {
		return err
	}

	if len(ph) == 0 {
		return nil
	}

	needReset := false
	var newPh []int64
	for _, h := range ph {
		if _, ok := roots[h]; ok {
			newPh = append(newPh, h)
		} else {
			needReset = true
		}
	}
	rs.pruneHeights = newPh

	if needReset {
		if rs.logger != nil {
			msg := fmt.Sprintf("Detected pruned heights length <%d>, reset to <%d>",
				len(ph), len(rs.pruneHeights))
			rs.logger.Info(msg)
		}
		batch := rs.db.NewBatch()
		setPruningHeights(batch, newPh)
		batch.Write()
		batch.Close()
	}

	return nil
}
func (rs *Store) getCommitID(infos map[string]storeInfo, name string) types.CommitID {
	info, ok := infos[name]
	if !ok {
		return types.CommitID{Version: tmtypes.GetStartBlockHeight()}
	}
	return info.Core.CommitID
}

func deleteKVStore(kv types.KVStore) error {
	// Note that we cannot write while iterating, so load all keys here, delete below
	var keys [][]byte
	itr := kv.Iterator(nil, nil)
	for itr.Valid() {
		keys = append(keys, itr.Key())
		itr.Next()
	}
	itr.Close()

	for _, k := range keys {
		kv.Delete(k)
	}
	return nil
}

// we simulate move by a copy and delete
func moveKVStoreData(oldDB types.KVStore, newDB types.KVStore) error {
	// we read from one and write to another
	itr := oldDB.Iterator(nil, nil)
	for itr.Valid() {
		newDB.Set(itr.Key(), itr.Value())
		itr.Next()
	}
	itr.Close()

	// then delete the old store
	return deleteKVStore(oldDB)
}

// SetInterBlockCache sets the Store's internal inter-block (persistent) cache.
// When this is defined, all CommitKVStores will be wrapped with their respective
// inter-block cache.
func (rs *Store) SetInterBlockCache(c types.MultiStorePersistentCache) {
	rs.interBlockCache = c
}

// SetTracer sets the tracer for the MultiStore that the underlying
// stores will utilize to trace operations. A MultiStore is returned.
func (rs *Store) SetTracer(w io.Writer) types.MultiStore {
	rs.traceWriter = w
	return rs
}

// SetTracingContext updates the tracing context for the MultiStore by merging
// the given context with the existing context by key. Any existing keys will
// be overwritten. It is implied that the caller should update the context when
// necessary between tracing operations. It returns a modified MultiStore.
func (rs *Store) SetTracingContext(tc types.TraceContext) types.MultiStore {
	if rs.traceContext != nil {
		for k, v := range tc {
			rs.traceContext[k] = v
		}
	} else {
		rs.traceContext = tc
	}

	return rs
}

// TracingEnabled returns if tracing is enabled for the MultiStore.
func (rs *Store) TracingEnabled() bool {
	return rs.traceWriter != nil
}

//----------------------------------------
// +CommitStore

// Implements Committer/CommitStore.
func (rs *Store) LastCommitID() types.CommitID {
	return rs.lastCommitInfo.CommitID()
}

func (rs *Store) CommitterCommit(*iavltree.TreeDelta) (_ types.CommitID, _ *iavltree.TreeDelta) {
	return
}

// Implements Committer/CommitStore.
func (rs *Store) CommitterCommitMap(inputDeltaMap iavltree.TreeDeltaMap) (types.CommitID, iavltree.TreeDeltaMap) {
	previousHeight := rs.lastCommitInfo.Version
	version := previousHeight + 1

	var outputDeltaMap iavltree.TreeDeltaMap
	rs.lastCommitInfo, outputDeltaMap = commitStores(version, rs.stores, inputDeltaMap, rs.commitHeightFilterPipeline(version))

	if !iavltree.EnableAsyncCommit {
		// Determine if pruneHeight height needs to be added to the list of heights to
		// be pruned, where pruneHeight = (commitHeight - 1) - KeepRecent.
		if int64(rs.pruningOpts.KeepRecent) < previousHeight {
			pruneHeight := previousHeight - int64(rs.pruningOpts.KeepRecent)
			// We consider this height to be pruned iff:
			//
			// - KeepEvery is zero as that means that all heights should be pruned.
			// - KeepEvery % (height - KeepRecent) != 0 as that means the height is not
			// a 'snapshot' height.
			if rs.pruningOpts.KeepEvery == 0 || pruneHeight%int64(rs.pruningOpts.KeepEvery) != 0 {
				rs.pruneHeights = append(rs.pruneHeights, pruneHeight)
				for k, v := range rs.versions {
					if v == pruneHeight {
						rs.versions = append(rs.versions[:k], rs.versions[k+1:]...)
						break
					}
				}
			}
		}

		if uint64(len(rs.versions)) > rs.pruningOpts.MaxRetainNum {
			rs.pruneHeights = append(rs.pruneHeights, rs.versions[:uint64(len(rs.versions))-rs.pruningOpts.MaxRetainNum]...)
			rs.versions = rs.versions[uint64(len(rs.versions))-rs.pruningOpts.MaxRetainNum:]
		}

		// batch prune if the current height is a pruning interval height
		if rs.pruningOpts.Interval > 0 && version%int64(rs.pruningOpts.Interval) == 0 {
			rs.pruneStores()
		}

		rs.versions = append(rs.versions, version)
	}
	flushMetadata(rs.db, version, rs.lastCommitInfo, rs.pruneHeights, rs.versions)

	return types.CommitID{
		Version: version,
		Hash:    rs.lastCommitInfo.Hash(),
	}, outputDeltaMap
}

// pruneStores will batch delete a list of heights from each mounted sub-store.
// Afterwards, pruneHeights is reset.
func (rs *Store) pruneStores() {
	pruneCnt := len(rs.pruneHeights)
	if pruneCnt == 0 {
		return
	}

	if rs.logger != nil {
		rs.logger.Info("pruning start", "pruning-count", pruneCnt, "curr-height", rs.lastCommitInfo.Version+1)
		rs.logger.Debug("pruning", "pruning-heights", rs.pruneHeights)
	}
	defer func() {
		if rs.logger != nil {
			rs.logger.Info("pruning end")
		}
	}()
	stores := rs.getFilterStores(rs.lastCommitInfo.Version + 1)
	//stores = rs.stores
	for key, store := range stores {
		if store.GetStoreType() == types.StoreTypeIAVL {
			// If the store is wrapped with an inter-block cache, we must first unwrap
			// it to get the underlying IAVL store.
			store = rs.GetCommitKVStore(key)
			if err := store.(*iavl.Store).DeleteVersions(rs.pruneHeights...); err != nil {
				if errCause := errors.Cause(err); errCause != nil && errCause != iavltree.ErrVersionDoesNotExist {
					panic(err)
				}
			}
		}
	}

	rs.pruneHeights = make([]int64, 0)
}

func (rs *Store) FlushPruneHeights(pruneHeights []int64, versions []int64) {
	flushMetadata(rs.db, rs.lastCommitInfo.Version, rs.lastCommitInfo, pruneHeights, versions)
}

// Implements CacheWrapper/Store/CommitStore.
func (rs *Store) CacheWrap() types.CacheWrap {
	return rs.CacheMultiStore().(types.CacheWrap)
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (rs *Store) CacheWrapWithTrace(_ io.Writer, _ types.TraceContext) types.CacheWrap {
	return rs.CacheWrap()
}

//----------------------------------------
// +MultiStore

// CacheMultiStore cache-wraps the multi-store and returns a CacheMultiStore.
// It implements the MultiStore interface.
func (rs *Store) CacheMultiStore() types.CacheMultiStore {
	stores := make(map[types.StoreKey]types.CacheWrapper)
	for k, v := range rs.stores {
		stores[k] = v
	}

	return cachemulti.NewStore(rs.db, stores, rs.keysByName, rs.traceWriter, rs.traceContext)
}

// CacheMultiStoreWithVersion is analogous to CacheMultiStore except that it
// attempts to load stores at a given version (height). An error is returned if
// any store cannot be loaded. This should only be used for querying and
// iterating at past heights.
func (rs *Store) CacheMultiStoreWithVersion(version int64) (types.CacheMultiStore, error) {
	cachedStores := make(map[types.StoreKey]types.CacheWrapper)
	for key, store := range rs.stores {
		switch store.GetStoreType() {
		case types.StoreTypeIAVL:
			// If the store is wrapped with an inter-block cache, we must first unwrap
			// it to get the underlying IAVL store.
			store = rs.GetCommitKVStore(key)

			// Attempt to lazy-load an already saved IAVL store version. If the
			// version does not exist or is pruned, an error should be returned.
			iavlStore, err := store.(*iavl.Store).GetImmutable(version)
			if err != nil {
				return nil, err
			}

			cachedStores[key] = iavlStore

		default:
			cachedStores[key] = store
		}
	}

	return cachemulti.NewStore(rs.db, cachedStores, rs.keysByName, rs.traceWriter, rs.traceContext), nil
}

// GetStore returns a mounted Store for a given StoreKey. If the StoreKey does
// not exist, it will panic. If the Store is wrapped in an inter-block cache, it
// will be unwrapped prior to being returned.
//
// TODO: This isn't used directly upstream. Consider returning the Store as-is
// instead of unwrapping.
func (rs *Store) GetStore(key types.StoreKey) types.Store {
	store := rs.GetCommitKVStore(key)
	if store == nil {
		panic(fmt.Sprintf("store does not exist for key: %s", key.Name()))
	}

	return store
}

// GetKVStore returns a mounted KVStore for a given StoreKey. If tracing is
// enabled on the KVStore, a wrapped TraceKVStore will be returned with the root
// store's tracer, otherwise, the original KVStore will be returned.
//
// NOTE: The returned KVStore may be wrapped in an inter-block cache if it is
// set on the root store.
func (rs *Store) GetKVStore(key types.StoreKey) types.KVStore {
	store := rs.stores[key].(types.KVStore)

	if rs.TracingEnabled() {
		store = tracekv.NewStore(store, rs.traceWriter, rs.traceContext)
	}

	return store
}

// getStoreByName performs a lookup of a StoreKey given a store name typically
// provided in a path. The StoreKey is then used to perform a lookup and return
// a Store. If the Store is wrapped in an inter-block cache, it will be unwrapped
// prior to being returned. If the StoreKey does not exist, nil is returned.
func (rs *Store) getStoreByName(name string) types.Store {
	key := rs.keysByName[name]
	if key == nil {
		return nil
	}

	return rs.GetCommitKVStore(key)
}

//---------------------- Query ------------------

// Query calls substore.Query with the same `req` where `req.Path` is
// modified to remove the substore prefix.
// Ie. `req.Path` here is `/<substore>/<path>`, and trimmed to `/<path>` for the substore.
// TODO: add proof for `multistore -> substore`.
func (rs *Store) Query(req abci.RequestQuery) abci.ResponseQuery {
	path := req.Path
	str := string(req.Data)
	storeName, subpath, err := parsePath(path)
	if err != nil {
		return sdkerrors.QueryResult(err)
	}

	store := rs.getStoreByName(storeName)
	if store == nil {
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "no such store: %s", storeName))
	}

	queryable, ok := store.(types.Queryable)
	if !ok {
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "store %s (type %T) doesn't support queries", storeName, store))
	}

	// trim the path and make the query
	req.Path = subpath
	res := queryable.Query(req)
	if strings.Contains(str, "connections") {
		defer func() {
			if res.Proof == nil || len(res.Proof.Ops) == 0 {
				panic("asd")
			}
		}()
	}

	if !req.Prove || !RequireProof(subpath) {
		return res
	}

	if res.Proof == nil || len(res.Proof.Ops) == 0 {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proof is unexpectedly empty; ensure height has not been pruned"))
	}

	// If the request's height is the latest height we've committed, then utilize
	// the store's lastCommitInfo as this commit info may not be flushed to disk.
	// Otherwise, we query for the commit info from disk.
	var commitInfo commitInfo

	if res.Height == rs.lastCommitInfo.Version {
		commitInfo = rs.lastCommitInfo
	} else {
		commitInfo, err = getCommitInfo(rs.db, res.Height)
		if err != nil {
			return sdkerrors.QueryResult(err)
		}
	}

	if tmtypes.HigherThanVenus1(req.Height) {
		queryIbcProof(&res, &commitInfo, storeName)
	} else {
		// Restore origin path and append proof op.
		res.Proof.Ops = append(res.Proof.Ops, NewMultiStoreProofOp(
			[]byte(storeName),
			NewMultiStoreProof(commitInfo.StoreInfos),
		).ProofOp())
	}

	// TODO: handle in another TM v0.26 update PR
	// res.Proof = buildMultiStoreProof(res.Proof, storeName, commitInfo.StoreInfos)
	return res
}

// parsePath expects a format like /<storeName>[/<subpath>]
// Must start with /, subpath may be empty
// Returns error if it doesn't start with /
func parsePath(path string) (storeName string, subpath string, err error) {
	if !strings.HasPrefix(path, "/") {
		return storeName, subpath, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid path: %s", path)
	}

	paths := strings.SplitN(path[1:], "/", 2)
	storeName = paths[0]

	if len(paths) == 2 {
		subpath = "/" + paths[1]
	}

	return storeName, subpath, nil
}

func (rs *Store) loadCommitStoreFromParams(key types.StoreKey, id types.CommitID, params storeParams) (types.CommitKVStore, error) {
	var db dbm.DB

	if params.db != nil {
		db = dbm.NewPrefixDB(params.db, []byte("s/_/"))
	} else {
		prefix := "s/k:" + params.key.Name() + "/"
		db = dbm.NewPrefixDB(rs.db, []byte(prefix))
	}

	switch params.typ {
	case types.StoreTypeMulti:
		panic("recursive MultiStores not yet supported")

	case types.StoreTypeIAVL:
		var store types.CommitKVStore
		var err error
		prefix := "s/k:" + params.key.Name() + "/"
		var prefixDB dbm.DB
		if rs.flatKVDB != nil {
			prefixDB = dbm.NewPrefixDB(rs.flatKVDB, []byte(prefix))
		}
		if params.initialVersion == 0 {
			store, err = iavl.LoadStore(db, prefixDB, id, rs.lazyLoading, tmtypes.GetStartBlockHeight())
		} else {
			store, err = iavl.LoadStoreWithInitialVersion(db, prefixDB, id, rs.lazyLoading, params.initialVersion)
		}

		if err != nil {
			return nil, err
		}

		if rs.interBlockCache != nil {
			// Wrap and get a CommitKVStore with inter-block caching. Note, this should
			// only wrap the primary CommitKVStore, not any store that is already
			// cache-wrapped as that will create unexpected behavior.
			store = rs.interBlockCache.GetStoreCache(key, store)
		}

		return store, err

	case types.StoreTypeDB:
		return commitDBStoreAdapter{Store: dbadapter.Store{DB: db}}, nil

	case types.StoreTypeTransient:
		_, ok := key.(*types.TransientStoreKey)
		if !ok {
			return nil, fmt.Errorf("invalid StoreKey for StoreTypeTransient: %s", key.String())
		}

		return transient.NewStore(), nil
	case types.StoreTypeMemory:
		if _, ok := key.(*types.MemoryStoreKey); !ok {
			return nil, fmt.Errorf("unexpected key type for a MemoryStoreKey; got: %s", key.String())
		}

		return mem.NewStore(), nil

	default:
		panic(fmt.Sprintf("unrecognized store type %v", params.typ))
	}
}
func (rs *Store) GetDBReadTime() int {
	count := 0
	for _, store := range rs.stores {
		count += store.GetDBReadTime()
	}
	return count
}

func (rs *Store) getCommitVersionFromParams(params storeParams) (int64, error) {
	var db dbm.DB

	if params.db != nil {
		db = dbm.NewPrefixDB(params.db, []byte("s/_/"))
	} else {
		prefix := "s/k:" + params.key.Name() + "/"
		db = dbm.NewPrefixDB(rs.db, []byte(prefix))
	}

	return iavl.GetCommitVersion(db)
}

func (rs *Store) GetDBWriteCount() int {
	count := 0
	for _, store := range rs.stores {
		count += store.GetDBWriteCount()
	}
	return count
}

func (rs *Store) GetDBReadCount() int {
	count := 0
	for _, store := range rs.stores {
		count += store.GetDBReadCount()
	}
	return count
}

func (rs *Store) GetNodeReadCount() int {
	count := 0
	for _, store := range rs.stores {
		count += store.GetNodeReadCount()
	}
	return count
}

func (rs *Store) ResetCount() {
	for _, store := range rs.stores {
		store.ResetCount()
	}
}

func (rs *Store) GetFlatKVReadTime() int {
	rt := 0
	for _, store := range rs.stores {
		rt += store.GetFlatKVReadTime()
	}
	return rt
}

func (rs *Store) GetFlatKVWriteTime() int {
	wt := 0
	for _, store := range rs.stores {
		wt += store.GetFlatKVWriteTime()
	}
	return wt
}

func (rs *Store) GetFlatKVReadCount() int {
	count := 0
	for _, store := range rs.stores {
		count += store.GetFlatKVReadCount()
	}
	return count
}

func (rs *Store) GetFlatKVWriteCount() int {
	count := 0
	for _, store := range rs.stores {
		count += store.GetFlatKVWriteCount()
	}
	return count
}

//----------------------------------------
// storeParams

type storeParams struct {
	key            types.StoreKey
	db             dbm.DB
	typ            types.StoreType
	initialVersion uint64
}

//----------------------------------------
// commitInfo

// NOTE: Keep commitInfo a simple immutable struct.
type commitInfo struct {

	// Version
	Version int64

	// Store info for
	StoreInfos []storeInfo
}

// Hash returns the simple merkle root hash of the stores sorted by name.
func (ci commitInfo) Hash() []byte {
	if tmtypes.HigherThanVenus1(ci.Version) {
		return ci.ibcHash()
	}
	return ci.originHash()
}

func (ci commitInfo) originHash() []byte {
	// TODO: cache to ci.hash []byte
	m := make(map[string][]byte, len(ci.StoreInfos))
	for _, storeInfo := range ci.StoreInfos {
		m[storeInfo.Name] = storeInfo.Hash()
	}
	return merkle.SimpleHashFromMap(m)
}

// Hash returns the simple merkle root hash of the stores sorted by name.
func (ci commitInfo) ibcHash() []byte {
	m := ci.toMap()
	rootHash, _, _ := sdkmaps.ProofsFromMap(m)
	return rootHash
}

func (ci commitInfo) CommitID() types.CommitID {
	return types.CommitID{
		Version: ci.Version,
		Hash:    ci.Hash(),
	}
}

//----------------------------------------
// storeInfo

// storeInfo contains the name and core reference for an
// underlying store.  It is the leaf of the Stores top
// level simple merkle tree.
type storeInfo struct {
	Name string
	Core storeCore
}

type storeCore struct {
	// StoreType StoreType
	CommitID types.CommitID
	// ... maybe add more state
}

// Implements merkle.Hasher.
func (si storeInfo) Hash() []byte {
	// Doesn't write Name, since merkle.SimpleHashFromMap() will
	// include them via the keys.
	bz := si.Core.CommitID.Hash
	hasher := tmhash.New()

	_, err := hasher.Write(bz)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}

	return hasher.Sum(nil)
}

//----------------------------------------
// Misc.

func getLatestVersion(db dbm.DB) int64 {
	var latest int64
	latestBytes, err := db.Get([]byte(latestVersionKey))
	if err != nil {
		panic(err)
	} else if latestBytes == nil {
		return 0
	}

	err = cdc.UnmarshalBinaryLengthPrefixed(latestBytes, &latest)
	if err != nil {
		panic(err)
	}

	return latest
}

type StoreSorts []StoreSort

func (s StoreSorts) Len() int {
	return len(s)
}

func (s StoreSorts) Less(i, j int) bool {
	return s[i].key.Name() < s[j].key.Name()
}

func (s StoreSorts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type StoreSort struct {
	key types.StoreKey
	v   types.CommitKVStore
}

// Commits each store and returns a new commitInfo.
func commitStores(version int64, storeMap map[types.StoreKey]types.CommitKVStore,
	inputDeltaMap iavltree.TreeDeltaMap, f func(str string) bool) (commitInfo, iavltree.TreeDeltaMap) {
	var storeInfos []storeInfo
	outputDeltaMap := iavltree.TreeDeltaMap{}

	for key, store := range storeMap {
		if f(key.Name()) {
			continue
		}
		if tmtypes.GetVenus1Height()+1 == version {
			//init store tree version with block height
			store.UpgradeVersion(version)
		}

		commitID, outputDelta := store.CommitterCommit(inputDeltaMap[key.Name()]) // CommitterCommit
		if store.GetStoreType() == types.StoreTypeTransient {
			continue
		}

		si := storeInfo{}
		si.Name = key.Name()
		si.Core.CommitID = commitID
		si.Core.CommitID.Version = version
		storeInfos = append(storeInfos, si)
		outputDeltaMap[key.Name()] = outputDelta
	}

	return commitInfo{
		Version:    version,
		StoreInfos: storeInfos,
	}, outputDeltaMap
}

// Gets commitInfo from disk.
func getCommitInfo(db dbm.DB, ver int64) (commitInfo, error) {
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, ver)

	cInfoBytes, err := db.Get([]byte(cInfoKey))
	if err != nil {
		return commitInfo{}, fmt.Errorf("failed to get commit info: %v", err)
	} else if cInfoBytes == nil {
		return commitInfo{}, fmt.Errorf("failed to get commit info: no data")
	}

	var cInfo commitInfo

	err = cdc.UnmarshalBinaryLengthPrefixed(cInfoBytes, &cInfo)
	if err != nil {
		return commitInfo{}, fmt.Errorf("failed to get Store: %v", err)
	}

	return cInfo, nil
}

func setCommitInfo(batch dbm.Batch, version int64, cInfo commitInfo) {
	cInfoBytes := cdc.MustMarshalBinaryLengthPrefixed(cInfo)
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, version)
	batch.Set([]byte(cInfoKey), cInfoBytes)
}

func setLatestVersion(batch dbm.Batch, version int64) {
	latestBytes := cdc.MustMarshalBinaryLengthPrefixed(version)
	batch.Set([]byte(latestVersionKey), latestBytes)
}

func setPruningHeights(batch dbm.Batch, pruneHeights []int64) {
	bz := cdc.MustMarshalBinaryBare(pruneHeights)
	batch.Set([]byte(pruneHeightsKey), bz)
}

func SetPruningHeights(db dbm.DB, pruneHeights []int64) {
	batch := db.NewBatch()
	setPruningHeights(batch, pruneHeights)
	batch.Write()
	batch.Close()
}

func GetPruningHeights(db dbm.DB) ([]int64, error) {
	return getPruningHeights(db, true)
}

func getPruningHeights(db dbm.DB, reportZeroLengthErr bool) ([]int64, error) {
	bz, err := db.Get([]byte(pruneHeightsKey))
	if err != nil {
		return nil, fmt.Errorf("failed to get pruned heights: %w", err)
	}
	if len(bz) == 0 {
		if reportZeroLengthErr {
			return nil, errors.New("no pruned heights found")
		} else {
			return nil, nil
		}
	}

	var prunedHeights []int64
	if err := cdc.UnmarshalBinaryBare(bz, &prunedHeights); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pruned heights: %w", err)
	}

	return prunedHeights, nil
}

func flushMetadata(db dbm.DB, version int64, cInfo commitInfo, pruneHeights []int64, versions []int64) {
	batch := db.NewBatch()
	defer batch.Close()

	setCommitInfo(batch, version, cInfo)
	setLatestVersion(batch, version)
	setPruningHeights(batch, pruneHeights)
	setVersions(batch, versions)

	if err := batch.Write(); err != nil {
		panic(fmt.Errorf("error on batch write %w", err))
	}
}

func setVersions(batch dbm.Batch, versions []int64) {
	bz := cdc.MustMarshalBinaryBare(versions)
	batch.Set([]byte(versionsKey), bz)
}

func getVersions(db dbm.DB) ([]int64, error) {
	bz, err := db.Get([]byte(versionsKey))
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	if len(bz) == 0 {
		return nil, nil
	}

	var versions []int64
	if err := cdc.UnmarshalBinaryBare(bz, &versions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pruned heights: %w", err)
	}

	return versions, nil
}

// Snapshot implements snapshottypes.Snapshotter. The snapshot output for a given format must be
// identical across nodes such that chunks from different sources fit together. If the output for a
// given format changes (at the byte level), the snapshot format must be bumped - see
// TestMultistoreSnapshot_Checksum test.
func (rs *Store) Export(to *Store, initVersion int64) error {
	curVersion := rs.lastCommitInfo.Version
	// Collect stores to snapshot (only IAVL stores are supported)
	type namedStore struct {
		fromStore *iavl.Store
		toStore   *iavl.Store
		name      string
	}
	stores := []namedStore{}
	for key := range rs.stores {
		switch store := rs.GetCommitKVStore(key).(type) {
		case *iavl.Store:
			var toKVStore types.CommitKVStore
			for toKey, toValue := range to.stores {
				if key.Name() == toKey.Name() {
					toKVStore = toValue
				}
			}
			toStore, _ := toKVStore.(*iavl.Store)
			stores = append(stores, namedStore{name: key.Name(), fromStore: store, toStore: toStore})
		case *transient.Store:
			// Non-persisted stores shouldn't be snapshotted
			continue
		default:
			return fmt.Errorf(
				"don't know how to snapshot store %q of type %T", key.Name(), store)
		}
	}
	sort.Slice(stores, func(i, j int) bool {
		return strings.Compare(stores[i].name, stores[j].name) == 1
	})

	// Export each IAVL store. Stores are serialized as a stream of SnapshotItem Protobuf
	// messages. The first item contains a SnapshotStore with store metadata (i.e. name),
	// and the following messages contain a SnapshotNode (i.e. an ExportNode). Store changes
	// are demarcated by new SnapshotStore items.
	for _, store := range stores {
		log.Println("--------- export ", store.name, " start ---------")
		exporter, err := store.fromStore.Export(curVersion)
		if err != nil {
			panic(err)
		}
		defer exporter.Close()

		importer, err := store.toStore.Import(initVersion)
		if err != nil {
			panic(err)
		}
		defer importer.Close()

		var totalCnt uint64
		var totalSize uint64
		for {
			node, err := exporter.Next()
			if err == iavltree.ExportDone {
				break
			}

			err = importer.Add(node)
			if err != nil {
				panic(err)
			}
			nodeSize := len(node.Key) + len(node.Value)
			totalCnt++
			totalSize += uint64(nodeSize)
			if totalCnt%10000 == 0 {
				log.Println("--------- total node count ", totalCnt, " ---------")
				log.Println("--------- total node size ", totalSize, " ---------")
			}
		}

		exporter.Close()
		err = importer.Commit()
		if err != nil {
			panic(err)
		}
		importer.Close()
		log.Println("--------- export ", store.name, " end ---------")
	}

	flushMetadata(to.db, initVersion, rs.buildCommitInfo(initVersion), []int64{}, []int64{})

	return nil
}

func (rs *Store) buildCommitInfo(version int64) commitInfo {
	storeInfos := []storeInfo{}
	for key, store := range rs.stores {
		if store.GetStoreType() == types.StoreTypeTransient {
			continue
		}
		storeInfos = append(storeInfos, storeInfo{
			Name: key.Name(),
			Core: storeCore{
				store.LastCommitID(),
			},
		})
	}
	return commitInfo{
		Version:    version,
		StoreInfos: storeInfos,
	}
}

func (src Store) Copy() *Store {
	dst := &Store{
		db:           src.db,
		pruningOpts:  src.pruningOpts,
		storesParams: make(map[types.StoreKey]storeParams, len(src.storesParams)),
		stores:       make(map[types.StoreKey]types.CommitKVStore, len(src.stores)),
		keysByName:   make(map[string]types.StoreKey, len(src.keysByName)),
		lazyLoading:  src.lazyLoading,
		pruneHeights: make([]int64, 0),
		versions:     make([]int64, 0),

		traceWriter:     src.traceWriter,
		traceContext:    src.traceContext,
		interBlockCache: src.interBlockCache,
		upgradeVersion:  src.upgradeVersion,
	}

	dst.lastCommitInfo = commitInfo{
		Version:    src.lastCommitInfo.Version,
		StoreInfos: make([]storeInfo, 0),
	}

	for _, info := range src.lastCommitInfo.StoreInfos {
		dst.lastCommitInfo.StoreInfos = append(dst.lastCommitInfo.StoreInfos, info)
	}

	for key, value := range src.storesParams {
		dst.storesParams[key] = value
	}

	for key, value := range src.stores {
		dst.stores[key] = value
	}

	for key, value := range src.keysByName {
		dst.keysByName[key] = value
	}

	for _, value := range src.pruneHeights {
		dst.pruneHeights = append(dst.pruneHeights, value)
	}

	for _, value := range src.versions {
		dst.versions = append(dst.versions, value)
	}

	return dst
}

func (rs *Store) StopStore() {
	for _, store := range rs.stores {
		switch store.GetStoreType() {
		case types.StoreTypeIAVL:
			s := store.(*iavl.Store)
			s.StopStore()
		case types.StoreTypeDB:
			panic("unexpected db store")
		case types.StoreTypeMulti:
			panic("unexpected multi store")
		case types.StoreTypeTransient:
		default:
		}
	}

}

func (rs *Store) SetLogger(log tmlog.Logger) {
	rs.logger = log.With("module", "root-multi")
}

func (rs *Store) UpgradeVersion(version int64) {

	rs.upgradeVersion = version
}
