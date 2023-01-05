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

			// check oldContract
			oldAcc := ad.ak.GetAccount(ctx, proposalType.OldContractAddr)
			oldEthAcc, ok := oldAcc.(*ethermint.EthAccount)
			if !ok {
				return ctx, fmt.Errorf("acc:%s not EthAccount", proposalType.OldContractAddr)
			}
			if !oldEthAcc.IsContract() {
				return ctx, evmtypes.ErrNotContracAddress(fmt.Errorf(proposalType.OldContractAddr.String()))
			}

			//check newContract
			newAcc := ad.ak.GetAccount(ctx, proposalType.NewContractAddr)
			newEthAcc, ok := newAcc.(*ethermint.EthAccount)
			if !ok {
				return ctx, fmt.Errorf("acc:%s not EthAccount", proposalType.NewContractAddr)
			}
			if !newEthAcc.IsContract() {
				return ctx, evmtypes.ErrNotContracAddress(fmt.Errorf(proposalType.NewContractAddr.String()))
			}
		}

	default:
		return next(ctx, tx, simulate)
	}
	return ctx, nil
}
