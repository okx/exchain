package keeper

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okexchain/x/gov/types"
)

func TestKeeper_IterateProposals(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	_, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	var findProposal sdkGovTypes.Proposal
	keeper.IterateProposals(ctx, func(proposal sdkGovTypes.Proposal) (stop bool) {
		if proposal.ProposalID == 1 {
			findProposal = proposal
			return true
		}
		return
	})
	require.Equal(t, uint64(1), findProposal.ProposalID)
}

func TestKeeper_IterateWaitingProposalsQueue(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// no proposal
	keeper.InsertWaitingProposalQueue(ctx, 100, 1)
	require.Panics(t, func() {
		keeper.IterateWaitingProposalsQueue(ctx, 100, func(proposal types.Proposal) (stop bool) {
			return
		})
	})

	content := types.NewTextProposal("Test", "description")
	_, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	keeper.InsertWaitingProposalQueue(ctx, 101, 2)
	keeper.InsertWaitingProposalQueue(ctx, 103, 3)

	var proposals sdkGovTypes.Proposals
	keeper.IterateWaitingProposalsQueue(ctx, 103, func(proposal types.Proposal) (stop bool) {
		if proposal.Status == types.StatusDepositPeriod {
			proposals = append(proposals, proposal)
		}
		return
	})
	require.Equal(t, 3, len(proposals))

	proposals = sdkGovTypes.Proposals{}
	keeper.IterateWaitingProposalsQueue(ctx, 103, func(proposal types.Proposal) (stop bool) {
		if proposal.ProposalID == 1 {
			proposals = append(proposals, proposal)
			return true
		}
		return
	})
	require.Equal(t, 1, len(proposals))
}

func TestKeeper_IterateAllWaitingProposals(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// no proposal
	keeper.InsertWaitingProposalQueue(ctx, 100, 1)
	require.Panics(t, func() {
		keeper.IterateWaitingProposalsQueue(ctx, 100, func(proposal types.Proposal) (stop bool) {
			return
		})
	})

	content := types.NewTextProposal("Test", "description")
	_, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	keeper.InsertWaitingProposalQueue(ctx, 101, 2)
	keeper.InsertWaitingProposalQueue(ctx, 103, 3)

	var proposals sdkGovTypes.Proposals
	keeper.IterateAllWaitingProposals(ctx, func(proposal types.Proposal, _, _ uint64) (stop bool) {
		if proposal.Status == types.StatusDepositPeriod {
			proposals = append(proposals, proposal)
		}
		return
	})
	require.Equal(t, 3, len(proposals))

	proposals = sdkGovTypes.Proposals{}
	keeper.IterateAllWaitingProposals(ctx, func(proposal types.Proposal, _, _ uint64) (stop bool) {
		if proposal.ProposalID == 1 {
			proposals = append(proposals, proposal)
			return true
		}
		return
	})
	require.Equal(t, 1, len(proposals))
}

func TestKeeper_IterateActiveProposalsQueue(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	baseTime := time.Now()
	// no proposal
	keeper.InsertActiveProposalQueue(ctx, 1, baseTime)
	require.Panics(t, func() {
		keeper.IterateActiveProposalsQueue(ctx, baseTime.Add(time.Second*1),
			func(proposal types.Proposal) (stop bool) {
				return
			})
	})

	content := types.NewTextProposal("Test", "description")
	_, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	keeper.InsertActiveProposalQueue(ctx, 2, baseTime.Add(time.Second*1))
	keeper.InsertActiveProposalQueue(ctx, 3, baseTime.Add(time.Second*2))

	var proposals sdkGovTypes.Proposals
	keeper.IterateActiveProposalsQueue(ctx, baseTime.Add(time.Second*2), func(proposal types.Proposal) (stop bool) {
		if proposal.Status == types.StatusDepositPeriod {
			proposals = append(proposals, proposal)
		}
		return
	})
	require.Equal(t, 3, len(proposals))

	proposals = sdkGovTypes.Proposals{}
	keeper.IterateActiveProposalsQueue(ctx, baseTime.Add(time.Second*2), func(proposal types.Proposal) (stop bool) {
		if proposal.ProposalID == 1 {
			proposals = append(proposals, proposal)
			return true
		}
		return
	})
	require.Equal(t, 1, len(proposals))
}

func TestKeeper_IterateInactiveProposalsQueue(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	baseTime := time.Now()
	// no proposal
	keeper.InsertInactiveProposalQueue(ctx, 1, baseTime)
	require.Panics(t, func() {
		keeper.IterateInactiveProposalsQueue(ctx, baseTime.Add(time.Second*1),
			func(proposal types.Proposal) (stop bool) {
				return
			})
	})
	keeper.RemoveFromInactiveProposalQueue(ctx, 1, baseTime)

	content := types.NewTextProposal("Test", "description")
	_, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	_, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	var proposals sdkGovTypes.Proposals
	keeper.IterateInactiveProposalsQueue(ctx, baseTime.Add(time.Second*2), func(proposal types.Proposal) (stop bool) {
		if proposal.Status == types.StatusDepositPeriod {
			proposals = append(proposals, proposal)
		}
		return
	})
	require.Equal(t, 3, len(proposals))

	proposals = sdkGovTypes.Proposals{}
	keeper.IterateInactiveProposalsQueue(ctx, baseTime.Add(time.Second*2), func(proposal types.Proposal) (stop bool) {
		if proposal.ProposalID == 1 {
			proposals = append(proposals, proposal)
			return true
		}
		return
	})
	require.Equal(t, 1, len(proposals))
}

func TestKeeper_IterateVotes(t *testing.T) {
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

	var findVote types.Vote
	keeper.IterateVotes(ctx, proposalID, func(vote types.Vote) (stop bool) {
		if vote.Voter.Equals(Addrs[1]) {
			findVote = vote
			return true
		}
		return
	})
	require.Equal(t, Addrs[1], findVote.Voter)
}

func TestKeeper_IterateDeposits(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	err = keeper.AddDeposit(ctx, proposalID, Addrs[1],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	var findDeposit types.Deposit
	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) (stop bool) {
		if deposit.Depositor.Equals(Addrs[1]) {
			findDeposit = deposit
			return true
		}
		return
	})
	require.Equal(t, Addrs[1], findDeposit.Depositor)
}

func TestKeeper_IterateAllDeposits(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	// proposal 1
	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	err = keeper.AddDeposit(ctx, proposalID, Addrs[1],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	// proposal 2
	content = types.NewTextProposal("Test", "description")
	proposal, err = keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID = proposal.ProposalID

	err = keeper.AddDeposit(ctx, proposalID, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	err = keeper.AddDeposit(ctx, proposalID, Addrs[1],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	var findDeposit types.Deposit
	keeper.IterateAllDeposits(ctx, func(deposit types.Deposit) (stop bool) {
		if deposit.ProposalID == 1 && deposit.Depositor.Equals(Addrs[1]) {
			findDeposit = deposit
			return true
		}
		return
	})
	require.Equal(t, Addrs[1], findDeposit.Depositor)
	require.Equal(t, uint64(1), findDeposit.ProposalID)
}

func TestKeeper_GetTallyParams(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	tallyParams := keeper.GetTallyParams(ctx)
	expectedParams := types.TallyParams{
		Quorum:          sdk.NewDecWithPrec(334, 3),
		Threshold:       sdk.NewDecWithPrec(5, 1),
		Veto:            sdk.NewDecWithPrec(334, 3),
		YesInVotePeriod: sdk.NewDecWithPrec(667, 3),
	}
	require.Equal(t, expectedParams, tallyParams)
}

func TestKeeper_RemoveFromWaitingProposalQueue(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID := proposal.ProposalID

	var proposals sdkGovTypes.Proposals
	keeper.InsertWaitingProposalQueue(ctx, 100, proposalID)
	keeper.IterateWaitingProposalsQueue(ctx, 100, func(proposal types.Proposal) (stop bool) {
		proposals = append(proposals, proposal)
		return
	})
	require.Equal(t, 1, len(proposals))

	proposals = sdkGovTypes.Proposals{}
	keeper.RemoveFromWaitingProposalQueue(ctx, 100, proposalID)
	keeper.IterateWaitingProposalsQueue(ctx, 100, func(proposal types.Proposal) (stop bool) {
		proposals = append(proposals, proposal)
		return
	})
	require.Equal(t, 0, len(proposals))
}

func TestKeeper_CheckMsgSubmitProposal(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	content := types.ContentFromProposalType("text", "text", types.ProposalTypeText)

	// not satisfy initial deposit
	amount, err := sdk.ParseDecCoins(fmt.Sprintf("1%s", sdk.DefaultBondDenom))
	require.Nil(t, err)
	msg := types.NewMsgSubmitProposal(content, amount, Addrs[0])
	err = keeper.CheckMsgSubmitProposal(ctx, msg)
	require.NotNil(t, err)
}
