package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/order/types"
	"github.com/okex/exchain/x/token"
)

// RegisterInvariants registers all order invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(keeper))
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// locks amounts held on store
func ModuleAccountInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var lockedCoins, lockedFees, orderLockedFees sdk.SysCoins

		for _, accCoins := range keeper.tokenKeeper.GetAllLockedCoins(ctx) {
			lockedCoins = lockedCoins.Add2(accCoins.Coins)
		}

		// lock fee
		keeper.tokenKeeper.IterateLockedFees(ctx, func(acc sdk.AccAddress, coins sdk.SysCoins) bool {
			lockedFees = lockedFees.Add2(coins)
			return false
		})

		// get open orders lock fee
		products := keeper.GetProductsFromDepthBookMap()
		for _, product := range products {
			depthBook := keeper.GetDepthBookCopy(product)
			for _, item := range depthBook.Items {
				buyKey := types.FormatOrderIDsKey(product, item.Price, types.BuyOrder)
				orderIDList := keeper.GetProductPriceOrderIDs(buyKey)
				sellKey := types.FormatOrderIDsKey(product, item.Price, types.SellOrder)
				orderIDList = append(orderIDList, keeper.GetProductPriceOrderIDs(sellKey)...)
				for _, orderID := range orderIDList {
					order := keeper.GetOrder(ctx, orderID)
					orderLockedFees = orderLockedFees.Add2(GetOrderNewFee(order))
				}
			}
		}

		if !lockedFees.IsEqual(orderLockedFees) {
			return sdk.FormatInvariant(types.ModuleName, "locks",
				fmt.Sprintf("\ttoken LockedFee coins: %s\n\tsum of order locked fee amounts:  %s\n",
					lockedFees, orderLockedFees)), true
		}

		macc := keeper.supplyKeeper.GetModuleAccount(ctx, token.ModuleName)
		broken := !macc.GetCoins().IsEqual(lockedCoins.Add2(lockedFees))
		return sdk.FormatInvariant(types.ModuleName, "locks",
			fmt.Sprintf("\ttoken ModuleAccount coins: %s\n\tsum of locks amounts:  %s\n",
				macc.GetCoins(), lockedCoins.Add2(lockedFees))), broken
	}
}
