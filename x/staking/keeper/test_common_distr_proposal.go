package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	//distr "github.com/okex/exchain/x/distribution"
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
