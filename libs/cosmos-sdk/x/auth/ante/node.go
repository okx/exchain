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
	res := 1
	// -1 for failure, 1 for success
	ctx = ctx.WithVerifyResult(res)

	defer fmt.Printf("NodeSignatureDecorator VerifyResult %d: %T\n", ctx.VerifyResult(), tx)

	if res < 0 {
		return ctx, fmt.Errorf("Invalid Node Signature Decorator")
	} else {
		return next(ctx, tx, simulate)
	}
}
