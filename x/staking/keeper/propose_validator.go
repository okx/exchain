package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/staking/types"
)

// GetProposeValidators gets proposed validators from db
func (k Keeper) GetProposeValidators(ctx sdk.Context) map[[sdk.AddrLen]byte]bool {
	proposeValidators := make(map[[sdk.AddrLen]byte]bool)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProposeValidatorsKey)
	defer iterator.Close()
	// iterate over the propose validators
	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		var isAdd bool
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		k.cdcMarshl.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &isAdd)
		proposeValidators[valAddr] = isAdd
	}
	return proposeValidators
}

// SetProposeValidator sets proposed validators to db
func (k Keeper) SetProposeValidator(ctx sdk.Context, operator sdk.ValAddress, isAdd bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdcMarshl.GetCdc().MustMarshalBinaryLengthPrefixed(isAdd)
	store.Set(types.GetProposeValidatorKey(operator), bz)
}

// DeleteProposeValidators deletes proposed validators to db
func (k Keeper) DeleteProposeValidators(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProposeValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		valAddr := make([]byte, sdk.AddrLen)
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		store.Delete(types.GetProposeValidatorKey(sdk.ValAddress(valAddr)))
	}
}
