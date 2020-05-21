package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/dex/types"
	ordertypes "github.com/okex/okchain/x/order/types"
)

// IsTokenPairLocked return true if token pair locked
func (k Keeper) IsTokenPairLocked(ctx sdk.Context, product string) bool {
	return k.getProductLock(ctx, product) != nil
}

func (k Keeper) getProductLock(ctx sdk.Context, product string) *ordertypes.ProductLock {
	store := ctx.KVStore(k.storeKey)
	productInfo := store.Get(types.GetLockProductKey(product))

	if productInfo == nil {
		return nil
	}
	productLock := &ordertypes.ProductLock{}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(productInfo, productLock)
	return productLock

	//l, ok := k.cache.lockMap.Data[product]
	//if ok {
	//	return l
	//}
	//return nil
}

// LockTokenPair locks token pair
func (k Keeper) LockTokenPair(ctx sdk.Context, product string, lock *ordertypes.ProductLock) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLockProductKey(product), k.cdc.MustMarshalBinaryLengthPrefixed(*lock))
	//k.cache.lockMap.Data[product] = lock
}

// UnlockTokenPair unlocks token pair
func (k Keeper) UnlockTokenPair(ctx sdk.Context, product string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLockProductKey(product))
	//delete(k.cache.lockMap.Data, product)
}

// LoadProductLocks loads product locked
func (k Keeper) LoadProductLocks(ctx sdk.Context) *ordertypes.ProductLockMap {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenPairLockKeyPrefix)
	defer iter.Close()

	lockMap := ordertypes.NewProductLockMap()

	for iter.Valid() {
		lock := &ordertypes.ProductLock{}
		lockBytes := iter.Value()
		k.cdc.MustUnmarshalBinaryLengthPrefixed(lockBytes, &lock)
		lockMap.Data[types.GetKey(iter)] = lock
		iter.Next()
	}
	return lockMap
}

// GetLockedProductsCopy returns deep copy of product locked
func (k Keeper) GetLockedProductsCopy(ctx sdk.Context) *ordertypes.ProductLockMap {
	//source := k.cache.lockMap
	//copy := ordertypes.NewProductLockMap()
	//
	//for k, v := range source.Data {
	//	copiedValue := *v
	//	copy.Data[k] = &copiedValue
	//}
	//return copy

	return k.LoadProductLocks(ctx)
}

// IsAnyProductLocked checks if any product is locked
func (k Keeper) IsAnyProductLocked(ctx sdk.Context) bool {
	return len(k.LoadProductLocks(ctx).Data) > 0
	//return len(k.cache.lockMap.Data) > 0
}
