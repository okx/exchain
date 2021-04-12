package farm

import (
	"testing"

	swaptypes "github.com/okex/exchain/x/ammswap/types"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"
	govtypes "github.com/okex/exchain/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestProposalHandlerPassed(t *testing.T) {
	ctx, k := keeper.GetKeeper(t)

	poolName := "pool1"
	hdlr := NewManageWhiteListProposalHandler(&k.Keeper)
	quoteSymbol := types.DefaultParams().QuoteSymbol

	proposal1 := govtypes.Proposal{Content: keeper.MockContent{}}
	err := hdlr(ctx, &proposal1)
	require.NotNil(t, err)

	proposal2 := govtypes.Proposal{Content: types.NewManageWhiteListProposal(
		"Test",
		"description",
		poolName,
		true,
	)}
	err = hdlr(ctx, &proposal2)
	require.NotNil(t, err)

	pool := types.FarmPool{
		Name:          poolName,
		MinLockAmount: sdk.NewDecCoinFromDec("xxb", sdk.ZeroDec()),
	}
	k.SetFarmPool(ctx, pool)
	err = hdlr(ctx, &proposal2)
	require.NotNil(t, err)

	keeper.SetSwapTokenPair(ctx, k.Keeper, pool.MinLockAmount.Denom, quoteSymbol)
	err = hdlr(ctx, &proposal2)
	require.Nil(t, err)
	require.True(t, inWhiteList(k.GetWhitelist(ctx), pool.Name))

	// test add LPT
	poolName = "pool2"
	baseSymbol := "okb"
	lockSymbol := swaptypes.PoolTokenPrefix + baseSymbol + "_" + quoteSymbol
	pool = types.FarmPool{
		Name:          poolName,
		MinLockAmount: sdk.NewDecCoinFromDec(lockSymbol, sdk.ZeroDec()),
	}
	k.SetFarmPool(ctx, pool)
	proposal3 := govtypes.Proposal{Content: types.NewManageWhiteListProposal(
		"Test",
		"description",
		poolName,
		true,
	)}
	err = hdlr(ctx, &proposal3)
	require.NotNil(t, err)

	keeper.SetSwapTokenPair(ctx, k.Keeper, baseSymbol, quoteSymbol)
	err = hdlr(ctx, &proposal3)
	require.Nil(t, err)
	require.True(t, inWhiteList(k.GetWhitelist(ctx), pool.Name))
}

func inWhiteList(list types.PoolNameList, name string) bool {
	for _, poolName := range list {
		if poolName == name {
			return true
		}
	}
	return false
}

//
//func TestProposalHandlerFailed(t *testing.T) {
//	ctx,k := keeper.GetKeeper(t)
//
//	account := accountKeeper.NewAccountWithAddress(ctx, recipient)
//	require.True(t, account.GetCoins().IsZero())
//	accountKeeper.SetAccount(ctx, account)
//
//	tp := testProposal(recipient, amount)
//	hdlr := NewCommunityPoolSpendProposalHandler(k)
//	require.Error(t, hdlr(ctx, &tp))
//	require.True(t, accountKeeper.GetAccount(ctx, recipient).GetCoins().IsZero())
//}
//
//func TestNewManageWhiteListProposalHandler(t *testing.T) {
//	ctx, k := keeper.GetKeeper(t)
//
//}
