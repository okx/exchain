package v0_11

import (
	"github.com/okex/okexchain/x/order/legacy/v0_10"
)

func Migrate(oldGenState v0_10.GenesisState) GenesisState {
	params := Params{
		OrderExpireBlocks:     oldGenState.Params.OrderExpireBlocks,
		MaxDealsPerBlock:      oldGenState.Params.MaxDealsPerBlock,
		FeePerBlock:           oldGenState.Params.FeePerBlock,
		TradeFeeRate:          oldGenState.Params.TradeFeeRate,
		NewOrderMsgGasUnit:    DefaultNewOrderMsgGasUnit,
		CancelOrderMsgGasUnit: DefaultCancelOrderMsgGasUnit,
	}
	return GenesisState{
		Params:     params,
		OpenOrders: oldGenState.OpenOrders,
	}
}
