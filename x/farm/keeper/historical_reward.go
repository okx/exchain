package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/farm/types"
)

func (k Keeper) GetPoolHistoricalRewards(
	ctx sdk.Context, poolName string, period uint64,
) (rewards types.PoolHistoricalRewards) {
	store := ctx.KVStore(k.StoreKey())
	bz := store.Get(types.GetPoolHistoricalRewardsKey(poolName, period))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rewards)
	return rewards
}

func (k Keeper) SetPoolHistoricalRewards(
	ctx sdk.Context, poolName string, period uint64, rewards types.PoolHistoricalRewards,
) {
	store := ctx.KVStore(k.StoreKey())
	store.Set(types.GetPoolHistoricalRewardsKey(poolName, period), k.cdc.MustMarshalBinaryLengthPrefixed(rewards))
	return
}

// DeletePoolHistoricalReward deletes a historical reward
func (k Keeper) DeletePoolHistoricalReward(ctx sdk.Context, poolName string, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolHistoricalRewardsKey(poolName, period))
}

// IteratePoolHistoricalRewards deletes historical rewards for a pool
func (k Keeper) IteratePoolHistoricalRewards(
	ctx sdk.Context, poolName string, handler func(store sdk.KVStore, key []byte, value []byte) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetPoolHistoricalRewardsPrefix(poolName))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if handler(store, iter.Key(), iter.Value()) {
			break
		}
	}
}

func (k Keeper) GetPoolCurrentRewards(ctx sdk.Context, poolName string) (period types.PoolCurrentRewards) {
	store := ctx.KVStore(k.StoreKey())
	bz := store.Get(types.GetPoolCurrentRewardsKey(poolName))
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &period)
	return period
}

func (k Keeper) SetPoolCurrentRewards(ctx sdk.Context, poolName string, rewards types.PoolCurrentRewards) {
	store := ctx.KVStore(k.StoreKey())
	store.Set(types.GetPoolCurrentRewardsKey(poolName), k.cdc.MustMarshalBinaryLengthPrefixed(rewards))
	return
}

// delete current rewards for a pool
func (k Keeper) DeletePoolCurrentRewards(ctx sdk.Context, poolName string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetPoolCurrentRewardsKey(poolName))
}

// HasPoolCurrentRewards check existence of the pool_current_period associated with a pool_name
func (k Keeper) HasPoolCurrentRewards(ctx sdk.Context, poolName string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetPoolCurrentRewardsKey(poolName))
}

// Iterate over historical rewards
func (k Keeper) IterateAllPoolHistoricalRewards(
	ctx sdk.Context, handler func(poolName string, period uint64, rewards types.PoolHistoricalRewards) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PoolHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.PoolHistoricalRewards
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &rewards)
		poolName, period := GetPoolHistoricalRewardsPoolNamePeriod(iter.Key())
		if handler(poolName, period, rewards) {
			break
		}
	}
}

// gets the address & period from a validator's historical rewards key
func GetPoolHistoricalRewardsPoolNamePeriod(key []byte) (poolName string, period uint64) {
	name := key[1 : len(key)-types.PeriodByteArrayLength]
	if len(name) > types.MaxPoolNameLength {
		panic("unexpected key length")
	}

	if len(key) <= types.PeriodByteArrayLength+len(types.PoolHistoricalRewardsPrefix) {
		panic("unexpected key length")
	}
	b := key[len(key)-types.PeriodByteArrayLength:]
	period = binary.LittleEndian.Uint64(b)
	return string(name), period
}

// Iterate over current rewards
func (k Keeper) IterateAllPoolCurrentRewards(
	ctx sdk.Context, handler func(poolName string, curRewards types.PoolCurrentRewards) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.PoolCurrentRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.PoolCurrentRewards
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &rewards)
		poolName := string(iter.Key()[1:])
		if handler(poolName, rewards) {
			break
		}
	}
}
