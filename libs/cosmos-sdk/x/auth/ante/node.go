package ante

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
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

	var wtx authtypes.WrappedTx
	var ok bool
	if wtx, ok = tx.(authtypes.WrappedTx); !ok {
		return ctx, fmt.Errorf("Invalid WrappedTx")
	}

	_ = wtx
	// load whitelist to verify the signature
	// wtx.Signature == verify(wtx.NodeKey, wtx.Payload+wtx.Metadata)

	res := -1
	// -1 for failure, 1 for success
	ctx = ctx.WithNodeSigVerifyResult(res)

	defer n.logger.Info("NodeSignatureDecorator AnteHandle",
		"NodeSigVerifyResult", ctx.NodeSigVerifyResult(),
		"tx-type", fmt.Sprintf("%T", tx),
		)

	if res < 0 {
		return ctx, fmt.Errorf("Invalid Node Signature Decorator")
	} else {
		return next(ctx, tx, simulate)
	}
}
