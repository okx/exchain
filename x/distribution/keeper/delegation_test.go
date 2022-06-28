package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking"
	"testing"
	"time"

	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/require"
)

func TestCalculateRewardsBasic(t *testing.T) {

	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitAllocateValidator(ctx)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// set new rate 0.5
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	DoEditValidator(t, ctx, sk, valOpAddr1, &newRate)
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator
	val := sk.Validator(ctx, valOpAddr1)

	// historical count should be 1 (once for validator init)
	require.Equal(t, uint64(1), dk.GetValidatorHistoricalReferenceCount(ctx))

	// end period
	dk.incrementValidatorPeriod(ctx, val)

	// historical count should be 1 still
	require.Equal(t, uint64(1), dk.GetValidatorHistoricalReferenceCount(ctx))

	// allocate some rewards
	initial := int64(10)
	tokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}
	dk.AllocateTokensToValidator(ctx, val, tokens)

	// end period
	dk.incrementValidatorPeriod(ctx, val)

	// commission should be the other half
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial / 2)}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
	require.Equal(t, sdk.SysCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}, dk.GetValidatorOutstandingRewards(ctx, valOpAddr1))
}

func TestCalculateRewardsMultiDelegator(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitAllocateValidator(ctx)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// set new rate 0.5
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	DoEditValidator(t, ctx, sk, valOpAddr1, &newRate)
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, valOpAddr1)

	// allocate some rewards
	initial := int64(20)
	tokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}
	dk.AllocateTokensToValidator(ctx, val, tokens)

	//first delegation
	DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(1), dk.GetValidatorHistoricalReferenceCount(ctx))
	valOpAddrs := []sdk.ValAddress{valOpAddr1}
	DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)
	// historical count should be 2(first is init validator)
	require.Equal(t, uint64(2), dk.GetValidatorHistoricalReferenceCount(ctx))

	// fetch updated validator
	val = sk.Validator(ctx, valOpAddr1)

	// end block
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	dk.AllocateTokensToValidator(ctx, val, tokens)

	// end period
	endingPeriod := dk.incrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards := dk.calculateDelegationRewards(ctx, val, delAddr1, endingPeriod)

	// rewards for del1 should be close to 1/2 initial
	require.True(t, rewards[0].Amount.LT(sdk.NewDec(initial/2)))
	require.True(t, rewards[0].Amount.GT(sdk.NewDec((initial/2)-1)))

	// calculate delegation rewards for del1
	rewards = dk.calculateDelegationRewards(ctx, val, delAddr1, endingPeriod)

	// rewards for del1 should be close to 1/2 initial
	require.True(t, rewards[0].Amount.LT(sdk.NewDec(initial/2)))
	require.True(t, rewards[0].Amount.GT(sdk.NewDec((initial/2)-1)))

	// commission should be equal to initial (50% twice)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
}

func TestWithdrawDelegationRewardsBasic(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, ak, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitAllocateValidator(ctx)

	balanceTokens := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), sdk.TokensFromConsensusPower(int64(1000))))

	//set module account coins
	distrAcc := dk.GetDistributionAccount(ctx)
	distrAcc.SetCoins(balanceTokens)
	dk.supplyKeeper.SetModuleAccount(ctx, distrAcc)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// set new rate 0.5
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	DoEditValidator(t, ctx, sk, valOpAddr1, &newRate)
	staking.EndBlocker(ctx, sk)
	valTokens := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sk.ParamsMinSelfDelegation(ctx))}
	// assert correct initial balance
	expTokens := balanceTokens.Sub(valTokens)
	require.Equal(t, expTokens, ak.GetAccount(ctx, sdk.AccAddress(valOpAddr1)).GetCoins())

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, valOpAddr1)

	initial := int64(20)
	tokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}

	dk.AllocateTokensToValidator(ctx, val, tokens)

	// historical count should be 1 (initial)
	require.Equal(t, uint64(1), dk.GetValidatorHistoricalReferenceCount(ctx))

	//assert correct balance
	exp := balanceTokens.Sub(valTokens)
	require.Equal(t, exp, ak.GetAccount(ctx, sdk.AccAddress(valOpAddr1)).GetCoins())

	// withdraw commission
	_, err := dk.WithdrawValidatorCommission(ctx, valOpAddr1)
	require.Nil(t, err)

	// assert correct balance
	exp = balanceTokens.Sub(valTokens).Add(tokens.QuoDec(sdk.NewDec(int64(2)))[0])
	require.Equal(t, exp, ak.GetAccount(ctx, sdk.AccAddress(valOpAddr1)).GetCoins())
}

func TestCalculateRewardsMultiDelegatorMultWithdraw(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, ak, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	tmtypes.UnittestOnlySetMilestoneSaturn1Height(-1)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitAllocateValidator(ctx)

	balanceTokens := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), sdk.TokensFromConsensusPower(int64(1000))))

	//set module account coins
	distrAcc := dk.GetDistributionAccount(ctx)
	distrAcc.SetCoins(balanceTokens)
	dk.supplyKeeper.SetModuleAccount(ctx, distrAcc)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// set new rate 0.5
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	DoEditValidator(t, ctx, sk, valOpAddr1, &newRate)
	staking.EndBlocker(ctx, sk)
	valTokens := sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sk.ParamsMinSelfDelegation(ctx))}
	// assert correct initial balance
	expTokens := balanceTokens.Sub(valTokens)
	require.Equal(t, expTokens, ak.GetAccount(ctx, sdk.AccAddress(valOpAddr1)).GetCoins())

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator
	val := sk.Validator(ctx, valOpAddr1)

	// allocate some rewards
	initial := int64(20)
	tokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}
	dk.AllocateTokensToValidator(ctx, val, tokens)

	//historical count should be 1 (validator init)
	require.Equal(t, uint64(1), dk.GetValidatorHistoricalReferenceCount(ctx))

	//first delegation
	DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	// historical count should be 1
	require.Equal(t, uint64(1), dk.GetValidatorHistoricalReferenceCount(ctx))
	valOpAddrs := []sdk.ValAddress{valOpAddr1}
	DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)
	// historical count should be 2 (first delegation init)
	require.Equal(t, uint64(2), dk.GetValidatorHistoricalReferenceCount(ctx))
	// end block
	staking.EndBlocker(ctx, sk)

	//second delegation
	DoDeposit(t, ctx, sk, delAddr2, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	// historical count should be 2
	require.Equal(t, uint64(2), dk.GetValidatorHistoricalReferenceCount(ctx))
	DoAddShares(t, ctx, sk, delAddr2, valOpAddrs)
	// historical count should be 3 (second delegation init)
	require.Equal(t, uint64(3), dk.GetValidatorHistoricalReferenceCount(ctx))
	// end block
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// fetch updated validator
	val = sk.Validator(ctx, valOpAddr1)

	// allocate some more rewards
	dk.AllocateTokensToValidator(ctx, val, tokens)

	// first delegator withdraws
	dk.WithdrawDelegationRewards(ctx, sdk.AccAddress(delAddr1), valOpAddr1)

	// second delegator withdraws
	dk.WithdrawDelegationRewards(ctx, sdk.AccAddress(delAddr2), valOpAddr1)

	// historical count should be 3 (two delegations)
	require.Equal(t, uint64(3), dk.GetValidatorHistoricalReferenceCount(ctx))

	// validator withdraws commission
	dk.WithdrawValidatorCommission(ctx, valOpAddr1)

	// end period
	endingPeriod := dk.incrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards := dk.calculateDelegationRewards(ctx, val, delAddr1, endingPeriod)

	// rewards for del1 should be zero
	require.True(t, rewards.IsZero())

	// calculate delegation rewards for del2
	rewards = dk.calculateDelegationRewards(ctx, val, delAddr2, endingPeriod)

	// rewards for del2 should be zero
	require.True(t, rewards.IsZero())

	// commission should be zero
	require.True(t, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1).IsZero())

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	dk.AllocateTokensToValidator(ctx, val, tokens)

	// first delegator withdraws again
	dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr1)

	// end period
	endingPeriod = dk.incrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards = dk.calculateDelegationRewards(ctx, val, delAddr1, endingPeriod)

	// rewards for del1 should be zero
	require.True(t, rewards.IsZero())

	// calculate delegation rewards for del2
	rewards = dk.calculateDelegationRewards(ctx, val, delAddr2, endingPeriod)

	// rewards for del2 should be close to 1/4 initial
	require.True(t, rewards[0].Amount.LT(sdk.NewDec(initial/4)))
	require.True(t, rewards[0].Amount.GT(sdk.NewDec((initial/4)-1)))

	// commission should be half initial
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial / 2)}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	dk.AllocateTokensToValidator(ctx, val, tokens)

	// withdraw commission
	dk.WithdrawValidatorCommission(ctx, valOpAddr1)

	// end period
	endingPeriod = dk.incrementValidatorPeriod(ctx, val)
	// calculate delegation rewards for del1
	rewards = dk.calculateDelegationRewards(ctx, val, delAddr1, endingPeriod)

	// rewards for del1 should be 1/4 initial
	require.True(t, rewards[0].Amount.LT(sdk.NewDec(initial/4)))
	require.True(t, rewards[0].Amount.GT(sdk.NewDec((initial/4)-1)))

	// calculate delegation rewards for del2
	rewards = dk.calculateDelegationRewards(ctx, val, delAddr2, endingPeriod)

	// rewards for del2 should be 1/4 + 1/4 initial
	// rewards for del1 should be close to 1/2 initial
	require.True(t, rewards[0].Amount.LT(sdk.NewDec(initial/2)))
	require.True(t, rewards[0].Amount.GT(sdk.NewDec((initial/2)-1)))

	// commission should be zero
	require.True(t, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1).IsZero())
}
