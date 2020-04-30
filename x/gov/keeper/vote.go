package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/okex/okchain/x/gov/types"
)

// AddVote adds a vote on a specific proposal
func (keeper Keeper) AddVote(
	ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, option types.VoteOption,
) (sdk.Error, string) {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return types.ErrUnknownProposal(keeper.Codespace(), proposalID), ""
	}
	if proposal.Status != types.StatusVotingPeriod {
		return types.ErrInvalidateProposalStatus(keeper.Codespace(),
			fmt.Sprintf("The status of proposal %d is in %s can not be voted.",
				proposal.ProposalID, proposal.Status)), ""
	}

	if !types.ValidVoteOption(option) {
		return types.ErrInvalidVote(keeper.Codespace(), option), ""
	}

	voteFeeStr := ""
	vote := types.Vote{
		ProposalID: proposalID, Voter: voterAddr, Option: option,
	}
	if keeper.ProposalHandlerRouter().HasRoute(proposal.ProposalRoute()) {
		var err sdk.Error
		voteFeeStr, err = keeper.ProposalHandlerRouter().GetRoute(proposal.ProposalRoute()).VoteHandler(ctx, proposal, vote)
		if err != nil {
			return err, ""
		}
	}

	keeper.SetVote(ctx, proposalID, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdkGovTypes.EventTypeProposalVote,
			sdk.NewAttribute(sdkGovTypes.AttributeKeyOption, option.String()),
			sdk.NewAttribute(sdkGovTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil, voteFeeStr
}

// SetVote stores the vote of a specific voter on a specific proposal
func (keeper Keeper) SetVote(ctx sdk.Context, proposalID uint64, vote types.Vote) {
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(vote)
	store.Set(types.VoteKey(proposalID, vote.Voter), bz)
}

// DeleteVotes deletes the votes of a specific proposal
func (keeper Keeper) DeleteVotes(ctx sdk.Context, proposalID uint64) {
	votes := keeper.GetVotes(ctx, proposalID)
	for _, vote := range votes {
		keeper.deleteVote(ctx, vote.ProposalID, vote.Voter)
	}
}

func (keeper Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey())
	store.Delete(types.VoteKey(proposalID, voterAddr))
}
