package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// getEarnings gets the earnings info by a given user address and a specific pool name
func (k Keeper) getEarnings(ctx sdk.Context, poolName string, accAddr sdk.AccAddress) (types.Earnings, sdk.Error) {
	var earnings types.Earnings
	lockInfo, found := k.GetLockInfo(ctx, accAddr, poolName)
	if !found {
		return earnings, types.ErrNoLockInfoFound(types.DefaultCodespace, accAddr.String())
	}

	pool, found := k.GetFarmPool(ctx, poolName)
	if !found {
		return earnings, types.ErrNoFarmPoolFound(types.DefaultCodespace, poolName)
	}

	_, endingPeriod := k.IncrementPoolPeriod(ctx, pool)
	rewards := k.calculateRewards(ctx, poolName, accAddr, endingPeriod)
	// calculate the yield amount of an account
	earnings.TargetBlockHeight = ctx.BlockHeight()
	earnings.AmountLocked = lockInfo.Amount
	earnings.AmountYielded = rewards
	return earnings, nil
}
