package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"
)

func queryTickerFromMarketKeeperV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	params := types.QueryTickerParams{}
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	keeper.marketKeeper.InitTokenPairMap(ctx, keeper.dexKeeper)
	tickers, err := keeper.marketKeeper.GetTickers()
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	result := types.DefaultTickerV2(params.Product)
	notExist := true
	for _, t := range tickers {
		if params.Product == t["product"] {
			notExist = false
			result.Last = t["price"]
			result.Open24H = t["open"]
			result.High24H = t["high"]
			result.Low24H = t["low"]
			result.BaseVolume24H = t["volume"]
			result.QuoteVolume24H = t["quote_volume_24h"]
			result.Timestamp = t["timestamp"]
			break
		}
	}

	if notExist {
		return nil, nil
	}

	bestBid, bestAsk := keeper.OrderKeeper.GetBestBidAndAsk(ctx, params.Product)
	result.BestBid = bestBid.String()
	result.BestAsk = bestAsk.String()

	res, err := json.Marshal(result)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryTickerListFromMarketKeeperV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	keeper.marketKeeper.InitTokenPairMap(ctx, keeper.dexKeeper)
	tickers, err := keeper.marketKeeper.GetTickers()
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	var tickerList []types.TickerV2
	for _, t := range tickers {
		var ticker types.TickerV2
		ticker.Last = t["price"]
		ticker.Open24H = t["open"]
		ticker.High24H = t["high"]
		ticker.Low24H = t["low"]
		ticker.BaseVolume24H = t["volume"]
		ticker.QuoteVolume24H = t["quote_volume_24h"]
		ticker.Timestamp = t["timestamp"]
		bestBid, bestAsk := keeper.OrderKeeper.GetBestBidAndAsk(ctx, t["price"])
		ticker.BestBid = bestBid.String()
		ticker.BestAsk = bestAsk.String()
		tickerList = append(tickerList, ticker)
	}

	if len(tickerList) == 0 {
		return nil, nil
	}

	res, err := json.Marshal(tickerList)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryTickerV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	params := types.QueryTickerParams{}
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data, ", err.Error()))
	}

	products := []string{params.Product}
	tickers := keeper.GetTickers(products, 1)

	result := types.DefaultTickerV2(params.Product)
	notExist := true
	for _, t := range tickers {
		if params.Product == t.Product {
			notExist = false
			result.Last = fmt.Sprintf("%f", t.Price)
			result.Open24H = fmt.Sprintf("%f", t.Open)
			result.High24H = fmt.Sprintf("%f", t.High)
			result.Low24H = fmt.Sprintf("%f", t.Low)
			result.BaseVolume24H = fmt.Sprintf("%f", t.Volume)
			result.QuoteVolume24H = fmt.Sprintf("%f", t.Volume)
			result.Timestamp = time.Unix(t.Timestamp, 0).UTC().Format("2006-01-02T15:04:05.000Z")
			break
		}
	}

	if notExist {
		return nil, nil
	}

	bestBid, bestAsk := keeper.OrderKeeper.GetBestBidAndAsk(ctx, params.Product)
	result.BestBid = bestBid.String()
	result.BestAsk = bestAsk.String()

	res, err := json.Marshal(result)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryTickerListV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	tickers := keeper.getAllTickers()
	var tickerList []types.TickerV2
	for _, t := range tickers {
		var ticker types.TickerV2
		ticker.Last = fmt.Sprintf("%f", t.Price)
		ticker.Open24H = fmt.Sprintf("%f", t.Open)
		ticker.High24H = fmt.Sprintf("%f", t.High)
		ticker.Low24H = fmt.Sprintf("%f", t.Low)
		ticker.BaseVolume24H = fmt.Sprintf("%f", t.Volume)
		ticker.QuoteVolume24H = fmt.Sprintf("%f", t.Volume)
		ticker.Timestamp = time.Unix(t.Timestamp, 0).UTC().Format("2006-01-02T15:04:05.000Z")
		bestBid, bestAsk := keeper.OrderKeeper.GetBestBidAndAsk(ctx, t.Product)
		ticker.BestBid = bestBid.String()
		ticker.BestAsk = bestAsk.String()
		tickerList = append(tickerList, ticker)
	}
	if len(tickerList) == 0 {
		return nil, nil
	}
	res, err := common.JSONMarshalV2(tickerList)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return res, nil
}

func queryInstrumentsV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	tokenPairs := keeper.dexKeeper.GetTokenPairs(ctx)

	var result []*types.InstrumentV2
	for _, t := range tokenPairs {
		if t == nil {
			panic("the nil pointer is not expected")
		}
		result = append(result, types.ConvertTokenPairToInstrumentV2(t))
	}

	res, err := json.Marshal(result)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryOrderListV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryOrderParamsV2
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	orders := keeper.getOrderListV2(ctx, params.Product, params.Address, params.Side, params.IsOpen, params.After, params.Before, params.Limit)

	var result []types.OrderV2

	for _, o := range orders {
		tmp := types.ConvertOrderToOrderV2(o)
		result = append(result, tmp)
	}

	if len(result) == 0 {
		return nil, nil
	}

	res, err := json.Marshal(result)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryOrderV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryOrderParamsV2
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	order := keeper.getOrderByIDV2(ctx, params.OrderID)

	if order == nil {
		return nil, nil
	}

	result := types.ConvertOrderToOrderV2(*order)

	res, err := json.Marshal(result)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryCandleListFromMarketKeeperV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryKlinesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	ctx.Logger().Debug(fmt.Sprintf("queryCandleList : %+v", params))
	// should init token pair map here
	keeper.marketKeeper.InitTokenPairMap(ctx, keeper.dexKeeper)
	restData, err := keeper.getCandlesByMarketKeeper(params.Product, params.Granularity, params.Size)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	res, err := json.Marshal(restData)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryCandleListV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryKlinesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	ctx.Logger().Debug(fmt.Sprintf("queryCandleList : %+v", params))
	restData, err := keeper.GetCandles(params.Product, params.Granularity, params.Size)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	res, err := json.Marshal(restData)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryMatchResultsV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryMatchParamsV2
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	matches := keeper.getMatchResultsV2(ctx, params.Product, params.After, params.Before, params.Limit)
	if len(matches) == 0 {
		return nil, nil
	}

	res, err := common.JSONMarshalV2(matches)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return res, nil
}

func queryFeeDetailsV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeeDetailsParamsV2
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	_, err = sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("invalid address", err.Error()))
	}

	feeDetails := keeper.getFeeDetailsV2(ctx, params.Address, params.After, params.Before, params.Limit)
	if len(feeDetails) == 0 {
		return nil, nil
	}

	res, err := common.JSONMarshalV2(feeDetails)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryDealsV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryDealsParamsV2
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	deals := keeper.getDealsV2(ctx, params.Address, params.Product, params.Side, params.After, params.Before, params.Limit)
	if len(deals) == 0 {
		return nil, nil
	}

	res, err := common.JSONMarshalV2(deals)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}

func queryTxListV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTxListParamsV2
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	if _, err := sdk.AccAddressFromBech32(params.Address); err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("invalid address", err.Error()))
	}

	txs := keeper.getTransactionListV2(ctx, params.Address, params.TxType, params.After, params.Before, params.Limit)
	if len(txs) == 0 {
		return nil, nil
	}

	res, err := common.JSONMarshalV2(txs)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	return res, nil
}
