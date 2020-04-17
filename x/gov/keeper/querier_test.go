package keeper

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli/flags"

	"github.com/okex/okchain/x/gov/types"
	"github.com/okex/okchain/x/staking"
)

const custom = "custom"

func getQueriedDepositParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.DepositParams {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryParams, types.ParamDeposit}, "/"),
		Data: []byte{},
	}
	bz, err := querier(ctx, []string{types.QueryParams, types.ParamDeposit}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)
	var params types.DepositParams
	err2 := cdc.UnmarshalJSON(bz, &params)
	require.Nil(t, err2)

	return params
}

func getQueriedVotingParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.VotingParams {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryParams, types.ParamVoting}, "/"),
		Data: []byte{},
	}
	bz, err := querier(ctx, []string{types.QueryParams, types.ParamVoting}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)
	var params types.VotingParams
	err2 := cdc.UnmarshalJSON(bz, &params)
	require.Nil(t, err2)

	return params
}

func getQueriedTallyParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.TallyParams {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryParams, types.ParamTallying}, "/"),
		Data: []byte{},
	}
	bz, err := querier(ctx, []string{types.QueryParams, types.ParamTallying}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)
	var params types.TallyParams
	err2 := cdc.UnmarshalJSON(bz, &params)
	require.Nil(t, err2)

	return params
}

func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier) types.Params {
	params := types.Params{
		VotingParams:  getQueriedVotingParams(t, ctx, cdc, querier),
		TallyParams:   getQueriedTallyParams(t, ctx, cdc, querier),
		DepositParams: getQueriedDepositParams(t, ctx, cdc, querier),
	}
	return params
}

func getQueriedProposals(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, depositor, voter sdk.AccAddress,
	status types.ProposalStatus, limit uint64,
) types.Proposals {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryProposals}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryProposalsParams(status, limit, voter, depositor)),
	}

	bz, err := querier(ctx, []string{types.QueryProposals}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var proposals types.Proposals
	err2 := cdc.UnmarshalJSON(bz, &proposals)
	require.Nil(t, err2)
	return proposals
}

func getQueriedDeposit(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64, depositor sdk.AccAddress,
) types.Deposit {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryDeposit}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryDepositParams(proposalID, depositor)),
	}

	bz, err := querier(ctx, []string{types.QueryDeposit}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var deposit types.Deposit
	err2 := cdc.UnmarshalJSON(bz, &deposit)
	require.Nil(t, err2)
	return deposit
}

func getQueriedDeposits(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64,
) types.Deposits {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryDeposits}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{types.QueryDeposits}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var deposits types.Deposits
	err2 := cdc.UnmarshalJSON(bz, &deposits)
	require.Nil(t, err2)
	return deposits
}

func getQueriedVote(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64, voter sdk.AccAddress,
) types.Vote {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryVote}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryVoteParams(proposalID, voter)),
	}

	bz, err := querier(ctx, []string{types.QueryVote}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var vote types.Vote
	err2 := cdc.UnmarshalJSON(bz, &vote)
	require.Nil(t, err2)
	return vote
}

func getQueriedVotes(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64,
) types.Votes {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryVotes}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{types.QueryVotes}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var votes types.Votes
	err2 := cdc.UnmarshalJSON(bz, &votes)
	require.Nil(t, err2)
	return votes
}

func getQueriedTally(
	t *testing.T, ctx sdk.Context, cdc *codec.Codec, querier sdk.Querier, proposalID uint64,
) types.TallyResult {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryTally}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{types.QueryTally}, query)
	require.Nil(t, err)
	require.NotNil(t, bz)

	var tally types.TallyResult
	err2 := cdc.UnmarshalJSON(bz, &tally)
	require.Nil(t, err2)
	return tally
}

func TestQueries(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	skHandler := staking.NewHandler(sk)
	querier := NewQuerier(keeper)
	cdc := keeper.Cdc()

	validator := sdk.ValAddress(Addrs[4])
	CreateValidators(t, skHandler, ctx, []sdk.ValAddress{validator}, []int64{500})
	staking.EndBlocker(ctx, sk)
	params := getQueriedParams(t, ctx, cdc, querier)

	// submit 3 proposals
	content := types.NewTextProposal("Test", "description")
	proposal1, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID1 := proposal1.ProposalID
	err = keeper.AddDeposit(ctx, proposalID1, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	proposal2, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID2 := proposal2.ProposalID
	err = keeper.AddDeposit(ctx, proposalID2, Addrs[0],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	content = types.NewTextProposal("Test", "description")
	proposal3, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)
	proposalID3 := proposal3.ProposalID
	err = keeper.AddDeposit(ctx, proposalID3, Addrs[1],
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}, "")
	require.Nil(t, err)

	// Addrs[1] deposits on proposals #2 & #3
	err = keeper.AddDeposit(ctx, proposalID2, Addrs[1], params.DepositParams.MinDeposit, "")
	require.Nil(t, err)
	err = keeper.AddDeposit(ctx, proposalID3, Addrs[1], params.DepositParams.MinDeposit, "")
	require.Nil(t, err)

	// only Addrs[0] deposits on proposalID1 initially
	deposits := getQueriedDeposits(t, ctx, cdc, querier, proposalID1)
	require.Len(t, deposits, 1)
	deposit := getQueriedDeposit(t, ctx, cdc, querier, proposalID1, Addrs[0])
	require.Equal(t, deposit, deposits[0])

	// Addrs[0] initially deposits on proposalID2 and Addrs[1] deposits on proposalID2 later
	deposits = getQueriedDeposits(t, ctx, cdc, querier, proposalID2)
	require.Len(t, deposits, 2)
	deposit = getQueriedDeposit(t, ctx, cdc, querier, proposalID2, Addrs[0])
	require.True(t, deposit.Equals(deposits[0]))
	deposit = getQueriedDeposit(t, ctx, cdc, querier, proposalID2, Addrs[1])
	require.True(t, deposit.Equals(deposits[1]))

	// only Addrs[1] deposits on proposalID1 initially
	deposits = getQueriedDeposits(t, ctx, cdc, querier, proposalID3)
	require.Len(t, deposits, 1)
	deposit = getQueriedDeposit(t, ctx, cdc, querier, proposalID3, Addrs[1])
	require.Equal(t, deposit.String(), deposits[0].String())

	// only proposal #1 should be in Deposit Period
	proposals := getQueriedProposals(t, ctx, cdc, querier, nil, nil, types.StatusDepositPeriod, 0)
	require.Len(t, proposals, 1)
	require.Equal(t, proposalID1, proposals[0].ProposalID)

	// proposals #2 and #3 should be in Voting Period
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, nil, types.StatusVotingPeriod, 0)
	require.Len(t, proposals, 2)
	require.Equal(t, proposalID2, proposals[0].ProposalID)
	require.Equal(t, proposalID3, proposals[1].ProposalID)

	// Addrs[0] and Addrs[1] vote on proposal #2
	err, _ = keeper.AddVote(ctx, proposalID2, Addrs[0], types.OptionYes)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposalID2, Addrs[1], types.OptionYes)
	require.Nil(t, err)

	// Addrs[0] and Addrs[1] votes on proposal #3
	err, _ = keeper.AddVote(ctx, proposalID3, Addrs[0], types.OptionYes)
	require.Nil(t, err)
	err, _ = keeper.AddVote(ctx, proposalID3, Addrs[1], types.OptionYes)
	require.Nil(t, err)

	// Test query voted by Addrs[0]
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, Addrs[0], types.StatusNil, 0)
	require.Equal(t, 2, len(proposals))
	require.Equal(t, proposalID2, (proposals[0]).ProposalID)

	// Test query votes on Proposal 2
	votes := getQueriedVotes(t, ctx, cdc, querier, proposalID2)
	require.Len(t, votes, 2)

	// Test query votes on Proposal 3
	votes = getQueriedVotes(t, ctx, cdc, querier, proposalID3)
	require.Len(t, votes, 2)
	require.True(t, Addrs[0].String() == votes[0].Voter.String())
	require.True(t, Addrs[1].String() == votes[1].Voter.String())
	vote := getQueriedVote(t, ctx, cdc, querier, proposalID3, Addrs[0])
	require.Equal(t, vote, votes[0])

	// Test proposals queries with filters

	// Test query all proposals
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, nil, types.StatusNil, 0)
	require.Equal(t, proposalID1, (proposals[0]).ProposalID)
	require.Equal(t, proposalID2, (proposals[1]).ProposalID)
	require.Equal(t, proposalID3, (proposals[2]).ProposalID)

	// Test query voted by Addrs[1]
	proposals = getQueriedProposals(t, ctx, cdc, querier, nil, Addrs[1], types.StatusNil, 0)
	require.Equal(t, proposalID2, (proposals[0]).ProposalID)

	// Test query deposited by Addrs[0]
	proposals = getQueriedProposals(t, ctx, cdc, querier, Addrs[0], nil, types.StatusNil, 0)
	require.Equal(t, proposalID1, (proposals[0]).ProposalID)

	// Test query deposited by Addrs[1]
	proposals = getQueriedProposals(t, ctx, cdc, querier, Addrs[1], nil, types.StatusNil, 0)
	require.Len(t, proposals, 2)
	require.Equal(t, proposalID3, (proposals[1]).ProposalID)

	// Test Tally Query
	status, dist, tallyResults := Tally(ctx, keeper, proposal2, true)
	require.True(t, dist)
	require.Equal(t, types.StatusRejected, status)
	proposal2.FinalTallyResult = tallyResults
	keeper.SetProposal(ctx, proposal2)
	tally := getQueriedTally(t, ctx, cdc, querier, proposalID2)
	require.Equal(t, tallyResults, tally)

	bz, err := querier(ctx, []string{""}, abci.RequestQuery{})
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryTally(t *testing.T) {
	ctx, _, keeper, sk, _ := CreateTestInput(t, false, 100000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)
	cdc := keeper.Cdc()

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

	err, _ = keeper.AddVote(ctx, proposal.ProposalID, Addrs[0], types.OptionYes)
	require.Nil(t, err)

	// no query params
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryTally}, "/"),
		Data: nil,
	}
	bz, err := querier(ctx, []string{types.QueryTally}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)

	// query proposal whose ID is 0
	query = abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryTally}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryProposalParams(0)),
	}
	bz, err = querier(ctx, []string{types.QueryTally}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)

	expectedTally := newTallyResult(t, "1", "1", "0.0", "0.0", "0.0", "2")
	require.Equal(t, expectedTally, getQueriedTally(t, ctx, cdc, querier, proposal.ProposalID))

	// proposal passed
	proposal.Status = types.StatusPassed
	proposal.FinalTallyResult = expectedTally
	keeper.SetProposal(ctx, proposal)
	require.Equal(t, expectedTally, getQueriedTally(t, ctx, cdc, querier, proposal.ProposalID))
}

func TestQueryParams(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryParams, "test"}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{types.QueryParams, "test"}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryVotes(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryVotes}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QueryVotes}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryVote(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryVote}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QueryVote}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryDeposits(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryDeposits}, "/"),
		Data: nil,
	}

	bz, err := querier(ctx, []string{types.QueryDeposits}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryDeposit(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryDeposit}, "/"),
		Data: nil,
	}
	bz, err := querier(ctx, []string{types.QueryDeposit}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestQueryProposal(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)
	cdc := keeper.Cdc()

	content := types.NewTextProposal("Test", "description")
	proposal1, err := keeper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	// no query params
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryProposal}, "/"),
		Data: nil,
	}
	bz, err := querier(ctx, []string{types.QueryProposal}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)

	// query proposal whose ID is 0
	query = abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryProposal}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryProposalParams(0)),
	}
	bz, err = querier(ctx, []string{types.QueryProposal}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)

	query = abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryProposal}, "/"),
		Data: cdc.MustMarshalJSON(types.NewQueryProposalParams(1)),
	}
	var queriedProposal types.Proposal
	bz, err = querier(ctx, []string{types.QueryProposal}, query)
	require.Nil(t, err)
	cdc.MustUnmarshalJSON(bz, &queriedProposal)
	require.Equal(t, proposal1.ProposalID, queriedProposal.ProposalID)
}

func TestQueryProposals(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	querier := NewQuerier(keeper)

	// no query params
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, types.QuerierRoute, types.QueryProposals}, "/"),
		Data: nil,
	}
	bz, err := querier(ctx, []string{types.QueryProposals}, query)
	require.NotNil(t, err)
	require.Nil(t, bz)
}
