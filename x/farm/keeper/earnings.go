package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// GetEarnings gets the earnings info by a given user address and a specific pool name
func (k Keeper) GetEarnings(ctx sdk.Context, poolName string, accAddr sdk.AccAddress) (types.Earnings, sdk.Error) {
	var earnings types.Earnings
	lockInfo, found := k.GetLockInfo(ctx, accAddr, poolName)
	if !found {
		return earnings, types.ErrNoLockInfoFound(types.DefaultCodespace, accAddr.String())
	}

	pool, found := k.GetFarmPool(ctx, poolName)
	if !found {
		return earnings, types.ErrNoFarmPoolFound(types.DefaultCodespace, poolName)
	}

	// calculate the yield amount of an account
	height := ctx.BlockHeight()
	updatedPool := k.LiquidateYieldTokenInfo(height, pool)
	selfAmountYielded, _ := k.calcYieldedAmount(height, updatedPool, lockInfo)

	// build return value
	earnings.TargetBlockHeight = height
	earnings.AmountLocked = lockInfo.Amount
	earnings.AmountYielded = selfAmountYielded
	return earnings, nil
}
