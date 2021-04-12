package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/gov/types"
)

// AddVote adds a vote on a specific proposal
func (keeper Keeper) AddVote(
	ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, option types.VoteOption,
) (sdk.Error, string) {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return types.ErrUnknownProposal(proposalID), ""
	}
	if proposal.Status != types.StatusVotingPeriod {
		return types.ErrInvalidateProposalStatus(), ""
	}

	if !types.ValidVoteOption(option) {
		return types.ErrInvalidVote(option), ""
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
			types.EventTypeProposalVote,
			sdk.NewAttribute(types.AttributeKeyOption, option.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil, voteFeeStr
}

// SetVote stores the vote of a specific voter on a specific proposal
func (keeper Keeper) SetVote(ctx sdk.Context, proposalID uint64, vote types.Vote) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
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
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.VoteKey(proposalID, voterAddr))
}

// GetAllVotes returns all the votes from the store
func (keeper Keeper) GetAllVotes(ctx sdk.Context) (votes types.Votes) {
	keeper.IterateAllVotes(ctx, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVotes returns all the votes from a proposal
func (keeper Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes types.Votes) {
	keeper.IterateVotes(ctx, proposalID, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVote gets the vote from an address on a specific proposal
func (keeper Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (vote types.Vote, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.VoteKey(proposalID, voterAddr))
	if bz == nil {
		return vote, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &vote)
	return vote, true
}

func (keeper Keeper) setVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, vote types.Vote) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(types.VoteKey(proposalID, voterAddr), bz)
}

// GetVotesIterator gets all the votes on a specific proposal as an sdk.Iterator
func (keeper Keeper) GetVotesIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.VotesKey(proposalID))
}
