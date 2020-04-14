package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/types"
	"github.com/okex/okchain/x/staking/exported"
)

// RegisterInvariants registers all distribution invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "nonnegative-commission", NonNegativeCommissionsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "can-withdraw", CanWithdrawInvariant(k))
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(k))
}

// NonNegativeCommissionsInvariant checks that accumulated commissions unwithdrawned fees are never negative
func NonNegativeCommissionsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int
		var commission sdk.DecCoins

		k.IterateValidatorAccumulatedCommissions(ctx,
			func(addr sdk.ValAddress, c types.ValidatorAccumulatedCommission) (stop bool) {
				commission = c
				if commission.IsAnyNegative() {
					count++
					msg += fmt.Sprintf("\t%v has negative accumulated commission coins: %v\n", addr, commission)
				}
				return false
			})
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "nonnegative accumulated commission",
			fmt.Sprintf("found %d validators with negative accumulated commission\n%s", count, msg)), broken
	}
}

// CanWithdrawInvariant checks that current commission can be completely withdrawn
func CanWithdrawInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		// iterate over all validators
		k.stakingKeeper.IterateValidators(ctx, func(_ int64, val exported.ValidatorI) (stop bool) {
			valAddr := val.GetOperator()
			accumCommission := k.GetValidatorAccumulatedCommission(ctx, valAddr)
			if accumCommission.IsZero() {
				return false
			}
			if _, err := k.WithdrawValidatorCommission(ctx, valAddr); err != nil {
				count++
				msg += fmt.Sprintf("\t%v failed to withdraw accumulated commission coins: %v. error: %v\n",
					valAddr, accumCommission, err)
			}
			return false
		})

		broken := count != 0
		return sdk.FormatInvariant(types.ModuleName, "withdraw commission", msg), broken
	}
}

// ModuleAccountInvariant checks that the coins held by the distr ModuleAccount
// is consistent with the sum of accumulated commissions
func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var accumulatedCommission sdk.DecCoins
		k.IterateValidatorAccumulatedCommissions(ctx,
			func(_ sdk.ValAddress, commission types.ValidatorAccumulatedCommission) (stop bool) {
				accumulatedCommission = accumulatedCommission.Add(commission)
				return false
			})
		communityPool := k.GetFeePoolCommunityCoins(ctx)
		macc := k.GetDistributionAccount(ctx)
		broken := !macc.GetCoins().IsEqual(communityPool.Add(accumulatedCommission))
		return sdk.FormatInvariant(types.ModuleName, "ModuleAccount coins",
			fmt.Sprintf("\texpected distribution ModuleAccount coins:     %s\n"+
				"\tacutal distribution ModuleAccount coins: %s\n",
				accumulatedCommission, macc.GetCoins())), broken
	}
}
