package gov

import (
	"testing"
	"time"

	"github.com/okex/okexchain/x/gov/keeper"
	"github.com/okex/okexchain/x/gov/types"
	"github.com/okex/okexchain/x/params"
	paramsTypes "github.com/okex/okexchain/x/params/types"
	"github.com/okex/okexchain/x/staking"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func newTextProposal(t *testing.T, ctx sdk.Context, initialDeposit sdk.DecCoins, govHandler sdk.Handler) sdk.Result {
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, initialDeposit, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	return res
}

func TestTickPassedVotingPeriod(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)
	govHandler := NewHandler(gk)

	inactiveQueue := gk.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	activeQueue := gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()

	proposalCoins := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 500)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, proposalCoins, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
	ctx = ctx.WithBlockHeader(newHeader)

	newDepositMsg := NewMsgDeposit(keeper.Addrs[1], proposalID, proposalCoins)
	res = govHandler(ctx, newDepositMsg)
	require.False(t, res.IsOK())

	newHeader = ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(gk.GetDepositParams(ctx).MaxDepositPeriod).
		Add(gk.GetVotingParams(ctx).VotingPeriod)
	ctx = ctx.WithBlockHeader(newHeader)

	inactiveQueue = gk.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()

	activeQueue = gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, activeQueue.Valid())
	var activeProposalID uint64
	err := gk.Cdc().UnmarshalBinaryLengthPrefixed(activeQueue.Value(), &activeProposalID)
	require.Nil(t, err)
	proposal, ok := gk.GetProposal(ctx, activeProposalID)
	require.True(t, ok)
	require.Equal(t, StatusVotingPeriod, proposal.Status)
	depositsIterator := gk.GetDeposits(ctx, proposalID)
	require.NotEqual(t, depositsIterator, []Deposit{})
	activeQueue.Close()

	EndBlocker(ctx, gk)

	activeQueue = gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()
}

// test deposit is not enough when expire max deposit period
func TestEndBlockerIterateInactiveProposalsQueue(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 10)}
	newTextProposal(t, ctx, initialDeposit, NewHandler(gk))

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(gk.GetMaxDepositPeriod(ctx, nil))
	ctx = ctx.WithBlockHeader(newHeader)
	inactiveQueue := gk.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, inactiveQueue.Valid())
	inactiveQueue.Close()
	EndBlocker(ctx, gk)
	inactiveQueue = gk.InactiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, inactiveQueue.Valid())
	inactiveQueue.Close()
}

func TestEndBlockerIterateActiveProposalsQueue1(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 150)}
	newTextProposal(t, ctx, initialDeposit, NewHandler(gk))

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(gk.GetVotingPeriod(ctx, nil))
	ctx = ctx.WithBlockHeader(newHeader)
	activeQueue := gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, activeQueue.Valid())
	activeQueue.Close()
	EndBlocker(ctx, gk)
	activeQueue = gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()
}

// test distribute
func TestEndBlockerIterateActiveProposalsQueue2(t *testing.T) {
	ctx, _, gk, sk, _ := keeper.CreateTestInput(t, false, 100000)
	govHandler := NewHandler(gk)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	skHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(keeper.Addrs[:3]))
	for i, addr := range keeper.Addrs[:3] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	keeper.CreateValidators(t, skHandler, ctx, valAddrs, []int64{10, 10, 10})
	staking.EndBlocker(ctx, sk)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 150)}
	res := newTextProposal(t, ctx, initialDeposit, NewHandler(gk))

	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	require.Equal(t, initialDeposit, gk.SupplyKeeper().
		GetModuleAccount(ctx, types.ModuleName).GetCoins())
	newVoteMsg := NewMsgVote(keeper.Addrs[0], proposalID, types.OptionNoWithVeto)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(gk.GetVotingPeriod(ctx, nil))
	ctx = ctx.WithBlockHeader(newHeader)
	activeQueue := gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, activeQueue.Valid())
	activeQueue.Close()
	EndBlocker(ctx, gk)
	activeQueue = gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()

	require.Equal(t, sdk.Coins(nil), gk.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName).GetCoins())
}

// test passed
func TestEndBlockerIterateActiveProposalsQueue3(t *testing.T) {
	ctx, _, gk, sk, _ := keeper.CreateTestInput(t, false, 100000)
	govHandler := NewHandler(gk)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	skHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(keeper.Addrs[:4]))
	for i, addr := range keeper.Addrs[:4] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	keeper.CreateValidators(t, skHandler, ctx, valAddrs, []int64{10, 10, 10, 10})
	staking.EndBlocker(ctx, sk)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 150)}
	res := newTextProposal(t, ctx, initialDeposit, NewHandler(gk))
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	require.Equal(t, initialDeposit, gk.SupplyKeeper().
		GetModuleAccount(ctx, types.ModuleName).GetCoins())
	newVoteMsg := NewMsgVote(keeper.Addrs[0], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())
	newVoteMsg = NewMsgVote(keeper.Addrs[1], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())

	newHeader := ctx.BlockHeader()
	newHeader.Time = ctx.BlockHeader().Time.Add(gk.GetVotingPeriod(ctx, nil))
	ctx = ctx.WithBlockHeader(newHeader)
	activeQueue := gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.True(t, activeQueue.Valid())
	activeQueue.Close()
	EndBlocker(ctx, gk)
	activeQueue = gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, activeQueue.Valid())
	activeQueue.Close()

	require.Equal(t, sdk.Coins(nil), gk.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName).GetCoins())
}

func TestEndBlockerIterateWaitingProposalsQueue(t *testing.T) {
	ctx, _, gk, sk, _ := keeper.CreateTestInput(t, false, 100000)
	govHandler := NewHandler(gk)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	skHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(keeper.Addrs[:4]))
	for i, addr := range keeper.Addrs[:4] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	keeper.CreateValidators(t, skHandler, ctx, valAddrs, []int64{10, 10, 10, 10})
	staking.EndBlocker(ctx, sk)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 150)}
	paramsChanges := []params.ParamChange{{Subspace: "staking", Key: "MaxValidators", Value: "105"}}
	height := uint64(ctx.BlockHeight() + 1000)
	content := paramsTypes.NewParameterChangeProposal("Test", "", paramsChanges, height)
	newProposalMsg := NewMsgSubmitProposal(content, initialDeposit, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newVoteMsg := NewMsgVote(keeper.Addrs[0], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())
	newVoteMsg = NewMsgVote(keeper.Addrs[1], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())
	newVoteMsg = NewMsgVote(keeper.Addrs[2], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())

	ctx = ctx.WithBlockHeight(int64(height))
	waitingQueue := gk.WaitingProposalQueueIterator(ctx, uint64(ctx.BlockHeight()))
	require.True(t, waitingQueue.Valid())
	waitingQueue.Close()
	EndBlocker(ctx, gk)
	waitingQueue = gk.ActiveProposalQueueIterator(ctx, ctx.BlockHeader().Time)
	require.False(t, waitingQueue.Valid())
	waitingQueue.Close()
}
