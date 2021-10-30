package farm

import (
	"testing"

	"github.com/stretchr/testify/require"

	swap "github.com/okex/exchain/x/ammswap/types"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func TestBeginBlocker(t *testing.T) {
	// init
	ctx, mk := keeper.GetKeeper(t)
	k := mk.Keeper
	farmParams := types.DefaultParams()
	farmParams.YieldNativeToken = true
	k.SetParams(ctx, farmParams)

	tests := []struct {
		name string
		run  func(ctx sdk.Context, k keeper.Keeper)
	}{
		{
			name: "empty MintFramingAccount",
			run: func(ctx sdk.Context, k keeper.Keeper) {
				require.NotPanics(t, func() {
					BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: 1}}, k)
				})
			},
		},
		{
			name: "no pools",
			run: func(ctx sdk.Context, k keeper.Keeper) {
				// mint native token
				coins := sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, sdk.NewDec(10000))
				k.SupplyKeeper().MintCoins(ctx, MintFarmingAccount, coins)
				require.NotPanics(t, func() {
					BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: 1}}, k)
				})
			},
		},
		{
			name: "MintFarmingAccount balance:10000, and three pools:poolA(50%),poolB(30%),poolC(20%)",
			run: func(ctx sdk.Context, k keeper.Keeper) {
				// mint native token
				coins := sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, sdk.NewDec(10000))
				k.SupplyKeeper().MintCoins(ctx, MintFarmingAccount, coins)
				moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, MintFarmingAccount)
				yieldedNativeTokenAmt := moduleAcc.GetCoins().AmountOf(sdk.DefaultBondDenom)

				// init swap pair
				lockSymbol := "xxb"
				quoteSymbol := k.GetParams(ctx).QuoteSymbol
				swapTokenPair := swap.NewSwapPair(lockSymbol, quoteSymbol)
				swapTokenPair.QuotePooledCoin.Amount = sdk.NewDec(10000)
				swapTokenPair.BasePooledCoin.Amount = sdk.NewDec(10000)
				k.SwapKeeper().SetSwapTokenPair(ctx, swapTokenPair.TokenPairName(), swapTokenPair)
				k.SwapKeeper().SetParams(ctx, swap.DefaultParams())

				// create pools
				valueLockedPoolA := yieldedNativeTokenAmt.MulInt64(50).QuoInt64(100)
				poolA := types.FarmPool{
					Name:             "pool-a",
					TotalValueLocked: sdk.NewDecCoinFromDec(lockSymbol, valueLockedPoolA),
				}
				k.SetFarmPool(ctx, poolA)
				k.SetWhitelist(ctx, poolA.Name)
				poolCurrentRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), 3, sdk.SysCoins{})
				k.SetPoolCurrentRewards(ctx, poolA.Name, poolCurrentRewards)

				valueLockedPoolB := yieldedNativeTokenAmt.MulInt64(30).QuoInt64(100)
				poolB := types.FarmPool{
					Name:             "pool-b",
					TotalValueLocked: sdk.NewDecCoinFromDec(lockSymbol, valueLockedPoolB),
				}
				k.SetFarmPool(ctx, poolB)
				k.SetWhitelist(ctx, poolB.Name)
				k.SetPoolCurrentRewards(ctx, poolB.Name, poolCurrentRewards)

				valueLockedPoolC := yieldedNativeTokenAmt.MulInt64(20).QuoInt64(100)
				poolC := types.FarmPool{
					Name:             "pool-c",
					TotalValueLocked: sdk.NewDecCoinFromDec(lockSymbol, valueLockedPoolC),
				}
				k.SetFarmPool(ctx, poolC)
				k.SetWhitelist(ctx, poolC.Name)
				k.SetPoolCurrentRewards(ctx, poolC.Name, poolCurrentRewards)

				// execute BeginBlocker
				require.NotPanics(t, func() {
					BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: 3}}, k)
				})

				// check pools
				var found bool
				poolA, found = k.GetFarmPool(ctx, poolA.Name)
				require.True(t, found)
				require.True(t, poolA.TotalAccumulatedRewards.IsEqual(sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, valueLockedPoolA)))
				poolB, found = k.GetFarmPool(ctx, poolB.Name)
				require.True(t, found)
				require.True(t, poolB.TotalAccumulatedRewards.IsEqual(sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, valueLockedPoolB)))
				poolC, found = k.GetFarmPool(ctx, poolC.Name)
				require.True(t, found)
				require.True(t, poolC.TotalAccumulatedRewards.IsEqual(sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, valueLockedPoolC)))

			},
		},
	}

	// run test
	for _, test := range tests {
		test.run(ctx, k)
	}

}

func TestYieldNativeTokenDisabled(t *testing.T) {
	ctx, mk := keeper.GetKeeper(t)
	k := mk.Keeper

	coins := sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, sdk.NewDec(10000))
	err := k.SupplyKeeper().MintCoins(ctx, MintFarmingAccount, coins)
	require.NoError(t, err)
	retCoins := k.SupplyKeeper().GetModuleAccount(ctx, MintFarmingAccount).GetCoins()
	require.Equal(t, coins, retCoins)
	BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: 1}}, k)
	retCoins = k.SupplyKeeper().GetModuleAccount(ctx, MintFarmingAccount).GetCoins()
	require.Nil(t, retCoins)
}
