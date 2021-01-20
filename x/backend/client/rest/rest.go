package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/okex/okexchain/x/backend/client/cli"
	"github.com/okex/okexchain/x/backend/types"
	"github.com/okex/okexchain/x/common"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/candles/{product}", candleHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/tickers", tickerHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/tickers/{product}", tickerHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/matches", matchHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/deals", dealHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/fees", feeDetailListHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/order/list/{openOrClosed}", orderListHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/orders/{orderID}", orderHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/accounts/{address}/orders", accountOrdersHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/block_tx_hashes/{blockHeight}", blockTxHashesHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/transactions", txListHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/latestheight", latestHeightHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/dex/fees", dexFeesHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/swap/watchlist", swapWatchlistHandler(cliCtx)).Methods("GET")

	// register farm rest
	registerFarmQueryRoutes(cliCtx, r)
}

func candleHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		product := vars["product"]

		strGranularity := r.URL.Query().Get("granularity")
		strSize := r.URL.Query().Get("size")

		if len(strSize) == 0 {
			strSize = "100"
		}

		if len(strGranularity) == 0 {
			strGranularity = "60"
		}

		size, err0 := strconv.Atoi(strSize)
		if err0 != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, fmt.Sprintf("parameter size %s not correct", strSize))
			return
		}
		granularity, err1 := strconv.Atoi(strGranularity)
		if err1 != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, fmt.Sprintf("parameter granularity %s not correct", strGranularity))
			return
		}

		params := types.NewQueryKlinesParams(product, granularity, size)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryCandleList), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func tickerHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		product := vars["product"]

		strCount := r.URL.Query().Get("count")
		strSort := r.URL.Query().Get("sort")

		if strCount == "" {
			strCount = "100"
		}

		if len(strSort) == 0 {
			strSort = "true"
		}

		sort, errSort := strconv.ParseBool(strSort)
		count, errCnt := strconv.Atoi(strCount)
		mErr := types.NewErrorsMerged(errSort, errCnt)
		if mErr != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, mErr.Error())
			return
		}

		params := types.QueryTickerParams{
			Product: product,
			Sort:    sort,
			Count:   count,
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryTickerList), bz)

		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)

	}
}

func matchHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := r.URL.Query().Get("product")
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")

		// validate request
		if product == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeProductIsRequired, "invalid params: product is required")
			return
		}
		var start, end int64
		var err error
		if startStr != "" {
			if start, err = strconv.ParseInt(startStr, 10, 64); err != nil {
				common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
				return
			}
		}
		if endStr != "" {
			if end, err = strconv.ParseInt(endStr, 10, 64); err != nil {
				common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
				return
			}
		}

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}

		params := types.NewQueryMatchParams(product, start, end, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryMatchResults), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func dealHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := r.URL.Query().Get("address")
		product := r.URL.Query().Get("product")
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")
		sideStr := r.URL.Query().Get("side")

		// validate request
		if addr == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressIsRequired, "bad request: address is required")
			return
		}
		var start, end int64
		var err error
		if startStr != "" {
			if start, err = strconv.ParseInt(startStr, 10, 64); err != nil {
				common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
				return
			}
		}
		if endStr != "" {
			if end, err = strconv.ParseInt(endStr, 10, 64); err != nil {
				common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
				return
			}
		}

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
			return
		}

		params := types.NewQueryDealsParams(addr, product, start, end, page, perPage, sideStr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryDealList), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func feeDetailListHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := r.URL.Query().Get("address")
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")

		// validate request
		if addr == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressIsRequired, "bad request: address is required")
			return
		}
		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}
		params := types.NewQueryFeeDetailsParams(addr, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryFeeDetails), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func orderListHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		openOrClosed := vars["openOrClosed"]
		if openOrClosed != "open" && openOrClosed != "closed" {
			common.HandleErrorMsg(w, cliCtx, types.CodeOrderStatusMustBeOpenOrClosed, fmt.Sprintf("order status should be open/closed"))
			return
		}
		addr := r.URL.Query().Get("address")
		product := r.URL.Query().Get("product")
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")
		sideStr := r.URL.Query().Get("side")
		hideNoFillStr := r.URL.Query().Get("hide_no_fill")

		// validate request
		if addr == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressIsRequired, "bad request: address is required")
			return
		}
		var start, end int64
		var err error
		if startStr == "" {
			startStr = "0"
		}
		if endStr == "" {
			endStr = "0"
		}
		start, errStart := strconv.ParseInt(startStr, 10, 64)
		end, errEnd := strconv.ParseInt(endStr, 10, 64)
		mErr := types.NewErrorsMerged(errStart, errEnd)
		if mErr != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, mErr.Error())
			return
		}

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}

		hideNoFill := hideNoFillStr == "1"

		params := types.NewQueryOrderListParams(
			addr, product, sideStr, page, perPage, start, end, hideNoFill)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s/%s", types.QueryOrderList, openOrClosed), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func orderHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["orderID"]
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s/%s", types.QueryOrderByID, orderID), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func accountOrdersHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")

		// validate request
		if address == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressIsRequired, "bad request: address is required")
			return
		}

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
			return
		}

		var start, end int64
		if startStr == "" {
			startStr = "0"
		}
		if endStr == "" {
			endStr = "0"
		}
		start, errStart := strconv.ParseInt(startStr, 10, 64)
		end, errEnd := strconv.ParseInt(endStr, 10, 64)
		mErr := types.NewErrorsMerged(errStart, errEnd)
		if mErr != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, mErr.Error())
			return
		}

		params := types.NewQueryAccountOrdersParams(address, start, end, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryAccountOrders), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func txListHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := r.URL.Query().Get("address")
		txTypeStr := r.URL.Query().Get("type")
		startStr := r.URL.Query().Get("start")
		endStr := r.URL.Query().Get("end")
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")

		// validate request
		if addr == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressIsRequired, "bad request: address is required")
			return
		}
		var txType, start, end int64
		var err error
		if txTypeStr != "" {
			if txType, err = strconv.ParseInt(txTypeStr, 10, 64); err != nil {
				common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
				return
			}
		}
		if startStr != "" {
			if start, err = strconv.ParseInt(startStr, 10, 64); err != nil {
				common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
				return
			}
		}
		if endStr != "" {
			if end, err = strconv.ParseInt(endStr, 10, 64); err != nil {
				common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
				return
			}
		}

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}
		params := types.NewQueryTxListParams(addr, txType, start, end, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryTxList), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func blockTxHashesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blockHeightStr := vars["blockHeight"]
		blockHeight, err := strconv.ParseInt(blockHeightStr, 10, 64)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
			return
		}
		res, err := cli.GetBlockTxHashes(cliCtx, blockHeight)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, types.CodeGetBlockTxHashesFailed,
				fmt.Sprintf("failed to get block tx hash: %s", err.Error()))
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func latestHeightHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h, err := rpc.GetChainHeight(cliCtx)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, types.CodeGetChainHeightFailed,
				fmt.Sprintf("failed to get chain height: %s", err.Error()))
			return
		}
		res := common.GetBaseResponse(h)
		bz, err := json.Marshal(res)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
		}
		rest.PostProcessResponse(w, cliCtx, bz)
	}
}

func dexFeesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		baseAsset := r.URL.Query().Get("base_asset")
		quoteAsset := r.URL.Query().Get("quote_asset")
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")
		if address == "" && baseAsset == "" && quoteAsset == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressAndProductRequired, "bad request: address、base_asset and quote_asset could not be empty at the same time")
			return
		}

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}

		params := types.NewQueryDexFeesParams(address, baseAsset, quoteAsset, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryDexFeesList), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func swapWatchlistHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")
		sortColumn := r.URL.Query().Get("sort_column")
		sortDirection := r.URL.Query().Get("sort_direction")

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}
		params := types.NewQuerySwapWatchlistParams(sortColumn, sortDirection, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QuerySwapWatchlist), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
