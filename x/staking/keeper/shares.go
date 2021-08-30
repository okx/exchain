package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking/types"
)

// GetShares gets the shares entity
func (k Keeper) GetShares(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (types.Shares, bool) {
	store := ctx.KVStore(k.storeKey)
	sharesBytes := store.Get(types.GetSharesKey(valAddr, delAddr))
	var shares types.Shares
	// the delegator never adds shares to this val before
	if sharesBytes == nil {
		return shares, false
	}

	shares = types.MustUnmarshalShares(k.cdc, sharesBytes)
	return shares, true
}

// SetShares sets the shares that added to validators to store
func (k Keeper) SetShares(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares types.Shares) {
	key := types.GetSharesKey(valAddr, delAddr)
	sharesBytes := k.cdc.MustMarshalBinaryLengthPrefixed(shares)
	ctx.KVStore(k.storeKey).Set(key, sharesBytes)
}

// DeleteShares deletes shares entire from store
func (k Keeper) DeleteShares(ctx sdk.Context, valAddr sdk.ValAddress, delAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetSharesKey(valAddr, delAddr))
}

// GetValidatorAllShares returns all shares added to a specific validator and it's useful for querier
func (k Keeper) GetValidatorAllShares(ctx sdk.Context, valAddr sdk.ValAddress) types.SharesResponses {
	store := ctx.KVStore(k.storeKey)

	var sharesResps types.SharesResponses
	iterator := sdk.KVStorePrefixIterator(store, types.GetSharesToValidatorsKey(valAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// 1.get the delegator address
		delAddr := sdk.AccAddress(iterator.Key()[1+sdk.AddrLen:])

		// 2.get the shares
		shares := types.MustUnmarshalShares(k.cdc, iterator.Value())

		// 3.assemble the result
		sharesResps = append(sharesResps, types.NewSharesResponse(delAddr, shares))
	}

	return sharesResps
}

// IterateShares iterates through all of the shares from store
func (k Keeper) IterateShares(ctx sdk.Context, fn func(index int64, delAddr sdk.AccAddress, valAddr sdk.ValAddress,
	shares types.Shares) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SharesKey)
	defer iterator.Close()

	boundIndex := sdk.AddrLen + 1
	for i := int64(0); iterator.Valid(); iterator.Next() {
		// 1.get delegator/validator address from the key
		key := iterator.Key()
		valAddr, delAddr := sdk.ValAddress(key[1:boundIndex]), sdk.AccAddress(key[boundIndex:])

		// 2.get the shares
		var shares types.Shares
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &shares)

		// 3.call back the function
		if stop := fn(i, delAddr, valAddr, shares); stop {
			break
		}
		i++
	}
}

// GetDelegatorsByProxy returns all delegator addresses binding a proxy and it's useful for querier
func (k Keeper) GetDelegatorsByProxy(ctx sdk.Context, proxyAddr sdk.AccAddress) (delAddrs []sdk.AccAddress) {
	k.IterateProxy(ctx, proxyAddr, false, func(_ int64, delAddr, _ sdk.AccAddress) (stop bool) {
		delAddrs = append(delAddrs, delAddr)
		return false
	})

	return
}
