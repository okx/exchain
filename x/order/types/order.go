package types

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OrderStatus int

const (
	Open OrderStatus = iota
	Filled
	Cancelled
	Expired
	PartialFilledCancelled
	PartialFilledExpired
)

func (p OrderStatus) String() string {
	switch p {
	case Open:
		return "Open"
	case Filled:
		return "Filled"
	case Cancelled:
		return "Cancelled"
	case Expired:
		return "Expired"
	case PartialFilledCancelled:
		return "PartialFilledCancelled"
	case PartialFilledExpired:
		return "PartialFilledExpired"
	default:
		return "Unknown"
	}
}

const (
	OrderStatusOpen                   = 0
	OrderStatusFilled                 = 1
	OrderStatusCancelled              = 2
	OrderStatusExpired                = 3
	OrderStatusPartialFilledCancelled = 4
	OrderStatusPartialFilledExpired   = 5
	OrderStatusPartialFilled          = 6
)

const (
	OrderExtraInfoKeyNewFee     = "newFee"
	OrderExtraInfoKeyCancelFee  = "cancelFee"
	OrderExtraInfoKeyExpireFee  = "expireFee"
	OrderExtraInfoKeyDealFee    = "dealFee"
	OrderExtraInfoKeyReceiveFee = "receiveFee"
)

type Order struct {
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

func NewOrder(txHash string, sender sdk.AccAddress, product, side string, price, quantity sdk.Dec,
	timestamp int64, orderExpireBlocks int64, feePerBlock sdk.DecCoin) *Order {
	order := &Order{
		TxHash:            txHash,
		Sender:            sender,
		Product:           product,
		Side:              side,
		Price:             price,
		Quantity:          quantity,
		Status:            OrderStatusOpen,
		RemainQuantity:    quantity,
		Timestamp:         timestamp,
		OrderExpireBlocks: orderExpireBlocks,
		FeePerBlock:       feePerBlock,
	}
	if side == BuyOrder {
		order.RemainLocked = price.Mul(quantity)
	} else {
		order.RemainLocked = quantity
	}
	return order
}

func (order *Order) String() string {
	if orderJSON, err := json.Marshal(order); err != nil {
		panic(err)
	} else {
		return string(orderJSON)
	}
}

func (order *Order) SetExtraInfo(extra map[string]string) {
	if extra != nil {
		bz, err := json.Marshal(extra)
		if err != nil {
			panic(err)
		}
		order.ExtraInfo = string(bz)
	}
}

func (order *Order) GetExtraInfo() map[string]string {
	extra := make(map[string]string)
	if order.ExtraInfo != "" && order.ExtraInfo != "{}" {
		if err := json.Unmarshal([]byte(order.ExtraInfo), &extra); err != nil {
			log.Printf("Unmarshal order extra info(%s) failed\n", order.ExtraInfo)
		}
	}
	return extra
}

func (order *Order) SetExtraInfoWithKeyValue(key, value string) {
	extra := order.GetExtraInfo()
	extra[key] = value
	order.SetExtraInfo(extra)
}

func (order *Order) GetExtraInfoWithKey(key string) string {
	extra := order.GetExtraInfo()
	if value, ok := extra[key]; ok {
		return value
	}
	return ""
}

func (order *Order) RecordOrderNewFee(fee sdk.DecCoins) {
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyNewFee, fee.String())
}

func (order *Order) RecordOrderCancelFee(fee sdk.DecCoins) {
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyCancelFee, fee.String())
}

func (order *Order) RecordOrderExpireFee(fee sdk.DecCoins) {
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyExpireFee, fee.String())
}

func (order *Order) RecordOrderReceiveFee(fee sdk.DecCoins) {
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyReceiveFee, fee.String())
}

// An order may have several deals
func (order *Order) RecordOrderDealFee(fee sdk.DecCoins) {
	oldValue := order.GetExtraInfoWithKey(OrderExtraInfoKeyDealFee)
	if oldValue == "" {
		order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyDealFee, fee.String())
		return
	}
	oldFee, err := sdk.ParseDecCoins(oldValue)
	if err != nil {
		log.Println(err)
		return
	}
	newFee := oldFee.Add(fee)
	order.SetExtraInfoWithKeyValue(OrderExtraInfoKeyDealFee, newFee.String())
}

func (order *Order) Fill(price, fillAmount sdk.Dec) {
	filledSum := order.FilledAvgPrice.Mul(order.Quantity.Sub(order.RemainQuantity))
	newFilledSum := filledSum.Add(price.Mul(fillAmount))
	order.RemainQuantity = order.RemainQuantity.Sub(fillAmount)
	order.FilledAvgPrice = newFilledSum.Quo(order.Quantity.Sub(order.RemainQuantity))
	if order.Side == BuyOrder {
		order.RemainLocked = order.RemainLocked.Sub(price.Mul(fillAmount))
	} else {
		order.RemainLocked = order.RemainLocked.Sub(fillAmount)
	}
	if order.RemainQuantity.IsZero() {
		order.Status = OrderStatusFilled
	}
}

func (order *Order) Cancel() {
	if order.RemainQuantity.Equal(order.Quantity) {
		order.Status = OrderStatusCancelled
	} else {
		order.Status = OrderStatusPartialFilledCancelled
	}
}

func (order *Order) Expire() {
	if order.RemainQuantity.Equal(order.Quantity) {
		order.Status = OrderStatusExpired
	} else {
		order.Status = OrderStatusPartialFilledExpired
	}
}

// when place a new order, we should lock the coins of sender
func (order *Order) NeedLockCoins() sdk.DecCoins {
	if order.Side == BuyOrder {
		token := strings.Split(order.Product, "_")[1]
		amount := order.Price.Mul(order.Quantity)
		return sdk.DecCoins{{Denom: token, Amount: amount}}
	}
	token := strings.Split(order.Product, "_")[0]
	amount := order.Quantity
	return sdk.DecCoins{{Denom: token, Amount: amount}}

}

// when order be cancelled/expired, we should unlock the coins of sender
func (order *Order) NeedUnlockCoins() sdk.DecCoins {
	if order.Side == BuyOrder {
		token := strings.Split(order.Product, "_")[1]
		return sdk.DecCoins{{Denom: token, Amount: order.RemainLocked}}
	}
	token := strings.Split(order.Product, "_")[0]
	return sdk.DecCoins{{Denom: token, Amount: order.RemainLocked}}

}

func (order *Order) Unlock() {
	order.RemainLocked = sdk.ZeroDec()
}

func MockOrder(orderID, product, side, price, quantity string) *Order {
	order := &Order{
		OrderID:           orderID,
		Product:           product,
		Side:              side,
		Price:             sdk.MustNewDecFromStr(price),
		FilledAvgPrice:    sdk.ZeroDec(),
		Quantity:          sdk.MustNewDecFromStr(quantity),
		RemainQuantity:    sdk.MustNewDecFromStr(quantity),
		Status:            OrderStatusOpen,
		OrderExpireBlocks: DefaultOrderExpireBlocks,
		FeePerBlock:       DefaultFeePerBlock,
	}
	if side == BuyOrder {
		order.RemainLocked = order.Price.Mul(order.Quantity)
	} else {
		order.RemainLocked = order.Quantity
	}
	return order
}

func FormatOrderID(blockHeight, orderNum int64) string {
	format := "ID%010d-%d"
	if blockHeight > 9999999999 {
		format = "ID%d-%d"
	}
	return fmt.Sprintf(format, blockHeight, orderNum)
}

func GetBlockHeightFromOrderID(orderID string) int64 {
	var blockHeight int64
	var id int64
	format := "ID%d-%d"
	fmt.Sscanf(orderID, format, &blockHeight, &id)

	return blockHeight
}
