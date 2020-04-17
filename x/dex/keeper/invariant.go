package keeper

import (
	"fmt"

	"github.com/okex/okchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/dex/types"
)

// RegisterInvariants registers all dex invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper IKeeper, supplyKeeper SupplyKeeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(keeper, supplyKeeper))
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// locks amounts held on store
func ModuleAccountInvariant(keeper IKeeper, supplyKeeper SupplyKeeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var depositsCoins = sdk.NewDecCoin(common.NativeToken, sdk.NewInt(0))
		withdrawCoins := depositsCoins

		// get product deposits
		for _, product := range keeper.GetTokenPairs(ctx) {
			depositsCoins = depositsCoins.Add(product.Deposits)
		}

		keeper.IterateWithdrawInfo(ctx, func(_ int64, withdrawInfo types.WithdrawInfo) (stop bool) {
			withdrawCoins = withdrawCoins.Add(withdrawInfo.Deposits)
			return false
		})

		moduleAcc := supplyKeeper.GetModuleAccount(ctx, types.ModuleName)

		broken := !moduleAcc.GetCoins().IsEqual(sdk.DecCoins{depositsCoins.Add(withdrawCoins)})

		return sdk.FormatInvariant(types.ModuleName, "module coins",
			fmt.Sprintf("\tdex ModuleAccount coins: %s\n\tsum of deposits coins: %s\tsum of withdraw coins: %s\n",
				moduleAcc.GetCoins(), depositsCoins, withdrawCoins)), broken
	}
}
