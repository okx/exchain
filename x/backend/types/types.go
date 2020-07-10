// nolint
package types

import (
	"fmt"
	orderTypes "github.com/okex/okchain/x/order/types"
	"time"
)

const (
	TxTypeTransfer    = 1
	TxTypeOrderNew    = 2
	TxTypeOrderCancel = 3

	TxSideBuy  = 1
	TxSideSell = 2
	TxSideFrom = 3
	TxSideTo   = 4

	BuyOrder      = orderTypes.BuyOrder
	SellOrder     = orderTypes.SellOrder
	TestTokenPair = orderTypes.TestTokenPair

	FeeTypeOrderNew     = orderTypes.FeeTypeOrderNew
	FeeTypeOrderCancel  = orderTypes.FeeTypeOrderCancel
	FeeTypeOrderExpire  = orderTypes.FeeTypeOrderExpire
	FeeTypeOrderDeal    = orderTypes.FeeTypeOrderDeal
	FeeTypeOrderReceive = orderTypes.FeeTypeOrderReceive
)

type Ticker struct {
	Symbol           string  `json:"symbol"`
	Product          string  `json:"product"`
	Timestamp        int64   `json:"timestamp"`
	Open             float64 `json:"open"`  // Open In 24h
	Close            float64 `json:"close"` // Close in 24h
	High             float64 `json:"high"`  // High in 24h
	Low              float64 `json:"low"`   // Low in 24h
	Price            float64 `json:"price"`
	Volume           float64 `json:"volume"`            // Volume in 24h
	Change           float64 `json:"change"`            // (Close - Open)
	ChangePercentage string  `json:"change_percentage"` // Change / Open * 100%
}

func (t *Ticker) GetChannelInfo() (channel, filter string, err error) {
	channel = "dex_spot/ticker"
	filter = t.Product
	return
}

//func (b *BaseKline) GetBrifeInfo() []string {
//	m := []string{
//		time.Unix(b.GetTimestamp(), 0).UTC().Format("2006-01-02T15:04:05.000Z"),
//		fmt.Sprintf("%.4f", b.GetOpen()),
//		fmt.Sprintf("%.4f", b.GetHigh()),
//		fmt.Sprintf("%.4f", b.GetLow()),
//		fmt.Sprintf("%.4f", b.GetClose()),
//		fmt.Sprintf("%.8f", b.GetVolume()),
//	}
//	return m
//}

func (t *Ticker) FormatResult() interface{} {
	result := map[string]string{
		"product":   t.Product,
		"symbol":    t.Symbol,
		"timestamp": time.Unix(t.Timestamp, 0).UTC().Format("2006-01-02T15:04:05.000Z"),
		"open":      fmt.Sprintf("%.4f", t.Open),
		"high":      fmt.Sprintf("%.4f", t.High),
		"low":       fmt.Sprintf("%.4f", t.Low),
		"close":     fmt.Sprintf("%.4f", t.Close),
		"volume":    fmt.Sprintf("%.8f", t.Volume),
		"price":     fmt.Sprintf("%.4f", t.Price),
	}
	return result
}

// PrettyString return string of ticker data
func (t *Ticker) PrettyString() string {
	return fmt.Sprintf("[Ticker] Symbol: %s, Price: %f, TStr: %s, Timestamp: %d, OCHLV(%f, %f, %f, %f, %f) [%f, %s])",
		t.Symbol, t.Price, TimeString(t.Timestamp), t.Timestamp, t.Open, t.Close, t.High, t.Low, t.Volume, t.Change, t.ChangePercentage)
}

type Tickers []Ticker

func (tickers Tickers) Len() int {
	return len(tickers)
}

func (c Tickers) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (tickers Tickers) Less(i, j int) bool {
	return tickers[i].Change < tickers[j].Change
}

type Order struct {
	TxHash         string `gorm:"type:varchar(80)" json:"txhash" v2:"txhash"`
	OrderID        string `gorm:"PRIMARY_KEY;type:varchar(30)" json:"order_id" v2:"order_id"`
	Sender         string `gorm:"index;type:varchar(80)" json:"sender" v2:"sender"`
	Product        string `gorm:"index;type:varchar(20)" json:"product" v2:"product"`
	Side           string `gorm:"type:varchar(10)" json:"side" v2:"side"`
	Price          string `gorm:"type:varchar(40)" json:"price" v2:"price"`
	Quantity       string `gorm:"type:varchar(40)" json:"quantity" v2:"quantity"`
	Status         int64  `gorm:"index;" json:"status" v2:"status"`
	FilledAvgPrice string `gorm:"type:varchar(40)" json:"filled_avg_price" v2:"filled_avg_price"`
	RemainQuantity string `gorm:"type:varchar(40)" json:"remain_quantity" v2:"remain_quantity"`
	Timestamp      int64  `gorm:"index;" json:"timestamp" v2:"timestamp"`
}

type Transaction struct {
	TxHash    string `gorm:"type:varchar(80)" json:"txhash" v2:"txhash"`
	Type      int64  `gorm:"index;" json:"type" v2:"type"` // 1:Transfer, 2:NewOrder, 3:CancelOrder
	Address   string `gorm:"index;type:varchar(80)" json:"address" v2:"address"`
	Symbol    string `gorm:"type:varchar(20)" json:"symbol" v2:"symbol"`
	Side      int64  `gorm:"" json:"side"` // 1:buy, 2:sell, 3:from, 4:to
	Quantity  string `gorm:"type:varchar(40)" json:"quantity" v2:"quantity"`
	Fee       string `gorm:"type:varchar(40)" json:"fee" v2:"fee"`
	Timestamp int64  `gorm:"index" json:"timestamp" v2:"timestamp"`
}
