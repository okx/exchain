package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/okex/okchain/x/order"

	"github.com/okex/okchain/x/stream/common"

	"github.com/okex/okchain/x/stream/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/backend"
	ordertype "github.com/okex/okchain/x/order/types"
	"github.com/okex/okchain/x/token"
)

type Writer interface {
	WriteSync(b *RedisBlock) (map[string]int, error) // atomic operation
}

type RedisBlock struct {
	Height        int64                      `json:"height"`     //blockHeight
	OrdersMap     map[string][]backend.Order `json:"orders"`     //key: address
	DepthBooksMap map[string]BookRes         `json:"depthBooks"` //key: product

	AccountsMap map[string]token.CoinInfo      `json:"accounts"`    //key: instrument_id:<address>
	Instruments map[string]struct{}            `json:"instruments"` //P3K:spot:instruments
	MatchesMap  map[string]backend.MatchResult `json:"matches"`     //key: product
}

func NewRedisBlock() *RedisBlock {
	return &RedisBlock{
		Height:        -1,
		OrdersMap:     make(map[string][]backend.Order),
		DepthBooksMap: make(map[string]BookRes),

		AccountsMap: make(map[string]token.CoinInfo),
		Instruments: make(map[string]struct{}),
		MatchesMap:  make(map[string]backend.MatchResult),
	}
}
func (rb RedisBlock) String() string {
	b, _ := json.Marshal(rb)
	return string(b)
}
func (rb *RedisBlock) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, cache *common.Cache) {
	rb.Height = ctx.BlockHeight()
	isTokenPairChanged := cache.GetTokenPairChanged()
	rb.storeInstrumentsWhenChanged(ctx, isTokenPairChanged, dexKeeper)
	rb.storeNewOrders(ctx, orderKeeper, rb.Height)
	rb.updateOrders(ctx, orderKeeper)
	rb.storeDepthBooks(ctx, orderKeeper, 200)
	updatedAccount := cache.GetUpdatedAccAddress()
	rb.storeAccount(ctx, updatedAccount, tokenKeeper)
	rb.storeMatches(ctx, orderKeeper)
}

func (rb *RedisBlock) storeInstrumentsWhenChanged(ctx sdk.Context, isTokenPairChanged bool, dexKeeper types.DexKeeper) {
	if isTokenPairChanged {
		rb.storeInstruments(ctx, dexKeeper)
	}
}

func (rb *RedisBlock) Empty() bool {
	if rb.Height == -1 && len(rb.DepthBooksMap) == 0 &&
		len(rb.OrdersMap) == 0 && len(rb.AccountsMap) == 0 &&
		len(rb.Instruments) == 0 && len(rb.MatchesMap) == 0 {
		return true
	}
	return false
}

func (rb *RedisBlock) Clear() {
	rb.Height = -1
	rb.OrdersMap = make(map[string][]backend.Order)
	rb.DepthBooksMap = make(map[string]BookRes)
	rb.Instruments = make(map[string]struct{})
	rb.AccountsMap = make(map[string]token.CoinInfo)
	rb.MatchesMap = make(map[string]backend.MatchResult)
}

func (b *RedisBlock) storeInstruments(ctx sdk.Context, dexKeeper types.DexKeeper) {
	logger := ctx.Logger().With("module", "stream")
	pairs := dexKeeper.GetTokenPairs(ctx)
	for _, pair := range pairs {
		product := fmt.Sprintf("%s_%s", pair.BaseAssetSymbol, pair.QuoteAssetSymbol)
		b.Instruments[product] = struct{}{}
		b.Instruments[pair.BaseAssetSymbol] = struct{}{}
		b.Instruments[pair.QuoteAssetSymbol] = struct{}{}
	}
	//curs := k.GetCurrencysInfo(ctx)
	//for _, curs := range curs {
	//	b.Instruments[curs.Symbol] = struct{}{}
	//}
	logger.Debug("storeInstruments",
		"instruments", b.Instruments,
	)
}

func getAddressProductPrefix(s1, s2 string) string {
	return s1 + ":" + s2
}

func (b *RedisBlock) storeNewOrders(ctx sdk.Context, orderKeeper types.OrderKeeper, blockHeight int64) {

	logger := ctx.Logger().With("module", "stream")
	orders, _ := backend.GetNewOrdersAtEndBlock(ctx, orderKeeper)
	for _, o := range orders {
		key := getAddressProductPrefix(o.Product, o.Sender)
		//key := o.Sender
		b.OrdersMap[key] = append(b.OrdersMap[key], *o)
		logger.Debug("storeNewOrders", "order", o)
	}
}

func (b *RedisBlock) updateOrders(ctx sdk.Context, orderKeeper types.OrderKeeper) {
	logger := ctx.Logger().With("module", "stream")
	orders := backend.GetUpdatedOrdersAtEndBlock(ctx, orderKeeper)
	for _, o := range orders {
		key := getAddressProductPrefix(o.Product, o.Sender)
		//key := o.Sender
		if _, ok := b.OrdersMap[key]; ok {
			if i, found := find(b.OrdersMap[key], *o); found {
				b.OrdersMap[key][i] = *o
			} else {
				b.OrdersMap[key] = append(b.OrdersMap[key], *o)
			}
		} else {
			b.OrdersMap[key] = append(b.OrdersMap[key], *o)
		}
		logger.Debug("updateOrders", "order", o)
	}
}

func (b *RedisBlock) storeMatches(ctx sdk.Context, orderKeeper types.OrderKeeper) {
	logger := ctx.Logger().With("module", "stream")
	_, matches, _ := backend.GetNewDealsAndMatchResultsAtEndBlock(ctx, orderKeeper)
	for _, m := range matches {
		b.MatchesMap[m.Product] = *m
		logger.Debug("storeMatches", "match", m)
	}
}

type BookResItem struct {
	Price      string `json:"price"`
	Quantity   string `json:"quantity"`
	OrderCount string `json:"order_count"`
}

type BookRes struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Product   string     `json:"instrument_id"`
	Timestamp string     `json:"timestamp"`
}

func (bri *BookResItem) toJsonList() []string {
	return []string{bri.Price, bri.Quantity, bri.OrderCount}
}

// nolint: unparam
//ask: small -> big, bids: big -> small
func (b *RedisBlock) storeDepthBooks(ctx sdk.Context, orderKeeper types.OrderKeeper, size int) {
	logger := ctx.Logger().With("module", "stream")

	products := orderKeeper.GetUpdatedDepthbookKeys()
	if len(products) == 0 {
		return
	}

	for _, product := range products {
		depthBook := orderKeeper.GetDepthBookCopy(product)
		bookRes := ConvertBookRes(product, orderKeeper, depthBook, size)
		b.DepthBooksMap[product] = bookRes
		logger.Debug("storeDepthBooks", "product", product, "depthBook", bookRes)
	}
}

func ConvertBookRes(product string, orderKeeper types.OrderKeeper, depthBook *order.DepthBook, size int) BookRes {
	asks := [][]string{}
	bids := [][]string{}
	for _, item := range depthBook.Items {
		if item.SellQuantity.IsPositive() {
			key := ordertype.FormatOrderIDsKey(product, item.Price, ordertype.SellOrder)
			ids := orderKeeper.GetProductPriceOrderIDs(key)
			bri := BookResItem{item.Price.String(), item.SellQuantity.String(), fmt.Sprintf("%d", len(ids))}
			asks = append([][]string{bri.toJsonList()}, asks...)

		}
		if item.BuyQuantity.IsPositive() {
			key := ordertype.FormatOrderIDsKey(product, item.Price, ordertype.BuyOrder)
			ids := orderKeeper.GetProductPriceOrderIDs(key)
			bri := BookResItem{item.Price.String(), item.BuyQuantity.String(), fmt.Sprintf("%d", len(ids))}
			// bids = append([][]string{bri.toJsonList()}, bids...)
			bids = append(bids, bri.toJsonList())
		}
	}
	if len(asks) > size {
		asks = asks[:size]
	}
	if len(bids) > size {
		bids = bids[:size]
	}

	bookRes := BookRes{
		asks,
		bids,
		product,
		time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}
	return bookRes
}

func (b *RedisBlock) storeAccount(ctx sdk.Context, updatedAccount []sdk.AccAddress, tokenKeeper types.TokenKeeper) {
	logger := ctx.Logger().With("module", "stream")

	for _, acc := range updatedAccount {
		coinsInfo := tokenKeeper.GetCoinsInfo(ctx, acc)
		for _, coinInfo := range coinsInfo {
			// key := acc.String() + ":" + coinInfo.Symbol
			key := coinInfo.Symbol + ":" + acc.String()
			b.AccountsMap[key] = coinInfo
		}
		logger.Debug("storeAccount",
			"address", acc.String(),
			"Currencies", coinsInfo,
		)
	}
}

var _ types.IStreamData = (*RedisBlock)(nil)

//impl IsStreamData interface
func (rb RedisBlock) BlockHeight() int64 {
	return rb.Height
}

func (rb RedisBlock) DataType() types.StreamDataKind {
	return types.StreamDataNotifyKind
}
