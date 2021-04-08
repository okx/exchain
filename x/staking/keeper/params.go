package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/params"
	"github.com/okex/okexchain/x/staking/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamKeyTable returns param table for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// UnbondingTime returns the param UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// MaxValidators returns the param Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// BondDenom returns the default the denomination of staking token
func (k Keeper) BondDenom(_ sdk.Context) string {
	return sdk.DefaultBondDenom
}

// GetParams gets all params as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.UnbondingTime(ctx),
		k.MaxValidators(ctx),
		k.ParamsEpoch(ctx),
		k.ParamsMaxValsToAddShares(ctx),
		k.ParamsMinDelegation(ctx),
		k.ParamsMinSelfDelegation(ctx),
	)
}

// SetParams sets the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// ParamsEpoch returns epoch from paramstore, only update the KeyEpoch after last epoch ends
func (k Keeper) ParamsEpoch(ctx sdk.Context) (res uint16) {
	k.paramstore.Get(ctx, types.KeyEpoch, &res)
	return
}

// GetEpoch returns the epoch for validators updates
func (k Keeper) GetEpoch(ctx sdk.Context) (epoch uint16) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyEpoch)
	if b == nil {
		return types.DefaultEpoch
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &epoch)
	return
}

// SetEpoch set epoch into keystore
func (k Keeper) SetEpoch(ctx sdk.Context, epoch uint16) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(epoch)
	store.Set(types.KeyEpoch, b)
}

// IsEndOfEpoch checks whether an epoch is end
func (k Keeper) IsEndOfEpoch(ctx sdk.Context) bool {
	return false
	blockInterval := ctx.BlockHeight() - k.GetTheEndOfLastEpoch(ctx)
	return blockInterval%int64(k.GetEpoch(ctx)) == 0
}

// GetTheEndOfLastEpoch returns the deadline of the current epoch
func (k Keeper) GetTheEndOfLastEpoch(ctx sdk.Context) (height int64) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyTheEndOfLastEpoch)
	if b == nil {
		return int64(0)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &height)
	return
}

// SetTheEndOfLastEpoch sets the deadline of the current epoch
func (k Keeper) SetTheEndOfLastEpoch(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(ctx.BlockHeight())
	store.Set(types.KeyTheEndOfLastEpoch, b)
}

// ParamsMaxValsToAddShares returns the param MaxValsToAddShares
func (k Keeper) ParamsMaxValsToAddShares(ctx sdk.Context) (num uint16) {
	k.paramstore.Get(ctx, types.KeyMaxValsToAddShares, &num)
	return
}

// ParamsMinDelegation returns the param MinDelegateAmount
func (k Keeper) ParamsMinDelegation(ctx sdk.Context) (num sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMinDelegation, &num)
	return
}

// ParamsMinSelfDelegation returns the param MinSelfDelegateAmount
func (k Keeper) ParamsMinSelfDelegation(ctx sdk.Context) (num sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMinSelfDelegation, &num)
	return
}
