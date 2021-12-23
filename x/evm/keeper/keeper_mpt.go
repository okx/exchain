package keeper

import (
	"encoding/binary"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	KeyPrefixLatestHeight       = []byte{0x01}
	KeyPrefixRootMptHash        = []byte{0x02}
	KeyPrefixLatestStoredHeight = []byte{0x03}
)

// GetLatestBlockHeight get latest mpt storage height
func (k *Keeper) GetLatestBlockHeight() uint64 {
	rst, err := k.db.TrieDB().DiskDB().Get(KeyPrefixLatestHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestBlockHeight sets the latest storage height
func (k *Keeper) SetLatestBlockHeight(height uint64) {
	hhash := sdk.Uint64ToBigEndian(height)
	k.db.TrieDB().DiskDB().Put(KeyPrefixLatestHeight, hhash)
}

// GetRootMptHash gets root mpt hash from block height
func (k *Keeper) GetRootMptHash(height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := k.db.TrieDB().DiskDB().Get(append(KeyPrefixRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}

// SetRootMptHash sets the mapping from block height to root mpt hash
func (k *Keeper) SetRootMptHash(height uint64, hash ethcmn.Hash) {
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
	//types3.GetStartBlockHeight() // start height of oec
	latestHeight := k.GetLatestBlockHeight()
	latestRootHash := k.GetRootMptHash(latestHeight)

	tr, err := k.db.OpenTrie(latestRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}
	k.rootTrie = tr
}

func (k *Keeper) OnStop() error {
	for !k.triegc.Empty() {
		k.db.TrieDB().Dereference(k.triegc.PopItem().(ethcmn.Hash))
	}

	return nil
}

func (k *Keeper) PushData2Database(ctx sdk.Context, root ethcmn.Hash) {
	triedb := k.db.TrieDB()
	// Full but not archive node, do proper garbage collection
	triedb.Reference(root, ethcmn.Hash{}) // metadata reference to keep trie alive
	k.triegc.Push(root, -int64(ctx.BlockHeight()))

	if sdk.TrieDirtyDisabled {
		if err := triedb.Commit(root, false, nil); err != nil {
			panic("fail to commit mpt data: " + err.Error())
		}
		k.SetLatestStoredBlockHeight(uint64(ctx.BlockHeight()))
	} else {
		if ctx.BlockHeight() > core.TriesInMemory {
			// If we exceeded our memory allowance, flush matured singleton nodes to disk
			var (
				nodes, imgs = triedb.Size()
				limit       = ethcmn.StorageSize(256) * 1024 * 1024
			)

			if nodes > limit || imgs > 4*1024*1024 {
				triedb.Cap(limit - ethdb.IdealBatchSize)
			}
			// Find the next state trie we need to commit
			chosen := ctx.BlockHeight() - core.TriesInMemory

			// If the header is missing (canonical chain behind), we're reorging a low
			// diff sidechain. Suspend committing until this operation is completed.
			chRoot := k.GetRootMptHash(uint64(chosen))
			if chRoot == (ethcmn.Hash{}) {
				k.Logger(ctx).Debug("Reorg in progress, trie commit postponed", "number", chosen)
			} else {
				k.SetLatestStoredBlockHeight(uint64(chosen))

				// Flush an entire trie and restart the counters, it's not a thread safe process,
				// cannot use a go thread to run, or it will lead 'fatal error: concurrent map read and map write' error
				if err := triedb.Commit(chRoot, true, nil); err != nil {
					panic("fail to commit mpt data: " + err.Error())
				}
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

	latestHeight := uint64(ctx.BlockHeight())
	k.SetRootMptHash(latestHeight, root)
	k.SetLatestBlockHeight(latestHeight)

	k.PushData2Database(ctx, root)
}
