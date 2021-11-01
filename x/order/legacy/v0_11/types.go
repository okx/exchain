package v0_11

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/order/legacy/v0_10"
)

const (
	ModuleName                   = "order"
	DefaultNewOrderMsgGasUnit    = 40000
	DefaultCancelOrderMsgGasUnit = 30000
)

type (
	// GenesisState - all order state that must be provided at genesis
	GenesisState struct {
		Params     Params         `json:"params"`
		OpenOrders []*v0_10.Order `json:"open_orders"`
	}

	// nolint : order parameters
	Params struct {
		OrderExpireBlocks     int64       `json:"order_expire_blocks"`
		MaxDealsPerBlock      int64       `json:"max_deals_per_block"`
		FeePerBlock           sdk.SysCoin `json:"fee_per_block"`
		TradeFeeRate          sdk.Dec     `json:"trade_fee_rate"`
		NewOrderMsgGasUnit    uint64      `json:"new_order_msg_gas_unit"`
		CancelOrderMsgGasUnit uint64      `json:"cancel_order_msg_gas_unit"`
	}
)
