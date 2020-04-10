package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/types"
	stakingtypes "github.com/okex/okchain/x/staking/types"
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
	// force-withdraw commission
	commission := h.k.GetValidatorAccumulatedCommission(ctx, valAddr)
	if !commission.IsZero() {
		// split into integral & remainder
		coins, remainder := commission.TruncateDecimal()
		if !remainder.IsZero() {
			err := h.k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, h.k.feeCollectorName, remainder)
			if err != nil {
				panic(err)
			}
		}

		accAddr := sdk.AccAddress(valAddr)
		withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, accAddr)
		// add to validator account
		if !coins.IsZero() {
			err := h.k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
			if err != nil {
				panic(err)
			}
		}
	}

	// remove commission record
	h.k.DeleteValidatorAccumulatedCommission(ctx, valAddr)
}

// AfterValidatorDestroyed nothing to do
func (h Hooks) AfterValidatorDestroyed(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {

}

// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                         {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}
