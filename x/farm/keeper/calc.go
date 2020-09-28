package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm"
	"github.com/okex/okexchain/x/farm/types"
)

// LiquidateYieldTokenInfo is used for calculating how many tokens haven been yielding from LastClaimedBlockHeight to CurrentHeight
// Then transfer YieldedTokenInfos[i].RemainingAmount -> AmountYielded
func (k Keeper) LiquidateYieldTokenInfo(height int64, pool types.FarmPool) types.FarmPool {
	if height <= pool.LastClaimedBlockHeight { // TODO: is there any neccessary to make a height comparison?
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

// calcYieldAmount calculates the yielded amount which belongs to an account on a giving block height
func (k Keeper) calcYieldAmount(blockHeight int64, pool types.FarmPool, lockInfo types.LockInfo) (
	selfAmountYielded sdk.DecCoins, numerator sdk.Dec) {
	currentHeight := sdk.NewDec(blockHeight)
	/* 1.1 Calculate its own weight during these blocks
	   (curHeight - Height1) * Amount1
	*/
	oldWeight := sdk.NewDec(lockInfo.StartBlockHeight).MulTruncate(lockInfo.Amount.Amount)
	currentWeight := currentHeight.MulTruncate(lockInfo.Amount.Amount)
	numerator = currentWeight.Sub(oldWeight)

	/* 1.2 Calculate all weight during these blocks
	    (curHeight - Height1) * Amount1 + (curHeight - Height2) * Amount2 + (curHeight - Height3) * Amount3
												||
	                                            \/
	   curHeight * (Amount1 + Amount2 + Amount3) - (Height1*Amount1 + Height2*Amount2 + Height3*Amount3)
												||
	                                            \/
	ctx.BlockHeight()  *  pool.TotalValueLocked.Amount  -  ( pool.TotalLockedWeight )
	*/
	denominator := currentHeight.MulTruncate(pool.TotalValueLocked.Amount).Sub(pool.TotalLockedWeight)

	// 1.3 Calculate how many yielded tokens to return
	selfAmountYielded = pool.AmountYielded.MulDecTruncate(numerator).QuoDecTruncate(denominator)

	return
}

func (k Keeper) ClaimRewards(ctx sdk.Context, pool types.FarmPool, lockInfo types.LockInfo,
	address sdk.AccAddress, changedAmount sdk.Dec) sdk.Error {
	// 1. calculation
	height := ctx.BlockHeight()
	currentHeight := sdk.NewDec(height)
	selfAmountYielded, numerator := k.calcYieldAmount(height, pool, lockInfo)
	// 2. Transfer yielded tokens to personal account
	if !selfAmountYielded.IsZero() {
		if err := k.SupplyKeeper().SendCoinsFromModuleToAccount(ctx, farm.ModuleName, address, selfAmountYielded); err != nil {
			return err
		}
	}

	// 3. Update the pool data
	pool.AmountYielded = pool.AmountYielded.Sub(selfAmountYielded)
	if !changedAmount.IsZero() {
		pool.TotalValueLocked.Amount = pool.TotalValueLocked.Amount.Add(changedAmount)
		numerator = numerator.Add(currentHeight.MulTruncate(changedAmount))
	}
	pool.TotalLockedWeight = pool.TotalLockedWeight.Add(numerator)
	pool.LastClaimedBlockHeight = height
	// Set the updated pool into store
	k.SetFarmPool(ctx, pool)

	return nil
}
