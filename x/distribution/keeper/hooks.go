package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"github.com/okex/exchain/x/distribution/types"
	stakingtypes "github.com/okex/exchain/x/staking/types"
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
	if tmtypes.HigherThanVenus2(ctx.BlockHeight()) {
		h.newAfterValidatorRemoved(ctx, nil, valAddr)
		return
	}

	// force-withdraw commission
	commission := h.k.GetValidatorAccumulatedCommission(ctx, valAddr)
	if !commission.IsZero() {
		// split into integral & remainder
		coins, remainder := commission.TruncateDecimal()
		// remainder to community pool
		if !remainder.IsZero() {
			feePool := h.k.GetFeePool(ctx)
			feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
			h.k.SetFeePool(ctx, feePool)
		}
		// add to validator account
		if !coins.IsZero() {
			accAddr := sdk.AccAddress(valAddr)
			withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, accAddr)
			err := h.k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
			if err != nil {
				panic(err)
			}
		}
	}

	// remove commission record
	h.k.deleteValidatorAccumulatedCommission(ctx, valAddr)
}

func (h Hooks) newAfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {
	// fetch outstanding
	outstanding := h.k.GetValidatorOutstandingRewards(ctx, valAddr)

	// force-withdraw commission
	commission := h.k.GetValidatorAccumulatedCommission(ctx, valAddr)
	if !commission.IsZero() {
		// subtract from outstanding
		outstanding = outstanding.Sub(commission)

		// split into integral & remainder
		coins, remainder := commission.TruncateDecimal()
		// remainder to community pool
		if !remainder.IsZero() {
			feePool := h.k.GetFeePool(ctx)
			feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
			h.k.SetFeePool(ctx, feePool)
		}
		// add to validator account
		if !coins.IsZero() {
			accAddr := sdk.AccAddress(valAddr)
			withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, accAddr)
			err := h.k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
			if err != nil {
				panic(err)
			}
		}
	}

	// add outstanding to community pool
	feePool := h.k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(outstanding...)
	h.k.SetFeePool(ctx, feePool)

	// delete outstanding
	h.k.DeleteValidatorOutstandingRewards(ctx, valAddr)

	// remove commission record
	h.k.deleteValidatorAccumulatedCommission(ctx, valAddr)

	// clear slashes
	h.k.DeleteValidatorSlashEvents(ctx, valAddr)

	// clear historical rewards
	h.k.DeleteValidatorHistoricalRewards(ctx, valAddr)

	// clear current rewards
	h.k.DeleteValidatorCurrentRewards(ctx, valAddr)
}

// increment period
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	for _, valAddr := range valAddrs {
		val := h.k.stakingKeeper.Validator(ctx, valAddr)
		h.k.incrementValidatorPeriod(ctx, val)
	}
}

// withdraw delegation rewards (which also increments period)
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	if !tmtypes.HigherThanVenus2(ctx.BlockHeight()) {
		return
	}

	//del := h.k.stakingKeeper.Delegator(ctx, delAddr)
	for _, valAddr := range valAddrs {
		val := h.k.stakingKeeper.Validator(ctx, valAddr)
		if _, err := h.k.withdrawDelegationRewards(ctx, val, delAddr); err != nil {
			panic(err)
		}
	}
}

// create new delegation period record
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress) {
	if !tmtypes.HigherThanVenus2(ctx.BlockHeight()) {
		return
	}
	for _, valAddr := range valAddrs {
		h.k.initializeDelegation(ctx, valAddr, delAddr)
	}
}

//// record the slash event
//func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
//	h.k.updateValidatorSlashFraction(ctx, valAddr, fraction)
//}

// AfterValidatorDestroyed nothing to do
func (h Hooks) AfterValidatorDestroyed(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {

}

// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                         {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)       {}
