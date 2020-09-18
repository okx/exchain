package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/common"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	mapp, _ := GetTestInput(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	// querier with wrong path
	querier := NewQuerier(keeper)
	path0 := []string{"any", types.TestBasePooledToken}
	tokenpair, err := querier(ctx, path0, abci.RequestQuery{})
	require.NotNil(t, err)
	require.Nil(t, tokenpair)

	// querier with wrong token
	path := []string{types.QuerySwapTokenPair, types.TestBasePooledToken, types.TestQuotePooledToken}
	tokenpair, err = querier(ctx, path, abci.RequestQuery{})
	require.NotNil(t, err)
	require.Nil(t, tokenpair)

	// add new tokenpair and querier
	tokenPair := types.GetSwapTokenPairName(types.TestBasePooledToken, types.TestQuotePooledToken)
	swapTokenPair := initTokenPair(types.TestBasePooledToken)
	keeper.SetSwapTokenPair(ctx, tokenPair, swapTokenPair)
	tokenpair, err = querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	require.NotNil(t, tokenpair)

	// check the value
	result := &types.SwapTokenPair{}
	keeper.cdc.MustUnmarshalJSON(tokenpair, result)
	require.EqualValues(t, result.BasePooledCoin.Denom, types.TestBasePooledToken)

	// delete tokenpair and querier
	keeper.DeleteSwapTokenPair(ctx, tokenPair)
	tokenpair, err = querier(ctx, path, abci.RequestQuery{})
	require.NotNil(t, err)
	require.Nil(t, tokenpair)
}

func initTokenPair(token string) types.SwapTokenPair {
	poolName := types.GetPoolTokenName(token, types.TestQuotePooledToken)
	baseToken := sdk.NewDecCoinFromDec(token, sdk.ZeroDec())
	quoteToken := sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.ZeroDec())

	swapTokenPair := types.SwapTokenPair{
		BasePooledCoin:  baseToken,
		QuotePooledCoin: quoteToken,
		PoolTokenName:   poolName,
	}

	return swapTokenPair
}
