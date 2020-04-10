package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/types"
)

// GetDelegatorWithdrawAddr returns the delegator withdraw address, defaulting to the delegator address
func (k Keeper) GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetDelegatorWithdrawAddrKey(delAddr))
	if b == nil {
		return delAddr
	}
	return sdk.AccAddress(b)
}

// SetDelegatorWithdrawAddr sets the delegator withdraw address
func (k Keeper) SetDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDelegatorWithdrawAddrKey(delAddr), withdrawAddr.Bytes())
}

// DeleteDelegatorWithdrawAddr deletes a delegator withdraw addr
func (k Keeper) DeleteDelegatorWithdrawAddr(ctx sdk.Context, delAddr, _ sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetDelegatorWithdrawAddrKey(delAddr))
}

// IterateDelegatorWithdrawAddrs iterates over delegator withdraw addrs
func (k Keeper) IterateDelegatorWithdrawAddrs(ctx sdk.Context,
	handler func(del sdk.AccAddress, addr sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, DelegatorWithdrawAddrPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Value())
		del := GetDelegatorWithdrawInfoAddress(iter.Key())
		if handler(del, addr) {
			break
		}
	}
}

// GetPreviousProposerConsAddr returns the proposer public key for this block
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) (consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &consAddr)
	return consAddr
}

// SetPreviousProposerConsAddr sets the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(consAddr)
	store.Set(ProposerKey, b)
}

// GetValidatorAccumulatedCommission returns accumulated commission for a validator
func (k Keeper) GetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress) (
	commission types.ValidatorAccumulatedCommission) {

	store := ctx.KVStore(k.storeKey)
	b := store.Get(GetValidatorAccumulatedCommissionKey(val))
	if b == nil {
		return types.ValidatorAccumulatedCommission{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &commission)
	return commission
}

// SetValidatorAccumulatedCommission sets accumulated commission for a validator
func (k Keeper) SetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress,
	commission types.ValidatorAccumulatedCommission) {

	var bz []byte
	store := ctx.KVStore(k.storeKey)
	if commission.IsZero() {
		bz = k.cdc.MustMarshalBinaryLengthPrefixed(types.InitialValidatorAccumulatedCommission())
	} else {
		bz = k.cdc.MustMarshalBinaryLengthPrefixed(commission)
	}
	store.Set(GetValidatorAccumulatedCommissionKey(val), bz)
}

// DeleteValidatorAccumulatedCommission deletes accumulated commission for a validator
func (k Keeper) DeleteValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetValidatorAccumulatedCommissionKey(val))
}

// IterateValidatorAccumulatedCommissions iterates over accumulated commissions
func (k Keeper) IterateValidatorAccumulatedCommissions(ctx sdk.Context,
	handler func(val sdk.ValAddress, commission types.ValidatorAccumulatedCommission) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, ValidatorAccumulatedCommissionPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var commission types.ValidatorAccumulatedCommission
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &commission)
		addr := GetValidatorAccumulatedCommissionAddress(iter.Key())
		if handler(addr, commission) {
			break
		}
	}
}
