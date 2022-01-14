package keeper

import (
	"encoding/binary"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	types2 "github.com/okex/exchain/x/evm/types"
)

var (
	KeyPrefixRootMptHash        = []byte{0x01}
	KeyPrefixLatestStoredHeight = []byte{0x02}
)

const TriesInMemory = 100

// GetMptRootHash gets root mpt hash from block height
func (k *Keeper) GetMptRootHash(height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := k.db.TrieDB().DiskDB().Get(append(KeyPrefixRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}

// SetMptRootHash sets the mapping from block height to root mpt hash
func (k *Keeper) SetMptRootHash(height uint64, hash ethcmn.Hash) {
	hhash := sdk.Uint64ToBigEndian(height)
	k.db.TrieDB().DiskDB().Put(append(KeyPrefixRootMptHash, hhash...), hash.Bytes())
}

// GetLatestStoredBlockHeight get latest stored mpt storage height
func (k *Keeper) GetLatestStoredBlockHeight() uint64 {
	rst, err := k.db.TrieDB().DiskDB().Get(KeyPrefixLatestStoredHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestStoredBlockHeight sets the latest stored storage height
func (k *Keeper) SetLatestStoredBlockHeight(height uint64) {
	hhash := sdk.Uint64ToBigEndian(height)
	k.db.TrieDB().DiskDB().Put(KeyPrefixLatestStoredHeight, hhash)
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
	// Ensure the state of a recent block is also stored to disk before exiting.
	// We're writing three different states to catch different restart scenarios:
	//  - HEAD:     So we don't need to reprocess any blocks in the general case
	//  - HEAD-1:   So we don't do large reorgs if our HEAD becomes an uncle
	//  - HEAD-127: So we have a hard limit on the number of blocks reexecuted
	if !sdk.TrieDirtyDisabled {
		triedb := k.db.TrieDB()
		oecStartHeight := uint64(tmtypes.GetStartBlockHeight()) // start height of oec

		for _, offset := range []uint64{0, 1, TriesInMemory - 1} {
			if number := uint64(ctx.BlockHeight()); number > offset {
				recent := number - offset
				if recent <= oecStartHeight || recent <= k.startHeight {
					break
				}

				recentMptRoot := k.GetMptRootHash(recent)
				if recentMptRoot == (ethcmn.Hash{}) {
					k.Logger(ctx).Debug("Reorg in progress, trie commit postponed", "block", recent)
				} else {
					k.Logger(ctx).Info("Writing cached state to disk", "block", recent, "trieHash", recentMptRoot)
					if err := triedb.Commit(recentMptRoot, true, nil); err != nil {
						k.Logger(ctx).Error("Failed to commit recent state trie", "err", err)
					}
				}
			}
		}

		for !k.triegc.Empty() {
			k.db.TrieDB().Dereference(k.triegc.PopItem().(ethcmn.Hash))
		}
	}

	return nil
}

func (k *Keeper) PushData2Database(ctx sdk.Context) {
	curHeight := ctx.BlockHeight()
	curMptRoot := k.GetMptRootHash(uint64(curHeight))

	triedb := k.db.TrieDB()
	if sdk.TrieDirtyDisabled {
		if err := triedb.Commit(curMptRoot, false, nil); err != nil {
			panic("fail to commit mpt data: " + err.Error())
		}
		k.SetLatestStoredBlockHeight(uint64(curHeight))
		k.Logger(ctx).Info("sync push data to db", "block", curHeight, "trieHash", curMptRoot)
	} else {
		// Full but not archive node, do proper garbage collection
		triedb.Reference(curMptRoot, ethcmn.Hash{}) // metadata reference to keep trie alive
		k.triegc.Push(curMptRoot, -int64(curHeight))

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

			if chosen <= int64(k.startHeight) {
				return
			}

			// If the header is missing (canonical chain behind), we're reorging a low
			// diff sidechain. Suspend committing until this operation is completed.
			chRoot := k.GetMptRootHash(uint64(chosen))
			if chRoot == (ethcmn.Hash{}) {
				k.Logger(ctx).Debug("Reorg in progress, trie commit postponed", "number", chosen)
			} else {
				// Flush an entire trie and restart the counters, it's not a thread safe process,
				// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
				if err := triedb.Commit(chRoot, true, nil); err != nil {
					panic("fail to commit mpt data: " + err.Error())
				}
				k.SetLatestStoredBlockHeight(uint64(chosen))
				k.Logger(ctx).Info("async push data to db", "block", chosen, "trieHash", chRoot)
			}

			// Garbage collect anything below our required write retention
			for !k.triegc.Empty() {
				root, number := k.triegc.Pop()
				if int64(-number) > chosen {
					k.triegc.Push(root, number)
					break
				}
				triedb.Dereference(root.(ethcmn.Hash))
			}
		}
	}
}

func (k *Keeper) Commit(ctx sdk.Context) {
	// commit contract storage mpt trie
	k.EvmStateDb.WithContext(ctx).Commit(true)

	// The onleaf func is called _serially_, so we can reuse the same account
	// for unmarshalling every time.
	var storageRoot ethcmn.Hash
	root, _ := k.rootTrie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		_, content, _, err := rlp.Split(leaf)
		if err != nil {
			k.EvmStateDb.SetError(err)
		}
		storageRoot.SetBytes(content)
		if storageRoot != types.EmptyRootHash {
			k.db.TrieDB().Reference(storageRoot, parent)
		}

		return nil
	})
	k.SetMptRootHash(uint64(ctx.BlockHeight()), root)
}
