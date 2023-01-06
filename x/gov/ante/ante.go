package ante

import (
	"fmt"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/gov/types"
	stakingkeeper "github.com/okex/exchain/x/staking"
)

type AnteDecorator struct {
	k  stakingkeeper.Keeper
	ak auth.AccountKeeper
}

func NewAnteDecorator(k stakingkeeper.Keeper, ak auth.AccountKeeper) AnteDecorator {
	return AnteDecorator{k: k, ak: ak}
}

func (ad AnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	m := tx.GetMsgs()[0]
	switch msg := m.(type) {
	case types.MsgSubmitProposal:
		switch proposalType := msg.Content.(type) {
		case evmtypes.ManagerContractByteCodeProposal:
			if !ad.k.IsValidator(ctx, msg.GetSigners()[0]) {
				return ctx, evmtypes.ErrCodeProposerMustBeValidator()
			}

			// check operation contract
			contract := ad.ak.GetAccount(ctx, proposalType.Contract)
			contractAcc, ok := contract.(*ethermint.EthAccount)
			if !ok {
				return ctx, fmt.Errorf("contract: %s not EthAccount", proposalType.Contract)
			}
			if !contractAcc.IsContract() {
				return ctx, evmtypes.ErrNotContracAddress(fmt.Errorf(proposalType.Contract.String()))
			}

			//check substitute contract
			substitute := ad.ak.GetAccount(ctx, proposalType.SubstituteContract)
			substituteAcc, ok := substitute.(*ethermint.EthAccount)
			if !ok {
				return ctx, fmt.Errorf("substitute contract:%s not EthAccount", proposalType.SubstituteContract)
			}
			if !substituteAcc.IsContract() {
				return ctx, evmtypes.ErrNotContracAddress(fmt.Errorf(proposalType.SubstituteContract.String()))
			}
		}

	default:
		return next(ctx, tx, simulate)
	}
	return ctx, nil
}
