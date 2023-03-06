package store

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/tendermint/go-amino"

	"github.com/pkg/errors"

	db "github.com/okx/okbchain/libs/tm-db"
	dbm "github.com/okx/okbchain/libs/tm-db"

	"github.com/okx/okbchain/libs/tendermint/types"
)

/*
BlockStore is a simple low level store for blocks.

There are three types of information stored:
  - BlockMeta:   Meta information about each block
  - Block part:  Parts of each block, aggregated w/ PartSet
  - Commit:      The commit part of each block, for gossiping precommit votes

Currently the precommit signatures are duplicated in the Block parts as
well as the Commit.  In the future this may change, perhaps by moving
the Commit data outside the Block. (TODO)

The store can be assumed to contain all contiguous blocks between base and height (inclusive).

// NOTE: BlockStore methods will panic if they encounter errors
// deserializing loaded data, indicating probable corruption on disk.
*/
type BlockStore struct {
	db dbm.DB

	mtx    sync.RWMutex
	base   int64
	height int64
}

// NewBlockStore returns a new BlockStore with the given DB,
// initialized to the last height that was committed to the DB.
func NewBlockStore(db dbm.DB) *BlockStore {
	bsjson := LoadBlockStoreStateJSON(db)
	return &BlockStore{
		base:   bsjson.Base,
		height: bsjson.Height,
		db:     db,
	}
}

// Base returns the first known contiguous block height, or 0 for empty block stores.
func (bs *BlockStore) Base() int64 {
	bs.mtx.RLock()
	defer bs.mtx.RUnlock()
	return bs.base
}

// Height returns the last known contiguous block height, or 0 for empty block stores.
func (bs *BlockStore) Height() int64 {
	bs.mtx.RLock()
	defer bs.mtx.RUnlock()
	return bs.height
}

// Size returns the number of blocks in the block store.
func (bs *BlockStore) Size() int64 {
	bs.mtx.RLock()
	defer bs.mtx.RUnlock()
	if bs.height == 0 {
		return 0
	}
	return bs.height - bs.base + 1
}

var blockLoadBufPool = &sync.Pool{
	New: func() interface{} {
		return &[2]bytes.Buffer{}
	},
}

// LoadBlock returns the block with the given height.
// If no block is found for that height, it returns nil.
func (bs *BlockStore) LoadBlock(height int64) *types.Block {
	b, _ := bs.LoadBlockWithExInfo(height)
	return b
}

// LoadBlockWithExInfo returns the block with the given height.
// and the BlockPartInfo is used to make block parts
func (bs *BlockStore) LoadBlockWithExInfo(height int64) (*types.Block, *types.BlockExInfo) {
	bufs := blockLoadBufPool.Get().(*[2]bytes.Buffer)
	defer blockLoadBufPool.Put(bufs)

	loadBuf, uncompressedBuf := &bufs[0], &bufs[1]

	loadBuf.Reset()
	uncompressedBuf.Reset()

	info := bs.loadBlockPartsBytesTo(height, loadBuf, uncompressedBuf)
	if loadBuf.Len() == 0 {
		return nil, nil
	}
	if !info.IsCompressed() {
		return bs.unmarshalBlockByBytes(loadBuf.Bytes()), &info
	} else {
		return bs.unmarshalBlockByBytes(uncompressedBuf.Bytes()), &info
	}
}

// unmarshalBlockByBytes returns the block with the given block parts bytes
func (bs *BlockStore) unmarshalBlockByBytes(blockBytes []byte) *types.Block {
	var block = new(types.Block)
	bz, err := amino.GetBinaryBareFromBinaryLengthPrefixed(blockBytes)
	if err == nil {
		err = block.UnmarshalFromAmino(cdc, bz)
	}
	if err != nil {
		block = new(types.Block)
		err = cdc.UnmarshalBinaryLengthPrefixed(blockBytes, block)
		if err != nil {
			// NOTE: The existence of meta should imply the existence of the
			// block. So, make sure meta is only saved after blocks are saved.
			panic(errors.Wrap(err, fmt.Sprintf("Error reading block, height:%d", block.Height)))
		}
	}
	return block
}

// LoadBlockByHash returns the block with the given hash.
// If no block is found for that hash, it returns nil.
// Panics if it fails to parse height associated with the given hash.
func (bs *BlockStore) LoadBlockByHash(hash []byte) *types.Block {
	bz, err := bs.db.Get(calcBlockHashKey(hash))
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		return nil
	}

	s := string(bz)
	height, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		panic(errors.Wrapf(err, "failed to extract height from %s", s))
	}
	return bs.LoadBlock(height)
}

func loadBlockPartFromBytes(bz []byte) *types.Part {
	if len(bz) == 0 {
		return nil
	}
	var part = new(types.Part)
	err := part.UnmarshalFromAmino(cdc, bz)
	if err != nil {
		part = new(types.Part)
		err = cdc.UnmarshalBinaryBare(bz, part)
		if err != nil {
			panic(errors.Wrap(err, "Error reading block part"))
		}
	}
	return part
}

// LoadBlockPart returns the Part at the given index
// from the block at the given height.
// If no part is found for the given height and index, it returns nil.
func (bs *BlockStore) LoadBlockPart(height int64, index int) *types.Part {
	v, err := bs.db.GetUnsafeValue(calcBlockPartKey(height, index), func(bz []byte) (interface{}, error) {
		return loadBlockPartFromBytes(bz), nil
	})
	if err != nil {
		panic(err)
	}
	return v.(*types.Part)
}

func loadBlockPartBytesFromBytesTo(bz []byte, buf *bytes.Buffer) {
	if len(bz) == 0 {
		return
	}
	lenBefore := buf.Len()
	err := unmarshalBlockPartBytesTo(bz, buf)
	if err == nil {
		return
	}
	part := loadBlockPartFromBytes(bz)
	buf.Truncate(lenBefore)
	buf.Write(part.Bytes)
}

func (bs *BlockStore) loadBlockPartBytesTo(height int64, index int, buf *bytes.Buffer) {
	_, err := bs.db.GetUnsafeValue(calcBlockPartKey(height, index), func(bz []byte) (interface{}, error) {
		loadBlockPartBytesFromBytesTo(bz, buf)
		return nil, nil
	})
	if err != nil {
		panic(err)
	}
}

// loadBlockPartsBytesTo load all block parts bytes to the given buffer,
// buf *Buffer stores the original block parts bytes,
// uncompressed *Buffer stores the uncompressed block parts bytes if block is compressed
func (bs *BlockStore) loadBlockPartsBytesTo(height int64, buf *bytes.Buffer, uncompressed *bytes.Buffer) types.BlockExInfo {
	var blockMeta = bs.LoadBlockMeta(height)
	if blockMeta == nil {
		return types.BlockExInfo{}
	}
	blockPartSize, bufBeforeLen := 0, buf.Len()
	for i := 0; i < blockMeta.BlockID.PartsHeader.Total; i++ {
		bs.loadBlockPartBytesTo(height, i, buf)
		if i == 0 {
			blockPartSize = buf.Len() - bufBeforeLen
		}
	}

	// uncompress if the block part bytes was created by compress block
	compressSign, err := types.UncompressBlockFromBytesTo(buf.Bytes(), uncompressed)
	if err != nil {
		panic(errors.Wrap(err, "failed to uncompress block"))
	}

	return types.BlockExInfo{
		BlockCompressType: compressSign / types.CompressDividing,
		BlockCompressFlag: compressSign % types.CompressDividing,
		BlockPartSize:     blockPartSize,
	}
}

func decodeBlockMeta(bz []byte) (*types.BlockMeta, error) {
	if len(bz) == 0 {
		return nil, nil
	}
	var blockMeta = new(types.BlockMeta)
	err := blockMeta.UnmarshalFromAmino(cdc, bz)
	if err != nil {
		err = cdc.UnmarshalBinaryBare(bz, blockMeta)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading block meta")
		}
	}
	return blockMeta, nil
}

// LoadBlockMeta returns the BlockMeta for the given height.
// If no block is found for the given height, it returns nil.
func (bs *BlockStore) LoadBlockMeta(height int64) *types.BlockMeta {
	v, err := bs.db.GetUnsafeValue(calcBlockMetaKey(height), func(bz []byte) (interface{}, error) {
		return decodeBlockMeta(bz)
	})
	if err != nil {
		panic(err)
	}
	return v.(*types.BlockMeta)
}

// LoadBlockCommit returns the Commit for the given height.
// This commit consists of the +2/3 and other Precommit-votes for block at `height`,
// and it comes from the block.LastCommit for `height+1`.
// If no commit is found for the given height, it returns nil.
func (bs *BlockStore) LoadBlockCommit(height int64) *types.Commit {
	var commit = new(types.Commit)
	bz, err := bs.db.Get(calcBlockCommitKey(height))
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		return nil
	}
	err = cdc.UnmarshalBinaryBare(bz, commit)
	if err != nil {
		panic(errors.Wrap(err, "Error reading block commit"))
	}
	return commit
}

// LoadSeenCommit returns the locally seen Commit for the given height.
// This is useful when we've seen a commit, but there has not yet been
// a new block at `height + 1` that includes this commit in its block.LastCommit.
func (bs *BlockStore) LoadSeenCommit(height int64) *types.Commit {
	var commit = new(types.Commit)
	bz, err := bs.db.Get(calcSeenCommitKey(height))
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		return nil
	}
	err = cdc.UnmarshalBinaryBare(bz, commit)
	if err != nil {
		panic(errors.Wrap(err, "Error reading block seen commit"))
	}
	return commit
}

// PruneBlocks removes block up to (but not including) a height. It returns number of blocks pruned.
func (bs *BlockStore) PruneBlocks(height int64) (uint64, error) {
	return bs.deleteBatch(height, false)
}

// DeleteBlocksFromTop removes block down to (but not including) a height. It returns number of blocks deleted.
func (bs *BlockStore) DeleteBlocksFromTop(height int64) (uint64, error) {
	return bs.deleteBatch(height, true)
}

func (bs *BlockStore) deleteBatch(height int64, deleteFromTop bool) (uint64, error) {
	if height <= 0 {
		return 0, fmt.Errorf("height must be greater than 0")
	}

	bs.mtx.RLock()
	top := bs.height
	base := bs.base
	bs.mtx.RUnlock()
	if height > top {
		return 0, fmt.Errorf("cannot delete beyond the latest height %v, delete from top %t", top, deleteFromTop)
	}
	if height < base {
		return 0, fmt.Errorf("cannot delete to height %v, it is lower than base height %v, delete from top %t",
			height, base, deleteFromTop)
	}

	deleted := uint64(0)
	batch := bs.db.NewBatch()
	defer batch.Close()
	flush := func(batch db.Batch, height int64) error {
		// We can't trust batches to be atomic, so update base first to make sure noone
		// tries to access missing blocks.
		bs.mtx.Lock()
		if deleteFromTop {
			bs.height = height
		} else {
			bs.base = height
		}
		bs.mtx.Unlock()
		bs.saveState()

		err := batch.WriteSync()
		if err != nil {
			batch.Close()
			return fmt.Errorf("failed to delete to height %v, delete from top %t: %w", height, deleteFromTop, err)
		}
		batch.Close()
		return nil
	}

	deleteFn := func(h int64) error {
		meta := bs.LoadBlockMeta(h)
		if meta == nil { // assume already deleted
			return nil
		}
		batch.Delete(calcBlockMetaKey(h))
		batch.Delete(calcBlockHashKey(meta.BlockID.Hash))
		batch.Delete(calcBlockCommitKey(h))
		batch.Delete(calcSeenCommitKey(h))
		for p := 0; p < meta.BlockID.PartsHeader.Total; p++ {
			batch.Delete(calcBlockPartKey(h, p))
		}
		deleted++

		// flush every 1000 blocks to avoid batches becoming too large
		if deleted%1000 == 0 && deleted > 0 {
			err := flush(batch, h)
			if err != nil {
				return err
			}
			batch = bs.db.NewBatch()
		}
		return nil
	}

	if deleteFromTop {
		for h := top; h > height; h-- {
			err := deleteFn(h)
			if err != nil {
				return 0, err
			}
		}
	} else {
		for h := base; h < height; h++ {
			err := deleteFn(h)
			if err != nil {
				return 0, err
			}
		}
	}

	err := flush(batch, height)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// SaveBlock persists the given block, blockParts, and seenCommit to the underlying db.
// blockParts: Must be parts of the block
// seenCommit: The +2/3 precommits that were seen which committed at height.
//
//	If all the nodes restart after committing a block,
//	we need this to reload the precommits to catch-up nodes to the
//	most recent height.  Otherwise they'd stall at H-1.
func (bs *BlockStore) SaveBlock(block *types.Block, blockParts *types.PartSet, seenCommit *types.Commit) {
	batch := bs.db.NewBatch()
	defer batch.Close()

	if block == nil {
		panic("BlockStore can only save a non-nil block")
	}

	height := block.Height
	hash := block.Hash()

	if g, w := height, bs.Height()+1; bs.Base() > 0 && g != w {
		panic(fmt.Sprintf("BlockStore can only save contiguous blocks. Wanted %v, got %v", w, g))
	}
	if !blockParts.IsComplete() {
		panic(fmt.Sprintf("BlockStore can only save complete block part sets"))
	}

	// Save block meta
	blockMeta := types.NewBlockMeta(block, blockParts)
	metaBytes := cdc.MustMarshalBinaryBare(blockMeta)
	batch.Set(calcBlockMetaKey(height), metaBytes)
	batch.Set(calcBlockHashKey(hash), []byte(fmt.Sprintf("%d", height)))

	// Save block parts
	for i := 0; i < blockParts.Total(); i++ {
		part := blockParts.GetPart(i)
		bs.saveBlockPart(batch, height, i, part)
	}

	// Save block commit (duplicate and separate from the Block)
	blockCommitBytes := cdc.MustMarshalBinaryBare(block.LastCommit)
	batch.Set(calcBlockCommitKey(height-1), blockCommitBytes)

	// Save seen commit (seen +2/3 precommits for block)
	// NOTE: we can delete this at a later height
	seenCommitBytes := cdc.MustMarshalBinaryBare(seenCommit)
	batch.Set(calcSeenCommitKey(height), seenCommitBytes)

	// Done!
	bs.mtx.Lock()
	bs.height = height
	if bs.base == 0 {
		bs.base = height
	}
	bs.mtx.Unlock()

	// Save new BlockStoreStateJSON descriptor
	bs.saveStateBatch(batch)

	// Flush
	batch.WriteSync()
}

func (bs *BlockStore) saveBlockPart(batch db.Batch, height int64, index int, part *types.Part) {
	partBytes := cdc.MustMarshalBinaryBare(part)
	batch.Set(calcBlockPartKey(height, index), partBytes)
}

func (bs *BlockStore) saveState() {
	bs.mtx.RLock()
	bsJSON := BlockStoreStateJSON{
		Base:   bs.base,
		Height: bs.height,
	}
	bs.mtx.RUnlock()
	bsJSON.Save(bs.db)
}

func (bs *BlockStore) saveStateBatch(batch db.Batch) {
	bs.mtx.RLock()
	bsJSON := BlockStoreStateJSON{
		Base:   bs.base,
		Height: bs.height,
	}
	bs.mtx.RUnlock()
	bsJSON.saveBatch(batch)
}

//-----------------------------------------------------------------------------

func calcBlockMetaKey(height int64) []byte {
	return amino.StrToBytes(strings.Join([]string{"H", strconv.FormatInt(height, 10)}, ":"))
}

func calcBlockPartKey(height int64, partIndex int) []byte {
	return amino.StrToBytes(strings.Join([]string{"P", strconv.FormatInt(height, 10), strconv.Itoa(partIndex)}, ":"))
}

func calcBlockCommitKey(height int64) []byte {
	return amino.StrToBytes(strings.Join([]string{"C", strconv.FormatInt(height, 10)}, ":"))
}

func calcSeenCommitKey(height int64) []byte {
	return amino.StrToBytes(strings.Join([]string{"SC", strconv.FormatInt(height, 10)}, ":"))
}

func calcBlockHashKey(hash []byte) []byte {
	return amino.StrToBytes(strings.Join([]string{"BH", amino.HexEncodeToString(hash)}, ":"))
}

//-----------------------------------------------------------------------------

var blockStoreKey = []byte("blockStore")

// BlockStoreStateJSON is the block store state JSON structure.
type BlockStoreStateJSON struct {
	Base   int64 `json:"base"`
	Height int64 `json:"height"`
}

// Save persists the blockStore state to the database as JSON.
func (bsj BlockStoreStateJSON) Save(db dbm.DB) {
	bytes, err := cdc.MarshalJSON(bsj)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal state bytes: %v", err))
	}
	db.SetSync(blockStoreKey, bytes)
}

// Save persists the blockStore state to the database as JSON.
func (bsj BlockStoreStateJSON) saveBatch(batch dbm.Batch) {
	bytes, err := cdc.MarshalJSON(bsj)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal state bytes: %v", err))
	}
	batch.Set(blockStoreKey, bytes)
}

// LoadBlockStoreStateJSON returns the BlockStoreStateJSON as loaded from disk.
// If no BlockStoreStateJSON was previously persisted, it returns the zero value.
func LoadBlockStoreStateJSON(db dbm.DB) BlockStoreStateJSON {
	bytes, err := db.Get(blockStoreKey)
	if err != nil {
		panic(err)
	}
	if len(bytes) == 0 {
		return BlockStoreStateJSON{
			Base:   0,
			Height: types.GetStartBlockHeight(),
		}
	}
	bsj := BlockStoreStateJSON{}
	err = cdc.UnmarshalJSON(bytes, &bsj)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal bytes: %X", bytes))
	}
	// Backwards compatibility with persisted data from before Base existed.
	if bsj.Height > 0 && bsj.Base == 0 {
		bsj.Base = 1
	}
	return bsj
}
