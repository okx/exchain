package ante

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)


type NodeSignatureDecorator struct{
	logger log.Logger
}

func NewNodeSignatureDecorator(l log.Logger) NodeSignatureDecorator {
	return NodeSignatureDecorator{
		logger: l,
	}
}

func (n NodeSignatureDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// load whitelist and verify the signature
	res := 1
	// -1 for failure, 1 for success
	ctx = ctx.WithVerifyResult(res)

	defer n.logger.Info("NodeSignatureDecorator AnteHandle",
		"VerifyResult", ctx.VerifyResult(),
		"tx-type", fmt.Sprintf("%T", tx),
		)

	if res < 0 {
		return ctx, fmt.Errorf("Invalid Node Signature Decorator")
	} else {
		return next(ctx, tx, simulate)
	}
}
