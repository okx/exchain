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
				// Currently, we support only one token of yield farming at the same time,
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

}
