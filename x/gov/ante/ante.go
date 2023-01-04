package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/gov/types"
	stakingkeeper "github.com/okex/exchain/x/staking"
)

type AnteDecorator struct {
	k stakingkeeper.Keeper
}

func NewAnteDecorator(k stakingkeeper.Keeper) AnteDecorator {
	return AnteDecorator{k: k}
}

func (ad AnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	m := tx.GetMsgs()[0]
	switch msg := m.(type) {
	case types.MsgSubmitProposal:
		if msg.Content.ProposalType() == "ManageContractByteCode" {
			if !ad.k.IsValidator(ctx, msg.GetSigners()[0]) {
				return ctx, evmtypes.ErrCodeProposerMustBeValidator()
			}
		}
	default:
		return next(ctx, tx, simulate)
	}
	return ctx, nil
}
