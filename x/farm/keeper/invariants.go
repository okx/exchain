package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// RegisterInvariants registers all distribution invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "yield-farming-account", YieldFarmingAccountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "mint-farming-account", MintFarmingAccountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "matched-address", MatchedAddressInvariant(k))
}

// ModuleAccountInvariant checks if farm ModuleAccount is consistent with the sum of deposit amount
func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// iterate all pools, then calculate the total deposit amount
		totalDepositAmount := sdk.DecCoins{}
		pools := k.GetFarmPools(ctx)
		for _, pool := range pools {
			totalDepositAmount = totalDepositAmount.Add(pool.DepositAmount.ToCoins())
		}

		// get farm module account
		moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName)
		amountInAmoduleAcc := moduleAcc.GetCoins()

		// make a comparison
		broken := !(totalDepositAmount.IsEqual(amountInAmoduleAcc))

		return sdk.FormatInvariant(types.ModuleName, "ModuleAccount coins",
			fmt.Sprintf("\texpected farm ModuleAccount coins: %s\n"+
				"\tacutal farm ModuleAccount coins: %s\n",
				totalDepositAmount, amountInAmoduleAcc)), broken
	}
}

// YieldFarmingAccountInvariant checks if yield_farming_account ModuleAccount is consistent with the total accumulated rewards
func YieldFarmingAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// iterate all pools, then calculate the total deposit amount
		totalRemainingYieldedAmount := sdk.DecCoins{}
		pools := k.GetFarmPools(ctx)
		for _, pool := range pools {
			totalRemainingYieldedAmount = totalRemainingYieldedAmount.Add(pool.TotalAccumulatedRewards)
		}

		// get yield_farming_account module account
		moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.YieldFarmingAccount)
		amountInAmoduleAcc := moduleAcc.GetCoins()

		// make a comparison
		broken := !(totalRemainingYieldedAmount.IsEqual(amountInAmoduleAcc))

		return sdk.FormatInvariant(types.ModuleName, "yield_farming_account coins",
			fmt.Sprintf("\texpected yield_farming_account coins: %s\n"+
				"\tacutal yield_farming_account coins: %s\n",
				totalRemainingYieldedAmount, amountInAmoduleAcc)), broken
	}
}

// MintFarmingAccountInvariant checks if mint_farming_account ModuleAccount is consistent with the sum of yielded native tokens
func MintFarmingAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		broken := false

		// get mint_farming_account module account
		moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.MintFarmingAccount)

		// get white_lists
		whiteLists := k.GetWhitelist(ctx)
		if len(whiteLists) != 0 {
			if !moduleAcc.GetCoins().IsZero() {
				broken = true
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "mint_farming_account coins",
			fmt.Sprintf("\texpected mint_farming_account coins should be zero\n"+
				"\tacutal mint_farming_account coins: %s\n"+
				"\twhite lists: %s\n",
				moduleAcc.GetCoins(), whiteLists)), broken
	}
}

// MatchedAddressInvariant
func MatchedAddressInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		return "", false
	}
}