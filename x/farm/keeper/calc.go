package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// LiquidateYieldedTokenInfo is used for calculating how many tokens haven been yielding from LastClaimedBlockHeight to CurrentHeight
// Then transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
func (k Keeper) LiquidateYieldedTokenInfo(height int64, pool types.FarmPool) types.FarmPool {
	if height <= pool.LastClaimedBlockHeight {
		return pool
	}

	// TODO: there are too many operations about MulTruncate, check the amount carefully, and write checking codes in invariants.go !!!
	// 1. Transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
	for i := 0; i < len(pool.YieldedTokenInfos); i++ {
		startBlockHeightToYield := pool.YieldedTokenInfos[i].StartBlockHeightToYield
		if height > startBlockHeightToYield {
			// calculate the exact interval
			var blockInterval sdk.Dec
			if startBlockHeightToYield > pool.LastClaimedBlockHeight {
				blockInterval = sdk.NewDec(height - startBlockHeightToYield)
			} else {
				blockInterval = sdk.NewDec(height - pool.LastClaimedBlockHeight)
			}

			// calculate how many coin have been yielded till the current block
			amountYielded := blockInterval.MulTruncate(pool.YieldedTokenInfos[i].AmountYieldedPerBlock)
			remainingAmount := pool.YieldedTokenInfos[i].RemainingAmount
			if amountYielded.LT(remainingAmount.Amount) {
				// add yielded amount
				pool.AmountYielded = pool.AmountYielded.Add(sdk.NewDecCoinsFromDec(remainingAmount.Denom, amountYielded))
				// subtract yielded_coin amount
				pool.YieldedTokenInfos[i].RemainingAmount.Amount = remainingAmount.Amount.Sub(amountYielded)
			} else {
				// add yielded amount
				pool.AmountYielded = pool.AmountYielded.Add(sdk.NewCoins(remainingAmount))

				// initialize yieldedTokenInfo
				pool.YieldedTokenInfos[i] = types.NewYieldedTokenInfo(sdk.NewDecCoin(remainingAmount.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec())

				// TODO: remove the YieldedTokenInfo when its amount become zero
				// Currently, we support only one token of	 yield farming at the same time,
				// so, it is unnecessary to remove the element in slice
			}
		}
	}

	return pool
}

func (k Keeper) ClaimRewards(ctx sdk.Context, pool types.FarmPool, lockInfo types.LockInfo, address sdk.AccAddress,
	changedAmount sdk.Dec) sdk.Error {
	// 1. calculation
	currentHeight := sdk.NewDec(ctx.BlockHeight())
	claimedAmount, selfChangedWeight := calculateYieldedAmount(currentHeight, pool, lockInfo)
	// 2. Transfer yielded tokens to personal account
	if !claimedAmount.IsZero() {
		if err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.YieldFarmingAccount, address, claimedAmount); err != nil {
			return err
		}
		// 3. Update the pool data
		pool.AmountYielded = pool.AmountYielded.Sub(claimedAmount)
	}
	if !changedAmount.IsZero() {
		pool.TotalValueLocked.Amount = pool.TotalValueLocked.Amount.Add(changedAmount)
		selfChangedWeight = selfChangedWeight.Add(currentHeight.MulTruncate(changedAmount))
	}
	pool.TotalLockedWeight = pool.TotalLockedWeight.Add(selfChangedWeight)
	pool.LastClaimedBlockHeight = ctx.BlockHeight()
	// Set the updated pool into store
	k.SetFarmPool(ctx, pool)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeClaim,
		sdk.NewAttribute(types.AttributeKeyClaimed, claimedAmount.String()),
	))

	return nil
}

// calculateYieldedAmount calculates the yielded amount which belongs to an account on a given block height
func calculateYieldedAmount(currentHeight sdk.Dec, pool types.FarmPool, lockInfo types.LockInfo) (sdk.DecCoins, sdk.Dec) {
	startBlockHeight := sdk.NewDec(lockInfo.StartBlockHeight)
	if currentHeight.LTE(startBlockHeight) {
		return sdk.NewCoins(), sdk.ZeroDec()
	}

	/* 1.1 Calculate its own weight during these blocks
	   (curHeight - Height1) * Amount1
	*/
	oldWeight := startBlockHeight.MulTruncate(lockInfo.Amount.Amount)
	currentWeight := currentHeight.MulTruncate(lockInfo.Amount.Amount)
	selfChangedWeight := currentWeight.Sub(oldWeight)

	/* 1.2 Calculate all weight during these blocks
	    (curHeight - Height1) * Amount1 + (curHeight - Height2) * Amount2 + (curHeight - Height3) * Amount3
												||
	                                            \/
	   curHeight * (Amount1 + Amount2 + Amount3) - (Height1*Amount1 + Height2*Amount2 + Height3*Amount3)
												||
	                                            \/
	ctx.BlockHeight()  *  pool.TotalValueLocked.Amount  -  ( pool.TotalLockedWeight )
	*/
	totalChangedWeight := currentHeight.MulTruncate(pool.TotalValueLocked.Amount).Sub(pool.TotalLockedWeight)

	if selfChangedWeight.IsZero() || totalChangedWeight.IsZero() {
		return sdk.NewCoins(), selfChangedWeight
	}
	// 1.3 Calculate how many yielded tokens to return
	claimedAmount := pool.AmountYielded.MulDecTruncate(selfChangedWeight).QuoDecTruncate(totalChangedWeight)
	return claimedAmount, selfChangedWeight
}

func calculateRewards(height int64, pool types.FarmPool, period types.PoolCurrentPeriod) types.FarmPool {
	for i := 0; i < len(pool.YieldedTokenInfos); i++ {
		startBlockHeightToYield := pool.YieldedTokenInfos[i].StartBlockHeightToYield
		if height > startBlockHeightToYield {
			// calculate the exact interval
			var blockInterval sdk.Dec
			if startBlockHeightToYield > period.StartBlockHeight {
				blockInterval = sdk.NewDec(height - startBlockHeightToYield)
			} else {
				blockInterval = sdk.NewDec(height - period.StartBlockHeight)
			}

			// calculate how many coin have been yielded till the current block
			amountYielded := blockInterval.MulTruncate(pool.YieldedTokenInfos[i].AmountYieldedPerBlock)
			remainingAmount := pool.YieldedTokenInfos[i].RemainingAmount
			if amountYielded.LT(remainingAmount.Amount) {
				// add yielded amount
				pool.AmountYieldedNativeToken = pool.AmountYieldedNativeToken.Add(sdk.NewDecCoinsFromDec(remainingAmount.Denom, amountYielded))
				// subtract yielded_coin amount
				pool.YieldedTokenInfos[i].RemainingAmount.Amount = remainingAmount.Amount.Sub(amountYielded)
			} else {
				// add yielded amount
				pool.AmountYieldedNativeToken = pool.AmountYieldedNativeToken.Add(sdk.NewCoins(remainingAmount))

				// initialize yieldedTokenInfo
				pool.YieldedTokenInfos[i] = types.NewYieldedTokenInfo(sdk.NewDecCoin(remainingAmount.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec())

				// TODO: remove the YieldedTokenInfo when its amount become zero
				// Currently, we support only one token of yield farming at the same time,
				// so, it is unnecessary to remove the element in slice
			}
		}
	}
	return types.FarmPool{}
}

func (k Keeper) withdrawRewards(ctx sdk.Context, poolName string) (sdk.DecCoins, sdk.Error) {
	// check existence of delegator starting info
	if !k.HasPoolCurrentPeriod(ctx, poolName) {
		return nil, types.ErrNoPoolCurrentPeriodFound(types.DefaultCodespace, poolName)
	}
	pool, _ := k.GetFarmPool(ctx, poolName)

	// end current period and calculate rewards
	endingPeriod := k.incrementValidatorPeriod(ctx, pool)
	rewardsRaw := k.calculateRewards(ctx, poolName, endingPeriod)
	outstanding := k.GetValidatorOutstandingRewards(ctx, del.GetValidatorAddr())

	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logger := k.Logger(ctx)
		logger.Info(fmt.Sprintf("missing rewards rounding error, delegator %v"+
			"withdrawing rewards from validator %v, should have received %v, got %v",
			val.GetOperator(), del.GetDelegatorAddr(), rewardsRaw, rewards))
	}

	// truncate coins, return remainder to community pool
	coins, remainder := rewards.TruncateDecimal()

	// add coins to user account
	if !coins.IsZero() {
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, del.GetDelegatorAddr())
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
		if err != nil {
			return nil, err
		}
	}

	// update the outstanding rewards and the community pool only if the
	// transaction was successful
	k.SetValidatorOutstandingRewards(ctx, del.GetValidatorAddr(), outstanding.Sub(rewards))
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remainder)
	k.SetFeePool(ctx, feePool)

	// decrement reference count of starting period
	startingInfo := k.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
	startingPeriod := startingInfo.PreviousPeriod
	k.decrementReferenceCount(ctx, poolName, startingPeriod)

	// remove delegator starting info
	k.DeleteDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())

	return coins, nil
}

// increment pool period, returning the period just ended
func (k Keeper) incrementPoolPeriod(ctx sdk.Context, pool types.FarmPool) uint64 {
	// fetch current rewards
	rewards := k.GetPoolCurrentPeriod(ctx, pool.Name)

	// calculate current ratio
	var current sdk.DecCoins
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	current = rewards.LastAmountYieldedNativeToken.Amount.QuoTruncate(pool.)

	// fetch historical rewards for last period
	historical := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), rewards.Period-1).CumulativeRewardRatio

	// decrement reference count
	k.decrementReferenceCount(ctx, val.GetOperator(), rewards.Period-1)

	// set new historical rewards with reference count of 1
	k.SetValidatorHistoricalRewards(ctx, val.GetOperator(), rewards.Period, types.NewValidatorHistoricalRewards(historical.Add(current), 1))

	// set current rewards, incrementing period by 1
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), types.NewValidatorCurrentRewards(sdk.DecCoins{}, rewards.Period+1))

	return rewards.Period
}

func (k Keeper) calculateRewards(ctx sdk.Context, poolName string, endingPeriod uint64)  sdk.DecCoins {
	// fetch current period
	currentPeriod, found := k.GetPoolCurrentPeriod(ctx, poolName)
	if !found {
		return nil
	}
	startingPeriod := currentPeriod.Period
	lastAmountYielded := currentPeriod.LastAmountYieldedNativeToken
	// calculate rewards for final period
	return k.calculateDelegationRewardsBetween(ctx, poolName, startingPeriod, endingPeriod, lastAmountYielded)
}

// calculate the rewards accrued by a pool between two periods
func (k Keeper) calculateDelegationRewardsBetween(ctx sdk.Context, poolName string, startingPeriod, endingPeriod uint64,
	amount sdk.DecCoin) (rewards sdk.DecCoins) {

	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// return amount * (ending - starting)
	starting := k.GetPoolHistoricalRewards(ctx, poolName, startingPeriod)
	ending := k.GetPoolHistoricalRewards(ctx, poolName, endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards = difference.MulDecTruncate(amount.Amount)
	return
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k Keeper) decrementReferenceCount(ctx sdk.Context, poolName string, period uint64) {
	historical := k.GetPoolHistoricalRewards(ctx, poolName, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeletePoolHistoricalReward(ctx, poolName, period)
	} else {
		k.SetPoolHistoricalRewards(ctx, poolName, period, historical)
	}
}