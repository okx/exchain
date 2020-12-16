package keeper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/okex/okexchain/x/ammswap"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend/cache"
	"github.com/okex/okexchain/x/backend/config"
	"github.com/okex/okexchain/x/backend/orm"
	"github.com/okex/okexchain/x/backend/types"
	"github.com/okex/okexchain/x/token"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	OrderKeeper  types.OrderKeeper  // The reference to the OrderKeeper to get deals
	TokenKeeper  types.TokenKeeper  // The reference to the TokenKeeper to get fee details
	marketKeeper types.MarketKeeper // The reference to the MarketKeeper to get fee details
	dexKeeper    types.DexKeeper    // The reference to the DexKeeper to get tokenpair
	swapKeeper   types.SwapKeeper
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
	Orm          *orm.ORM
	stopChan     chan struct{}
	Config       *config.Config
	Logger       log.Logger
	wsChan       chan types.IWebsocket // Websocket channel, it's only available when websocket config enabled
	ticker3sChan chan types.IWebsocket // Websocket channel, it's used by tickers merge triggered 3s once
	Cache        *cache.Cache          // Memory cache
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, swapKeeper types.SwapKeeper, marketKeeper types.MarketKeeper, cdc *codec.Codec, logger log.Logger, cfg *config.Config) Keeper {
	k := Keeper{
		OrderKeeper:  orderKeeper,
		TokenKeeper:  tokenKeeper,
		marketKeeper: marketKeeper,
		dexKeeper:    dexKeeper,
		swapKeeper:   swapKeeper,
		cdc:          cdc,
		Logger:       logger.With("module", "backend"),
		Config:       cfg,
		wsChan:       nil,
	}

	if k.Config.EnableBackend {
		k.Cache = cache.NewCache()
		orm, err := orm.New(k.Config.LogSQL, &k.Config.OrmEngine, &k.Logger)
		if err != nil {
			panic(fmt.Sprintf("backend new orm error:%s", err.Error()))
		}
		k.Orm = orm
		k.stopChan = make(chan struct{})

		if k.Config.EnableMktCompute {
			// websocket channel
			k.wsChan = make(chan types.IWebsocket, types.WebsocketChanCapacity)
			k.ticker3sChan = make(chan types.IWebsocket, types.WebsocketChanCapacity)
			go generateKline1M(k)
			// init ticker buffer
			ts := time.Now().Unix()

			k.UpdateTickersBuffer(ts-types.SecondsInADay*14, ts, nil)

			go k.mergeTicker3SecondEvents()

			// set observer keeper
			k.swapKeeper.SetObserverKeeper(k)
		}

	}
	logger.Debug(fmt.Sprintf("%+v", k.Config))
	return k
}

func (k Keeper) pushWSItem(obj types.IWebsocket) {
	if k.wsChan != nil {
		k.Logger.Debug("pushWSItem", "typeof(obj)", reflect.TypeOf(obj))
		k.wsChan <- obj
	}
}

func (k Keeper) pushTickerItems(obj types.IWebsocket) {
	if k.ticker3sChan != nil {
		k.ticker3sChan <- obj
	}
}

// Emit all of the WSItems as tendermint events
func (k Keeper) EmitAllWsItems(ctx sdk.Context) {
	if k.wsChan == nil {
		return
	}

	k.Logger.Debug("EmitAllWsItems", "eventCnt", len(k.wsChan))

	// TODO: Add filter to reduce events to send
	allChannelNotifies := map[string]int64{}
	updatedChannels := map[string]bool{}
	for len(k.wsChan) > 0 {
		item, ok := <-k.wsChan
		if ok {
			channel, _, err := item.GetChannelInfo()
			fullchannel := item.GetFullChannel()

			formatedResult := item.FormatResult()
			if formatedResult == nil {
				allChannelNotifies[channel] = item.GetTimestamp()
				continue
			}

			jstr, jerr := json.Marshal(formatedResult)
			if jerr == nil && err == nil {
				k.Logger.Debug("EmitAllWsItems Item[#1]", "type", reflect.TypeOf(item), "channel", fullchannel, "data", string(jstr))
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						"backend",
						sdk.NewAttribute("channel", fullchannel),
						sdk.NewAttribute("data", string(jstr))))

				updatedChannels[fullchannel] = true

			} else {
				k.Logger.Error("failed to EmitAllWsItems[#1] ", "Json Error", jerr, "GetChannelInfo Error", err)
				break
			}

		} else {
			break
		}
	}

	// Push All product kline when trigger by FakeEvent
	for klineType, ts := range allChannelNotifies {

		freq := types.GetFreqByKlineType(klineType)

		tokenPairs := k.getAllProducts(ctx)
		for _, tp := range tokenPairs {

			klines, err := k.getCandlesWithTimeFromORM(tp, freq, 1, ts)
			if err != nil || len(klines) == 0 {
				k.Logger.Error("EmitAllWsItems[#2] failed to getCandlesWithTimeFromORM", "error", err)
				continue
			}
			lastKline := klines[len(klines)-1]
			item := lastKline.(types.IWebsocket)

			fullchannel := item.GetFullChannel()
			bSkip, ok := updatedChannels[fullchannel]
			if bSkip || ok {
				continue
			}

			formatedResult := item.FormatResult()
			jstr, jerr := json.Marshal(formatedResult)
			if jerr == nil {
				k.Logger.Debug("EmitAllWsItems Item[#2]", "type", reflect.TypeOf(item), "data", string(jstr))
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						"backend",
						sdk.NewAttribute("channel", fullchannel),
						sdk.NewAttribute("data", string(jstr))))
			} else {
				k.Logger.Error("EmitAllWsItems[#2] failed to EmitAllWsItems ", "Json Error", jerr, "GetChannelInfo Error", err)
				break
			}
		}

	}

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
	if k.Config.EnableBackend && k.Config.EnableMktCompute {
		k.Logger.Debug(fmt.Sprintf("[backend] get new tx, txHash: %s", txHash))
		txs := types.GenerateTx(tx, txHash, ctx, k.OrderKeeper, timestamp)
		k.Cache.AddTransaction(txs)
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

// nolint
func (k Keeper) GetDexFees(ctx sdk.Context, dexHandlingAddr, product string, offset, limit int) ([]types.DexFees, int) {
	return k.Orm.GetDexFees(dexHandlingAddr, product, offset, limit)
}

func (k Keeper) getAllProducts(ctx sdk.Context) []string {
	tokenPairs := k.dexKeeper.GetTokenPairs(ctx)
	products := make([]string, len(tokenPairs))
	for i := 0; i < len(tokenPairs); i++ {
		products[i] = tokenPairs[i].Name()
	}
	return products
}

// nolint
func (k Keeper) getCandlesWithTimeFromORM(product string, granularity, size int, ts int64) (r []types.IKline, err error) {
	if !k.Config.EnableBackend {
		return nil, types.ErrBackendPluginNotEnabled()
	}

	m := types.GetAllKlineMap()
	candleType := m[granularity]
	if candleType == "" || len(candleType) == 0 || (size < 0 || size > 1000) {
		return nil, types.ErrParamNotCorrect(size, granularity)
	}

	klines, err := types.NewKlinesFactory(candleType)
	if err == nil {
		err := k.Orm.GetLatestKlinesByProduct(product, size, ts, klines)
		iklines := types.ToIKlinesArray(klines, ts, true)
		return iklines, err
	}
	return nil, err

}

// nolint
func (k Keeper) GetCandlesWithTime(product string, granularity, size int, ts int64) (r [][]string, err error) {

	iklines, err := k.getCandlesWithTimeFromORM(product, granularity, size, ts)
	if err == nil {
		restData := types.ToRestfulData(&iklines, size)
		return restData, err
	}
	return nil, err
}

func (k Keeper) getCandlesByMarketKeeper(productID uint64, granularity, size int) (r [][]string, err error) {
	if !k.Config.EnableBackend {
		return nil, types.ErrBackendPluginNotEnabled()
	}

	if k.marketKeeper == nil {
		return nil, types.ErrMarketkeeperNotInitialized()
	}

	m := types.GetAllKlineMap()
	candleType := m[granularity]
	if candleType == "" || len(candleType) == 0 || (size < 0 || size > 1000) {
		return nil, types.ErrParamNotCorrect(size, granularity)
	}

	klines, err := k.marketKeeper.GetKlineByProductID(productID, granularity, size)
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
				if ticker != nil {
					tickers = append(tickers, *ticker)
				}
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

	k.Orm.Debug(fmt.Sprintf("[backend] entering UpdateTickersBuffer, latestTickers: %+v, TickerTimeRange: [%d, %d)=[%s, %s), productList: %v",
		k.Cache.LatestTicker, startTS, endTS, types.TimeString(startTS), types.TimeString(endTS), productList))

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
			k.pushWSItem(ticker)
			k.pushTickerItems(ticker)
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
			if previousTicker != nil && (endTS > previousTicker.Timestamp+types.SecondsInADay) {
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
		if ticker != nil {
			tickers = append(tickers, *ticker)
		}
	}
	return tickers
}

func (k Keeper) mergeTicker3SecondEvents() (err error) {

	sysTicker := time.NewTicker(3 * time.Second)

	merge := func() *types.MergedTickersEvent {
		tickersMap := map[string]types.IWebsocket{}
		for len(k.ticker3sChan) > 0 {
			ticker, ok := <-k.ticker3sChan
			if !ok {
				break
			}

			if ticker.FormatResult() == nil {
				continue
			}

			tickersMap[ticker.GetFullChannel()] = ticker
		}

		allTickers := []interface{}{}
		for _, ticker := range tickersMap {
			allTickers = append(allTickers, ticker.FormatResult())
		}

		if len(allTickers) > 0 {
			return types.NewMergedTickersEvent(time.Now().Unix(), 3, allTickers)
		}

		return nil

	}

	for {
		select {
		case <-sysTicker.C:
			mEvt := merge()
			if mEvt != nil {
				k.pushWSItem(mEvt)
			}
		case <-k.stopChan:
			break
		}
	}
}

func (k Keeper) OnSwapToken(ctx sdk.Context, address sdk.AccAddress, swapTokenPair ammswap.SwapTokenPair, sellAmount sdk.SysCoin, buyAmount sdk.SysCoin) {
	swapInfo := &types.SwapInfo{
		Address:          address.String(),
		TokenPairName:    swapTokenPair.TokenPairName(),
		BaseTokenAmount:  swapTokenPair.BasePooledCoin.String(),
		QuoteTokenAmount: swapTokenPair.QuotePooledCoin.String(),
		SellAmount:       sellAmount.String(),
		BuysAmount:       buyAmount.String(),
		Price:            swapTokenPair.BasePooledCoin.Amount.Quo(swapTokenPair.QuotePooledCoin.Amount).String(),
		Timestamp:        ctx.BlockTime().Unix(),
	}
	k.Cache.AddSwapInfo(swapInfo)
}

func (k Keeper) OnSwapCreateExchange(ctx sdk.Context, swapTokenPair ammswap.SwapTokenPair) {
}
