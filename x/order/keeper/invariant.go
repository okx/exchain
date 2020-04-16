package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/order/types"
	"github.com/okex/okchain/x/token"
)

// RegisterInvariants registers all order invariants
func RegisterInvariants(ir sdk.InvariantRegistry, keeper Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-account", ModuleAccountInvariant(keeper))
}

// ModuleAccountInvariant checks that the module account coins reflects the sum of
// locks amounts held on store
func ModuleAccountInvariant(keeper Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var lockedCoins, lockedFee, orderLockedFee sdk.DecCoins

		for _, accCoins := range keeper.tokenKeeper.GetAllLockedCoins(ctx) {
			lockedCoins = lockedCoins.Add(accCoins.Coins)
		}

		// lock fee
		keeper.tokenKeeper.IterateLockedFee(ctx, func(acc sdk.AccAddress, coins sdk.DecCoins) bool {
			lockedFee = lockedFee.Add(coins)
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
					orderLockedFee = orderLockedFee.Add(GetOrderNewFee(order))
				}
			}
		}

		if !lockedFee.IsEqual(orderLockedFee) {
			return sdk.FormatInvariant(types.ModuleName, "locks",
				fmt.Sprintf("\ttoken LockedFee coins: %s\n\tsum of order locked fee amounts:  %s\n",
					lockedFee, orderLockedFee)), true
		}

		macc := keeper.supplyKeeper.GetModuleAccount(ctx, token.ModuleName)
		broken := !macc.GetCoins().IsEqual(lockedCoins.Add(lockedFee))
		return sdk.FormatInvariant(types.ModuleName, "locks",
			fmt.Sprintf("\ttoken ModuleAccount coins: %s\n\tsum of locks amounts:  %s\n",
				macc.GetCoins(), lockedCoins.Add(lockedFee))), broken
	}
}
