package gov

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli/flags"

	"github.com/okex/okchain/x/gov/keeper"
	"github.com/okex/okchain/x/gov/types"
)

func TestNewHandler(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)
	govHandler := NewHandler(gk)

	res := govHandler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
}

func TestHandleMsgDeposit(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)
	govHandler := NewHandler(gk)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 50)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, initialDeposit, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newDepositMsg := NewMsgDeposit(keeper.Addrs[0], proposalID,
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)})
	res = govHandler(ctx, newDepositMsg)
	require.True(t, res.IsOK())

	// nil address deposit on proposal
	newDepositMsg = NewMsgDeposit(sdk.AccAddress{}, proposalID,
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 1000)})
	res = govHandler(ctx, newDepositMsg)
	require.False(t, res.IsOK())

	// deposit on proposal whose proposal id is 0
	newDepositMsg = NewMsgDeposit(keeper.Addrs[0], 0,
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 1000)})
	res = govHandler(ctx, newDepositMsg)
	require.False(t, res.IsOK())
}

func TestHandleMsgVote(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)
	govHandler := NewHandler(gk)

	proposalCoins := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 500)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, proposalCoins, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	newVoteMsg := NewMsgVote(keeper.Addrs[4], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())

	newVoteMsg = NewMsgVote(keeper.Addrs[4], 0, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.False(t, res.IsOK())

	newVoteMsg = NewMsgVote(sdk.AccAddress{}, proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.False(t, res.IsOK())
}

func TestHandleMsgVote2(t *testing.T) {
	ctx, _, gk, sk, _ := keeper.CreateTestInput(t, false, 100000)
	govHandler := NewHandler(gk)

	proposalCoins := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 500)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, proposalCoins, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	skHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(keeper.Addrs[:2]))
	for i, addr := range keeper.Addrs[:2] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	keeper.CreateValidators(t, skHandler, ctx, valAddrs, []int64{10, 10})
	staking.EndBlocker(ctx, sk)

	newVoteMsg := NewMsgVote(keeper.Addrs[0], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())

	newVoteMsg = NewMsgVote(keeper.Addrs[1], proposalID, types.OptionYes)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())
}

// test distribute deposits after voting
func TestHandleMsgVote3(t *testing.T) {
	ctx, _, gk, sk, _ := keeper.CreateTestInput(t, false, 100000)
	govHandler := NewHandler(gk)

	proposalCoins := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 500)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, proposalCoins, keeper.Addrs[0])
	res := govHandler(ctx, newProposalMsg)
	require.True(t, res.IsOK())
	var proposalID uint64
	gk.Cdc().MustUnmarshalBinaryLengthPrefixed(res.Data, &proposalID)

	ctx = ctx.WithBlockHeight(int64(sk.GetEpoch(ctx)))
	skHandler := staking.NewHandler(sk)
	valAddrs := make([]sdk.ValAddress, len(keeper.Addrs[:2]))
	for i, addr := range keeper.Addrs[:2] {
		valAddrs[i] = sdk.ValAddress(addr)
	}
	keeper.CreateValidators(t, skHandler, ctx, valAddrs, []int64{10, 10})
	staking.EndBlocker(ctx, sk)

	require.Equal(t, proposalCoins, gk.SupplyKeeper().
		GetModuleAccount(ctx, types.ModuleName).GetCoins())
	newVoteMsg := NewMsgVote(keeper.Addrs[0], proposalID, types.OptionNoWithVeto)
	res = govHandler(ctx, newVoteMsg)
	require.True(t, res.IsOK())
	require.Equal(t, sdk.Coins(nil), gk.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName).GetCoins())
}

func TestHandleMsgSubmitProposal(t *testing.T) {
	ctx, _, gk, _, _ := keeper.CreateTestInput(t, false, 1000)
	log, err := flags.ParseLogLevel("*:error", ctx.Logger(), "error")
	require.Nil(t, err)
	ctx = ctx.WithLogger(log)
	handler := NewHandler(gk)

	proposalCoins := sdk.DecCoins{sdk.NewInt64DecCoin("xxx", 500)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg := NewMsgSubmitProposal(content, proposalCoins, keeper.Addrs[0])
	res := handler(ctx, newProposalMsg)
	require.False(t, res.IsOK())

	proposalCoins = sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 500)}
	content = types.NewTextProposal("Test", "description")
	newProposalMsg = NewMsgSubmitProposal(content, proposalCoins, sdk.AccAddress{})
	res = handler(ctx, newProposalMsg)
	require.False(t, res.IsOK())

	//content = tokenTypes.NewDexListProposal("Test", "", keeper.Addrs[0],
	//	"btc-123", common.NativeToken, sdk.NewDec(1000), 0,
	//	4, 4, sdk.NewDec(1))
	//newProposalMsg = NewMsgSubmitProposal(content, proposalCoins, keeper.Addrs[0])
	//res = handler(ctx, newProposalMsg)
	//require.False(t, res.IsOK())
}
