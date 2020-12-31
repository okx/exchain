// nolint
package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/okex/okexchain/x/dex"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MatchResult struct {
	Timestamp   int64   `gorm:"index;" json:"timestamp" v2:"timestamp"`
	BlockHeight int64   `gorm:"PRIMARY_KEY;type:bigint" json:"block_height" v2:"block_height"`
	Product     string  `gorm:"PRIMARY_KEY;type:varchar(20)" json:"product" v2:"product"`
	Price       float64 `gorm:"type:DOUBLE" json:"price" v2:"price"`
	Quantity    float64 `gorm:"type:DOUBLE" json:"volume" v2:"volume"`
}

type Deal struct {
	Timestamp   int64   `gorm:"index;" json:"timestamp" v2:"timestamp"`
	BlockHeight int64   `gorm:"PRIMARY_KEY;type:bigint" json:"block_height" v2:"block_height"`
	OrderID     string  `gorm:"PRIMARY_KEY;type:varchar(30)" json:"order_id" v2:"order_id"`
	Sender      string  `gorm:"index;type:varchar(80)" json:"sender" v2:"sender"`
	Product     string  `gorm:"index;type:varchar(20)" json:"product" v2:"product"`
	Side        string  `gorm:"type:varchar(10)" json:"side" v2:"side"`
	Price       float64 `gorm:"type:DOUBLE" json:"price" v2:"price"`
	Quantity    float64 `gorm:"type:DOUBLE" json:"volume" v2:"volume"`
	Fee         string  `gorm:"type:varchar(20)" json:"fee" v2:"fee"`
	FeeReceiver string  `gorm:"index;type:varchar(80)" json:"fee_receiver" v2:"fee_receiver"`
}

type TickerV2 struct {
	InstrumentID   string `json:"instrument_id"` // name of token pair
	Last           string `json:"last"`
	BestBid        string `json:"best_bid"`
	BestAsk        string `json:"best_ask"`
	Open24H        string `json:"open_24h"`
	High24H        string `json:"high_24h"`
	Low24H         string `json:"low_24h"`
	BaseVolume24H  string `json:"base_volume_24h"`
	QuoteVolume24H string `json:"quote_volume_24h"`
	Timestamp      string `json:"timestamp"`
}

func DefaultTickerV2(instrumentID string) TickerV2 {
	return TickerV2{
		InstrumentID:   instrumentID,
		Last:           "-1",
		BestBid:        "0",
		BestAsk:        "0",
		Open24H:        "0",
		High24H:        "0",
		Low24H:         "0",
		BaseVolume24H:  "0",
		QuoteVolume24H: "0",
		Timestamp:      time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}
}

type InstrumentV2 struct {
	InstrumentID  string `json:"instrument_id"` // name of token pair
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	MinSize       string `json:"min_size"`
	SizeIncrement string `json:"size_increment"`
	TickSize      string `json:"tick_size"`
}

// ConvertTokenPairToInstrumentV2 convert TokenPair to InstrumentV2
func ConvertTokenPairToInstrumentV2(tokenPair *dex.TokenPair) *InstrumentV2 {
	res := &InstrumentV2{}
	res.InstrumentID = tokenPair.Name()
	res.BaseCurrency = tokenPair.BaseAssetSymbol
	res.QuoteCurrency = tokenPair.QuoteAssetSymbol
	res.MinSize = tokenPair.MinQuantity.String()
	// convert 4 to "0.0001"
	fSizeIncrement := 1 / math.Pow10(int(tokenPair.MaxQuantityDigit))
	res.SizeIncrement = strings.TrimRight(fmt.Sprintf("%.10f", fSizeIncrement), "0")

	fTickSize := 1 / math.Pow10(int(tokenPair.MaxPriceDigit))
	res.TickSize = strings.TrimRight(fmt.Sprintf("%.10f", fTickSize), "0")
	return res
}

type OrderV2 struct {
	OrderID        string `json:"order_id"`
	Price          string `json:"price"`
	Size           string `json:"size"`
	OrderType      string `json:"order_type"`
	Notional       string `json:"notional"`
	InstrumentID   string `json:"instrument_id"`
	Side           string `json:"side"`
	Type           string `json:"type"`
	Timestamp      string `json:"timestamp"`
	FilledSize     string `json:"filled_size"`
	FilledNotional string `json:"filled_notional"`
	State          string `json:"state"`
}

func ConvertOrderToOrderV2(order Order) OrderV2 {
	res := OrderV2{}
	res.OrderID = order.OrderID
	res.Price = order.Price
	res.Size = order.Quantity
	res.OrderType = "0"
	res.Notional = order.FilledAvgPrice
	res.InstrumentID = order.Product
	res.Side = order.Side
	res.Type = "limit"
	res.Timestamp = time.Unix(order.Timestamp, 0).UTC().Format("2006-01-02T15:04:05.000Z")
	res.State = strconv.FormatInt(order.Status, 10)

	filledSizeDec := sdk.MustNewDecFromStr(order.Quantity).Sub(sdk.MustNewDecFromStr(order.RemainQuantity))
	filledNotionalDec := filledSizeDec.Mul(sdk.MustNewDecFromStr(order.FilledAvgPrice))

	res.FilledSize = filledSizeDec.String()
	res.FilledNotional = filledNotionalDec.String()

	return res
}

type QueryOrderParamsV2 struct {
	OrderID string
	Product string
	Side    string
	Address string
	After   string
	Before  string
	Limit   int
	IsOpen  bool
}

type QueryFeeDetailsParamsV2 struct {
	Address string
	After   string
	Before  string
	Limit   int
}

type QueryMatchParamsV2 struct {
	Product string
	After   string
	Before  string
	Limit   int
}

type QueryDealsParamsV2 struct {
	Address string
	Product string
	Side    string
	After   string
	Before  string
	Limit   int
}

type QueryTxListParamsV2 struct {
	Address string
	TxType  int
	After   string
	Before  string
	Limit   int
}

type DexFees struct {
	Timestamp       int64  `json:"timestamp"`
	OrderID         string `json:"order_id"`
	Product         string `json:"product"`
	Fee             string `json:"fee"`
	HandlingFeeAddr string `json:"handling_fee_addr"`
}
