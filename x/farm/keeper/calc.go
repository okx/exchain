package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// CalculateAmountYieldedBetween is used for calculating how many tokens haven been yielded from
// startBlockHeight to endBlockHeight. And return the amount.
func (k Keeper) CalculateAmountYieldedBetween(ctx sdk.Context, pool types.FarmPool) (types.FarmPool, sdk.SysCoins) {
	currentPeriod := k.GetPoolCurrentRewards(ctx, pool.Name)
	endBlockHeight := ctx.BlockHeight()

	totalYieldedTokens := sdk.SysCoins{}
	for i := 0; i < len(pool.YieldedTokenInfos); i++ {
		startBlockHeightToYield := pool.YieldedTokenInfos[i].StartBlockHeightToYield
		var startBlockHeight int64
		if currentPeriod.StartBlockHeight <= startBlockHeightToYield {
			startBlockHeight = startBlockHeightToYield
		} else {
			startBlockHeight = currentPeriod.StartBlockHeight
		}

		// no tokens to yield
		if startBlockHeightToYield == 0 || startBlockHeight >= endBlockHeight {
			continue
		}

		yieldedTokens := sdk.SysCoins{}
		// calculate how many tokens to be yielded between startBlockHeight and endBlockHeight
		blockInterval := sdk.NewDec(endBlockHeight - startBlockHeight)
		amount := blockInterval.MulTruncate(pool.YieldedTokenInfos[i].AmountYieldedPerBlock)
		remaining := pool.YieldedTokenInfos[i].RemainingAmount
		if amount.LT(remaining.Amount) {
			pool.YieldedTokenInfos[i].RemainingAmount.Amount = remaining.Amount.Sub(amount)
			yieldedTokens = sdk.NewDecCoinsFromDec(remaining.Denom, amount)
		} else {
			pool.YieldedTokenInfos[i] = types.NewYieldedTokenInfo(
				sdk.NewDecCoin(remaining.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec(),
			)
			yieldedTokens = sdk.NewDecCoinsFromDec(remaining.Denom, remaining.Amount)
		}
		pool.TotalAccumulatedRewards = pool.TotalAccumulatedRewards.Add(yieldedTokens)
		totalYieldedTokens = totalYieldedTokens.Add(yieldedTokens)
	}
	return pool, totalYieldedTokens
}

func (k Keeper) WithdrawRewards(
	ctx sdk.Context, poolName string, totalValueLocked sdk.SysCoin, yieldedTokens sdk.SysCoins, addr sdk.AccAddress,
) (sdk.SysCoins, sdk.Error) {
	// 0. check existence of lock info
	lockInfo, found := k.GetLockInfo(ctx, addr, poolName)
	if !found {
		return nil, types.ErrNoLockInfoFound(types.DefaultCodespace, addr.String(), poolName)
	}

	// 1. end current period and calculate rewards
	endingPeriod := k.IncrementPoolPeriod(ctx, poolName, totalValueLocked, yieldedTokens)
	rewards := k.calculateRewards(ctx, poolName, addr, endingPeriod, lockInfo)

	// 2. transfer rewards to user account
	if !rewards.IsZero() {
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.YieldFarmingAccount, addr, rewards)
		if err != nil {
			return nil, err
		}
	}

	// 3. decrement reference count of lock info
	k.decrementReferenceCount(ctx, poolName, lockInfo.ReferencePeriod)

	return rewards, nil
}

// IncrementPoolPeriod increments pool period, returning the period just ended
func (k Keeper) IncrementPoolPeriod(
	ctx sdk.Context, poolName string, totalValueLocked sdk.SysCoin, yieldedTokens sdk.SysCoins,
) uint64 {
	// 1. fetch current period rewards
	rewards := k.GetPoolCurrentRewards(ctx, poolName)
	// 2. calculate current reward ratio
	rewards.Rewards = rewards.Rewards.Add(yieldedTokens)
	var currentRatio sdk.SysCoins
	if totalValueLocked.IsZero() {
		currentRatio = sdk.SysCoins{}
	} else {
		currentRatio = rewards.Rewards.QuoDecTruncate(totalValueLocked.Amount)
	}

	// 3.1 get the previous pool historical rewards
	historical := k.GetPoolHistoricalRewards(ctx, poolName, rewards.Period-1).CumulativeRewardRatio
	// 3.2 decrement reference count
	k.decrementReferenceCount(ctx, poolName, rewards.Period-1)
	// 3.3 create new pool historical rewards with reference count of 1, then set it into store
	newHistoricalRewards := types.NewPoolHistoricalRewards(historical.Add(currentRatio), 1)
	k.SetPoolHistoricalRewards(ctx, poolName, rewards.Period, newHistoricalRewards)

	// 4. set new current rewards into store, incrementing period by 1
	newCurRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), rewards.Period+1, sdk.SysCoins{})
	k.SetPoolCurrentRewards(ctx, poolName, newCurRewards)

	return rewards.Period
}

// incrementReferenceCount increments the reference count for a historical rewards value
func (k Keeper) incrementReferenceCount(ctx sdk.Context, poolName string, period uint64) {
	historical := k.GetPoolHistoricalRewards(ctx, poolName, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetPoolHistoricalRewards(ctx, poolName, period, historical)
}

// decrementReferenceCount decrements the reference count for a historical rewards value,
// and delete if zero references remain.
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

func (k Keeper) calculateRewards(
	ctx sdk.Context, poolName string, addr sdk.AccAddress, endingPeriod uint64, lockInfo types.LockInfo,
) (rewards sdk.SysCoins) {
	if lockInfo.StartBlockHeight == ctx.BlockHeight() {
		// started this height, no rewards yet
		return
	}

	startingPeriod := lockInfo.ReferencePeriod
	// calculate rewards for final period
	return k.calculateLockRewardsBetween(ctx, poolName, startingPeriod, endingPeriod, lockInfo.Amount)
}

// calculateLockRewardsBetween calculate the rewards accrued by a pool between two periods
func (k Keeper) calculateLockRewardsBetween(ctx sdk.Context, poolName string, startingPeriod, endingPeriod uint64,
	amount sdk.SysCoin) (rewards sdk.SysCoins) {

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
	rewards = difference.MulDecTruncate(amount.Amount)
	return
}

// UpdateLockInfo updates lock info for the modified lock info
func (k Keeper) UpdateLockInfo(ctx sdk.Context, addr sdk.AccAddress, poolName string, changedAmount sdk.Dec) {
	// period has already been incremented - we want to store the period ended by this lock action
	previousPeriod := k.GetPoolCurrentRewards(ctx, poolName).Period - 1

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
		// increment reference count for the period we're going to track
		k.incrementReferenceCount(ctx, poolName, previousPeriod)

		// set the updated lock info
		k.SetLockInfo(ctx, lockInfo)
		k.SetAddressInFarmPool(ctx, lockInfo.PoolName, lockInfo.Owner)
	}
}
