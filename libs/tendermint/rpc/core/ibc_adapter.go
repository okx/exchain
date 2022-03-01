package core

import (
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	rpctypes "github.com/okex/exchain/libs/tendermint/rpc/jsonrpc/types"
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

//func CommitIBC(ctx *rpctypes.Context, heightPtr *int64) (*coretypes.ResultCommit, error) {
//	height, err := getHeight(env.BlockStore.Height(), heightPtr)
//	if err != nil {
//		return nil, err
//	}
//
//	blockMeta := env.BlockStore.LoadBlockMeta(height)
//	if blockMeta == nil {
//		return nil, nil
//	}
//	header := blockMeta.Header
//
//	// If the next block has not been committed yet,
//	// use a non-canonical commit
//	if height == env.BlockStore.Height() {
//		commit := env.BlockStore.LoadSeenCommit(height)
//		return ConvResultCommitToTendermint(ctypes.NewResultCommit(&header, commit, false)), nil
//	}
//	// Return the canonical commit (comes from the block at height+1)
//	commit := env.BlockStore.LoadBlockCommit(height)
//	return ConvResultCommitToTendermint(ctypes.NewResultCommit(&header, commit, true)), nil
//}
