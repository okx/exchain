package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	govtypes "github.com/okex/exchain/x/gov/types"
)

type ParamsKeeper interface {
	CheckMsgSubmitProposal(ctx sdk.Context, msg govtypes.MsgSubmitProposal) sdk.Error
}

type AnteDecorator struct {
	pk ParamsKeeper
}

func NewAnteDecorator(pk ParamsKeeper) AnteDecorator {
	return AnteDecorator{pk: pk}
}

func (ad AnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, m := range tx.GetMsgs() {
		switch msg := m.(type) {
		case govtypes.MsgSubmitProposal:
			if err := ad.checkMsgSubmitProposal(ctx, simulate, msg); err != nil {
				return ctx, err
			}
		}
	}

	return next(ctx, tx, simulate)
}

func (ad AnteDecorator) checkMsgSubmitProposal(ctx sdk.Context, simulate bool, msg govtypes.MsgSubmitProposal) error {
	err := ad.pk.CheckMsgSubmitProposal(ctx, msg)
	if err != nil {
		if e, ok := err.(sdk.EnvelopedErr); ok {
			if e.ABCICode() == common.CodeUnknownProposalType {
				// it's not an error that we don't know the proposal's type here .
				return nil
			}
		}
	}

	return err
}
