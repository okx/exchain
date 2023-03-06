package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

// BeforeDelegationCreated - call hook if registered
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationCreated(ctx, delAddr, valAddrs)
	}
}

// BeforeDelegationSharesModified - call hook if registered
func (k Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddrs)
	}
}

// BeforeDelegationRemoved - call hook if registered
func (k Keeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	}
}

// AfterDelegationModified - call hook if registered
func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterDelegationModified(ctx, delAddr, valAddrs)
	}
}

//// BeforeValidatorSlashed - call hook if registered
//func (k Keeper) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
//	if k.hooks != nil {
//		k.hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
//	}
//}

// CheckEnabled - check modules enabled
func (k Keeper) CheckEnabled(ctx sdk.Context) bool {
	if k.hooks == nil {
		return true
	}

	return k.hooks.CheckEnabled(ctx)
}
