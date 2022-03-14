package keeper

import (
	"encoding/binary"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/mpt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	types3 "github.com/okex/exchain/libs/types"
	types2 "github.com/okex/exchain/x/evm/types"
)

// GetMptRootHash gets root mpt hash from block height
func (k *Keeper) GetMptRootHash(height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := k.db.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}
	return ethcmn.BytesToHash(rst)
}

// SetMptRootHash sets the mapping from block height to root mpt hash
func (k *Keeper) SetMptRootHash(ctx sdk.Context, hash ethcmn.Hash) {
	hhash := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	k.db.TrieDB().DiskDB().Put(append(mpt.KeyPrefixRootMptHash, hhash...), hash.Bytes())

	// put root hash to iavl and participate the process of calculate appHash
	if tmtypes.HigherThanMars(ctx.BlockHeight()) {
		store := k.Ada.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixRootMptHashForIavl)
		store.Set(hhash, hash.Bytes())
	}
}

// GetLatestStoredBlockHeight get latest stored mpt storage height
func (k *Keeper) GetLatestStoredBlockHeight() uint64 {
	rst, err := k.db.TrieDB().DiskDB().Get(mpt.KeyPrefixLatestStoredHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestStoredBlockHeight sets the latest stored storage height
func (k *Keeper) SetLatestStoredBlockHeight(height uint64) {
	hhash := sdk.Uint64ToBigEndian(height)
	k.db.TrieDB().DiskDB().Put(mpt.KeyPrefixLatestStoredHeight, hhash)
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
	k.EvmStateDb = types2.NewCommitStateDB(k.GenerateCSDBParams())
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (k *Keeper) OnStop(ctx sdk.Context) error {
	if !types3.TrieDirtyDisabled {
		triedb := k.db.TrieDB()
		oecStartHeight := uint64(tmtypes.GetStartBlockHeight()) // start height of oec

		latestVersion := uint64(ctx.BlockHeight())
		offset := uint64(mpt.TriesInMemory)
		for ; offset > 0; offset-- {
			if latestVersion > offset {
				version := latestVersion - offset
				if version <= oecStartHeight || version <= k.startHeight {
					continue
				}

				recentMptRoot := k.GetMptRootHash(version)
				if recentMptRoot == (ethcmn.Hash{}) || recentMptRoot == types.EmptyRootHash {
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
	if types3.TrieDirtyDisabled {
		if curMptRoot == (ethcmn.Hash{}) || curMptRoot == types.EmptyRootHash {
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

			if chosen % 10 == 0 {
				k.mptCommitMu.Lock()
				defer k.mptCommitMu.Unlock()
				// If the header is missing (canonical chain behind), we're reorging a low
				// diff sidechain. Suspend committing until this operation is completed.
				chRoot := k.GetMptRootHash(uint64(chosen))
				if chRoot == (ethcmn.Hash{}) || chRoot == types.EmptyRootHash {
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

	if tmtypes.HigherThanMars(ctx.BlockHeight()) || types3.EnableDoubleWrite {
		// The onleaf func is called _serially_, so we can reuse the same account
		// for unmarshalling every time.
		var storageRoot ethcmn.Hash
		root, _ := k.rootTrie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
			storageRoot.SetBytes(leaf)
			if storageRoot != types.EmptyRootHash {
				k.db.TrieDB().Reference(storageRoot, parent)
			}

			return nil
		})
		k.SetMptRootHash(ctx, root)
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
