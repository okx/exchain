package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFarmPools(t *testing.T) {
	tests := []struct {
		owner                   sdk.AccAddress
		name                    string
		lockedSymbol            string
		depositAmount           sdk.DecCoin
		totalValueLocked        sdk.DecCoin
		yieldedTokenInfos       YieldedTokenInfos
		totalAccumulatedRewards sdk.DecCoins
		isFinished              bool
	}{
		{
			owner:            sdk.AccAddress{0x1},
			name:             "pool",
			lockedSymbol:     "xxb",
			depositAmount:    sdk.DecCoin{},
			totalValueLocked: sdk.NewDecCoinFromDec("wwb", sdk.NewDec(100)),
			yieldedTokenInfos: YieldedTokenInfos{
				{
					RemainingAmount: sdk.NewDecCoinFromDec("wwb", sdk.NewDec(100)),
				},
			},
			totalAccumulatedRewards: sdk.DecCoins{},
			isFinished:              false,
		},
		{
			owner:            sdk.AccAddress{0x1},
			name:             "pool",
			lockedSymbol:     "xxb",
			depositAmount:    sdk.DecCoin{},
			totalValueLocked: sdk.NewDecCoinFromDec("wwb", sdk.NewDec(0)),
			yieldedTokenInfos: YieldedTokenInfos{
				{
					RemainingAmount: sdk.NewDecCoinFromDec("wwb", sdk.NewDec(100)),
				},
			},
			totalAccumulatedRewards: sdk.DecCoins{},
			isFinished:              false,
		},
		{
			owner:            sdk.AccAddress{0x1},
			name:             "pool",
			lockedSymbol:     "xxb",
			depositAmount:    sdk.DecCoin{},
			totalValueLocked: sdk.NewDecCoinFromDec("wwb", sdk.NewDec(0)),
			yieldedTokenInfos: YieldedTokenInfos{
				{
					RemainingAmount: sdk.NewDecCoinFromDec("wwb", sdk.NewDec(0)),
				},
			},
			totalAccumulatedRewards: sdk.DecCoins{},
			isFinished:              true,
		},
	}

	for _, test := range tests {
		pool := NewFarmPool(
			test.owner, test.name, sdk.NewDecCoin(test.lockedSymbol, sdk.ZeroInt()), test.depositAmount, test.totalValueLocked,
			test.yieldedTokenInfos, test.totalAccumulatedRewards,
		)
		require.Equal(t, test.isFinished, pool.Finished())
	}
}
