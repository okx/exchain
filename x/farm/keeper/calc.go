package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// CalculateAmountYieldedBetween is used for calculating how many tokens haven been yielding from LastClaimedBlockHeight to CurrentHeight
// Then transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
func CalculateAmountYieldedBetween(currentHeight int64, startBlockHeight int64, pool types.FarmPool) (types.FarmPool, sdk.DecCoins) {
	yieldedTokens := sdk.DecCoins{}
	for i := 0; i < len(pool.YieldedTokenInfos); i++ {
		startBlockHeightToYield := pool.YieldedTokenInfos[i].StartBlockHeightToYield
		if currentHeight > startBlockHeightToYield {
			// calculate the exact interval
			var blockInterval sdk.Dec
			if startBlockHeightToYield > startBlockHeight {
				blockInterval = sdk.NewDec(currentHeight - startBlockHeightToYield)
			} else {
				blockInterval = sdk.NewDec(currentHeight - startBlockHeight)
			}

			// calculate how many coin have been yielded till the current block
			amount := blockInterval.MulTruncate(pool.YieldedTokenInfos[i].AmountYieldedPerBlock)
			remaining := pool.YieldedTokenInfos[i].RemainingAmount
			if amount.LT(remaining.Amount) {
				// subtract yielded_coin amount
				pool.YieldedTokenInfos[i].RemainingAmount.Amount = remaining.Amount.Sub(amount)

				yieldedTokens = yieldedTokens.Add(sdk.NewDecCoinsFromDec(remaining.Denom, amount))
			} else {
				// TODO: remove the YieldedTokenInfo when its amount become zero
				// Currently, we support only one token of yield farming at the same time,
				// so, it is unnecessary to remove the element in slice

				// initialize yieldedTokenInfo
				pool.YieldedTokenInfos[i] = types.NewYieldedTokenInfo(sdk.NewDecCoin(remaining.Denom, sdk.ZeroInt()), 0, sdk.ZeroDec())

				yieldedTokens = yieldedTokens.Add(sdk.NewDecCoinsFromDec(remaining.Denom, remaining.Amount))
			}
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
	_ = k.incrementPoolPeriod(ctx, pool)

	// TODO get the rewards, then send rewards from module to account
	// TODO not sure where the rewards should be calculated?
	// TODO not sure if check the amount precision?

	return nil, nil
}

// increment pool period, returning the period just ended
func (k Keeper) incrementPoolPeriod(ctx sdk.Context, pool types.FarmPool) uint64 {
	// fetch current rewards status
	curReward := k.GetPoolCurrentRewards(ctx, pool.Name)

	// 1.1 calculate how many provided token has been yielded between start_block_height and current_height
	// TODO choose a perfect position to update pool, remember!
	updatedPool, yieldedTokens := CalculateAmountYieldedBetween(ctx.BlockHeight(), curReward.StartBlockHeight, pool)

	// 1.2 calculate how many native token has been yielded between start_block_height and current_height
	curReward.AccumulatedRewards = curReward.AccumulatedRewards.Add(yieldedTokens)

	currentRatio := sdk.DecCoins{}
	if !curReward.AccumulatedRewards.IsZero() { // warning: can't calculate ratio for zero-token
		// 2. calculate current reward ratio
		currentRatio = curReward.AccumulatedRewards.QuoDecTruncate(pool.TotalValueLocked.Amount)
	}

	// 3.1 get the previous pool_historical_rewards
	oldHistoricalRewards := k.GetPoolHistoricalRewards(ctx, pool.Name, curReward.Period-1)
	// 3.2 decrement reference count
	k.decrementReferenceCount(ctx, pool.Name, curReward.Period-1)
	// 3.3 create new pool_historical_rewards with reference count of 1, then set it into store
	newHistoricalRewards := types.NewPoolHistoricalRewards(oldHistoricalRewards.CumulativeRewardRatio.Add(currentRatio),1)
	k.SetPoolHistoricalRewards(ctx, pool.Name, curReward.Period, newHistoricalRewards)

	// 4. set new current newYieldedRewards into store, incrementing period by 1
	newRewards := types.NewPoolCurrentRewards(ctx.BlockHeight(), curReward.Period+1, sdk.DecCoins{})
	k.SetPoolCurrentRewards(ctx, pool.Name, newRewards)

	// 5. set updated pool
	k.SetFarmPool(ctx, updatedPool)

	return curReward.Period
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

func (k Keeper) calculateRewards(ctx sdk.Context, poolName string, endingPeriod uint64)  sdk.DecCoins {
	// fetch current period
	currentPeriod := k.GetPoolCurrentRewards(ctx, poolName)
	startingPeriod := currentPeriod.Period
	lastAmountYielded := currentPeriod.AccumulatedRewards
	// calculate rewards for final period
	return k.calculateDelegationRewardsBetween(ctx, poolName, startingPeriod, endingPeriod, lastAmountYielded)
}

// calculate the rewards accrued by a pool between two periods
func (k Keeper) calculateDelegationRewardsBetween(ctx sdk.Context, poolName string, startingPeriod, endingPeriod uint64,
	amount sdk.DecCoins) (rewards sdk.DecCoins) {

	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// return amount * (ending - starting)
	starting := k.GetPoolHistoricalRewards(ctx, poolName, startingPeriod)
	ending := k.GetPoolHistoricalRewards(ctx, poolName, endingPeriod)
	differences := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if differences.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	for _, difference := range differences {
		difference.Amount = difference.Amount.MulTruncate(amount.AmountOf(difference.Denom))
		rewards = append(rewards, difference)
	}
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