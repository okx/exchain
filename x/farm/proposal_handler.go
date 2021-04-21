package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/farm/types"
	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewManageWhiteListProposalHandler handles "gov" type message in "farm"
func NewManageWhiteListProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.ManageWhiteListProposal:
			return handleManageWhiteListProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(DefaultCodespace, content.ProposalType())
		}
	}
}

func handleManageWhiteListProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	// check
	manageWhiteListProposal, ok := proposal.Content.(types.ManageWhiteListProposal)
	if !ok {
		return types.ErrUnexpectedProposalType(proposal.Content.ProposalType())
	}
	if sdkErr := k.CheckMsgManageWhiteListProposal(ctx, manageWhiteListProposal); sdkErr != nil {
		return sdkErr
	}

	if manageWhiteListProposal.IsAdded {
		// add pool name into whitelist
		k.SetWhitelist(ctx, manageWhiteListProposal.PoolName)
		return nil
	}

	// remove pool name from whitelist
	k.DeleteWhiteList(ctx, manageWhiteListProposal.PoolName)
	return nil
}
