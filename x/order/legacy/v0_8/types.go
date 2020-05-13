// nolint
package v0_8

import (
	"github.com/okex/okchain/x/order/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
const (
	ModuleName = types.ModuleName
)

// nolint
type (
	Params struct {
		OrderExpireBlocks   int64   `json:"order_expire_blocks"`
		MaxDealsPerBlock    int64   `json:"max_deals_per_block"`
		NewOrder            sdk.Dec `json:"new_order"`
		CancelFeeRateNative sdk.Dec `json:"cancel_fee_rate_native"`
		ExpireFeeRateNative sdk.Dec `json:"expire_fee_rate_native"`
		TradeFeeRate        sdk.Dec `json:"trade_fee_rate"`
		TradeFeeRateNative  sdk.Dec `json:"trade_fee_rate_native"`
	}

	Order struct {
		TxHash         string         `json:"txHash"`         // txHash of the place order tx
		OrderID        string         `json:"orderId"`        // order id
		Sender         sdk.AccAddress `json:"sender"`         // order maker address
		Product        string         `json:"product"`        // product for trading pair
		Side           string         `json:"side"`           // BUY/SELL
		Price          sdk.Dec        `json:"price"`          // price of the order
		Quantity       sdk.Dec        `json:"quantity"`       // quantity of the order
		Status         int64          `json:"status"`         // order status, see OrderStatusXXX
		FilledAvgPrice sdk.Dec        `json:"filledAvgPrice"` // filled average price
		RemainQuantity sdk.Dec        `json:"remainQuantity"` // Remaining quantity of the order
		RemainLocked   sdk.Dec        `json:"remainLocked"`   // Remaining locked quantity of token
		Timestamp      int64          `json:"timestamp"`      // created timestamp
		ExtraInfo      string         `json:"extraInfo"`      // extra info of order in json format
	}

	GenesisState struct {
		Params     Params   `json:"params"`
		OpenOrders []*Order `json:"open_orders"`
	}
)
