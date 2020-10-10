package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

func (k Keeper) GetPoolHistoricalRewards(ctx sdk.Context, poolName string, period int64) (
	rewards types.PoolHistoricalRewards, found bool) {
	store := ctx.KVStore(k.StoreKey())
	bz := store.Get(types.GetPoolHistoricalRewardsKey(poolName, period))
	if bz == nil {
		return rewards, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewards)
	return rewards, true
}

func (k Keeper) SetPoolHistoricalRewards(ctx sdk.Context, poolName string, period int64, rewards types.PoolHistoricalRewards) {
	store := ctx.KVStore(k.StoreKey())
	store.Set(types.GetPoolHistoricalRewardsKey(poolName, period), k.cdc.MustMarshalBinaryLengthPrefixed(rewards))
	return
}

func (k Keeper) GetPoolCurrentPeriod(ctx sdk.Context, poolName string) (period types.PoolCurrentPeriod, found bool) {
	store := ctx.KVStore(k.StoreKey())
	bz := store.Get(types.GetPoolCurrentPeriodKey(poolName))
	if bz == nil {
		return period, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &period)
	return period, true
}

func (k Keeper) SetPoolCurrentPeriod(ctx sdk.Context, poolName string, period types.PoolCurrentPeriod) {
	store := ctx.KVStore(k.StoreKey())
	store.Set(types.GetPoolCurrentPeriodKey(poolName), k.cdc.MustMarshalBinaryLengthPrefixed(period))
	return
}
