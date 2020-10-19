package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// CalculateAmountYieldedBetween is used for calculating how many tokens haven been yielded from
// startBlockHeight to endBlockHeight. And return the amount.
func (k Keeper) CalculateAmountYieldedBetween(ctx sdk.Context, pool types.FarmPool) (types.FarmPool, sdk.DecCoins) {
	currentPeriod := k.GetPoolCurrentRewards(ctx, pool.Name)
	endBlockHeight := ctx.BlockHeight()

	// add native tokens in yieldedTokens
	yieldedTokens := sdk.DecCoins{}
	for i := 0; i < len(pool.YieldedTokenInfos); i++ {
		startBlockHeightToYield := pool.YieldedTokenInfos[i].StartBlockHeightToYield
		var startBlockHeight int64
		if currentPeriod.StartBlockHeight <= startBlockHeightToYield {
			startBlockHeight = startBlockHeightToYield
		} else {
			startBlockHeight = currentPeriod.StartBlockHeight
		}

		// if condition startBlockHeightToYield <= startBlockHeight < endBlockHeight is not satisfied, then continue
		if startBlockHeightToYield == 0 || startBlockHeight >= endBlockHeight {
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
			// Ⅱ. add yielded tokens in yieldedTokens
			yieldedTokens = yieldedTokens.Add(sdk.NewDecCoinsFromDec(remaining.Denom, amount))
		} else {
			// initialize yieldedTokenInfo
			pool.YieldedTokenInfos[i] = types.NewYieldedTokenInfo(sdk.NewDecCoin(remaining.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec())
			// add yielded tokens in yieldedTokens
			yieldedTokens = yieldedTokens.Add(sdk.NewDecCoinsFromDec(remaining.Denom, remaining.Amount))
		}
	}
	return pool, yieldedTokens
}

func (k Keeper) WithdrawRewards(
	ctx sdk.Context, poolName string, totalValue sdk.DecCoin, yieldedTokens sdk.DecCoins, addr sdk.AccAddress,
) (sdk.DecCoins, sdk.Error) {
	// 0. check existence of delegator starting info
	if !k.HasLockInfo(ctx, addr, poolName) {
		return nil, types.ErrNoLockInfoFound(types.DefaultCodespace, addr.String())
	}

	// 1. end current period and calculate rewards
	//endingPeriod := k.incrementPoolPeriod(ctx, pool)
	endingPeriod := k.IncrementPoolPeriod(ctx, poolName, totalValue, yieldedTokens)
	rewards := k.calculateRewards(ctx, poolName, addr, endingPeriod)

	// add rewards to user account
	if !rewards.IsZero() {
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.YieldFarmingAccount, addr, rewards)
		if err != nil {
			return nil, err
		}
	}

	// decrement reference count of lock info
	lockInfo, _ := k.GetLockInfo(ctx, addr, poolName)
	k.decrementReferenceCount(ctx, poolName, lockInfo.ReferencePeriod)

	return rewards, nil
}

// IncrementPoolPeriod increment pool period, returning the period just ended
func (k Keeper) IncrementPoolPeriod(ctx sdk.Context, poolName string, totalValue sdk.DecCoin, yieldedTokens sdk.DecCoins) uint64 {
	// 1. fetch current rewards status
	rewards := k.GetPoolCurrentRewards(ctx, poolName)
	// 2. calculate current reward ratio
	rewards.Rewards = rewards.Rewards.Add(yieldedTokens)
	var currentRatio sdk.DecCoins
	if totalValue.IsZero() {
		currentRatio = sdk.DecCoins{}
	} else {
		currentRatio = rewards.Rewards.QuoDecTruncate(totalValue.Amount)
	}

	// 3.1 get the previous pool_historical_rewards
	historical := k.GetPoolHistoricalRewards(ctx, poolName, rewards.Period-1).CumulativeRewardRatio
	// 3.2 decrement reference count
	k.decrementReferenceCount(ctx, poolName, rewards.Period-1)
	// 3.3 create new pool_historical_rewards with reference count of 1, then set it into store
	newHistoricalRewards := types.NewPoolHistoricalRewards(historical.Add(currentRatio), 1)
	k.SetPoolHistoricalRewards(ctx, poolName, rewards.Period, newHistoricalRewards)

	// 4. set new current newYieldedRewards into store, incrementing period by 1
	newCurRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), rewards.Period+1, sdk.DecCoins{})
	k.SetPoolCurrentRewards(ctx, poolName, newCurRewards)

	return rewards.Period
}

// incrementReferenceCount increment the reference count for a historical rewards value
func (k Keeper) incrementReferenceCount(ctx sdk.Context, poolName string, period uint64) {
	historical := k.GetPoolHistoricalRewards(ctx, poolName, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetPoolHistoricalRewards(ctx, poolName, period, historical)
}

// decrementReferenceCount decrement the reference count for a historical rewards value, and delete if zero references remain
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
	if lockInfo.StartBlockHeight == ctx.BlockHeight() {
		// started this height, no rewards yet
		return nil
	}

	startingPeriod := lockInfo.ReferencePeriod
	// calculate rewards for final period
	return k.calculateLockRewardsBetween(ctx, poolName, startingPeriod, endingPeriod, lockInfo.Amount)
}

// calculateLockRewardsBetween calculate the rewards accrued by a pool between two periods
func (k Keeper) calculateLockRewardsBetween(ctx sdk.Context, poolName string, startingPeriod, endingPeriod uint64,
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

// UpdateLockInfo updates lock info for a new lock info
func (k Keeper) UpdateLockInfo(ctx sdk.Context, addr sdk.AccAddress, poolName string, changedAmount sdk.Dec) {
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
	if lockInfo.Amount.IsZero() {
		k.DeleteLockInfo(ctx, lockInfo.Owner, lockInfo.PoolName)
		k.DeleteAddressInFarmPool(ctx, lockInfo.PoolName, lockInfo.Owner)
	} else {
		k.SetLockInfo(ctx, lockInfo)
		k.SetAddressInFarmPool(ctx, lockInfo.PoolName, lockInfo.Owner)
	}
}
