package keeper

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/backend/cache"
	"github.com/okex/okchain/x/backend/config"
	"github.com/okex/okchain/x/backend/orm"
	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/token"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	OrderKeeper  types.OrderKeeper  // The reference to the OrderKeeper to get deals
	TokenKeeper  types.TokenKeeper  // The reference to the TokenKeeper to get fee details
	marketKeeper types.MarketKeeper // The reference to the MarketKeeper to get fee details
	dexKeeper    types.DexKeeper    // The reference to the DexKeeper to get tokenpair
	cdc          *codec.Codec       // The wire codec for binary encoding/decoding.
	Orm          *orm.ORM
	stopChan     chan struct{}
	Config       *config.Config
	Logger       log.Logger

	// memory cache
	Cache *cache.Cache
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, marketKeeper types.MarketKeeper, cdc *codec.Codec, logger log.Logger, cfg *config.Config) Keeper {
	k := Keeper{
		OrderKeeper:  orderKeeper,
		TokenKeeper:  tokenKeeper,
		marketKeeper: marketKeeper,
		dexKeeper:    dexKeeper,
		cdc:          cdc,
		Logger:       logger,
		Config:       cfg,
	}

	if k.Config.EnableBackend {
		k.Cache = cache.NewCache()
		orm, err := orm.New(k.Config.LogSQL, &k.Config.OrmEngine, &logger)
		if err == nil {
			k.Orm = orm
			k.stopChan = make(chan struct{})

			if k.Config.EnableMktCompute {
				go generateKline1M(k.stopChan, k.Config, k.Orm, &k.Logger)
				// init ticker buffer
				ts := time.Now().Unix()

				k.UpdateTickersBuffer(ts-types.SecondsInADay*14, ts, nil)
			}
		}
	}

	logger.Debug(fmt.Sprintf("%+v", k.Config))
	return k
}

// Stop close database
func (k Keeper) Stop() {
	defer types.PrintStackIfPanic()
	if k.stopChan != nil {
		close(k.stopChan)
	}
	if k.Orm != nil {
		if err := k.Orm.Close(); err != nil {
			k.Orm.Error(fmt.Sprintf("failed to close orm because %s ", err.Error()))
		}
	}
}

// Flush temporary cache
func (k Keeper) Flush() {
	defer k.Cache.Flush()
}

// SyncTx generate transaction and add it to cache, called at DeliverTx
func (k Keeper) SyncTx(ctx sdk.Context, tx *auth.StdTx, txHash string, timestamp int64) {
	if tx == nil {
		panic("failed. a nil pointer appears")
	}
	if k.Config.EnableBackend && k.Config.EnableMktCompute {
		k.Logger.Debug(fmt.Sprintf("[backend] get new tx, txHash: %s", txHash))
		txs := types.GenerateTx(tx, txHash, ctx, k.OrderKeeper, timestamp)
		for _, tx := range txs {
			k.Cache.AddTransaction(tx)
		}
	}
}

func (k Keeper) getMatchResults(ctx sdk.Context, product string, start, end int64, offset, limit int) ([]types.MatchResult, int) {
	return k.Orm.GetMatchResults(product, start, end, offset, limit)
}

// nolint
func (k Keeper) GetDeals(ctx sdk.Context, sender, product, side string, start, end int64, offset, limit int) ([]types.Deal, int) {
	return k.Orm.GetDeals(sender, product, side, start, end, offset, limit)
}

// nolint
func (k Keeper) GetFeeDetails(ctx sdk.Context, addr string, offset, limit int) ([]token.FeeDetail, int) {
	return k.Orm.GetFeeDetails(addr, offset, limit)
}

// nolint
func (k Keeper) GetOrderList(ctx sdk.Context, addr, product, side string, open bool,
	offset, limit int, startTS, endTS int64, hideNoFill bool) ([]types.Order, int) {
	return k.Orm.GetOrderList(addr, product, side, open, offset, limit, startTS, endTS, hideNoFill)
}

// nolint
func (k Keeper) GetTransactionList(ctx sdk.Context, addr string, txType, startTime, endTime int64, offset, limit int) ([]types.Transaction, int) {
	return k.Orm.GetTransactionList(addr, txType, startTime, endTime, offset, limit)
}

func (k Keeper) getAllProducts(ctx sdk.Context) []string {
	products := []string{}
	tokenPairs := k.dexKeeper.GetTokenPairs(ctx)
	for _, tp := range tokenPairs {
		if tp != nil {
			products = append(products, fmt.Sprintf("%s_%s", tp.BaseAssetSymbol, tp.QuoteAssetSymbol))
		}
	}

	k.Cache.ProductsBuf = products

	return products
}

// nolint
func (k Keeper) GetCandlesWithTime(product string, granularity, size int, ts int64) (r [][]string, err error) {
	if !k.Config.EnableBackend {
		return nil, fmt.Errorf("backend is not enabled, no candle found, maintian.conf: %+v", k.Config)
	}

	m := types.GetAllKlineMap()
	candleType := m[granularity]
	if candleType == "" || len(candleType) == 0 || (size < 0 || size > 1000) {
		return nil, fmt.Errorf("parameter's not correct, size: %d, granularity: %d", size, granularity)
	}

	klines, err := types.NewKlinesFactory(candleType)
	if err == nil {
		err := k.Orm.GetLatestKlinesByProduct(product, size, ts, klines)
		iklines := types.ToIKlinesArray(klines, ts, true)
		restData := types.ToRestfulData(&iklines, size)
		return restData, err
	}

	return nil, err
}

func (k Keeper) getCandlesByMarketKeeper(product string, granularity, size int) (r [][]string, err error) {
	if !k.Config.EnableBackend {
		return nil, fmt.Errorf("backend is not enabled, no candle found, maintian.conf: %+v", k.Config)
	}

	if k.marketKeeper == nil {
		return nil, errors.New("Market keeper is not initialized properly")
	}

	m := types.GetAllKlineMap()
	candleType := m[granularity]
	if candleType == "" || len(candleType) == 0 || (size < 0 || size > 1000) {
		return nil, fmt.Errorf("parameter's not correct, size: %d, granularity: %d", size, granularity)
	}

	klines, err := k.marketKeeper.GetKlineByInstrument(product, granularity, size)
	if err == nil && klines != nil && len(klines) > 0 {
		return klines, err
	}

	return [][]string{}, err
}

// nolint
func (k Keeper) GetCandles(product string, granularity, size int) (r [][]string, err error) {
	return k.GetCandlesWithTime(product, granularity, size, time.Now().Unix())
}

// nolint
func (k Keeper) GetTickers(products []string, count int) []types.Ticker {
	tickers := []types.Ticker{}
	if len(k.Cache.LatestTicker) > 0 {

		if len(products) > 0 {
			for _, p := range products {
				t := k.Cache.LatestTicker[p]
				if t != nil {
					tickers = append(tickers, *t)
				}
			}
		} else {
			for _, ticker := range k.Cache.LatestTicker {
				tickers = append(tickers, *ticker)
			}
		}
	}

	maxUpper := count
	if len(tickers) > 0 {
		if len(tickers) < maxUpper {
			maxUpper = len(tickers)
		}
		return tickers[0:maxUpper]
	} else {
		return tickers
	}
}

// UpdateTickersBuffer calculate and update the products ticker
func (k Keeper) UpdateTickersBuffer(startTS, endTS int64, productList []string) {

	defer types.PrintStackIfPanic()

	k.Orm.Debug(fmt.Sprintf("[backend] entering UpdateTickersBuffer, latestTickers: %+v, TickerTimeRange: [%d, %d)=[%s, %s)",
		k.Cache.LatestTicker, startTS, endTS, types.TimeString(startTS), types.TimeString(endTS)))

	latestProducts := []string{}
	for p := range k.Cache.LatestTicker {
		latestProducts = append(latestProducts, p)
	}
	tickerMap, err := k.Orm.RefreshTickers(startTS, endTS, productList)
	if err != nil {
		k.Orm.Error(fmt.Sprintf("generateTicker error %+v, latestTickers %+v, returnTickers: %+v", err, k.Cache.LatestTicker, tickerMap))
		return
	}

	if len(tickerMap) > 0 {
		for product, ticker := range tickerMap {
			k.Cache.LatestTicker[product] = ticker
		}
		k.Orm.Debug(fmt.Sprintf("UpdateTickersBuffer LatestTickerMap: %+v", k.Cache.LatestTicker))
	} else {
		k.Orm.Debug(fmt.Sprintf("UpdateTickersBuffer No product's deal refresh in [%d, %d), latestTicker: %+v", startTS, endTS, k.Cache.LatestTicker))
	}

	// Case: No deals produced in last 24 hours.
	for _, p := range latestProducts {
		refreshedTicker := tickerMap[p]
		if refreshedTicker == nil {
			previousTicker := k.Cache.LatestTicker[p]
			if previousTicker != nil {
				previousTicker.Open = previousTicker.Close
				previousTicker.High = previousTicker.Close
				previousTicker.Low = previousTicker.Close
				previousTicker.Volume = 0
				previousTicker.Change = 0
				previousTicker.ChangePercentage = "0.00%"
			}

		}
	}
}

func (k Keeper) getOrderListV2(ctx sdk.Context, instrumentID string, address string, side string, open bool, after string, before string, limit int) []types.Order {
	return k.Orm.GetOrderListV2(instrumentID, address, side, open, after, before, limit)
}

func (k Keeper) getOrderByIDV2(ctx sdk.Context, orderID string) *types.Order {
	return k.Orm.GetOrderByID(orderID)
}

func (k Keeper) getMatchResultsV2(ctx sdk.Context, instrumentID string, after string, before string, limit int) []types.MatchResult {
	return k.Orm.GetMatchResultsV2(instrumentID, after, before, limit)
}

func (k Keeper) getFeeDetailsV2(ctx sdk.Context, addr string, after string, before string, limit int) []token.FeeDetail {
	return k.Orm.GetFeeDetailsV2(addr, after, before, limit)
}

func (k Keeper) getDealsV2(ctx sdk.Context, sender, product, side string, after string, before string, limit int) []types.Deal {
	return k.Orm.GetDealsV2(sender, product, side, after, before, limit)
}

func (k Keeper) getTransactionListV2(ctx sdk.Context, addr string, txType int, after string, before string, limit int) []types.Transaction {
	return k.Orm.GetTransactionListV2(addr, txType, after, before, limit)
}

func (k Keeper) getAllTickers() []types.Ticker {
	var tickers []types.Ticker
	for _, ticker := range k.Cache.LatestTicker {
		tickers = append(tickers, *ticker)
	}
	return tickers
}
