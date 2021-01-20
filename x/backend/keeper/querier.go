package keeper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend/types"
	"github.com/okex/okexchain/x/common"
	orderTypes "github.com/okex/okexchain/x/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		if !keeper.Config.EnableBackend {
			response := common.GetErrorResponse(types.CodeBackendPluginNotEnabled, "Backend Plugin's Not Enabled", "Backend Plugin's Not Enabled")
			res, eJSON := json.Marshal(response)
			if eJSON != nil {
				return nil, common.ErrMarshalJSONFailed(eJSON.Error())
			}
			return res, nil
		}

		defer func() {
			if e := recover(); e != nil {
				errMsg := fmt.Sprintf("%+v", e)
				response := common.GetErrorResponse(types.CodeGoroutinePanic, errMsg, errMsg)
				resJSON, eJSON := json.Marshal(response)
				if eJSON != nil {
					res = nil
					err = common.ErrMarshalJSONFailed(eJSON.Error())
				} else {
					res = resJSON
					err = nil
				}

			}
		}()

		switch path[0] {
		case types.QueryMatchResults:
			res, err = queryMatchResults(ctx, path[1:], req, keeper)
		case types.QueryDealList:
			res, err = queryDeals(ctx, path[1:], req, keeper)
		case types.QueryFeeDetails:
			res, err = queryFeeDetails(ctx, path[1:], req, keeper)
		case types.QueryOrderList:
			res, err = queryOrderList(ctx, path[1:], req, keeper)
		case types.QueryOrderByID:
			res, err = queryOrderByID(ctx, path[1:], req, keeper)
		case types.QueryAccountOrders:
			res, err = queryAccountOrders(ctx, path[1:], req, keeper)
		case types.QueryTxList:
			res, err = queryTxList(ctx, path[1:], req, keeper)
		case types.QueryCandleList:
			if keeper.Config.EnableMktCompute {
				res, err = queryCandleList(ctx, path[1:], req, keeper)
			} else {
				res, err = queryCandleListFromMarketKeeper(ctx, path[1:], req, keeper)
			}
		case types.QueryTickerList:
			if keeper.Config.EnableMktCompute {
				res, err = queryTickerList(ctx, path[1:], req, keeper)
			} else {
				res, err = queryTickerListFromMarketKeeper(ctx, path[1:], req, keeper)
			}
		case types.QueryDexFeesList:
			res, err = queryDexFees(ctx, path[1:], req, keeper)

		case types.QuerySwapWatchlist:
			res, err = querySwapWatchlist(ctx, req, keeper)
		case types.QuerySwapTokens:
			res, err = querySwapTokens(ctx, req, keeper)
		case types.QuerySwapTokenPairs:
			res, err = querySwapTokenPairs(ctx, path[1:], req, keeper)
		case types.QuerySwapLiquidityHistories:
			res, err = querySwapLiquidityHistories(ctx, req, keeper)
		case types.QueryTickerListV2:
			if keeper.Config.EnableMktCompute {
				res, err = queryTickerListV2(ctx, path[1:], req, keeper)
			} else {
				res, err = queryTickerListFromMarketKeeperV2(ctx, path[1:], req, keeper)
			}
		case types.QueryTickerV2:
			if keeper.Config.EnableMktCompute {
				res, err = queryTickerV2(ctx, path[1:], req, keeper)
			} else {
				res, err = queryTickerFromMarketKeeperV2(ctx, path[1:], req, keeper)
			}
		case types.QueryInstrumentsV2:
			res, err = queryInstrumentsV2(ctx, path[1:], req, keeper)
		case types.QueryOrderListV2:
			res, err = queryOrderListV2(ctx, path[1:], req, keeper)
		case types.QueryOrderV2:
			res, err = queryOrderV2(ctx, path[1:], req, keeper)
		case types.QueryCandleListV2:
			if keeper.Config.EnableMktCompute {
				res, err = queryCandleListV2(ctx, path[1:], req, keeper)
			} else {
				res, err = queryCandleListFromMarketKeeperV2(ctx, path[1:], req, keeper)
			}
		case types.QueryMatchResultsV2:
			res, err = queryMatchResultsV2(ctx, path[1:], req, keeper)
		case types.QueryFeeDetailsV2:
			res, err = queryFeeDetailsV2(ctx, path[1:], req, keeper)
		case types.QueryDealListV2:
			res, err = queryDealsV2(ctx, path[1:], req, keeper)
		case types.QueryTxListV2:
			res, err = queryTxListV2(ctx, path[1:], req, keeper)
		default:
			res, err = nil, types.ErrBackendModuleUnknownQueryType()
		}

		if err != nil {
			response := common.GetErrorResponse(types.CodeBackendModuleUnknownQueryType, "backend module unknown query type", err.Error())
			res, eJSON := json.Marshal(response)
			if eJSON != nil {
				return nil, common.ErrMarshalJSONFailed(eJSON.Error())
			}
			return res, err
		}

		return res, nil
	}
}

// nolint: unparam
func queryDeals(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryDealsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	if params.Address != "" {
		_, err := sdk.AccAddressFromBech32(params.Address)
		if err != nil {
			return nil, common.ErrCreateAddrFromBech32Failed(params.Address, err.Error())
		}
	}
	if params.Side != "" && params.Side != orderTypes.BuyOrder && params.Side != orderTypes.SellOrder {
		return nil, types.ErrOrderSideParamMustBuyOrSell(params.Side)
	}
	if params.Page < 0 || params.PerPage < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}

	offset, limit := common.GetPage(params.Page, params.PerPage)
	deals, total := keeper.GetDeals(ctx, params.Address, params.Product, params.Side, params.Start, params.End, offset, limit)
	var response *common.ListResponse
	if len(deals) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, deals)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// nolint: unparam
func queryMatchResults(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryMatchParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	if params.Page < 0 || params.PerPage < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}
	offset, limit := common.GetPage(params.Page, params.PerPage)
	matches, total := keeper.getMatchResults(ctx, params.Product, params.Start, params.End, offset, limit)
	var response *common.ListResponse
	if len(matches) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, matches)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// nolint: unparam
func queryFeeDetails(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryFeeDetailsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	_, err = sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, common.ErrCreateAddrFromBech32Failed(params.Address, err.Error())
	}
	if params.Page < 0 || params.PerPage < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}

	offset, limit := common.GetPage(params.Page, params.PerPage)
	feeDetails, total := keeper.GetFeeDetails(ctx, params.Address, offset, limit)
	var response *common.ListResponse
	if len(feeDetails) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, feeDetails)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryCandleList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {

	var params types.QueryKlinesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	if params.Product == "" {
		return nil, types.ErrProductIsRequired()
	}
	if keeper.dexKeeper.GetTokenPair(ctx, params.Product) == nil {
		return nil, types.ErrProductDoesNotExist(params.Product)
	}

	ctx.Logger().Debug(fmt.Sprintf("queryCandleList : %+v", params))
	restData, err := keeper.GetCandles(params.Product, params.Granularity, params.Size)

	var response *common.BaseResponse
	if err != nil {
		response = common.GetErrorResponse(types.CodeGetCandlesFailed, err.Error(), err.Error())
	} else {
		response = common.GetBaseResponse(restData)
	}

	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryCandleListFromMarketKeeper(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryKlinesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	if params.Product == "" {
		return nil, types.ErrProductIsRequired()
	}
	tokenPair := keeper.dexKeeper.GetTokenPair(ctx, params.Product)
	if tokenPair == nil {
		return nil, types.ErrProductDoesNotExist(params.Product)
	}

	ctx.Logger().Debug(fmt.Sprintf("queryCandleList : %+v", params))
	// should init token pair map here
	restData, err := keeper.getCandlesByMarketKeeper(tokenPair.ID, params.Granularity, params.Size)

	var response *common.BaseResponse
	if err != nil {
		response = common.GetErrorResponse(types.CodeGetCandlesByMarketFailed, err.Error(), err.Error())
	} else {
		response = common.GetBaseResponse(restData)
	}

	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryTickerList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	params := types.QueryTickerParams{}
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}

	products := []string{}
	if params.Product != "" {
		if tokenPair := keeper.dexKeeper.GetTokenPair(ctx, params.Product); tokenPair == nil {
			return nil, types.ErrProductDoesNotExist(params.Product)
		}
		products = append(products, params.Product)
	} else {
		products = keeper.getAllProducts(ctx)
	}

	// set default count to 10
	if params.Count <= 0 {
		params.Count = 10
	}

	addedTickers := []types.Ticker{}
	tickers := keeper.GetTickers(products, params.Count)
	for _, p := range products {

		exists := false
		for _, t := range tickers {
			if p == t.Product {
				exists = true
				break
			}
		}

		if !exists {
			//tmpPrice := keeper.orderKeeper.GetLastPrice(ctx, p)
			tmpTicker := types.Ticker{
				Price:            -1,
				Product:          p,
				Symbol:           p,
				Open:             0,
				Close:            0,
				High:             0,
				Low:              0,
				Volume:           0,
				Change:           0,
				ChangePercentage: "0.00%",
				Timestamp:        time.Now().Unix(),
			}
			addedTickers = append(addedTickers, tmpTicker)
		}

	}

	if len(addedTickers) > 0 {
		tickers = append(tickers, addedTickers...)
	}

	var sortedTickers types.Tickers = tickers
	sort.Sort(sortedTickers)

	response := common.GetBaseResponse(sortedTickers)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryTickerListFromMarketKeeper(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	params := types.QueryTickerParams{}
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}

	var products []string
	if params.Product != "" {
		if tokenPair := keeper.dexKeeper.GetTokenPair(ctx, params.Product); tokenPair == nil {
			return nil, types.ErrProductDoesNotExist(params.Product)
		}
		products = append(products, params.Product)
	} else {
		products = keeper.getAllProducts(ctx)
	}

	// set default count to 10
	if params.Count <= 0 {
		params.Count = 10
	}

	allTickers, err := keeper.marketKeeper.GetTickerByProducts(products)

	var filterTickers []map[string]string
	for _, p := range products {
		exists := false
		for _, t := range allTickers {
			if p == t["product"] {
				filterTickers = append(filterTickers, t)
				exists = true
				break
			}
		}

		if !exists {
			tmpTicker := map[string]string{
				"price":   "-1",
				"product": p,
				"symbol":  p,
				"open":    "0",
				"close":   "0",
				"high":    "0",
				"low":     "0",
				"volume":  "0",
				"change":  "0",
				//"changePercentage": "0.00%",
				"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			}
			filterTickers = append(filterTickers, tmpTicker)
		}

	}

	if len(filterTickers) > params.Count {
		filterTickers = filterTickers[0:params.Count]
	}

	var response *common.BaseResponse
	if err != nil {
		response = common.GetErrorResponse(types.CodeGetTickerByProductsFailed, "", err.Error())
	} else {
		response = common.GetBaseResponse(filterTickers)
	}

	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryOrderList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	isOpen := path[0] == "open"
	var params types.QueryOrderListParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	_, err = sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, common.ErrCreateAddrFromBech32Failed(params.Address, err.Error())
	}
	if params.Page < 0 || params.PerPage < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}
	offset, limit := common.GetPage(params.Page, params.PerPage)
	orders, total := keeper.GetOrderList(ctx, params.Address, params.Product, params.Side, isOpen,
		offset, limit, params.Start, params.End, params.HideNoFill)

	var response *common.ListResponse
	if len(orders) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, orders)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

func queryOrderByID(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) == 0 {
		return nil, types.ErrOrderIdIsRequired()
	}
	orderID := path[0]
	order := keeper.Orm.GetOrderByID(orderID)
	response := common.GetBaseResponse(order)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryAccountOrders(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryAccountOrdersParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	if params.Address != "" {
		_, err := sdk.AccAddressFromBech32(params.Address)
		if err != nil {
			return nil, common.ErrCreateAddrFromBech32Failed(params.Address, err.Error())
		}
	}
	if params.Page < 0 || params.PerPage < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}

	offset, limit := common.GetPage(params.Page, params.PerPage)
	orders, total := keeper.Orm.GetAccountOrders(params.Address, params.Start, params.End, offset, limit)
	var response *common.ListResponse
	if len(orders) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, orders)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryTxList(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTxListParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	_, err = sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, common.ErrCreateAddrFromBech32Failed(params.Address, err.Error())
	}
	if params.Page < 0 || params.PerPage < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}
	offset, limit := common.GetPage(params.Page, params.PerPage)
	txs, total := keeper.GetTransactionList(ctx, params.Address, params.TxType, params.StartTime, params.EndTime, offset, limit)

	var response *common.ListResponse
	if len(txs) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, txs)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func queryDexFees(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryDexFeesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	if params.Page < 0 || params.PerPage < 0 {
		return nil, common.ErrInvalidPaginateParam(params.Page, params.PerPage)
	}
	offset, limit := common.GetPage(params.Page, params.PerPage)

	var fees []types.DexFees
	var total int
	if params.BaseAsset == "" && params.QuoteAsset == "" {
		fees, total = keeper.GetDexFees(ctx, params.DexHandlingAddr, "", offset, limit)
	} else { // filter base asset and quote asset
		tokenPairs := keeper.dexKeeper.GetTokenPairs(ctx)
		for _, tokenPair := range tokenPairs {
			if params.BaseAsset != "" && !strings.Contains(tokenPair.BaseAssetSymbol, params.BaseAsset) {
				continue
			}
			if params.QuoteAsset != "" && !strings.Contains(tokenPair.QuoteAssetSymbol, params.QuoteAsset) {
				continue
			}
			partialFees, partial := keeper.GetDexFees(ctx, params.DexHandlingAddr, tokenPair.Name(), offset, limit)
			fees = append(fees, partialFees...)
			total += partial
		}
	}

	var response *common.ListResponse
	if len(fees) > 0 {
		response = common.GetListResponse(total, params.Page, params.PerPage, fees)
	} else {
		response = common.GetEmptyListResponse(total, params.Page, params.PerPage)
	}
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}
