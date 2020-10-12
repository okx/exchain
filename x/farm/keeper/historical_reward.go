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

// DeleteValidatorHistoricalReward deletes a historical reward
func (k Keeper) DeleteValidatorHistoricalReward(ctx sdk.Context, poolName string, period int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolHistoricalRewardsKey(poolName, period))
}

// DeleteValidatorHistoricalReward deletes historical rewards for a pool
func (k Keeper) DeleteValidatorHistoricalRewards(ctx sdk.Context, poolName string) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorHistoricalRewardsPrefix(poolName))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
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

// HasPoolCurrentPeriod check existence of the pool_current_period associated with a pool_name
func (k Keeper) HasPoolCurrentPeriod(ctx sdk.Context, poolName string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetPoolCurrentPeriodKey(poolName))
}