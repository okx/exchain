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

	// TODO update the query
	currentPeriod := k.GetPoolCurrentRewards(ctx, poolName)
	height := ctx.BlockHeight()
	_, yieldedTokens := CalculateAmountYieldedBetween(height, currentPeriod.StartBlockHeight, pool)

	// build return value
	earnings.TargetBlockHeight = height
	earnings.AmountLocked = lockInfo.Amount
	earnings.AmountYielded = currentPeriod.Rewards.Add(yieldedTokens)
	return earnings, nil
}
