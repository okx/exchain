package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/farm/types"
)

// GetEarnings gets the earnings info by a given user address and a specific pool name
func (k Keeper) GetEarnings(
	ctx sdk.Context, poolName string, accAddr sdk.AccAddress,
) (types.Earnings, sdk.Error) {
	var earnings types.Earnings
	lockInfo, found := k.GetLockInfo(ctx, accAddr, poolName)
	if !found {
		return earnings, types.ErrNoLockInfoFound(accAddr.String(), poolName)
	}

	pool, found := k.GetFarmPool(ctx, poolName)
	if !found {
		return earnings, types.ErrNoFarmPoolFound(poolName)
	}

	// 1.1 Calculate how many provided token & native token have been yielded
	// between start block height and current height
	updatedPool, yieldedTokens := k.CalculateAmountYieldedBetween(ctx, pool)

	endingPeriod := k.IncrementPoolPeriod(ctx, poolName, updatedPool.TotalValueLocked, yieldedTokens)
	rewards := k.calculateRewards(ctx, poolName, accAddr, endingPeriod, lockInfo)

	earnings = types.NewEarnings(ctx.BlockHeight(), lockInfo.Amount, rewards)
	return earnings, nil
}
