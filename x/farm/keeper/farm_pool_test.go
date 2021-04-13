package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	swaptypes "github.com/okex/exchain/x/ammswap/types"
	"github.com/okex/exchain/x/farm/types"
	"github.com/stretchr/testify/require"
)

type tokenPair struct {
	token0 string
	token1 string
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
			lockedSymbol:  quoteSymbol,
			lockedValue:   sdk.NewDec(100),
			isLPT:         false,
			expectedValue: sdk.NewDec(100),
		},
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
		{
			lockedSymbol:  token2Sym,
			lockedValue:   sdk.ZeroDec(),
			isLPT:         false,
			expectedValue: sdk.ZeroDec(),
		},
	}

	for _, test := range tests {
		pool := types.FarmPool{
			MinLockAmount:    sdk.NewDecCoinFromDec(test.lockedSymbol, sdk.ZeroDec()),
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

func TestIterateAllLockInfos(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	lockInfoList := []types.LockInfo{}
	lockInfoNum := 10
	for i := 0; i < lockInfoNum; i++ {
		keeper.Keeper.SetLockInfo(ctx, types.LockInfo{Owner: Addrs[i], PoolName: fmt.Sprintf("pool%d", i)})
	}
	handler := func(lockInfo types.LockInfo) (stop bool) {
		lockInfoList = append(lockInfoList, lockInfo)
		return false
	}
	keeper.IterateAllLockInfos(ctx, handler)
	require.Equal(t, lockInfoNum, len(lockInfoList))
}

func TestReadWriteFarmPool(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	poolName := "pool"
	_, found := keeper.Keeper.GetFarmPool(ctx, poolName)
	require.False(t, found)
	pool := types.FarmPool{
		Name: poolName,
	}
	keeper.Keeper.SetFarmPool(ctx, pool)
	_, found = keeper.Keeper.GetFarmPool(ctx, poolName)
	require.True(t, found)
	keeper.Keeper.DeleteFarmPool(ctx, poolName)
	_, found = keeper.Keeper.GetFarmPool(ctx, poolName)
	require.False(t, found)
}
