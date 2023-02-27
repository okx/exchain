package keeper

import (
	"encoding/binary"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt"
	snap "github.com/okex/exchain/libs/cosmos-sdk/store/mpt/snapshot"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types"
)

// GetMptRootHash gets root mpt hash from block height
func (k *Keeper) GetMptRootHash(height uint64) ethcmn.Hash {
	heightBytes := sdk.Uint64ToBigEndian(height)
	rst, err := k.db.TrieDB().DiskDB().Get(append(mpt.KeyPrefixEvmRootMptHash, heightBytes...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}
	return ethcmn.BytesToHash(rst)
}

// GetLatestStoredBlockHeight get latest stored mpt storage height
func (k *Keeper) GetLatestStoredBlockHeight() uint64 {
	rst, err := k.db.TrieDB().DiskDB().Get(mpt.KeyPrefixEvmLatestStoredHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestStoredBlockHeight sets the latest stored storage height
func (k *Keeper) SetLatestStoredBlockHeight(height uint64) {
	heightBytes := sdk.Uint64ToBigEndian(height)
	k.db.TrieDB().DiskDB().Put(mpt.KeyPrefixEvmLatestStoredHeight, heightBytes)
}

func (k *Keeper) OpenTrie() {

}

func (k *Keeper) SetTargetMptVersion(targetVersion int64) {

}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (k *Keeper) OnStop(ctx sdk.Context) error {
	if !mpt.TrieDirtyDisabled {
		k.cmLock.Lock()
		defer k.cmLock.Unlock()

		triedb := k.db.TrieDB()
		oecStartHeight := uint64(tmtypes.GetStartBlockHeight()) // start height of oec

		latestStoreVersion := k.GetLatestStoredBlockHeight()
		curVersion := uint64(ctx.BlockHeight())
		for version := latestStoreVersion; version <= curVersion; version++ {
			if version <= oecStartHeight || version <= k.startHeight {
				continue
			}

			recentMptRoot := k.GetMptRootHash(version)
			if recentMptRoot == (ethcmn.Hash{}) || recentMptRoot == ethtypes.EmptyRootHash {
				recentMptRoot = ethcmn.Hash{}
			} else {
				if err := triedb.Commit(recentMptRoot, true, nil); err != nil {
					k.Logger().Error("Failed to commit recent state trie", "err", err)
					break
				}
			}
			k.SetLatestStoredBlockHeight(version)
			k.Logger().Info("Writing evm cached state to disk", "block", version, "trieHash", recentMptRoot)
		}

		for !k.triegc.Empty() {
			k.db.TrieDB().Dereference(k.triegc.PopItem().(ethcmn.Hash))
		}
	}

	return nil
}

// PushData2Database writes all associated state in cache to the database
func (k *Keeper) PushData2Database(height int64, log log.Logger) {
	k.cmLock.Lock()
	defer k.cmLock.Unlock()

	curMptRoot := k.GetMptRootHash(uint64(height))
	if mpt.TrieDirtyDisabled {
		// If we're running an archive node, always flush
		k.fullNodePersist(curMptRoot, height, log)
	} else {
		k.otherNodePersist(curMptRoot, height, log)
	}
}

// fullNodePersist persist data without pruning
func (k *Keeper) fullNodePersist(curMptRoot ethcmn.Hash, curHeight int64, log log.Logger) {
	if curMptRoot == (ethcmn.Hash{}) || curMptRoot == ethtypes.EmptyRootHash {
		curMptRoot = ethcmn.Hash{}
	} else {
		// Commit all cached state changes into underlying memory database.
		if err := k.db.TrieDB().Commit(curMptRoot, false, nil); err != nil {
			panic("fail to commit mpt data: " + err.Error())
		}
	}
	k.SetLatestStoredBlockHeight(uint64(curHeight))
	log.Info("sync push evm data to db", "block", curHeight, "trieHash", curMptRoot)
}

// otherNodePersist persist data with pruning
func (k *Keeper) otherNodePersist(curMptRoot ethcmn.Hash, curHeight int64, log log.Logger) {
	triedb := k.db.TrieDB()

	// Full but not archive node, do proper garbage collection
	triedb.Reference(curMptRoot, ethcmn.Hash{}) // metadata reference to keep trie alive
	k.triegc.Push(curMptRoot, -int64(curHeight))

	if curHeight > mpt.TriesInMemory {
		// If we exceeded our memory allowance, flush matured singleton nodes to disk
		var (
			nodes, imgs = triedb.Size()
			nodesLimit  = ethcmn.StorageSize(mpt.TrieNodesLimit) * 1024 * 1024
			imgsLimit   = ethcmn.StorageSize(mpt.TrieImgsLimit) * 1024 * 1024
		)

		if nodes > nodesLimit || imgs > imgsLimit {
			triedb.Cap(nodesLimit - ethdb.IdealBatchSize)
		}
		// Find the next state trie we need to commit
		chosen := curHeight - mpt.TriesInMemory

		if chosen <= int64(k.startHeight) {
			return
		}

		if chosen%mpt.TrieCommitGap == 0 {
			// If the header is missing (canonical chain behind), we're reorging a low
			// diff sidechain. Suspend committing until this operation is completed.
			chRoot := k.GetMptRootHash(uint64(chosen))
			if chRoot == (ethcmn.Hash{}) || chRoot == ethtypes.EmptyRootHash {
				chRoot = ethcmn.Hash{}
			} else {
				// Flush an entire trie and restart the counters, it's not a thread safe process,
				// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
				if err := triedb.Commit(chRoot, true, nil); err != nil {
					panic("fail to commit mpt data: " + err.Error())
				}
			}
			k.SetLatestStoredBlockHeight(uint64(chosen))
			log.Info("async push evm data to db", "block", chosen, "trieHash", chRoot)
		}

		// Garbage collect anything below our required write retention
		for !k.triegc.Empty() {
			root, number := k.triegc.Pop()
			if -number > chosen {
				k.triegc.Push(root, number)
				break
			}
			triedb.Dereference(root.(ethcmn.Hash))
		}
	}
}

/*
 * Getters for keys in x/evm/types/keys.go
 * TODO: these interfaces are used for setting/getting data in rawdb, instead of iavl.
 * TODO: delete these if we decide persist data in iavl.
 */
func (k Keeper) getBlockHashInDiskDB(hash []byte) (int64, bool) {
	key := types.AppendBlockHashKey(hash)
	bz, err := k.db.TrieDB().DiskDB().Get(key)
	if err != nil {
		return 0, false
	}
	if len(bz) == 0 {
		return 0, false
	}

	height := binary.BigEndian.Uint64(bz)
	return int64(height), true
}

func (k Keeper) setBlockHashInDiskDB(hash []byte, height int64) {
	key := types.AppendBlockHashKey(hash)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	k.db.TrieDB().DiskDB().Put(key, bz)
}

func (k Keeper) iterateBlockHashInDiskDB(fn func(key []byte, value []byte) (stop bool)) {
	iterator := k.db.TrieDB().DiskDB().NewIterator(types.KeyPrefixBlockHash, nil)
	defer iterator.Release()
	for iterator.Next() {
		if !types.IsBlockHashKey(iterator.Key()) {
			continue
		}
		key, value := iterator.Key(), iterator.Value()
		if stop := fn(key, value); stop {
			break
		}
	}
}

func (k Keeper) getBlockBloomInDiskDB(height int64) ethtypes.Bloom {
	key := types.AppendBloomKey(height)
	bz, err := k.db.TrieDB().DiskDB().Get(key)
	if err != nil {
		return ethtypes.Bloom{}
	}

	return ethtypes.BytesToBloom(bz)
}

func (k Keeper) setBlockBloomInDiskDB(height int64, bloom ethtypes.Bloom) {
	key := types.AppendBloomKey(height)
	k.db.TrieDB().DiskDB().Put(key, bloom.Bytes())
}

func (k Keeper) iterateBlockBloomInDiskDB(fn func(key []byte, value []byte) (stop bool)) {
	iterator := k.db.TrieDB().DiskDB().NewIterator(types.KeyPrefixBloom, nil)
	defer iterator.Release()
	for iterator.Next() {
		if !types.IsBloomKey(iterator.Key()) {
			continue
		}
		key, value := iterator.Key(), iterator.Value()
		if stop := fn(key, value); stop {
			break
		}
	}
}

func (k *Keeper) openSnapshot() error {
	latestHeight := k.getLatestStoredBlockHeight()
	latestRootHash := k.getMptRootHash(latestHeight)
	var recovery bool
	if layer := rawdb.ReadSnapshotRecoveryNumber(snap.GetDiskDB()); layer != nil && *layer > latestHeight {
		k.logger.Error("Enabling snapshot recovery", "chainhead", layer, "diskbase", *layer)
		recovery = true
	}
	var err error
	k.snaps, err = snapshot.New(snap.GetDiskDB(), k.db.TrieDB(), 256, latestRootHash, false, true, recovery)
	if err != nil {
		k.logger.Error("open snapshot error ", "error", err)
		return fmt.Errorf("open snapshot error %v", err)
	}
	k.snap = k.snaps.Snapshot(latestRootHash)

	return nil
}

func (k *Keeper) getLatestStoredBlockHeight() uint64 {
	rst, err := k.db.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

func (k *Keeper) getMptRootHash(height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := k.db.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}
