package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
)

// GetDelegator gets Delegator entity from store
func (k Keeper) GetDelegator(ctx sdk.Context, delAddr sdk.AccAddress) (delegator types.Delegator, found bool) {
	bytes := ctx.KVStore(k.storeKey).Get(types.GetDelegatorKey(delAddr))
	if bytes == nil {
		return delegator, false
	}

	delegator = types.MustUnMarshalDelegator(k.cdc, bytes)
	return delegator, true
}

// SetDelegator sets Delegator info to store
func (k Keeper) SetDelegator(ctx sdk.Context, delegator types.Delegator) {
	key := types.GetDelegatorKey(delegator.DelegatorAddress)
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(delegator)
	ctx.KVStore(k.storeKey).Set(key, bytes)
}

// DeleteDelegator deletes Delegator info from store
func (k Keeper) DeleteDelegator(ctx sdk.Context, delAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetDelegatorKey(delAddr))
}

// IterateDelegator iterates through all of the delegators info from the store
func (k Keeper) IterateDelegator(ctx sdk.Context, fn func(index int64, delegator types.Delegator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DelegatorKey)
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		var delegator types.Delegator
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &delegator)
		if stop := fn(i, delegator); stop {
			break
		}
		i++
	}
}
