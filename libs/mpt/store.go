package mpt

import (
	"fmt"
	"github.com/VictoriaMetrics/fastcache"
	ethcmn "github.com/ethereum/go-ethereum/common"
	types3 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	types2 "github.com/okex/exchain/libs/types"
	"io"

	"github.com/ethereum/go-ethereum/common/prque"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/okex/exchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/tracekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

var cdc = codec.New()

var (
	_ types.KVStore       = (*MptStore)(nil)
	_ types.CommitStore   = (*MptStore)(nil)
	_ types.CommitKVStore = (*MptStore)(nil)
	_ types.Queryable     = (*MptStore)(nil)
)

// MptStore Implements types.KVStore and CommitKVStore.
type MptStore struct {
	trie    ethstate.Trie
	db      ethstate.Database
	triegc  *prque.Prque
	logger  tmlog.Logger
	kvCache *fastcache.Cache

	prefetcher   *TriePrefetcher
	originalRoot ethcmn.Hash
	exitSignal   chan struct{}

	version      int64
	startVersion int64
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
	triegc := prque.New(nil)

	mptStore := &MptStore{
		db:         db,
		triegc:     triegc,
		logger:     logger,
		kvCache:    fastcache.New(int(AccStoreCache) * 1024 * 1024),
		exitSignal: make(chan struct{}),
	}
	err := mptStore.openTrie(id)

	return mptStore, err
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

	if ms.logger != nil {
		ms.logger.Info("open acc mpt trie", "version", openHeight, "trieHash", openedRootHash)
	}

	ms.StartPrefetcher("mptStore")
	ms.prefetchData()

	return nil
}

/*
*  implement KVStore
 */
func (ms *MptStore) GetStoreType() types.StoreType {
	return StoreTypeMPT
}

func (ms *MptStore) CacheWrap() types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(ms)
}

func (ms *MptStore) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	//TODO implement me
	return cachekv.NewStore(tracekv.NewStore(ms, w, tc))
}

func (ms *MptStore) Get(key []byte) []byte {
	if enc := ms.kvCache.Get(nil, key); len(enc) > 0 {
		return enc
	}

	value, err := ms.trie.TryGet(key)
	if err != nil {
		return nil
	}
	ms.kvCache.Set(key, value)

	return value
}

func (ms *MptStore) Has(key []byte) bool {
	if ms.kvCache.Has(key) {
		return true
	}

	return ms.Get(key) != nil
}

func (ms *MptStore) Set(key, value []byte) {
	types.AssertValidValue(value)

	if ms.prefetcher != nil {
		ms.prefetcher.Used(ms.originalRoot, [][]byte{key})
	}

	ms.kvCache.Set(key, value)
	err := ms.trie.TryUpdate(key, value)
	if err != nil {
		return
	}
	return
}

func (ms *MptStore) Delete(key []byte) {
	if ms.prefetcher != nil {
		ms.prefetcher.Used(ms.originalRoot, [][]byte{key})
	}

	ms.kvCache.Del(key)
	err := ms.trie.TryDelete(key)
	if err != nil {
		return
	}
}

func (ms *MptStore) Iterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.trie, start, end)
}

func (ms *MptStore) ReverseIterator(start, end []byte) types.Iterator {
	return newMptIterator(ms.trie, start, end)
}

/*
*  implement CommitStore, CommitKVStore
 */
func (ms *MptStore) CommitterCommit(delta *iavl.TreeDelta) (types.CommitID, *iavl.TreeDelta) {
	ms.version++

	// stop pre round prefetch
	ms.StopPrefetcher()

	root, err := ms.trie.Commit(nil)
	if err != nil {
		panic("fail to commit trie data: " + err.Error())
	}
	ms.SetMptRootHash(uint64(ms.version), root)
	ms.originalRoot = root

	// TODO: use a thread to push data to database
	// push data to database
	ms.PushData2Database(ms.version)

	// start next found prefetch
	ms.StartPrefetcher("mptStore")

	return types.CommitID{
		Version: ms.version,
		Hash:    root.Bytes(),
	}, nil
}

func (ms *MptStore) LastCommitID() types.CommitID {
	return types.CommitID{
		Version: ms.version,
		Hash:    ms.trie.Hash().Bytes(),
	}
}

func (ms *MptStore) SetPruning(options types.PruningOptions) {
	panic("cannot set pruning options on an initialized MPT store")
}

func (ms *MptStore) GetDBWriteCount() int {
	return 0
}

func (ms *MptStore) GetDBReadCount() int {
	return 0
}

func (ms *MptStore) GetNodeReadCount() int {
	return 0
}

func (ms *MptStore) ResetCount() {
	return
}

func (ms *MptStore) GetDBReadTime() int {
	return 0
}

func (ms *MptStore) PushData2Database(curHeight int64) {
	curMptRoot := ms.GetMptRootHash(uint64(curHeight))

	triedb := ms.db.TrieDB()
	if types2.TrieDirtyDisabled {
		if curMptRoot == (ethcmn.Hash{}) || curMptRoot == types3.EmptyRootHash {
			curMptRoot = ethcmn.Hash{}
		} else {
			if err := triedb.Commit(curMptRoot, false, nil); err != nil {
				panic("fail to commit mpt data: " + err.Error())
			}
		}
		ms.SetLatestStoredBlockHeight(uint64(curHeight))
		if ms.logger != nil {
			ms.logger.Info("sync push acc data to db", "block", curHeight, "trieHash", curMptRoot)
		}
	} else {
		// Full but not archive node, do proper garbage collection
		triedb.Reference(curMptRoot, ethcmn.Hash{}) // metadata reference to keep trie alive
		ms.triegc.Push(curMptRoot, -int64(curHeight))

		if curHeight > TriesInMemory {
			// If we exceeded our memory allowance, flush matured singleton nodes to disk
			var (
				nodes, imgs = triedb.Size()
				limit       = ethcmn.StorageSize(256) * 1024 * 1024
			)

			if nodes > limit || imgs > 4*1024*1024 {
				triedb.Cap(limit - ethdb.IdealBatchSize)
			}
			// Find the next state trie we need to commit
			chosen := curHeight - TriesInMemory

			// we start at startVersion, but the chosen height may be startVersion - triesInMemory
			if chosen <= ms.startVersion {
				return
			}

			// If we exceeded out time allowance, flush an entire trie to disk
			if chosen%10 == 0 {
				// If the header is missing (canonical chain behind), we're reorging a low
				// diff sidechain. Suspend committing until this operation is completed.
				chRoot := ms.GetMptRootHash(uint64(chosen))
				if chRoot == (ethcmn.Hash{}) || chRoot == types3.EmptyRootHash {
					chRoot = ethcmn.Hash{}
				} else {
					// Flush an entire trie and restart the counters, it's not a thread safe process,
					// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
					if err := triedb.Commit(chRoot, true, nil); err != nil {
						panic("fail to commit mpt data: " + err.Error())
					}
				}
				ms.SetLatestStoredBlockHeight(uint64(chosen))
				if ms.logger != nil {
					ms.logger.Info("async push acc data to db", "block", chosen, "trieHash", chRoot)
				}
			}

			// Garbage collect anything below our required write retention
			for !ms.triegc.Empty() {
				root, number := ms.triegc.Pop()
				if int64(-number) > chosen {
					ms.triegc.Push(root, number)
					break
				}
				triedb.Dereference(root.(ethcmn.Hash))
			}
		}
	}
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (ms *MptStore) OnStop() error {
	ms.exitSignal <- struct{}{}
	ms.StopPrefetcher()

	if !tmtypes.HigherThanMars(ms.version) && !types2.EnableDoubleWrite {
		return nil
	}

	// Ensure the state of a recent block is also stored to disk before exiting.
	if !types2.TrieDirtyDisabled {
		triedb := ms.db.TrieDB()
		oecStartHeight := uint64(tmtypes.GetStartBlockHeight()) // start height of oec

		latestStoreVersion := ms.GetLatestStoredBlockHeight()
		curVersion := uint64(ms.version)
		for version := latestStoreVersion; version <= curVersion; version++ {
			if version <= oecStartHeight || version <= uint64(ms.startVersion) {
				continue
			}

			recentMptRoot := ms.GetMptRootHash(version)
			if recentMptRoot == (ethcmn.Hash{}) || recentMptRoot == types3.EmptyRootHash {
				recentMptRoot = ethcmn.Hash{}
			} else {
				if err := triedb.Commit(recentMptRoot, true, nil); err != nil {
					if ms.logger != nil {
						ms.logger.Error("Failed to commit recent state trie", "err", err)
					}
					break
				}
			}
			ms.SetLatestStoredBlockHeight(version)
			if ms.logger != nil {
				ms.logger.Info("Writing acc cached state to disk", "block", version, "trieHash", recentMptRoot)
			}
		}

		for !ms.triegc.Empty() {
			ms.db.TrieDB().Dereference(ms.triegc.PopItem().(ethcmn.Hash))
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

	height := ms.getHeight(req)
	res.Height = int64(height)

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	trie, err := ms.getTrieWithHeight(height)
	if err != nil {
		res.Log = fmt.Sprintf("trie of height %d doesn't exist: %s", err)
		return
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
			res.Value, _ = getVersioned(trie, key)
		}

	case "/subspace":
		var KVs []types.KVPair

		subspace := req.Data
		res.Key = subspace

		iterator := newMptIterator(trie, subspace, sdk.PrefixEndBytes(subspace))
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

func (ms *MptStore) getHeight(req abci.RequestQuery) uint64 {
	height := uint64(req.Height)
	latestStoredBlockHeight := ms.GetLatestStoredBlockHeight()
	if height == 0 || height > latestStoredBlockHeight {
		height = latestStoredBlockHeight
	}
	return height
}

// Handle gatest the latest height, if height is 0
func (ms *MptStore) getTrieWithHeight(height uint64) (ethstate.Trie, error) {
	latestRootHash := ms.GetMptRootHash(height)
	return ms.db.OpenTrie(latestRootHash)
}

func getVersioned(trie ethstate.Trie, key []byte) ([]byte, error) {
	return trie.TryGet(key)
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
	if !tmtypes.HigherThanMars(ms.version) {
		return
	}

	if ms.prefetcher != nil {
		ms.prefetcher.Close()
		ms.prefetcher = nil
	}

	ms.prefetcher = NewTriePrefetcher(ms.db, ms.originalRoot, namespace)
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
