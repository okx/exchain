package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/mock"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func initQurierTest(t *testing.T) (*TestInput, mock.AddrKeysSlice, sdk.Context, Keeper, sdk.Querier) {
	mapp, addrSlice := GetTestInput(t, 1)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	keeper := mapp.swapKeeper
	keeper.SetParams(ctx, types.DefaultParams())
	return mapp, addrSlice, ctx, keeper, NewQuerier(mapp.swapKeeper)
}

func TestNewQuerier(t *testing.T) {
	_, _, ctx, keeper, querier := initQurierTest(t)

	// querier with wrong path
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
	tokenPair := types.TestSwapTokenPairName
	swapTokenPair := types.GetTestSwapTokenPair()
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

func TestQueryParams(t *testing.T) {
	_, _, ctx, keeper, querier := initQurierTest(t)

	path0 := []string{types.QueryParams}
	resultBytes, err := querier(ctx, path0, abci.RequestQuery{})
	require.Nil(t, err)
	result := types.Params{}
	keeper.cdc.MustUnmarshalJSON(resultBytes, &result)
	require.Equal(t, types.DefaultParams(), result)
}

func TestQuerySwapTokenPairs(t *testing.T) {
	_, _, ctx, keeper, querier := initQurierTest(t)

	tokenPair := types.TestSwapTokenPairName
	swapTokenPair := types.GetTestSwapTokenPair()
	keeper.SetSwapTokenPair(ctx, tokenPair, swapTokenPair)

	path := []string{types.QuerySwapTokenPairs}
	resultBytes, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	var result []types.SwapTokenPair
	keeper.cdc.MustUnmarshalJSON(resultBytes, &result)
	expectedSwapTokenPairList := []types.SwapTokenPair{swapTokenPair}
	require.Equal(t, expectedSwapTokenPairList, result)
}

func initTestPool(t *testing.T, addrList mock.AddrKeysSlice, mapp *TestInput,
	ctx sdk.Context, keeper Keeper, baseTokenAmount, quoteTokenAmount sdk.SysCoin, poolTokenAmount sdk.Dec) types.SwapTokenPair{
	swapTokenPair := types.SwapTokenPair{
		QuotePooledCoin: quoteTokenAmount,
		BasePooledCoin: baseTokenAmount,
		PoolTokenName: types.GetPoolTokenName(baseTokenAmount.Denom, quoteTokenAmount.Denom),
	}
	keeper.SetSwapTokenPair(ctx, types.GetSwapTokenPairName(baseTokenAmount.Denom, quoteTokenAmount.Denom), swapTokenPair)
	poolToken := types.InitPoolToken(swapTokenPair.PoolTokenName)
	initPoolTokenAmount := sdk.NewDecCoinFromDec(swapTokenPair.PoolTokenName, poolTokenAmount)
	mapp.tokenKeeper.NewToken(ctx, poolToken)
	err := keeper.MintPoolCoinsToUser(ctx, sdk.SysCoins{initPoolTokenAmount}, addrList[0].Address)
	require.Nil(t, err)
	return swapTokenPair
}

func TestQueryRedeemableAssets(t *testing.T) {
	mapp, addrList, ctx, keeper, querier := initQurierTest(t)

	baseTokenAmount := sdk.NewDecCoinFromDec(types.TestBasePooledToken, sdk.NewDec(100))
	quoteTokenAmount := sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.NewDec(100))
	poolTokenAmount := sdk.NewDec(1)
	swapTokenPair := initTestPool(t, addrList, mapp, ctx, keeper, baseTokenAmount, quoteTokenAmount, poolTokenAmount)


	path := []string{types.QueryRedeemableAssets, swapTokenPair.BasePooledCoin.Denom, swapTokenPair.QuotePooledCoin.Denom, poolTokenAmount.String()}
	resultBytes, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	var result []sdk.SysCoin
	keeper.cdc.MustUnmarshalJSON(resultBytes, &result)
	expectedAmountList := []sdk.SysCoin{swapTokenPair.BasePooledCoin, swapTokenPair.QuotePooledCoin}
	require.Equal(t, expectedAmountList, result)
}

func TestQueryBuyAmount(t *testing.T) {
	mapp, addrList, ctx, keeper, querier := initQurierTest(t)

	baseTokenAmount := sdk.NewDecCoinFromDec(types.TestBasePooledToken, sdk.NewDec(100))
	quoteTokenAmount := sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.NewDec(100))
	poolTokenAmount := sdk.NewDec(1)
	swapTokenPair := initTestPool(t, addrList, mapp, ctx, keeper, baseTokenAmount, quoteTokenAmount, poolTokenAmount)

	path := []string{types.QueryBuyAmount}
	var queryParams types.QueryBuyAmountParams
	queryParams.SoldToken = swapTokenPair.QuotePooledCoin
	queryParams.TokenToBuy = swapTokenPair.BasePooledCoin.Denom
	requestBytes := keeper.cdc.MustMarshalJSON(queryParams)
	resultBytes, err := querier(ctx, path, abci.RequestQuery{Data: requestBytes})
	require.Nil(t, err)
	var result string
	keeper.cdc.MustUnmarshalJSON(resultBytes, &result)
	expectedToken := "49.92488733"
	require.Equal(t, expectedToken, result)

	baseTokenAmount2 := sdk.NewDecCoinFromDec(types.TestBasePooledToken2, sdk.NewDec(100))
	quoteTokenAmount2 := sdk.NewDecCoinFromDec(types.TestQuotePooledToken, sdk.NewDec(100))
	poolTokenAmount2 := sdk.NewDec(1)
	swapTokenPair2 := initTestPool(t, addrList, mapp, ctx, keeper, baseTokenAmount2, quoteTokenAmount2, poolTokenAmount2)

	queryParams.SoldToken = swapTokenPair2.BasePooledCoin
	queryParams.TokenToBuy = swapTokenPair.BasePooledCoin.Denom
	requestBytes = keeper.cdc.MustMarshalJSON(queryParams)
	resultBytes, err = querier(ctx, path, abci.RequestQuery{Data: requestBytes})
	require.Nil(t, err)
	keeper.cdc.MustUnmarshalJSON(resultBytes, &result)
	expectedToken = "33.23323333"
	require.Equal(t, expectedToken, result)
}