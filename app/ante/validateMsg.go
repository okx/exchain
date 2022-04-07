package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type ValidateMsgHandler func(ctx sdk.Context, msgs []sdk.Msg) error

type ValidateMsgHandlerDecorator struct {
	validateMsgHandler ValidateMsgHandler
}

func NewValidateMsgHandlerDecorator(validateHandler ValidateMsgHandler) ValidateMsgHandlerDecorator {
	return ValidateMsgHandlerDecorator{validateMsgHandler: validateHandler}
}

func (vmhd ValidateMsgHandlerDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// *ABORT* the tx in case of failing to validate it in checkTx mode
	if ctx.IsCheckTx() && !simulate && vmhd.validateMsgHandler != nil {
		err := vmhd.validateMsgHandler(ctx, tx.GetMsgs())
		if err != nil {
			return ctx, err
		}
	}

	return next(ctx, tx, simulate)
}
