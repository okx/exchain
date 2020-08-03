package v0_10

import sdk "github.com/cosmos/cosmos-sdk/types"

const ModuleName = "order"

type (
	// GenesisState - all slashing state that must be provided at genesis
	GenesisState struct {
		Params     Params   `json:"params"`
		OpenOrders []*Order `json:"open_orders"`
	}

	// nolint : order parameters
	Params struct {
		OrderExpireBlocks int64       `json:"order_expire_blocks"`
		MaxDealsPerBlock  int64       `json:"max_deals_per_block"`
		FeePerBlock       sdk.DecCoin `json:"fee_per_block"`
		TradeFeeRate      sdk.Dec     `json:"trade_fee_rate"`
	}

	// nolint
	Order struct {
		TxHash            string         `json:"txhash"`           // txHash of the place order tx
		OrderID           string         `json:"order_id"`         // order id
		Sender            sdk.AccAddress `json:"sender"`           // order maker address
		Product           string         `json:"product"`          // product for trading pair
		Side              string         `json:"side"`             // BUY/SELL
		Price             sdk.Dec        `json:"price"`            // price of the order
		Quantity          sdk.Dec        `json:"quantity"`         // quantity of the order
		Status            int64          `json:"status"`           // order status, see OrderStatusXXX
		FilledAvgPrice    sdk.Dec        `json:"filled_avg_price"` // filled average price
		RemainQuantity    sdk.Dec        `json:"remain_quantity"`  // Remaining quantity of the order
		RemainLocked      sdk.Dec        `json:"remain_locked"`    // Remaining locked quantity of token
		Timestamp         int64          `json:"timestamp"`        // created timestamp
		OrderExpireBlocks int64          `json:"order_expire_blocks"`
		FeePerBlock       sdk.DecCoin    `json:"fee_per_block"`
		ExtraInfo         string         `json:"extra_info"` // extra info of order in json format
	}
)
