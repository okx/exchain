package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
	govtypes "github.com/okex/okexchain/x/gov/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCheckMsgSubmitProposal(t *testing.T) {
	ctx, k := GetKeeper(t)

	proposal := govtypes.Proposal{Content: types.NewManageWhiteListProposal(
		"Test",
		"description",
		"pool",
		true,
	)}

	params := types.DefaultParams()
	k.Keeper.SetParams(ctx, params)
	require.Equal(t, sdk.DecCoins(nil), k.GetMinDeposit(ctx, MockContent{}))
	require.Equal(t, params.ManageWhiteListMinDeposit, k.GetMinDeposit(ctx, proposal.Content))

	require.Equal(t, time.Duration(0), k.GetMaxDepositPeriod(ctx, MockContent{}))
	require.Equal(t, params.ManageWhiteListMaxDepositPeriod, k.GetMaxDepositPeriod(ctx, proposal.Content))

	require.Equal(t, time.Duration(0), k.GetVotingPeriod(ctx, MockContent{}))
	require.Equal(t, params.ManageWhiteListVotingPeriod, k.GetVotingPeriod(ctx, proposal.Content))

	require.Error(t, k.CheckMsgSubmitProposal(ctx, govtypes.MsgSubmitProposal{Content: MockContent{}}))
	err := k.CheckMsgSubmitProposal(ctx, govtypes.MsgSubmitProposal{Content: proposal.Content})
	require.Error(t, err)
	require.Equal(t, types.CodeInvalidFarmPool, err.Code())
}

func TestCheckMsgManageWhiteListProposal(t *testing.T) {
	ctx, k := GetKeeper(t)
	quoteSymbol := types.DefaultParams().QuoteSymbol

	proposal := types.NewManageWhiteListProposal(
		"Test",
		"description",
		"pool",
		false,
	)

	err := k.CheckMsgManageWhiteListProposal(ctx, proposal)
	require.Error(t, err)
	require.Equal(t, types.CodePoolNameNotExistedInWhiteList, err.Code())

	k.SetWhitelist(ctx, proposal.PoolName)
	err = k.CheckMsgManageWhiteListProposal(ctx, proposal)
	require.NoError(t, err)

	proposal.IsAdded = true
	err = k.CheckMsgManageWhiteListProposal(ctx, proposal)
	require.Error(t, err)
	require.Equal(t, types.CodeInvalidFarmPool, err.Code())

	lockedSymbol := "xxb"
	pool := types.FarmPool{
		Name: proposal.PoolName,
		LockedSymbol: lockedSymbol,
	}
	k.SetFarmPool(ctx, pool)
	err = k.CheckMsgManageWhiteListProposal(ctx, proposal)
	require.Error(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	SetSwapTokenPair(ctx, k.Keeper, lockedSymbol, quoteSymbol)
	err = k.CheckMsgManageWhiteListProposal(ctx, proposal)
	require.NoError(t, err)
}




