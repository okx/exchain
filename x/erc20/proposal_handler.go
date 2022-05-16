package erc20

import (
	ethcmm "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/erc20/types"
	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewProposalHandler handles "gov" type message in "erc20"
func NewProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.TokenMappingProposal:
			return handleTokenMappingProposal(ctx, k, content)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleTokenMappingProposal(ctx sdk.Context, k *Keeper, p types.TokenMappingProposal) sdk.Error {
	if len(p.Contract) == 0 {
		// delete existing mapping
		k.DeleteExternalContractForDenom(ctx, p.Denom)
	} else {
		// update the mapping
		contract := ethcmm.HexToAddress(p.Contract)
		if err := k.SetExternalContractForDenom(ctx, p.Denom, contract); err != nil {
			return err
		}
	}
	return nil
}
