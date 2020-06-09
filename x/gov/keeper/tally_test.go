package keeper

import (
	"testing"
	"time"

	"github.com/okex/okchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okchain/x/gov/types"
	"github.com/okex/okchain/x/staking"
)

func newTallyResult(t *testing.T, totalVoted, yes, abstain, no, veto, totalVoting string) types.TallyResult {
	decTotalVoting, err := sdk.NewDecFromStr(totalVoting)
	require.Nil(t, err)
	decTotalVoted, err := sdk.NewDecFromStr(totalVoted)
	require.Nil(t, err)
	decYes, err := sdk.NewDecFromStr(yes)
	require.Nil(t, err)
	decAbstain, err := sdk.NewDecFromStr(abstain)
	require.Nil(t, err)
	decNo, err := sdk.NewDecFromStr(no)
	require.Nil(t, err)
	decNoWithVeto, err := sdk.NewDecFromStr(veto)
	require.Nil(t, err)
	return types.TallyResult{
		TotalPower:      decTotalVoting,
		TotalVotedPower: decTotalVoted,
		Yes:             decYes,
		Abstain:         decAbstain,
		No:              decNo,
		NoWithVeto:      decNoWithVeto,
	}
}

func TestTallyNoBondedTokens(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	// less quorum when expire VotingPeriod
	status, dist, tallyResults := Tally(ctx, keeper, proposal, true)
	require.False(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(types.EmptyTallyResult(keeper.totalPower(ctx))))

	// less quorum when in VotingPeriod
	status, dist, tallyResults = Tally(ctx, keeper, proposal, false)
	require.False(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(types.EmptyTallyResult(keeper.totalPower(ctx))))
}

func TestTallyNoOneVotes(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	stakingHandler := staking.NewHandler(sk)

	valAddrs := make([]sdk.ValAddress, len(Addrs[:2]))
	for i, addr := range Addrs[:2] {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	CreateValidators(t, stakingHandler, ctx, valAddrs, []int64{5, 5})
	staking.EndBlocker(ctx, sk)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	// less quorum when expire VotingPeriod
	status, dist, tallyResults := Tally(ctx, keeper, proposal, true)
	require.True(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(types.EmptyTallyResult(keeper.totalPower(ctx))))

	// less quorum when in VotingPeriod
	status, dist, tallyResults = Tally(ctx, keeper, proposal, false)
	require.False(t, dist)
	require.Equal(t, types.StatusVotingPeriod, status)
	require.True(t, tallyResults.Equals(types.EmptyTallyResult(keeper.totalPower(ctx))))
}

func TestTallyAllValidatorsVoteAbstain(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	stakingHandler := staking.NewHandler(sk)

	valAddrs := make([]sdk.ValAddress, len(Addrs[:2]))
	for i, addr := range Addrs[:2] {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	CreateValidators(t, stakingHandler, ctx, valAddrs, []int64{5, 5})
	staking.EndBlocker(ctx, sk)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[0], types.OptionAbstain)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[1], types.OptionAbstain)
	require.Nil(t, err)

	expectedTallyResult := newTallyResult(t, "2", "0.0", "2", "0.0", "0.0", "2")
	// when expire VotingPeriod
	status, dist, tallyResults := Tally(ctx, keeper, proposal, true)
	require.False(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))

	// when in VotingPeriod
	status, dist, tallyResults = Tally(ctx, keeper, proposal, false)
	require.False(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))
}

// test more than one third validator vote veto, in this test there are two validators
// and one vote veto.
func TestTallyAllValidatorsMoreThanOneThirdVeto(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	stakingHandler := staking.NewHandler(sk)

	valAddrs := make([]sdk.ValAddress, len(Addrs[:2]))
	for i, addr := range Addrs[:2] {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	CreateValidators(t, stakingHandler, ctx, valAddrs, []int64{5, 5})
	staking.EndBlocker(ctx, sk)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[0], types.OptionNoWithVeto)
	require.Nil(t, err)

	expectedTallyResult := newTallyResult(t, "1", "0.0", "0.0", "0.0", "1", "2")
	// when expire VotingPeriod
	status, dist, tallyResults := Tally(ctx, keeper, proposal, true)
	require.True(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))

	// when in VotingPeriod
	status, dist, tallyResults = Tally(ctx, keeper, proposal, false)
	require.True(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))
}

func TestTallyOtherCase(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)
	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	stakingHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(Addrs[:2]))
	for i, addr := range Addrs[:2] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	CreateValidators(t, stakingHandler, ctx, valAddrs, []int64{5, 5})
	staking.EndBlocker(ctx, sk)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	// one of two validators vote no, that is more than or equal to 1/2 of non-abstain vote not Yes
	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[0], types.OptionNo)
	require.Nil(t, err)

	expectedTallyResult := newTallyResult(t, "1", "0.0", "0.0", "1", "0.0", "2")
	// when expire VotingPeriod
	status, dist, tallyResults := Tally(ctx, keeper, proposal, true)
	require.False(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))

	// when in VotingPeriod
	status, dist, tallyResults = Tally(ctx, keeper, proposal, false)
	require.False(t, dist)
	require.Equal(t, types.StatusRejected, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))

	// all validators vote yes, that is more than to 1/2 of non-abstain vote Yes when expire VotingPeriod
	// and more than 2/3 of totalBonded vote Yes when in VotingPeriod
	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[0], types.OptionYes)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[1], types.OptionYes)
	require.Nil(t, err)

	expectedTallyResult = newTallyResult(t, "2", "2", "0.0", "0.0", "0.0", "2")
	// when expire VotingPeriod
	status, dist, tallyResults = Tally(ctx, keeper, proposal, true)
	require.False(t, dist)
	require.Equal(t, types.StatusPassed, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))

	// when in VotingPeriod
	status, dist, tallyResults = Tally(ctx, keeper, proposal, false)
	require.False(t, dist)
	require.Equal(t, types.StatusPassed, status)
	require.True(t, tallyResults.Equals(expectedTallyResult))
}

func TestTallyDelegatorInherit(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)
	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	ctx = ctx.WithBlockTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	stakingHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(Addrs[:3]))
	for i, addr := range Addrs[:3] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	CreateValidators(t, stakingHandler, ctx, valAddrs, []int64{5, 5, 5})
	staking.EndBlocker(ctx, sk)

	coin, err := sdk.ParseDecCoin("11000" + common.NativeToken)
	require.Nil(t, err)
	delegator1Msg := staking.NewMsgDeposit(Addrs[3], coin)
	stakingHandler(ctx, delegator1Msg)

	addSharesMsg := staking.NewMsgAddShares(Addrs[3], []sdk.ValAddress{sdk.ValAddress(Addrs[2])})
	stakingHandler(ctx, addSharesMsg)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)

	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[0], types.OptionNo)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[1], types.OptionNo)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[2], types.OptionYes)
	require.Nil(t, err)

	// there are 3 validators with 1 voting power for each one (0.001okt -> 1 power)
	//  2 vals -> OptionNo
	//  1 val -> OptionYes
	expectedTallyResult := newTallyResult(t, "11003", "11001", "0.0", "2", "0.0", "11003")
	status, dist, tallyResults := Tally(ctx, keeper, proposal, true)
	require.False(t, dist)
	require.Equal(t, types.StatusPassed, status)
	require.Equal(t, expectedTallyResult, tallyResults)
}

func TestTallyDelegatorOverride(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)
	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	ctx = ctx.WithBlockTime(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	stakingHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(Addrs[:3]))
	for i, addr := range Addrs[:3] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	CreateValidators(t, stakingHandler, ctx, valAddrs, []int64{5, 5, 5})
	staking.EndBlocker(ctx, sk)

	coin, err := sdk.ParseDecCoin("1" + common.NativeToken)
	require.Nil(t, err)
	delegator1Msg := staking.NewMsgDeposit(Addrs[3], coin)
	stakingHandler(ctx, delegator1Msg)

	addSharesMsg := staking.NewMsgAddShares(Addrs[3], []sdk.ValAddress{sdk.ValAddress(Addrs[2])})
	stakingHandler(ctx, addSharesMsg)

	content := types.NewTextProposal("Test", "description")
	proposal, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposal.Status = types.StatusVotingPeriod
	keeper.SetProposal(ctx, proposal)
	proposalID := proposal.ProposalID

	err, _ = keeper.AddVote(ctx, proposalID, Addrs[0], types.OptionYes)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposalID, Addrs[1], types.OptionYes)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposalID, Addrs[2], types.OptionYes)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposalID, Addrs[3], types.OptionNo)
	require.Nil(t, err)

	expectedTallyResult := newTallyResult(t, "4", "3", "0.0", "1", "0.0", "4")
	status, dist, tallyResults := Tally(ctx, keeper, proposal, true)
	require.False(t, dist)
	require.Equal(t, types.StatusPassed, status)
	require.Equal(t, expectedTallyResult, tallyResults)
}
