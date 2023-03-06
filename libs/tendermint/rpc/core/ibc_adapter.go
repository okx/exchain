package core

import (
	ctypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
	rpctypes "github.com/okx/okbchain/libs/tendermint/rpc/jsonrpc/types"
	"github.com/okx/okbchain/libs/tendermint/types"
	//coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

func CommitIBC(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.IBCResultCommit, error) {
	height, err := getHeight(env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}

	blockMeta := env.BlockStore.LoadBlockMeta(height)
	if blockMeta == nil {
		return nil, nil
	}
	header := blockMeta.Header

	// If the next block has not been committed yet,
	// use a non-canonical commit
	if height == env.BlockStore.Height() {
		commit := env.BlockStore.LoadSeenCommit(height)
		return ConvResultCommitTOIBC(ctypes.NewResultCommit(&header, commit, false)), nil
	}
	// Return the canonical commit (comes from the block at height+1)
	commit := env.BlockStore.LoadBlockCommit(height)
	return ConvResultCommitTOIBC(ctypes.NewResultCommit(&header, commit, true)), nil
}

func CM40Block(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.CM40ResultBlock, error) {
	height, err := getHeight(env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}

	block := env.BlockStore.LoadBlock(height)
	blockMeta := env.BlockStore.LoadBlockMeta(height)

	if blockMeta == nil {
		return &ctypes.CM40ResultBlock{BlockID: ConvBlockID2CM40BlockID(types.BlockID{}), Block: ConvBlock2CM40Block(block)}, nil
	}
	ret := &ctypes.CM40ResultBlock{BlockID: ConvBlockID2CM40BlockID(blockMeta.BlockID), Block: ConvBlock2CM40Block(block)}
	return ret, nil
}
