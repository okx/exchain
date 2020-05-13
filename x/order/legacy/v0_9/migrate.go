// nolint
package v0_9

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	v08order "github.com/okex/okchain/x/order/legacy/v0_8"
	"github.com/okex/okchain/x/order/types"
)

// Migrate : migrate from v08 to v09
func Migrate(oldGenState v08order.GenesisState) GenesisState {
	params := types.Params{
		OrderExpireBlocks: oldGenState.Params.OrderExpireBlocks,
		MaxDealsPerBlock:  oldGenState.Params.MaxDealsPerBlock,
		FeePerBlock:       types.DefaultFeePerBlock,
		TradeFeeRate:      oldGenState.Params.TradeFeeRate,
	}

	orders := make([]*types.Order, 0, len(oldGenState.OpenOrders))
	for _, order := range oldGenState.OpenOrders {
		orders = append(orders, &types.Order{
			TxHash:            order.TxHash,
			OrderID:           order.OrderID,
			Sender:            order.Sender,
			Product:           order.Product,
			Side:              order.Side,
			Price:             order.Price,
			Quantity:          order.Quantity,
			Status:            order.Status,
			FilledAvgPrice:    order.FilledAvgPrice,
			RemainQuantity:    order.RemainQuantity,
			RemainLocked:      order.RemainLocked,
			Timestamp:         order.Timestamp,
			OrderExpireBlocks: params.OrderExpireBlocks,
			FeePerBlock:       sdk.NewDecCoinFromDec(common.NativeToken, sdk.ZeroDec()),
			ExtraInfo:         order.ExtraInfo})
	}

	return GenesisState{
		Params:     params,
		OpenOrders: orders,
	}
}
