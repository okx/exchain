package ante

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)


type NodeSignatureDecorator struct{}

func NewNodeSignatureDecorator() NodeSignatureDecorator {
	return NodeSignatureDecorator{}
}

func (mfd NodeSignatureDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// load whitelist and verify the signature

	// -1 for failure, 1 for success
	ctx.WithChecked(1)

	return ctx, fmt.Errorf("Invalid Node Signature Decorator")
	return next(ctx, tx, simulate)
}
