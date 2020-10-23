package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInvariants(t *testing.T) {
	ctx, keeper := GetKeeper(t)
	keeper.swapKeeper.SetParams(ctx, swaptypes.DefaultParams())
	quoteSymbol := types.DefaultParams().QuoteSymbol
	token1Sym, token2Sym := initPoolsAndLockInfos(ctx, keeper)

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
