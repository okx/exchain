package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

func (k Keeper) GetPoolHistoricalRewards(ctx sdk.Context, poolName string, period uint64) (rewards types.PoolHistoricalRewards) {
	store := ctx.KVStore(k.StoreKey())
	bz := store.Get(types.GetPoolHistoricalRewardsKey(poolName, period))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewards)
	return rewards
}

func (k Keeper) SetPoolHistoricalRewards(ctx sdk.Context, poolName string, period uint64, rewards types.PoolHistoricalRewards) {
	store := ctx.KVStore(k.StoreKey())
	store.Set(types.GetPoolHistoricalRewardsKey(poolName, period), k.cdc.MustMarshalBinaryLengthPrefixed(rewards))
	return
}

// DeletePoolHistoricalReward deletes a historical reward
func (k Keeper) DeletePoolHistoricalReward(ctx sdk.Context, poolName string, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolHistoricalRewardsKey(poolName, period))
}

// DeletePoolHistoricalRewards deletes historical rewards for a pool
func (k Keeper) DeletePoolHistoricalRewards(ctx sdk.Context, poolName string) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorHistoricalRewardsPrefix(poolName))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

func (k Keeper) GetPoolCurrentPeriod(ctx sdk.Context, poolName string) (period types.PoolCurrentPeriod) {
	store := ctx.KVStore(k.StoreKey())
	bz := store.Get(types.GetPoolCurrentPeriodKey(poolName))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &period)
	return period
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