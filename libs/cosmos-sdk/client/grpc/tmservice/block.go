package tmservice

import (
	"context"

	cliContext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	tmproto "github.com/okex/exchain/libs/tendermint/proto/types"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
)

func getBlock(ctx context.Context, clientCtx cliContext.CLIContext, height *int64) (*coretypes.ResultBlock, error) {
	// get the node
	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.Block(height)
}

func GetProtoBlock(ctx context.Context, clientCtx cliContext.CLIContext, height *int64) (tmproto.BlockID, *tmproto.Block, error) {
	block, err := getBlock(ctx, clientCtx, height)
	if err != nil {
		return tmproto.BlockID{}, nil, err
	}
	protoBlock, err := block.Block.ToProto()
	if err != nil {
		return tmproto.BlockID{}, nil, err
	}
	protoBlockID := block.BlockID.ToProto()

	return protoBlockID, protoBlock, nil
}
