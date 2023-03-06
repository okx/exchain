package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

func (dk mockDistributionKeeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
}
func (dk mockDistributionKeeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
}
func (dk mockDistributionKeeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (dk mockDistributionKeeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
}
func (dk mockDistributionKeeper) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
}
func (dk mockDistributionKeeper) CheckEnabled(ctx sdk.Context) bool { return true }
