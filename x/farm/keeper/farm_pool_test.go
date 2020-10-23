package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type tokenPair struct {
	token0 string
	token1 string
}

func TestGetPoolLockedValue(t *testing.T) {

}

func TestCalculateLPTValueWithQuote(t *testing.T) {
	ctx, keeper := GetKeeper(t)

	token0Symbol := "xxb"
	token1Symbol := types.DefaultParams().QuoteSymbol

	tokenPairName := swaptypes.GetSwapTokenPairName(token0Symbol, token1Symbol)
	exchange := swaptypes.NewSwapPair(token0Symbol, token1Symbol)
	exchange.QuotePooledCoin.Amount = sdk.NewDec(10000)
	exchange.BasePooledCoin.Amount = sdk.NewDec(10000)

	keeper.swapKeeper.SetSwapTokenPair(ctx, tokenPairName, exchange)

	token0Amount := sdk.NewDecCoinFromDec(token0Symbol, sdk.NewDec(100))
	token1Amount := sdk.NewDecCoinFromDec(token1Symbol, sdk.NewDec(100))

	retValue1 := keeper.calculateLPTValueWithQuote(
		ctx, token0Amount, token1Amount, token1Symbol, swaptypes.DefaultParams(),
	)

	retValue2 := keeper.calculateLPTValueWithQuote(
		ctx, token1Amount, token0Amount, token1Symbol, swaptypes.DefaultParams(),
	)

	require.Equal(t, retValue1, retValue2)
}

func TestCalculateLPTValueWithoutQuote(t *testing.T) {
	ctx, keeper := GetKeeper(t)

	quoteSymbol := types.DefaultParams().QuoteSymbol
	token1Symbol, token2Symbol, _ := initSwapExchange(ctx, keeper, quoteSymbol)

	token1Amount := sdk.NewDecCoinFromDec(token1Symbol, sdk.NewDec(100))
	token2Amount := sdk.NewDecCoinFromDec(token2Symbol, sdk.NewDec(100))

	retValue1 := keeper.calculateLPTValueWithoutQuote(
		ctx, token1Amount, token2Amount, quoteSymbol, swaptypes.DefaultParams(),
	)

	retValue2 := keeper.calculateLPTValueWithoutQuote(
		ctx, token2Amount, token1Amount, quoteSymbol, swaptypes.DefaultParams(),
	)

	require.Equal(t, retValue1, retValue2)
}

func TestGetLockInfo(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	keeper.swapKeeper.SetParams(ctx, swaptypes.DefaultParams())
	quoteSymbol := types.DefaultParams().QuoteSymbol
	token1Sym, token2Sym, _ := initSwapExchange(ctx, keeper, quoteSymbol)

	tests := []struct {
		lockedSymbol  string
		lockedValue   sdk.Dec
		isLPT         bool
		expectedValue sdk.Dec
	}{
		{
			lockedSymbol:  swaptypes.PoolTokenPrefix + swaptypes.GetSwapTokenPairName(token1Sym, quoteSymbol),
			lockedValue:   sdk.NewDec(100),
			isLPT:         true,
			expectedValue: sdk.NewDec(200),
		},
		{
			lockedSymbol:  swaptypes.PoolTokenPrefix + swaptypes.GetSwapTokenPairName(token1Sym, token2Sym),
			lockedValue:   sdk.NewDec(100),
			isLPT:         true,
			expectedValue: sdk.NewDec(200),
		},
		{
			lockedSymbol:  token2Sym,
			lockedValue:   sdk.NewDec(100),
			isLPT:         false,
			expectedValue: sdk.NewDec(100),
		},
	}

	for _, test := range tests {
		pool := types.FarmPool{
			LockedSymbol:     test.lockedSymbol,
			TotalValueLocked: sdk.NewDecCoinFromDec(test.lockedSymbol, test.lockedValue),
		}
		if test.isLPT {
			retValue := keeper.calculateLockedLPTValue(ctx, pool, quoteSymbol, swaptypes.DefaultParams())
			require.Equal(t, test.expectedValue, retValue)
		}
		retValue := keeper.GetPoolLockedValue(ctx, pool)
		require.Equal(t, test.expectedValue, retValue)
	}
}

func initSwapExchange(
	ctx sdk.Context, keeper MockFarmKeeper, quoteSymbol string,
) (string, string, []tokenPair) {
	token1Symbol := "xxb"
	token2Symbol := "okb"

	tokenPairs := []tokenPair{
		{token1Symbol, quoteSymbol},
		{token2Symbol, quoteSymbol},
		{token1Symbol, token2Symbol},
	}

	for _, tokenPair := range tokenPairs {
		tokenPairName := swaptypes.GetSwapTokenPairName(tokenPair.token0, tokenPair.token1)
		exchange := swaptypes.NewSwapPair(tokenPair.token0, tokenPair.token1)
		exchange.QuotePooledCoin.Amount = sdk.NewDec(10000)
		exchange.BasePooledCoin.Amount = sdk.NewDec(10000)
		keeper.swapKeeper.SetSwapTokenPair(ctx, tokenPairName, exchange)
		err := keeper.supplyKeeper.MintCoins(
			ctx, swaptypes.ModuleName,
			sdk.NewDecCoinsFromDec(swaptypes.PoolTokenPrefix+tokenPairName, sdk.NewDec(10000)),
		)
		if err != nil {
			panic("should not happen")
		}
	}
	return token1Symbol, token2Symbol, tokenPairs
}
