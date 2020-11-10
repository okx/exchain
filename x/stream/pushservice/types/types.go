package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/okex/okexchain/x/order"

	"github.com/okex/okexchain/x/stream/common"

	"github.com/okex/okexchain/x/stream/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend"
	ordertype "github.com/okex/okexchain/x/order/types"
	"github.com/okex/okexchain/x/token"
)

type Writer interface {
	WriteSync(b *RedisBlock) (map[string]int, error) // atomic operation
}

type RedisBlock struct {
	Height        int64                      `json:"height"`
	OrdersMap     map[string][]backend.Order `json:"orders"`
	DepthBooksMap map[string]BookRes         `json:"depthBooks"`

	AccountsMap map[string]token.CoinInfo      `json:"accounts"`
	Instruments map[string]struct{}            `json:"instruments"`
	MatchesMap  map[string]backend.MatchResult `json:"matches"`
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
	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}
	return string(b)
}
func (rb *RedisBlock) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, swapKeeper types.SwapKeeper, cache *common.Cache) {
	rb.Height = ctx.BlockHeight()

	rb.storeInstruments(ctx, cache, dexKeeper, swapKeeper)
	rb.storeNewOrders(ctx, orderKeeper, rb.Height)
	rb.updateOrders(ctx, orderKeeper)
	rb.storeDepthBooks(ctx, orderKeeper, 200)
	updatedAccount := cache.GetUpdatedAccAddress()
	rb.storeAccount(ctx, updatedAccount, tokenKeeper)
	rb.storeMatches(ctx, orderKeeper)
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

func (rb *RedisBlock) storeInstruments(ctx sdk.Context, cache *common.Cache, dexKeeper types.DexKeeper, swapKeeper types.SwapKeeper) {
	logger := ctx.Logger().With("module", "stream")

	// store instruments when token pair changed
	isTokenPairChanged := cache.GetTokenPairChanged()
	newSwapPairs := cache.GetNewSwapTokenPairs()
	if isTokenPairChanged || len(newSwapPairs) > 0 {
		// store token in dex token pair
		tokenPairs := dexKeeper.GetTokenPairs(ctx)
		for _, tokenPair := range tokenPairs {
			rb.Instruments[tokenPair.Name()] = struct{}{}
			rb.Instruments[tokenPair.BaseAssetSymbol] = struct{}{}
			rb.Instruments[tokenPair.QuoteAssetSymbol] = struct{}{}
		}

		// store token in swap token pair
		swapTokenPairs := swapKeeper.GetSwapTokenPairs(ctx)
		for _, swapTokenPair := range swapTokenPairs {
			rb.Instruments[swapTokenPair.BasePooledCoin.Denom] = struct{}{}
			rb.Instruments[swapTokenPair.QuotePooledCoin.Denom] = struct{}{}
			rb.Instruments[swapTokenPair.PoolTokenName] = struct{}{}
		}
	}

	logger.Debug("storeInstruments",
		"instruments", rb.Instruments,
	)
}

func getAddressProductPrefix(s1, s2 string) string {
	return s1 + ":" + s2
}

// nolint
func (rb *RedisBlock) storeNewOrders(ctx sdk.Context, orderKeeper types.OrderKeeper, blockHeight int64) {
	logger := ctx.Logger().With("module", "stream")
	orders, err := backend.GetNewOrdersAtEndBlock(ctx, orderKeeper)
	if err != nil {
		logger.Error("RedisBlock storeNewOrders error", "msg", err.Error())
	}
	for _, o := range orders {
		// key := o.Sender
		key := getAddressProductPrefix(o.Product, o.Sender)
		rb.OrdersMap[key] = append(rb.OrdersMap[key], *o)
		logger.Debug("storeNewOrders", "order", o)
	}
}

// nolint
func (rb *RedisBlock) updateOrders(ctx sdk.Context, orderKeeper types.OrderKeeper) {
	logger := ctx.Logger().With("module", "stream")
	orders := backend.GetUpdatedOrdersAtEndBlock(ctx, orderKeeper)
	for _, o := range orders {
		// key := o.Sender
		key := getAddressProductPrefix(o.Product, o.Sender)
		if _, ok := rb.OrdersMap[key]; ok {
			if i, found := find(rb.OrdersMap[key], *o); found {
				rb.OrdersMap[key][i] = *o
			} else {
				rb.OrdersMap[key] = append(rb.OrdersMap[key], *o)
			}
		} else {
			rb.OrdersMap[key] = append(rb.OrdersMap[key], *o)
		}
		logger.Debug("updateOrders", "order", o)
	}
}

// nolint
func (rb *RedisBlock) storeMatches(ctx sdk.Context, orderKeeper types.OrderKeeper) {
	logger := ctx.Logger().With("module", "stream")
	_, matches, err := backend.GetNewDealsAndMatchResultsAtEndBlock(ctx, orderKeeper)
	if err != nil {
		logger.Error("RedisBlock storeMatches error", "msg", err.Error())
	}
	for _, m := range matches {
		rb.MatchesMap[m.Product] = *m
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

func (bri *BookResItem) toJSONList() []string {
	return []string{bri.Price, bri.Quantity, bri.OrderCount}
}

// nolint: unparam
//ask: small -> big, bids: big -> small
func (rb *RedisBlock) storeDepthBooks(ctx sdk.Context, orderKeeper types.OrderKeeper, size int) {
	logger := ctx.Logger().With("module", "stream")

	products := orderKeeper.GetUpdatedDepthbookKeys()
	if len(products) == 0 {
		return
	}

	for _, product := range products {
		depthBook := orderKeeper.GetDepthBookCopy(product)
		bookRes := ConvertBookRes(product, orderKeeper, depthBook, size)
		rb.DepthBooksMap[product] = bookRes
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
			asks = append([][]string{bri.toJSONList()}, asks...)

		}
		if item.BuyQuantity.IsPositive() {
			key := ordertype.FormatOrderIDsKey(product, item.Price, ordertype.BuyOrder)
			ids := orderKeeper.GetProductPriceOrderIDs(key)
			bri := BookResItem{item.Price.String(), item.BuyQuantity.String(), fmt.Sprintf("%d", len(ids))}
			// bids = append([][]string{bri.toJsonList()}, bids...)
			bids = append(bids, bri.toJSONList())
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

func (rb *RedisBlock) storeAccount(ctx sdk.Context, updatedAccount []sdk.AccAddress, tokenKeeper types.TokenKeeper) {
	logger := ctx.Logger().With("module", "stream")

	for _, acc := range updatedAccount {
		coinsInfo := tokenKeeper.GetCoinsInfo(ctx, acc)
		for _, coinInfo := range coinsInfo {
			// key := acc.String() + ":" + coinInfo.Symbol
			key := coinInfo.Symbol + ":" + acc.String()
			rb.AccountsMap[key] = coinInfo
		}
		logger.Debug("storeAccount",
			"address", acc.String(),
			"Currencies", coinsInfo,
		)
	}
}

var _ types.IStreamData = (*RedisBlock)(nil)

// BlockHeight impl IsStreamData interface
func (rb RedisBlock) BlockHeight() int64 {
	return rb.Height
}

func (rb RedisBlock) DataType() types.StreamDataKind {
	return types.StreamDataNotifyKind
}
