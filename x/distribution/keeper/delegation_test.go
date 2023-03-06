package keeper

import (
	"testing"
	"time"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/types"
	"github.com/okx/okbchain/x/staking"
	stakingexported "github.com/okx/okbchain/x/staking/exported"
	stakingtypes "github.com/okx/okbchain/x/staking/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/require"
)

func TestCalculateRewardsBasic(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitExistedValidatorFlag(ctx, true)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// set new rate 0.5
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	DoEditValidator(t, ctx, sk, valOpAddr1, newRate)
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
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitExistedValidatorFlag(ctx, true)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)

	// end block to bond validator
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// set new rate 0.5
	newRate, _ := sdk.NewDecFromStr("0.5")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	DoEditValidator(t, ctx, sk, valOpAddr1, newRate)
	staking.EndBlocker(ctx, sk)

	// next block
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// fetch validator and delegation
	val := sk.Validator(ctx, valOpAddr1)

	// allocate some rewards
	initial := int64(20)
	tokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}
	dk.AllocateTokensToValidator(ctx, val, tokens)

	valOpAddrs := []sdk.ValAddress{valOpAddr1}
	//first delegation
	DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	require.Equal(t, uint64(1), dk.GetValidatorHistoricalReferenceCount(ctx))
	DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)
	// historical count should be 2(first is init validator)
	require.Equal(t, uint64(2), dk.GetValidatorHistoricalReferenceCount(ctx))

	//second delegation
	DoDeposit(t, ctx, sk, delAddr2, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))
	DoAddShares(t, ctx, sk, delAddr2, valOpAddrs)
	require.Equal(t, uint64(3), dk.GetValidatorHistoricalReferenceCount(ctx))

	// fetch updated validator
	val = sk.Validator(ctx, valOpAddr1)

	// end block
	staking.EndBlocker(ctx, sk)
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)

	// allocate some more rewards
	dk.AllocateTokensToValidator(ctx, val, tokens)
	// end period
	endingPeriod := dk.incrementValidatorPeriod(ctx, val)

	// calculate delegation rewards for del1
	rewards1 := dk.calculateDelegationRewards(ctx, val, delAddr1, endingPeriod)
	require.True(t, rewards1[0].Amount.LT(sdk.NewDec(initial/4)))
	require.True(t, rewards1[0].Amount.GT(sdk.NewDec((initial/4)-1)))

	// calculate delegation rewards for del2
	rewards2 := dk.calculateDelegationRewards(ctx, val, delAddr2, endingPeriod)
	require.True(t, rewards2[0].Amount.LT(sdk.NewDec(initial/4)))
	require.True(t, rewards2[0].Amount.GT(sdk.NewDec((initial/4)-1)))

	// commission should be equal to initial (50% twice)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(initial)}}, dk.GetValidatorAccumulatedCommission(ctx, valOpAddr1))
}

func TestWithdrawDelegationRewardsBasic(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, ak, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitExistedValidatorFlag(ctx, true)

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
	DoEditValidator(t, ctx, sk, valOpAddr1, newRate)
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
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitExistedValidatorFlag(ctx, true)

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
	DoEditValidator(t, ctx, sk, valOpAddr1, newRate)
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

func TestInitExistedDelegationStartInfoSuite(t *testing.T) {
	suite.Run(t, new(InitExistedDelegationStartInfoestSuite))
}

type InitExistedDelegationStartInfoestSuite struct {
	suite.Suite
}

func changeDistribution(ctx sdk.Context, dk Keeper) {
	//change to distribution onchain
	dk.SetInitExistedValidatorFlag(ctx, true)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingexported.ValidatorI) (stop bool) {
		if validator != nil {
			dk.initExistedValidatorForDistrProposal(ctx, validator)
		}
		return false
	})
}

func (suite *InitExistedDelegationStartInfoestSuite) TestInitExistedDelegationStartInfo() {
	testCases := []struct {
		title                         string
		execute1                      func(ctx *sdk.Context, dk Keeper)
		execute2                      func(ctx *sdk.Context, dk Keeper)
		execute3                      func(ctx *sdk.Context, dk Keeper)
		execute4                      func(ctx *sdk.Context, dk Keeper)
		beforeAddSharesReferenceCount uint64
		afterAddSharesReferenceCount  uint64
		afterWithdrawReferenceCount   uint64
		coins                         sdk.Coins
		err                           error
	}{
		{
			"NO ERROR Before create validator",
			func(ctx *sdk.Context, dk Keeper) {
				changeDistribution(*ctx, dk)
			},
			func(ctx *sdk.Context, dk Keeper) {
				ctx.SetBlockTime(time.Now())
			},
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {},
			1,
			2,
			2,
			sdk.Coins(nil),
			nil,
		},
		{
			"NO ERROR Before Do Deposit",
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {
				ctx.SetBlockTime(time.Now())
				changeDistribution(*ctx, dk)
			},
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {},
			2,
			3,
			3,
			sdk.Coins(nil),
			nil,
		},
		{
			"NO ERROR Before Do Add Shares",
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {
				ctx.SetBlockTime(time.Now())
			},
			func(ctx *sdk.Context, dk Keeper) {
				changeDistribution(*ctx, dk)
			},
			func(ctx *sdk.Context, dk Keeper) {},
			2,
			3,
			3,
			sdk.Coins(nil),
			nil,
		},
		{
			"NO ERROR After Do Add Shares",
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {
				ctx.SetBlockTime(time.Now())
			},
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {
				changeDistribution(*ctx, dk)
			},
			0,
			0,
			3,
			sdk.Coins(nil),
			nil,
		},
		{
			"ERROR No Shares",
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {},
			func(ctx *sdk.Context, dk Keeper) {
				changeDistribution(*ctx, dk)
			},
			0,
			0,
			2,
			sdk.Coins(nil),
			types.ErrCodeZeroDelegationShares(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			communityTax := sdk.NewDecWithPrec(2, 2)
			ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(suite.T(), false, 1000, communityTax)
			balanceTokens := sdk.NewCoins(sdk.NewCoin(sk.BondDenom(ctx), sdk.TokensFromConsensusPower(int64(1000))))
			//set module account coins
			distrAcc := dk.GetDistributionAccount(ctx)
			distrAcc.SetCoins(balanceTokens)
			dk.supplyKeeper.SetModuleAccount(ctx, distrAcc)
			tc.execute1(&ctx, dk)
			// create validator
			DoCreateValidator(suite.T(), ctx, sk, valOpAddr1, valConsPk1)
			// end block to bond validator
			staking.EndBlocker(ctx, sk)

			ctx.SetBlockHeight(ctx.BlockHeight() + 1)
			//ctx.SetBlockTime(time.Now())
			tc.execute2(&ctx, dk)
			//delegation
			DoDeposit(suite.T(), ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(100)))

			coins, err := dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr1)
			require.Equal(suite.T(), types.ErrCodeEmptyDelegationVoteValidator(), err)

			tc.execute3(&ctx, dk)
			// historical count
			require.Equal(suite.T(), tc.beforeAddSharesReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
			valOpAddrs := []sdk.ValAddress{valOpAddr1}
			DoAddShares(suite.T(), ctx, sk, delAddr1, valOpAddrs)
			require.Equal(suite.T(), tc.afterAddSharesReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))

			tc.execute4(&ctx, dk)
			// end block
			staking.EndBlocker(ctx, sk)

			coins, err = dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr3)
			require.Equal(suite.T(), types.ErrCodeEmptyValidatorDistInfo(), err)

			coins, err = dk.WithdrawDelegationRewards(ctx, delAddr2, valOpAddr1)
			require.Equal(suite.T(), types.ErrCodeEmptyDelegationDistInfo(), err)

			coins, err = dk.WithdrawDelegationRewards(ctx, delAddr1, valOpAddr1)
			require.Equal(suite.T(), tc.afterWithdrawReferenceCount, dk.GetValidatorHistoricalReferenceCount(ctx))
			require.Equal(suite.T(), tc.coins, coins)
			require.Equal(suite.T(), tc.err, err)
		})
	}
}

func TestInvalidDelegation(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	// Cannot init when distribution proposal valid
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)
	dk.initializeDelegation(ctx, valOpAddr1, delAddr1)
	require.False(t, dk.HasDelegatorStartingInfo(ctx, valOpAddr1, delAddr1))

	// Cannot init when distribution proposal valid
	val := dk.stakingKeeper.Validator(ctx, valOpAddr1)
	del := dk.stakingKeeper.Delegator(ctx, delAddr1)
	dk.initExistedDelegationStartInfo(ctx, val, del)
	require.False(t, dk.HasDelegatorStartingInfo(ctx, valOpAddr1, delAddr1))

	// init delegation
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	ctx.SetBlockTime(time.Now())
	DoDeposit(t, ctx, sk, delAddr1, sdk.NewCoin(sk.BondDenom(ctx), sdk.NewInt(1)))
	valOpAddrs := []sdk.ValAddress{valOpAddr1}
	DoAddShares(t, ctx, sk, delAddr1, valOpAddrs)
	staking.EndBlocker(ctx, sk)

	// check calculateDelegationRewardsBetween startingPeriod > endingPeriod
	newDec, _ := sdk.NewDecFromStr("1")
	period := types.DelegatorStartingInfo{
		PreviousPeriod: uint64(100),
		Stake:          newDec,
		Height:         10,
	}
	dk.SetDelegatorStartingInfo(ctx, valOpAddr1, delAddr1, period)
	panicFunc := func() {
		dk.calculateDelegationRewards(ctx, val, delAddr1, 1)
	}
	assert.PanicsWithValue(t, "startingPeriod cannot be greater than endingPeriod", panicFunc)

	// check calculateDelegationRewards stake.GT(del.GetLastAddedShares())
	newDec, _ = sdk.NewDecFromStr("100000000000000")
	period = types.DelegatorStartingInfo{
		PreviousPeriod: uint64(100),
		Stake:          newDec,
		Height:         10,
	}
	dk.SetDelegatorStartingInfo(ctx, valOpAddr1, delAddr1, period)
	panicFunc = func() {
		dk.calculateDelegationRewards(ctx, val, delAddr1, 1)
	}
	assert.Panics(t, panicFunc)
}

func TestIncrementValidatorPeriod(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)
	val := sk.Validator(ctx, valOpAddr1)

	// distribution type invalid, No Panic
	noPanicFunc := func() {
		dk.incrementValidatorPeriod(ctx, val)
	}
	assert.NotPanics(t, noPanicFunc)
}

func TestRewardToCommunity(t *testing.T) {
	communityTax := sdk.NewDecWithPrec(2, 2)
	ctx, _, _, dk, sk, _, _ := CreateTestInputAdvanced(t, false, 1000, communityTax)
	dk.SetDistributionType(ctx, types.DistributionTypeOnChain)
	dk.SetInitExistedValidatorFlag(ctx, true)

	// create validator
	DoCreateValidator(t, ctx, sk, valOpAddr1, valConsPk1)
	newRate, _ := sdk.NewDecFromStr("0")
	ctx.SetBlockTime(time.Now().UTC().Add(48 * time.Hour))
	DoEditValidator(t, ctx, sk, valOpAddr1, newRate)
	val := sk.Validator(ctx, valOpAddr1)

	// allocate some rewards
	tokens := sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewDec(int64(20))}}
	dk.AllocateTokensToValidator(ctx, val, tokens)

	sk.SetValidator(ctx, stakingtypes.Validator{OperatorAddress: val.GetOperator(), DelegatorShares: sdk.NewDec(int64(0))})
	val = sk.Validator(ctx, valOpAddr1)

	beforeFeePool := dk.GetFeePool(ctx)
	dk.incrementValidatorPeriod(ctx, val)
	afterFeePool := dk.GetFeePool(ctx)
	require.Equal(t, tokens, afterFeePool.CommunityPool.Sub(beforeFeePool.CommunityPool))
}
