package staking

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/staking/types"

	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewProposalHandler handles "gov" type message in "staking"
func NewProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.ProposeValidatorProposal:
			return handleProposeValidatorProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleProposeValidatorProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	// check
	proposeValidatorProposal, ok := proposal.Content.(types.ProposeValidatorProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}

	if proposeValidatorProposal.IsAdd {
		//// add/update treasures into state
		//if err := k.UpdateTreasures(ctx, manageTreasuresProposal.Treasures); err != nil {
		//	return types.ErrTreasuresInternal(err)
		//}
		return nil
	}

	// delete treasures into state
	//if err := k.DeleteTreasures(ctx, manageTreasuresProposal.Treasures); err != nil {
	//	return types.ErrTreasuresInternal(err)
	//}
	return nil
}
