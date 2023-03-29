package mpt

import (
	"encoding/hex"
	"fmt"
	"github.com/okx/okbchain/libs/system/trace/persist"
	"io"
	"sync"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/cachekv"
	mpttype "github.com/okx/okbchain/libs/cosmos-sdk/store/mpt/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/iavl"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/libs/tendermint/crypto/merkle"
	tmlog "github.com/okx/okbchain/libs/tendermint/libs/log"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
)

const (
	FlagTrieAccStoreCache = "trie.account-store-cache"
)

var (
	TrieAccStoreCache uint = 32 // MB

	AccountStateRootRetriever StateRootRetriever = EmptyStateRootRetriever{}

	applyDelta   = false
	produceDelta = false
)

func SetApplyDelta(val bool) {
	applyDelta = val
}

func SetProduceDelta(val bool) {
	produceDelta = val
}

var cdc = codec.New()

var (
	_ types.KVStore       = (*MptStore)(nil)
	_ types.CommitStore   = (*MptStore)(nil)
	_ types.CommitKVStore = (*MptStore)(nil)
	_ types.Queryable     = (*MptStore)(nil)
)

// MptStore Implements types.KVStore and CommitKVStore.
// Its main purpose is to own the same interface as iavl store in libs/cosmos-sdk/store/iavl/iavl_store.go
type MptStore struct {
	trie                ethstate.Trie
	storageTrieForWrite map[ethcmn.Address]ethstate.Trie
	db                  ethstate.Database
	triegc              *prque.Prque
	logger              tmlog.Logger

	prefetcher   *TriePrefetcher
	originalRoot ethcmn.Hash
	exitSignal   chan struct{}

	version      int64
	startVersion int64
	cmLock       sync.Mutex

	retriever StateRootRetriever

	snaps         *snapshot.Tree
	snap          snapshot.Snapshot
	snapDestructs map[ethcmn.Hash]struct{}
	snapAccounts  map[ethcmn.Hash][]byte
	snapStorage   map[ethcmn.Hash]map[ethcmn.Hash][]byte
	snapRWLock    sync.RWMutex

	// for time statistics
	statisticsBeginTime time.Time
}

func (ms *MptStore) CommitterCommitMap(deltaMap iavl.TreeDeltaMap) (_ types.CommitID, _ iavl.TreeDeltaMap) {
	return
}

func (ms *MptStore) GetFlatKVReadTime() int {
	return 0
}

func (ms *MptStore) GetFlatKVWriteTime() int {
	return 0
}

func (ms *MptStore) GetFlatKVReadCount() int {
	return 0
}

func (ms *MptStore) GetFlatKVWriteCount() int {
	return 0
}

func NewMptStore(logger tmlog.Logger, id types.CommitID) (*MptStore, error) {
	db := InstanceOfMptStore()
	return generateMptStore(logger, id, db, AccountStateRootRetriever)
}

func generateMptStore(logger tmlog.Logger, id types.CommitID, db ethstate.Database, retriever StateRootRetriever) (*MptStore, error) {
	triegc := prque.New(nil)
	mptStore := &MptStore{
		storageTrieForWrite: make(map[ethcmn.Address]ethstate.Trie, 0),
		db:                  db,
		triegc:              triegc,
		logger:              logger,
		retriever:           retriever,
		exitSignal:          make(chan struct{}),
	}
	if mptStore.logger == nil {
		mptStore.logger = tmlog.NewNopLogger()
	}
	err := mptStore.openTrie(id)
	if logger != nil {
		gAsyncDB.SetLogger(logger.With("module", "asyncdb"))
	}
	if err != nil {
		return mptStore, err
	}
	if err := mptStore.openSnapshot(); err != nil {
		mptStore.logger.Error("open snapshot fail", "warn", err)
	} else {
		mptStore.logger.Info("open snapshot successfully", "snapshot", "ok")
	}

	return mptStore, err
}

func mockMptStore(logger tmlog.Logger, id types.CommitID) (*MptStore, error) {
	db := ethstate.NewDatabaseWithConfig(rawdb.NewMemoryDatabase(), &trie.Config{
		Cache:     int(TrieCacheSize),
		Journal:   "",
		Preimages: true,
	})
	return generateMptStore(logger, id, db, EmptyStateRootRetriever{})
}

func (ms *MptStore) openTrie(id types.CommitID) error {
	latestStoredHeight := ms.GetLatestStoredBlockHeight()
	openHeight := uint64(id.Version)
	if latestStoredHeight > 0 && openHeight > latestStoredHeight {
		return fmt.Errorf("fail to open mpt trie, the target version is: %d, the latest stored version is: %d, "+
			"please repair", openHeight, latestStoredHeight)
	}

	openedRootHash := ms.GetMptRootHash(openHeight)
	tr, err := ms.db.OpenTrie(openedRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}

	ms.trie = tr
	ms.version = id.Version
	ms.startVersion = id.Version
	ms.originalRoot = openedRootHash

	ms.logger.Info("open acc mpt trie", "version", openHeight, "trieHash", openedRootHash)

	ms.StartPrefetcher("mptStore")
	ms.prefetchData()

	return nil
}

func (ms *MptStore) GetImmutable(height int64) (*ImmutableMptStore, error) {
	rootHash := ms.GetMptRootHash(uint64(height))

	return NewImmutableMptStore(ms.db, rootHash)
}

/*
*  implement KVStore
 */
func (ms *MptStore) GetStoreType() types.StoreType {
	return StoreTypeMPT
}

func (ms *MptStore) CacheWrap() types.CacheWrap {
	stores := cachekv.NewStore(ms)
	stores.StatisticsCell = ms
	return stores
}

func (ms *MptStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(tracekv.NewStore(ms, w, tc))
}

func (ms *MptStore) Get(key []byte) []byte {
	switch mptKeyType(key) {
	case storageType:
		addr, stateRoot, realKey := decodeAddressStorageInfo(key)
		value, err := ms.getSnapStorage(addr, realKey)
		if err == nil {
			return value
		}
		t := ms.tryGetStorageTrie(addr, stateRoot, false)
		value, err = t.TryGet(realKey)
		if err != nil {
			return nil
		}
		return value
	case addressType:
		v, err := ms.getSnapAccount(key)
		if err == nil {
			return v
		}
		value, err := ms.db.CopyTrie(ms.trie).TryGet(key)
		if err != nil {
			return nil
		}

		return value
	default:
		panic(fmt.Errorf("not support key %s for mpt get", hex.EncodeToString(key)))
	}

}

func (ms *MptStore) tryGetStorageTrie(addr ethcmn.Address, stateRoot ethcmn.Hash, useCache bool) ethstate.Trie {
	if useCache {
		if t, ok := ms.storageTrieForWrite[addr]; ok {
			return t
		}
	}
	addrHash := mpttype.Keccak256HashWithSyncPool(addr[:])
	var t ethstate.Trie
	var err error
	t, err = ms.db.OpenStorageTrie(addrHash, stateRoot)
	if err != nil {
		t, err = ms.db.OpenStorageTrie(addrHash, ethcmn.Hash{})
		if err != nil {
			panic("unexcepted err")
		}
	}

	if useCache {
		ms.storageTrieForWrite[addr] = t
	}
	return t
}

func (ms *MptStore) Has(key []byte) bool {
	return ms.Get(key) != nil
}

func (ms *MptStore) Set(key, value []byte) {
	types.AssertValidValue(value)

	if ms.prefetcher != nil {
		ms.prefetcher.Used(ms.originalRoot, [][]byte{key})
	}
	switch mptKeyType(key) {
	case storageType:
		addr, stateRoot, realKey := decodeAddressStorageInfo(key)
		t := ms.tryGetStorageTrie(addr, stateRoot, true)
		t.TryUpdate(realKey, value)
		ms.updateSnapStorages(addr, realKey, value)
	case addressType:
		ms.trie.TryUpdate(key, value)
		ms.updateSnapAccounts(key, value)
	default:
		panic(fmt.Errorf("not support key %s for mpt set", hex.EncodeToString(key)))
	}

	return
}

func (ms *MptStore) Delete(key []byte) {
	if ms.prefetcher != nil {
		ms.prefetcher.Used(ms.originalRoot, [][]byte{key})
	}
	switch mptKeyType(key) {
	case storageType:
		addr, stateRoot, realKey := decodeAddressStorageInfo(key)
		t := ms.tryGetStorageTrie(addr, stateRoot, true)
		t.TryDelete(realKey)
		ms.updateSnapStorages(addr, realKey, nil)
	case addressType:
		ms.trie.TryDelete(key)
		ms.updateDestructs(key)
	default:
		panic(fmt.Errorf("not support key %s for mpt delete", hex.EncodeToString(key)))

	}
}

func (ms *MptStore) Iterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.db.CopyTrie(ms.trie), start, end)
}

func (ms *MptStore) ReverseIterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.db.CopyTrie(ms.trie), start, end)
}

/*
*  implement CommitStore, CommitKVStore
 */
func (ms *MptStore) commitStorageWithDelta(storageDelta []*trie.StorageDelta, nodeSets *trie.MergedNodeSet) {
	if storageDelta == nil || len(storageDelta) == 0 {
		return
	}
	for _, storage := range storageDelta {
		addr := storage.Addr
		addrHash := mpttype.Keccak256HashWithSyncPool(addr[:])
		stateRoot := ms.retriever.GetAccStateRoot(storage.PreAcc)
		t, err := ms.db.OpenStorageTrie(addrHash, stateRoot)
		if err != nil {
			panic(err)
		}
		_, set, err := t.CommitWithDelta(storage.NodeDelta, false)
		if set != nil {
			if err := nodeSets.Merge(set); err != nil {
				panic("fail to commit trie data(storage nodeSets merge): " + err.Error())
			}
		}
	}
}

func (ms *MptStore) commitStorageForDelta(nodeSets *trie.MergedNodeSet) []*trie.StorageDelta {
	storageDelta := make([]*trie.StorageDelta, 0, len(ms.storageTrieForWrite))
	for addr, v := range ms.storageTrieForWrite {
		stateR, set, outStorageDelta, err := v.CommitForDelta(false)
		if err != nil {
			panic(fmt.Errorf("unexcepted err:%v while commit storage tire ", err))
		}
		key := AddressStoreKey(addr.Bytes())
		preValue, err := ms.trie.TryGet(key)
		if err == nil { // maybe acc already been deleted
			newValue := ms.retriever.ModifyAccStateRoot(preValue, stateR)
			if err := ms.trie.TryUpdate(key, newValue); err != nil {
				panic(fmt.Errorf("unexcepted err:%v while update acc %s ", err, addr.String()))
			}
			ms.updateSnapAccounts(key, newValue)
		} else {
			panic(fmt.Errorf("unexcepted err:%v while update get acc %s ", err, addr.String()))
		}
		if set != nil {
			if err := nodeSets.Merge(set); err != nil {
				panic("fail to commit trie data(storage nodeSets merge): " + err.Error())
			}
		}
		storageDelta = append(storageDelta, &trie.StorageDelta{Addr: addr, PreAcc: preValue, NodeDelta: outStorageDelta})
		delete(ms.storageTrieForWrite, addr)
	}
	return storageDelta
}

func (ms *MptStore) commitStorage(nodeSets *trie.MergedNodeSet) {
	for addr, v := range ms.storageTrieForWrite {
		stateR, set, err := v.Commit(false)
		if err != nil {
			panic(fmt.Errorf("unexcepted err:%v while commit storage tire ", err))
		}
		key := AddressStoreKey(addr.Bytes())
		preValue, err := ms.trie.TryGet(key)
		if err == nil { // maybe acc already been deleted
			newValue := ms.retriever.ModifyAccStateRoot(preValue, stateR)
			if err := ms.trie.TryUpdate(key, newValue); err != nil {
				panic(fmt.Errorf("unexcepted err:%v while update acc %s ", err, addr.String()))
			}
			ms.updateSnapAccounts(key, newValue)
		} else {
			panic(fmt.Errorf("unexcepted err:%v while update get acc %s ", err, addr.String()))
		}
		if set != nil {
			if err := nodeSets.Merge(set); err != nil {
				panic("fail to commit trie data(storage nodeSets merge): " + err.Error())
			}
		}
		delete(ms.storageTrieForWrite, addr)
	}
}

func (ms *MptStore) CommitterCommit(inputDelta interface{}) (rootHash types.CommitID, outputDelta interface{}) {
	ms.version++

	// stop pre round prefetch
	ms.StopPrefetcher()
	nodeSets := trie.NewMergedNodeSet()

	var root ethcmn.Hash
	var set *trie.NodeSet
	var err error
	if applyDelta && inputDelta != nil {
		delta, ok := inputDelta.(*trie.MptDelta)
		if !ok {
			panic(fmt.Sprintf("wrong input delta of mpt. delta: %v", inputDelta))
		}
		ms.commitStorageWithDelta(delta.Storage, nodeSets)
		root, set, err = ms.trie.CommitWithDelta(delta.NodeDelta, true)
	} else if produceDelta {
		var outputNodeDelta []*trie.NodeDelta
		outStorage := ms.commitStorageForDelta(nodeSets)
		root, set, outputNodeDelta, err = ms.trie.CommitForDelta(true)
		outputDelta = &trie.MptDelta{NodeDelta: outputNodeDelta, Storage: outStorage}
	} else {
		ms.commitStorage(nodeSets)
		root, set, err = ms.trie.Commit(true)
	}
	if err != nil {
		panic("fail to commit trie data(acc_trie.Commit): " + err.Error())
	}

	if set != nil {
		if err := nodeSets.Merge(set); err != nil {
			panic("fail to commit trie data(acc nodeSets merge): " + err.Error())
		}
	}

	if err := ms.db.TrieDB().UpdateForOK(nodeSets, AccountStateRootRetriever.RetrieveStateRoot); err != nil {
		panic("fail to commit trie data (UpdateForOK): " + err.Error())
	}
	ms.SetMptRootHash(uint64(ms.version), root)
	ms.originalRoot = root

	// TODO: use a thread to push data to database
	// push data to database
	ms.PushData2Database(ms.version)

	ms.sprintDebugLog(ms.version)

	// start next found prefetch
	ms.StartPrefetcher("mptStore")

	return types.CommitID{
		Version: ms.version,
		Hash:    root.Bytes(),
	}, outputDelta
}

func (ms *MptStore) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: ms.version,
		Hash:    ms.trie.Hash().Bytes(),
	}
}

func (ms *MptStore) LastCommitVersion() int64 {
	return ms.version
}

func (ms *MptStore) SetPruning(options types.PruningOptions) {
	panic("cannot set pruning options on an initialized MPT store")
}

func (ms *MptStore) GetDBWriteCount() int {
	return gStatic.getDBWriteCount()
}

func (ms *MptStore) GetDBReadCount() int {
	return gStatic.getDBReadCount()
}

func (ms *MptStore) GetNodeReadCount() int {
	return ms.db.TrieDB().GetNodeReadCount()
}

func (ms *MptStore) GetCacheReadCount() int {
	return ms.db.TrieDB().GetCacheReadCount()
}

func (ms *MptStore) ResetCount() {
	gStatic.resetCount()
	ms.db.TrieDB().ResetCount()
}

func (ms *MptStore) GetDBReadTime() int {
	return gStatic.getDBReadTime()
}

func (ms *MptStore) StartTiming() {
	ms.statisticsBeginTime = time.Now()
}

func (ms *MptStore) EndTiming(tag string) {
	persist.GetStatistics().Accumulate(tag, ms.statisticsBeginTime)
}

// PushData2Database writes all associated state in cache to the database
func (ms *MptStore) PushData2Database(curHeight int64) {
	ms.cmLock.Lock()
	defer ms.cmLock.Unlock()

	curMptRoot := ms.GetMptRootHash(uint64(curHeight))
	if TrieDirtyDisabled {
		// If we're running an archive node, always flush
		ms.fullNodePersist(curMptRoot, curHeight)
	} else {
		ms.otherNodePersist(curMptRoot, curHeight)
	}
}

// fullNodePersist persist data without pruning
func (ms *MptStore) fullNodePersist(curMptRoot ethcmn.Hash, curHeight int64) {
	if curMptRoot == (ethcmn.Hash{}) || curMptRoot == ethtypes.EmptyRootHash {
		curMptRoot = ethcmn.Hash{}
	} else {
		if err := ms.db.TrieDB().Commit(curMptRoot, false, nil); err != nil {
			panic("fail to commit mpt data: " + err.Error())
		}
		ms.commitSnap(curMptRoot)
	}
	ms.SetLatestStoredBlockHeight(uint64(curHeight))
	ms.logger.Info("sync push acc data to db", "block", curHeight, "trieHash", curMptRoot)
}

// otherNodePersist persist data with pruning
func (ms *MptStore) otherNodePersist(curMptRoot ethcmn.Hash, curHeight int64) {
	triedb := ms.db.TrieDB()

	// Full but not archive node, do proper garbage collection
	triedb.Reference(curMptRoot, ethcmn.Hash{}) // metadata reference to keep trie alive
	ms.triegc.Push(curMptRoot, -int64(curHeight))

	// we start at startVersion, but the chosen height may be startVersion - triesInMemory
	if curHeight <= ms.startVersion {
		return
	}

	commitGap := TrieCommitGap
	if !EnableAsyncCommit {
		commitGap = 1
	}

	// If we exceeded out time allowance, flush an entire trie to disk
	if curHeight%commitGap == 0 {
		// If the header is missing (canonical chain behind), we're reorging a low
		// diff sidechain. Suspend committing until this operation is completed.
		chRoot := ms.GetMptRootHash(uint64(curHeight))
		if chRoot == (ethcmn.Hash{}) || chRoot == ethtypes.EmptyRootHash {
			chRoot = ethcmn.Hash{}
		} else {
			// Flush an entire trie and restart the counters, it's not a thread safe process,
			// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
			if err := triedb.Commit(chRoot, true, nil); err != nil {
				panic("fail to commit mpt data: " + err.Error())
			}
			ms.commitSnap(chRoot)
			gAsyncDB.LogStats()
		}
		ms.SetLatestStoredBlockHeight(uint64(curHeight))
		ms.logger.Info("async push acc data to db", "block", curHeight, "trieHash", chRoot)
	}
	gAsyncDB.Prune()
	// Garbage collect anything below our required write retention
	if curHeight > int64(TriesInMemory) {
		for !ms.triegc.Empty() {
			root, number := ms.triegc.Pop()
			if -number > curHeight-int64(TriesInMemory) {
				ms.triegc.Push(root, number)
				break
			}
			triedb.Dereference(root.(ethcmn.Hash))
		}
	}
}
func (ms *MptStore) CurrentVersion() int64 {
	return ms.version
}

func (ms *MptStore) OnStop() error {
	return ms.StopWithVersion(ms.version)
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (ms *MptStore) StopWithVersion(targetVersion int64) error {
	curVersion := uint64(targetVersion)
	ms.exitSignal <- struct{}{}
	ms.StopPrefetcher()

	ms.cmLock.Lock()
	defer ms.cmLock.Unlock()

	// Ensure the state of a recent block is also stored to disk before exiting.
	if !TrieDirtyDisabled {
		triedb := ms.db.TrieDB()
		okbcStartHeight := uint64(tmtypes.GetStartBlockHeight()) // start height of okbc

		latestStoreVersion := ms.GetLatestStoredBlockHeight()

		for version := latestStoreVersion; version <= curVersion; version++ {
			if version <= okbcStartHeight || version <= uint64(ms.startVersion) {
				continue
			}

			recentMptRoot := ms.GetMptRootHash(version)
			if recentMptRoot == (ethcmn.Hash{}) || recentMptRoot == ethtypes.EmptyRootHash {
				recentMptRoot = ethcmn.Hash{}
			} else {
				if err := triedb.Commit(recentMptRoot, true, nil); err != nil {
					ms.logger.Error("Failed to commit recent state trie", "err", err)
					break
				}
			}
			ms.SetLatestStoredBlockHeight(version)
			ms.logger.Info("Writing acc cached state to disk", "block", version, "trieHash", recentMptRoot)
		}

		for !ms.triegc.Empty() {
			ms.db.TrieDB().Dereference(ms.triegc.PopItem().(ethcmn.Hash))
		}
	}

	ts := time.Now()
	if err := ms.flattenPersistSnapshot(); err != nil {
		ms.logger.Error("Writing snapshot state to disk", "error", err)
	} else {
		ms.logger.Error("Writing snapshot successfully", "cost", time.Since(ts))
	}

	if gAsyncDB != nil {
		ms.logger.Error("Closing async database")
		err := gAsyncDB.Close()
		if err != nil {
			ms.logger.Error("Close async db", "err", err)
		} else {
			ms.logger.Error("Close async db OK.")
		}
	}

	return nil
}

/*
*  implement Queryable
 */
func (ms *MptStore) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(req.Data) == 0 {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrTxDecode, "query cannot be zero length"))
	}

	height := ms.getlatestHeight(uint64(req.Height))
	res.Height = int64(height)

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	trie, err := ms.getTrieByHeight(height)
	if err != nil {
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrInvalidVersion, "open trie failed: %s", err.Error()))
	}

	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes

		res.Key = key
		if req.Prove {
			value, proof, err := getVersionedWithProof(trie, key)
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
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{newProofOpMptValue(key, proof)}}
			} else {
				// value wasn't found
				res.Value = nil
				res.Proof = &merkle.Proof{Ops: []merkle.ProofOp{newProofOpMptAbsence(key, proof)}}
			}
		} else {
			res.Value, err = trie.TryGet(key)
			if err != nil {
				return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "failed to query in trie: %s", err.Error()))
			}
		}

	case "/subspace":
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "not supported query subspace path: %v in mptStore", req.Path))

	default:
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unexpected query path: %v", req.Path))
	}

	return res
}

// Handle latest the latest height - 1  (committed), if height is 0
func (ms *MptStore) getlatestHeight(height uint64) uint64 {
	if height == 0 {
		height = uint64(ms.version)
	}
	return height
}

func (ms *MptStore) getTrieByHeight(height uint64) (ethstate.Trie, error) {
	latestRootHash := ms.GetMptRootHash(height)
	if latestRootHash == NilHash {
		return nil, fmt.Errorf("header %d not found", height)
	}
	return ms.db.OpenTrie(latestRootHash)
}

// getVersionedWithProof returns the Merkle proof for given storage slot.
func getVersionedWithProof(trie ethstate.Trie, key []byte) ([]byte, [][]byte, error) {
	value, err := trie.TryGet(key)
	if err != nil {
		return nil, nil, err
	}

	var proof ProofList
	err = trie.Prove(crypto.Keccak256(key), 0, &proof)
	return value, proof, err
}

func (ms *MptStore) StartPrefetcher(namespace string) {

	if ms.prefetcher != nil {
		ms.prefetcher.Close()
		ms.prefetcher = nil
	}

	//ms.prefetcher = NewTriePrefetcher(ms.db, ms.originalRoot, namespace)
}

// StopPrefetcher terminates a running prefetcher and reports any leftover stats
// from the gathered metrics.
func (ms *MptStore) StopPrefetcher() {
	if ms.prefetcher != nil {
		ms.prefetcher.Close()
		ms.prefetcher = nil
	}
}

func (ms *MptStore) prefetchData() {
	go func() {
		for {
			select {
			case <-ms.exitSignal:
				return
			case <-GAccTryUpdateTrieChannel:
				if ms.prefetcher != nil {
					if trie := ms.prefetcher.Trie(ms.originalRoot); trie != nil {
						ms.trie = trie
					}
				}
				GAccTrieUpdatedChannel <- struct{}{}
			case addr := <-GAccToPrefetchChannel:
				if ms.prefetcher != nil {
					ms.prefetcher.Prefetch(ms.originalRoot, addr)
				}
			}
		}
	}()
}

func (ms *MptStore) SetUpgradeVersion(i int64) {}

var (
	keyPrefixStorageMpt  = byte(0)
	keyPrefixAddrMpt     = byte(1) // TODO auth.AddressStoreKeyPrefix
	prefixSizeInMpt      = 1
	storageKeySize       = prefixSizeInMpt + len(ethcmn.Address{}) + len(ethcmn.Hash{}) + len(ethcmn.Hash{})
	addrKeySize          = prefixSizeInMpt + sdk.AddrLen
	wasmContractKeySize  = prefixSizeInMpt + sdk.WasmContractAddrLen
	storageKeyPrefixSize = prefixSizeInMpt + len(ethcmn.Address{}) + len(ethcmn.Hash{})
)

func AddressStoragePrefixMpt(address ethcmn.Address, stateRoot ethcmn.Hash) []byte {
	t1 := append([]byte{keyPrefixStorageMpt}, address.Bytes()...)
	return append(t1, stateRoot.Bytes()...)
}

func decodeAddressStorageInfo(key []byte) (ethcmn.Address, ethcmn.Hash, []byte) {
	addr := ethcmn.BytesToAddress(key[prefixSizeInMpt : prefixSizeInMpt+20])
	storageRoot := ethcmn.BytesToHash(key[prefixSizeInMpt+20 : prefixSizeInMpt+20+32])
	updateKey := key[prefixSizeInMpt+20+32:]
	return addr, storageRoot, updateKey
}

func AddressStoreKey(addr []byte) []byte {
	return append([]byte{keyPrefixAddrMpt}, addr...)
}

var (
	storageType = 0
	addressType = 1
)

/*
storageType : 0x0 + addr + stateRoot + key
addressType : 0x1 + addr
*/

func mptKeyType(key []byte) int {
	switch key[0] {
	case keyPrefixAddrMpt:
		size := len(key)
		if size == wasmContractKeySize || size == addrKeySize {
			return addressType
		}
		return -1
	case keyPrefixStorageMpt:
		if len(key) == storageKeySize {
			return storageType
		}
		return -1
	default:
		return -1

	}
}

func IsStoragePrefix(prefix []byte) bool {
	return len(prefix) == storageKeyPrefixSize && prefix[0] == keyPrefixStorageMpt
}

func GetAddressFromStoragePrefix(prefix []byte) ethcmn.Address {
	return ethcmn.BytesToAddress(prefix[1:21])
}
