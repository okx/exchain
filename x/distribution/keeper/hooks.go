package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"

	stakingtypes "github.com/okx/okbchain/x/staking/types"
)

// Hooks is struct of keepers from other module
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Hooks creates new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// AfterValidatorCreated initializes validator distribution record
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	val := h.k.stakingKeeper.Validator(ctx, valAddr)
	h.k.initializeValidator(ctx, val)
}

// AfterValidatorRemoved cleans up for after validator is removed
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.afterValidatorRemovedForDistributionProposal(ctx, nil, valAddr)
	return
}

// AfterValidatorDestroyed nothing to do
func (h Hooks) AfterValidatorDestroyed(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {

}

// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                         {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}
