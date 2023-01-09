package ante

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/gov/types"
	stakingkeeper "github.com/okex/exchain/x/staking"
)

type AnteDecorator struct {
	sk stakingkeeper.Keeper
	ak auth.AccountKeeper
}

func NewAnteDecorator(k stakingkeeper.Keeper, ak auth.AccountKeeper) AnteDecorator {
	return AnteDecorator{sk: k, ak: ak}
}

func (ad AnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	for _, m := range tx.GetMsgs() {
		switch msg := m.(type) {
		case types.MsgSubmitProposal:
			switch proposalType := msg.Content.(type) {
			case evmtypes.ManagerContractByteCodeProposal:
				if !ad.sk.IsValidator(ctx, msg.Proposer) {
					return ctx, evmtypes.ErrCodeProposerMustBeValidator()
				}

				// check operation contract
				contract := ad.ak.GetAccount(ctx, proposalType.Contract)
				contractAcc, ok := contract.(*ethermint.EthAccount)
				if !ok || !contractAcc.IsContract() {
					return ctx, evmtypes.ErrNotContracAddress(fmt.Errorf(ethcmn.BytesToAddress(proposalType.Contract).String()))
				}

				//check substitute contract
				substitute := ad.ak.GetAccount(ctx, proposalType.SubstituteContract)
				substituteAcc, ok := substitute.(*ethermint.EthAccount)
				if !ok || !substituteAcc.IsContract() {
					return ctx, evmtypes.ErrNotContracAddress(fmt.Errorf(ethcmn.BytesToAddress(proposalType.SubstituteContract).String()))
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}
