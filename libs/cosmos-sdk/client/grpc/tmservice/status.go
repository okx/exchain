package tmservice

import (
	"context"

	cliContext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
)

func getNodeStatus(ctx context.Context, clientCtx cliContext.CLIContext) (*ctypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &ctypes.ResultStatus{}, err
	}
	return node.Status()
}
