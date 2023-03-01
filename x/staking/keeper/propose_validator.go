package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking/types"
)

// GetProposeValidators gets proposed validators from db
func (k Keeper) GetProposeValidators(ctx sdk.Context) (proposeValidators map[[sdk.AddrLen]byte]bool, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.ProposeValidatorsKey)
	if value == nil {
		return nil, false
	}
	k.cdcMarshl.GetCdc().MustUnmarshalBinaryLengthPrefixed(value, &proposeValidators)
	return proposeValidators, true
}

// SetProposeValidators sets proposed validators to db
func (k Keeper) SetProposeValidators(ctx sdk.Context, proposeValidators map[[sdk.AddrLen]byte]bool) {
	store := ctx.KVStore(k.storeKey)
	bytes := k.cdcMarshl.GetCdc().MustMarshalBinaryLengthPrefixed(proposeValidators)
	store.Set(types.ProposeValidatorsKey, bytes)
}

// DeleteProposeValidators deletes proposed validators to db
func (k Keeper) DeleteProposeValidators(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ProposeValidatorsKey)
}
