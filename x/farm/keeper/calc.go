package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// CalculateAmountYieldedBetween is used for calculating how many tokens haven been yielded from
// startBlockHeight to endBlockHeight. And return the amount.
func CalculateAmountYieldedBetween(
	endBlockHeight int64, startBlockHeight int64, pool types.FarmPool,
) (types.FarmPool, sdk.DecCoins) {
	yieldedTokens := sdk.DecCoins{}
	for i := 0; i < len(pool.YieldedTokenInfos); i++ {
		startBlockHeightToYield := pool.YieldedTokenInfos[i].StartBlockHeightToYield

		// if condition startBlockHeightToYield <= startBlockHeight < endBlockHeight is not satisfied, then continue
		if startBlockHeightToYield == 0 || startBlockHeight < startBlockHeightToYield || startBlockHeight >= endBlockHeight {
			continue
		}

		// calculate the exact interval
		blockInterval := sdk.NewDec(endBlockHeight - startBlockHeight)
		// calculate how many coin have been yielded till the current block
		amount := blockInterval.MulTruncate(pool.YieldedTokenInfos[i].AmountYieldedPerBlock)
		remaining := pool.YieldedTokenInfos[i].RemainingAmount
		if amount.LT(remaining.Amount) {
			// subtract yielded_coin amount
			pool.YieldedTokenInfos[i].RemainingAmount.Amount = remaining.Amount.Sub(amount)
			yieldedTokens = yieldedTokens.Add(sdk.NewDecCoinsFromDec(remaining.Denom, amount))
		} else {
			// initialize yieldedTokenInfo
			pool.YieldedTokenInfos[i] = types.NewYieldedTokenInfo(sdk.NewDecCoin(remaining.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec())
			yieldedTokens = yieldedTokens.Add(sdk.NewDecCoinsFromDec(remaining.Denom, remaining.Amount))
		}
	}
	return pool, yieldedTokens
}

func (k Keeper) WithdrawRewards(ctx sdk.Context, pool types.FarmPool, addr sdk.AccAddress) (sdk.DecCoins, sdk.Error) {
	// 0. check existence of delegator starting info
	if !k.HasLockInfo(ctx, addr, pool.Name) {
		return nil, types.ErrNoLockInfoFound(types.DefaultCodespace, addr.String())
	}

	// 1. end current period and calculate rewards
	//endingPeriod := k.incrementPoolPeriod(ctx, pool)
	endingPeriod := k.IncrementPoolPeriod(ctx, pool)
	rewards := k.calculateRewards(ctx, pool.Name, addr, endingPeriod)

	// add coins to user account
	if !rewards.IsZero() {
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.YieldFarmingAccount, addr, rewards)
		if err != nil {
			return nil, err
		}
	}

	// decrement reference count of lock info
	lockInfo, _ := k.GetLockInfo(ctx, addr, pool.Name)
	k.decrementReferenceCount(ctx, pool.Name, lockInfo.ReferencePeriod)

	// remove delegator starting info
	k.DeleteLockInfo(ctx, addr, pool.Name)

	return rewards, nil
}

// increment pool period, returning the period just ended
func (k Keeper) IncrementPoolPeriod(ctx sdk.Context, pool types.FarmPool) uint64 {
	// fetch current rewards status
	rewards := k.GetPoolCurrentRewards(ctx, pool.Name)

	// 1.1 calculate how many provided token has been yielded between start_block_height and current_height
	updatedPool, yieldedTokens := CalculateAmountYieldedBetween(ctx.BlockHeight(), rewards.StartBlockHeight, pool)

	// 1.2 calculate how many native token has been yielded between start_block_height and current_height
	rewards.Rewards = rewards.Rewards.Add(yieldedTokens)

	// 2. calculate current reward ratio
	var currentRatio sdk.DecCoins
	if pool.TotalValueLocked.IsZero() {
		currentRatio = sdk.DecCoins{}
	} else {
		currentRatio = rewards.Rewards.QuoDecTruncate(pool.TotalValueLocked.Amount)
	}

	// 3.1 get the previous pool_historical_rewards
	historical := k.GetPoolHistoricalRewards(ctx, pool.Name, rewards.Period-1).CumulativeRewardRatio
	// 3.2 decrement reference count
	k.decrementReferenceCount(ctx, pool.Name, rewards.Period-1)
	// 3.3 create new pool_historical_rewards with reference count of 1, then set it into store
	newHistoricalRewards := types.NewPoolHistoricalRewards(historical.Add(currentRatio), 1)
	k.SetPoolHistoricalRewards(ctx, pool.Name, rewards.Period, newHistoricalRewards)

	// 4. set new current newYieldedRewards into store, incrementing period by 1
	newCurRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), rewards.Period+1, sdk.DecCoins{})
	k.SetPoolCurrentRewards(ctx, pool.Name, newCurRewards)

	// 5. set updated pool
	k.SetFarmPool(ctx, updatedPool)

	return rewards.Period
}

// increment the reference count for a historical rewards value
func (k Keeper) incrementReferenceCount(ctx sdk.Context, poolName string, period uint64) {
	historical := k.GetPoolHistoricalRewards(ctx, poolName, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetPoolHistoricalRewards(ctx, poolName, period, historical)
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

func (k Keeper) calculateRewards(ctx sdk.Context, poolName string, addr sdk.AccAddress, endingPeriod uint64) sdk.DecCoins {
	// fetch lock info
	lockInfo, found := k.GetLockInfo(ctx, addr, poolName)
	if !found {
		panic("should not happen")
	}
	if lockInfo.StartBlockHeight <= ctx.BlockHeight() {
		// started this height, no rewards yet
		return nil
	}

	currentPeriod := k.GetPoolCurrentRewards(ctx, poolName)
	startingPeriod := currentPeriod.Period
	// calculate rewards for final period
	return k.calculateDelegationRewardsBetween(ctx, poolName, startingPeriod, endingPeriod, lockInfo.Amount)
}

// calculate the rewards accrued by a pool between two periods
func (k Keeper) calculateDelegationRewardsBetween(ctx sdk.Context, poolName string, startingPeriod, endingPeriod uint64,
	amount sdk.DecCoin) (rewards sdk.DecCoins) {

	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	if amount.Amount.LT(sdk.ZeroDec()) {
		panic("amount should not be negative")
	}

	// return amount * (ending - starting)
	starting := k.GetPoolHistoricalRewards(ctx, poolName, startingPeriod)
	ending := k.GetPoolHistoricalRewards(ctx, poolName, endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	rewards = difference.MulDecTruncate(amount.Amount)
	return
}

// initialize starting info for a new lock info
func (k Keeper) InitializeLockInfo(ctx sdk.Context, addr sdk.AccAddress, poolName string, changedAmount sdk.Dec) {
	// period has already been incremented - we want to store the period ended by this delegation action
	previousPeriod := k.GetPoolCurrentRewards(ctx, poolName).Period - 1

	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, poolName, previousPeriod)

	// get lock info, then set it into store
	lockInfo, found := k.GetLockInfo(ctx, addr, poolName)
	if !found {
		panic("the lock info can't be found")
	}
	lockInfo.StartBlockHeight = ctx.BlockHeight()
	lockInfo.ReferencePeriod = previousPeriod
	lockInfo.Amount.Amount = lockInfo.Amount.Amount.Add(changedAmount)
	if lockInfo.Amount.IsZero() { // TODO should the lockinfo be deleted when its amount is zero?
		k.DeleteLockInfo(ctx, lockInfo.Owner, lockInfo.PoolName)
	} else {
		k.SetLockInfo(ctx, lockInfo)
	}
}
