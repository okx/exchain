package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/dex/types"
)

// RegisterInvariants registers all dex invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper IKeeper, supplyKeeper SupplyKeeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(keeper, supplyKeeper))
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// locks amounts held on store
func ModuleAccountInvariant(keeper IKeeper, supplyKeeper SupplyKeeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var depositsCoins, withdrawCoins sdk.SysCoins

		// get product deposits
		for _, product := range keeper.GetTokenPairs(ctx) {
			if product == nil {
				panic("the nil pointer is not expected")
			}
			depositsCoins = depositsCoins.Add(sdk.SysCoins{product.Deposits})
		}

		keeper.IterateWithdrawInfo(ctx, func(_ int64, withdrawInfo types.WithdrawInfo) (stop bool) {
			withdrawCoins = withdrawCoins.Add(sdk.SysCoins{withdrawInfo.Deposits})
			return false
		})

		moduleAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)

		broken := !moduleAcc.GetCoins().IsEqual(depositsCoins.Add(withdrawCoins))

		return sdk.FormatInvariant(types.ModuleName, "module coins",
			fmt.Sprintf("\tdex ModuleAccount coins: %s\n\tsum of deposits coins: %s\tsum of withdraw coins: %s\n",
				moduleAcc.GetCoins(), depositsCoins, withdrawCoins)), broken
	}
}
