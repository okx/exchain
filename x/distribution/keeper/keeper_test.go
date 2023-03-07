package keeper

import (
	"testing"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func TestSetWithdrawAddr(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInputDefault(t, false, 1000)
	keeper.SetDistributionType(ctx, types.DistributionTypeOffChain)
	keeper.SetWithdrawRewardEnabled(ctx, true)
	keeper.SetRewardTruncatePrecision(ctx, 0)

	params := keeper.GetParams(ctx)
	params.WithdrawAddrEnabled = false
	keeper.SetParams(ctx, params)

	err := keeper.SetWithdrawAddr(ctx, delAddr1, delAddr2)
	require.NotNil(t, err)

	params.WithdrawAddrEnabled = true
	keeper.SetParams(ctx, params)

	err = keeper.SetWithdrawAddr(ctx, delAddr1, delAddr2)
	require.Nil(t, err)

	keeper.blacklistedAddrs[distrAcc.GetAddress().String()] = true
	require.Error(t, keeper.SetWithdrawAddr(ctx, delAddr1, distrAcc.GetAddress()))
}

func TestWithdrawValidatorCommission(t *testing.T) {
	ctx, ak, keeper, sk, _ := CreateTestInputDefault(t, false, 1000)

	valCommission := sdk.DecCoins{
		sdk.NewDecCoinFromDec("mytoken", sdk.NewDec(5).Quo(sdk.NewDec(4))),
		sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDec(3).Quo(sdk.NewDec(2))),
	}

	// set module account coins
	distrAcc := keeper.GetDistributionAccount(ctx)
	distrAcc.SetCoins(sdk.NewCoins(
		sdk.NewCoin("mytoken", sdk.NewInt(2)),
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2)),
	))
	keeper.supplyKeeper.SetModuleAccount(ctx, distrAcc)

	// check initial balance
	balance := ak.GetAccount(ctx, sdk.AccAddress(valOpAddr3)).GetCoins()
	expTokens := sdk.TokensFromConsensusPower(1000)
	subMsdCoin := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sk.ParamsMinSelfDelegation(ctx))
	expCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, expTokens).Sub(subMsdCoin))
	require.Equal(t, expCoins, balance)

	// withdraw commission error
	_, err := keeper.WithdrawValidatorCommission(ctx, valOpAddr3)
	require.Equal(t, err, types.ErrNoValidatorCommission())

	// set outstanding rewards
	keeper.SetValidatorOutstandingRewards(ctx, valOpAddr3, valCommission)

	// set commission
	keeper.SetValidatorAccumulatedCommission(ctx, valOpAddr3, valCommission)

	// withdraw commission
	_, err = keeper.WithdrawValidatorCommission(ctx, valOpAddr3)
	require.Equal(t, err, nil)

	// check balance increase
	balance = ak.GetAccount(ctx, sdk.AccAddress(valOpAddr3)).GetCoins()
	require.Equal(t, sdk.NewCoins(
		sdk.NewCoin("mytoken", sdk.NewInt(1)),
		sdk.NewCoin(sdk.DefaultBondDenom, expTokens.AddRaw(1)).Sub(subMsdCoin),
	), balance)

	// check remainder
	remainder := keeper.GetValidatorAccumulatedCommission(ctx, valOpAddr3)
	require.Equal(t, sdk.DecCoins{
		sdk.NewDecCoinFromDec("mytoken", sdk.NewDec(1).Quo(sdk.NewDec(4))),
		sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDec(1).Quo(sdk.NewDec(2))),
	}, remainder)

	require.True(t, true)
}

func TestGetTotalRewards(t *testing.T) {
	ctx, _, keeper, _, _ := CreateTestInputDefault(t, false, 1000)

	valCommission := sdk.DecCoins{
		sdk.NewDecCoinFromDec("mytoken", sdk.NewDec(5).Quo(sdk.NewDec(4))),
		sdk.NewDecCoinFromDec("stake", sdk.NewDec(3).Quo(sdk.NewDec(2))),
	}

	keeper.SetValidatorOutstandingRewards(ctx, valOpAddr1, valCommission)
	keeper.SetValidatorOutstandingRewards(ctx, valOpAddr2, valCommission)

	expectedRewards := valCommission.MulDec(sdk.NewDec(2))
	totalRewards := keeper.GetTotalRewards(ctx)

	require.Equal(t, expectedRewards, totalRewards)
}
