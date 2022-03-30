package keeper

import (
	"encoding/binary"
	"fmt"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/mpt"
	mpttypes "github.com/okex/exchain/libs/mpt/types"
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

// SetMptRootHash sets the mapping from block height to root mpt hash
func (k *Keeper) SetMptRootHash(ctx sdk.Context, hash ethcmn.Hash) {
	heightBytes := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	k.db.TrieDB().DiskDB().Put(append(mpt.KeyPrefixEvmRootMptHash, heightBytes...), hash.Bytes())

	// put root hash to iavl and participate the process of calculate appHash
	if tmtypes.HigherThanMars(ctx.BlockHeight()) {
		k.accountKeeper.SetEvmRootHash(ctx, hash.Bytes())
	}
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
	//startHeight := types2.GetStartBlockHeight() // start height of oec
	latestStoredHeight := k.GetLatestStoredBlockHeight()
	latestStoredRootHash := k.GetMptRootHash(latestStoredHeight)

	tr, err := k.db.OpenTrie(latestStoredRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}
	k.rootTrie = tr
	k.rootHash = latestStoredRootHash
	k.startHeight = latestStoredHeight

	if latestStoredHeight == 0 {
		k.startHeight = uint64(tmtypes.GetStartBlockHeight())
	}
}

func (k *Keeper) SetTargetMptVersion(targetVersion int64) {
	if !tmtypes.HigherThanMars(targetVersion) {
		return
	}

	latestStoredHeight := k.GetLatestStoredBlockHeight()
	if latestStoredHeight < uint64(targetVersion) {
		panic(fmt.Sprintf("The target mpt height is: %v, but the latest stored evm height is: %v", targetVersion, latestStoredHeight))
	}
	targetMptRootHash := k.GetMptRootHash(uint64(targetVersion))

	tr, err := k.db.OpenTrie(targetMptRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}
	k.rootTrie = tr
	k.rootHash = targetMptRootHash
	k.startHeight = uint64(targetVersion)
	if targetVersion == 0 {
		k.startHeight = uint64(tmtypes.GetStartBlockHeight())
	}

	k.EvmStateDb = types.NewCommitStateDB(k.GenerateCSDBParams())
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (k *Keeper) OnStop(ctx sdk.Context) error {
	if !mpttypes.TrieDirtyDisabled {
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
				k.mptCommitMu.Lock()
				if err := triedb.Commit(recentMptRoot, true, nil); err != nil {
					k.Logger(ctx).Error("Failed to commit recent state trie", "err", err)
					break
				}
				k.mptCommitMu.Unlock()
			}
			k.SetLatestStoredBlockHeight(version)
			k.Logger(ctx).Info("Writing evm cached state to disk", "block", version, "trieHash", recentMptRoot)
		}

		for !k.triegc.Empty() {
			k.db.TrieDB().Dereference(k.triegc.PopItem().(ethcmn.Hash))
		}
	}

	return nil
}

func (k *Keeper) PushData2Database(height int64, log log.Logger) {
	curHeight := height
	curMptRoot := k.GetMptRootHash(uint64(curHeight))

	triedb := k.db.TrieDB()
	if mpttypes.TrieDirtyDisabled {
		if curMptRoot == (ethcmn.Hash{}) || curMptRoot == ethtypes.EmptyRootHash {
			curMptRoot = ethcmn.Hash{}
		} else {
			k.mptCommitMu.Lock()
			if err := triedb.Commit(curMptRoot, false, nil); err != nil {
				panic("fail to commit mpt data: " + err.Error())
			}
			k.mptCommitMu.Unlock()
		}
		k.SetLatestStoredBlockHeight(uint64(curHeight))
		log.Info("sync push evm data to db", "block", curHeight, "trieHash", curMptRoot)
	} else {
		// Full but not archive node, do proper garbage collection
		triedb.Reference(curMptRoot, ethcmn.Hash{}) // metadata reference to keep trie alive
		k.triegc.Push(curMptRoot, -int64(curHeight))

		if curHeight > mpt.TriesInMemory {
			// If we exceeded our memory allowance, flush matured singleton nodes to disk
			var (
				nodes, imgs = triedb.Size()
				limit       = ethcmn.StorageSize(256) * 1024 * 1024
			)

			if nodes > limit || imgs > 4*1024*1024 {
				triedb.Cap(limit - ethdb.IdealBatchSize)
			}
			// Find the next state trie we need to commit
			chosen := curHeight - mpt.TriesInMemory

			if chosen <= int64(k.startHeight) {
				return
			}

			if chosen%10 == 0 {
				k.mptCommitMu.Lock()
				defer k.mptCommitMu.Unlock()
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
}

func (k *Keeper) Commit(ctx sdk.Context) {
	k.mptCommitMu.Lock()
	defer k.mptCommitMu.Unlock()

	// commit contract storage mpt trie
	k.EvmStateDb.WithContext(ctx).Commit(true)
	k.EvmStateDb.StopPrefetcher()

	if tmtypes.HigherThanMars(ctx.BlockHeight()) || mpttypes.EnableDoubleWrite {
		k.rootTrie = k.EvmStateDb.GetRootTrie()

		// The onleaf func is called _serially_, so we can reuse the same account
		// for unmarshalling every time.
		var storageRoot ethcmn.Hash
		root, _ := k.rootTrie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
			storageRoot.SetBytes(leaf)
			if storageRoot != ethtypes.EmptyRootHash {
				k.db.TrieDB().Reference(storageRoot, parent)
			}

			return nil
		})
		k.SetMptRootHash(ctx, root)
		k.rootHash = root
	}
}

func (k *Keeper) AddMptAsyncTask(height int64) {
	k.asyncChain <- height
}
func (k *Keeper) asyncCommit(logger log.Logger) {
	go func() {
		for {
			select {
			case height := <-k.asyncChain:
				k.PushData2Database(height, logger)
			}
		}
	}()

}

/*
 * Getters for keys in x/evm/types/keys.go
 */
func (k Keeper) getBlockHashInDiskDB(hash []byte) (int64, bool) {
	bz, err := k.db.TrieDB().DiskDB().Get(append(types.KeyPrefixBlockHash, hash...))
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
	bz := sdk.Uint64ToBigEndian(uint64(height))
	k.db.TrieDB().DiskDB().Put(append(types.KeyPrefixBlockHash, hash...), bz)
}

func (k Keeper) iterateBlockHashInDiskDB(fn func(key []byte, value []byte) (stop bool)) {
	iterator := k.db.TrieDB().DiskDB().NewIterator(types.KeyPrefixBlockHash, nil)
	defer iterator.Release()
	for iterator.Next() {
		if len(iterator.Key()) != len(types.KeyPrefixBlockHash)+ethcmn.HashLength {
			continue
		}
		key, value := iterator.Key(), iterator.Value()
		if stop := fn(key, value); stop {
			break
		}
	}
}

func (k Keeper) getBlockBloomInDiskDB(height int64) ethtypes.Bloom {
	bz, err := k.db.TrieDB().DiskDB().Get(append(types.KeyPrefixBloom, types.BloomKey(height)...))
	if err != nil {
		return ethtypes.Bloom{}
	}

	return ethtypes.BytesToBloom(bz)
}

func (k Keeper) setBlockBloomInDiskDB(height int64, bloom ethtypes.Bloom) {
	k.db.TrieDB().DiskDB().Put(append(types.KeyPrefixBloom, types.BloomKey(height)...), bloom.Bytes())
}

func (k Keeper) iterateBlockBloomInDiskDB(fn func(key []byte, value []byte) (stop bool)) {
	iterator := k.db.TrieDB().DiskDB().NewIterator(types.KeyPrefixBloom, nil)
	defer iterator.Release()
	for iterator.Next() {
		if len(iterator.Key()) != len(types.KeyPrefixBloom)+8 {
			continue
		}
		key, value := iterator.Key(), iterator.Value()
		if stop := fn(key, value); stop {
			break
		}
	}
}

func (k Keeper) getChainConfigInDiskDB() (types.ChainConfig, bool) {
	var config types.ChainConfig
	bz, err := k.db.TrieDB().DiskDB().Get(types.KeyPrefixChainConfig)
	if err != nil {
		return config, false
	}
	if err := config.UnmarshalFromAmino(k.cdc, bz[4:]); err != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &config)
	}
	return config, true
}

func (k Keeper) setChainConfigInDiskDB(config types.ChainConfig) {
	bz := k.cdc.MustMarshalBinaryBare(config)
	k.db.TrieDB().DiskDB().Put(types.KeyPrefixChainConfig, bz)
}
