package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGov "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/stretchr/testify/require"

	"github.com/okex/okchain/x/gov/types"
)

func TestKeeper_AddVote(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// vote on proposal which is not exist
	err, votefee := keeper.AddVote(ctx, 0, Addrs[0], types.OptionYes)
	require.NotNil(t, err)
	require.Equal(t, "", votefee)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	// nil address deposit
	err, votefee = keeper.AddVote(ctx, 0, sdk.AccAddress{}, types.OptionYes)
	require.NotNil(t, err)
	require.Equal(t, "", votefee)

	// vote on proposal whose status is not VotingPeriod
	proposal.Status = types.StatusPassed
	keeper.SetProposal(ctx, proposal)
	err, votefee = keeper.AddVote(ctx, proposalID, Addrs[0], types.OptionYes)
	require.NotNil(t, err)
	require.Equal(t, "", votefee)

	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	// vote invalid vote option
	err, votefee = keeper.AddVote(ctx, proposalID, Addrs[0], types.OptionEmpty)
	require.NotNil(t, err)
	require.Equal(t, "", votefee)

	// change old vote and activate proposal
	err, votefee = keeper.AddVote(ctx, proposalID, Addrs[0], types.OptionYes)
	require.Nil(t, err)
	require.Equal(t, "", votefee)
	vote, ok := keeper.GetVote(ctx, proposalID, Addrs[0])
	sdkGovVote := sdkGov.Vote{ProposalID: proposalID, Voter: Addrs[0], Option: types.OptionYes}
	expectedVote := types.Vote{Vote: sdkGovVote, VoteID: 0}
	require.True(t, ok)
	require.Equal(t, expectedVote, vote)

	err, votefee = keeper.AddVote(ctx, proposalID, Addrs[0], types.OptionNo)
	require.Nil(t, err)
	require.Equal(t, "", votefee)
	vote, ok = keeper.GetVote(ctx, proposalID, Addrs[0])
	sdkGovVote = sdkGov.Vote{ProposalID: proposalID, Voter: Addrs[0], Option: types.OptionNo}
	expectedVote = types.Vote{Vote: sdkGovVote, VoteID: 0}
	require.True(t, ok)
	require.Equal(t, expectedVote, vote)
}

func TestKeeper_GetVote(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 150)}, "")
	require.Nil(t, err)

	err, voteFee := keeper.AddVote(ctx, proposalID, Addrs[0], types.OptionYes)
	fmt.Println(err)
	require.Nil(t, err)
	require.Equal(t, "", voteFee)

	expectedVote := types.Vote{
		Vote: sdkGov.Vote{ProposalID: proposalID, Voter: Addrs[0],
			Option: types.OptionYes}, VoteID: 0,
	}
	vote, found := keeper.GetVote(ctx, proposalID, Addrs[0])
	require.True(t, found)
	require.True(t, vote.Equals(expectedVote))

	// get vote from db
	vote, found = keeper.GetVote(ctx, proposalID, Addrs[0])
	require.True(t, found)
	require.True(t, vote.Equals(expectedVote))
}

func TestKeeper_GetVotes(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 150)}, "")
	require.Nil(t, err)

	err, voteFee := keeper.AddVote(ctx, proposalID, Addrs[1], types.OptionYes)
	require.Nil(t, err)
	require.Equal(t, "", voteFee)

	err, voteFee = keeper.AddVote(ctx, proposalID, Addrs[2], types.OptionNo)
	require.Nil(t, err)
	require.Equal(t, "", voteFee)

	expectedVotes := types.Votes{
		{
			Vote: sdkGov.Vote{
				ProposalID: proposalID,
				Voter:      Addrs[1],
				Option:     types.OptionYes,
			},
			VoteID: 0,
		},
		{
			Vote: sdkGov.Vote{
				ProposalID: proposalID,
				Voter:      Addrs[2],
				Option:     types.OptionNo,
			},
			VoteID: 1,
		},
	}
	votes := keeper.GetVotes(ctx, proposalID)
	require.Equal(t, expectedVotes, votes)

	// get votes from db
	votes = keeper.GetVotes(ctx, proposalID)
	require.Equal(t, expectedVotes, votes)
}

func TestKeeper_DeleteVotes(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 150)}, "")
	require.Nil(t, err)

	err, voteFee := keeper.AddVote(ctx, proposalID, Addrs[1], types.OptionYes)
	require.Nil(t, err)
	require.Equal(t, "", voteFee)

	err, voteFee = keeper.AddVote(ctx, proposalID, Addrs[2], types.OptionNo)
	require.Nil(t, err)
	require.Equal(t, "", voteFee)

	votes := keeper.GetVotes(ctx, proposalID)
	require.Equal(t, 2, len(votes))
	keeper.DeleteVotes(ctx, proposalID)
	votes = keeper.GetVotes(ctx, proposalID)
	require.Equal(t, 0, len(votes))

	votes = keeper.GetVotes(ctx, proposalID)
	require.Equal(t, 0, len(votes))
}
