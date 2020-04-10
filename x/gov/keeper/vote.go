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
		Vote: sdkGovTypes.Vote{ProposalID: proposalID, Voter: voterAddr, Option: option},
	}
	if keeper.ProposalHandlerRouter().HasRoute(proposal.ProposalRoute()) {
		var err sdk.Error
		voteFeeStr, err = keeper.ProposalHandlerRouter().GetRoute(proposal.ProposalRoute()).VoteHandler(ctx, proposal, vote)
		if err != nil {
			return err, ""
		}
	}

	voteOld, found := keeper.GetVote(ctx, proposalID, voterAddr)
	if found {
		keeper.changeOldVote(ctx, proposalID, voteOld.VoteID, vote)
	} else {
		keeper.SetVote(ctx, proposalID, vote)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdkGovTypes.EventTypeProposalVote,
			sdk.NewAttribute(sdkGovTypes.AttributeKeyOption, option.String()),
			sdk.NewAttribute(sdkGovTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil, voteFeeStr
}

// GetVotes gets all the votes on a specific proposal
func (keeper Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes types.Votes) {
	var voteNum, i uint64
	voteNum = keeper.getProposalVoteCnt(ctx, proposalID)
	store := ctx.KVStore(keeper.StoreKey())
	for i = 0; i < voteNum; i++ {
		var vote types.Vote
		bz := store.Get(types.VoteKey(proposalID, i))
		if bz == nil {
			continue
		}
		keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &vote)
		votes = append(votes, vote)
	}
	return votes
}

// GetVote gets the vote of a specific voter on a specific proposal
func (keeper Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.Address) (types.Vote, bool) {
	var voteNum, i uint64
	voteNum = keeper.getProposalVoteCnt(ctx, proposalID)
	store := ctx.KVStore(keeper.StoreKey())
	for i = 0; i < voteNum; i++ {
		var vote types.Vote
		bz := store.Get(types.VoteKey(proposalID, i))
		if bz == nil {
			continue
		}
		keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &vote)
		if vote.Voter.Equals(voterAddr) {
			return vote, true
		}
	}
	return types.Vote{}, false
}

// SetVote stores the vote of a specific voter on a specific proposal
func (keeper Keeper) SetVote(ctx sdk.Context, proposalID uint64, vote types.Vote) {
	voteID := keeper.getProposalVoteCnt(ctx, proposalID)
	vote.VoteID = voteID
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(vote)
	store.Set(types.VoteKey(proposalID, voteID), bz)
	keeper.setProposalVoteCnt(ctx, proposalID)
}

func (keeper Keeper) changeOldVote(ctx sdk.Context, proposalID, voteID uint64, vote types.Vote) {
	vote.VoteID = voteID
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(vote)
	store.Set(types.VoteKey(proposalID, voteID), bz)
}

// setProposalDepositCnt save new count of vote for a specific proposalID
func (keeper Keeper) setProposalVoteCnt(ctx sdk.Context, proposalID uint64) {
	cnt := keeper.getProposalVoteCnt(ctx, proposalID)
	store := ctx.KVStore(keeper.StoreKey())
	bz := keeper.Cdc().MustMarshalBinaryLengthPrefixed(cnt + 1)
	store.Set(types.VoteCntKey(proposalID), bz)
}

func (keeper Keeper) getProposalVoteCnt(ctx sdk.Context, proposalID uint64) uint64 {
	var cnt uint64
	store := ctx.KVStore(keeper.StoreKey())
	bz := store.Get(types.VoteCntKey(proposalID))
	if bz == nil {
		bz = keeper.Cdc().MustMarshalBinaryLengthPrefixed(0)
		store.Set(types.VoteCntKey(proposalID), bz)
		return 0
	}
	keeper.Cdc().MustUnmarshalBinaryLengthPrefixed(bz, &cnt)
	return cnt
}

// DeleteVotes deletes the votes of a specific proposal
func (keeper Keeper) DeleteVotes(ctx sdk.Context, proposalID uint64) {
	votes := keeper.GetVotes(ctx, proposalID)
	for _, vote := range votes {
		keeper.deleteVote(ctx, vote.ProposalID, vote.VoteID)
	}
	keeper.deleteVoteCnt(ctx, proposalID)
}

func (keeper Keeper) deleteVote(ctx sdk.Context, proposalID, voteID uint64) {
	store := ctx.KVStore(keeper.StoreKey())
	store.Delete(types.VoteKey(proposalID, voteID))
}

func (keeper Keeper) deleteVoteCnt(ctx sdk.Context, proposalID uint64) {
	if votes := keeper.GetVotes(ctx, proposalID); len(votes) == 0 {
		store := ctx.KVStore(keeper.StoreKey())
		store.Delete(types.VoteCntKey(proposalID))
	}
}
