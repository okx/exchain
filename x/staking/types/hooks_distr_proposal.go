package types

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

func (h MultiStakingHooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	for i := range h {
		h[i].BeforeDelegationCreated(ctx, delAddr, valAddrs)
	}
}

func (h MultiStakingHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	for i := range h {
		h[i].BeforeDelegationSharesModified(ctx, delAddr, valAddrs)
	}
}

func (h MultiStakingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeDelegationRemoved(ctx, delAddr, valAddr)
	}
}

func (h MultiStakingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	for i := range h {
		h[i].AfterDelegationModified(ctx, delAddr, valAddrs)
	}
}

//func (h MultiStakingHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
//	for i := range h {
//		h[i].BeforeValidatorSlashed(ctx, valAddr, fraction)
//	}
//}

func (h MultiStakingHooks) CheckEnabled(ctx sdk.Context) bool {
	for i := range h {
		if !h[i].CheckEnabled(ctx) {
			return false
		}
	}

	return true
}
