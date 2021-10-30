package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/okex/exchain/x/backend/types"
	"github.com/okex/exchain/x/common"
)

const (
	defaultLimit = "100"
)

// RegisterRoutesV2 - Central function to define routes for interface version 2
func RegisterRoutesV2(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/block_tx_hashes/{blockHeight}", blockTxHashesHandler(cliCtx)).Methods("GET")

	r.HandleFunc("/instruments", instrumentsHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/instruments/ticker", tickerListHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/instruments/{instrument_id}/ticker", tickerHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/orders_pending", orderOpenListHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/orders/list/open", orderOpenListHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/orders/list/closed", orderClosedListHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/orders/{order_id}", orderHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/instruments/{instrument_id}/candles", candleHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/instruments/{instrument_id}matches", matchHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/fees", feesHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/deals", dealsHandlerV2(cliCtx)).Methods("GET")
	r.HandleFunc("/transactions", txListHandlerV2(cliCtx)).Methods("GET")
}

func txListHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		txType := r.URL.Query().Get("type")
		after := r.URL.Query().Get("after")
		before := r.URL.Query().Get("before")
		limit := r.URL.Query().Get("limit")

		// validate request
		if address == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}
		var typeInt int
		var err error
		if txType != "" {
			if typeInt, err = strconv.Atoi(txType); err != nil {
				common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
				return
			}
		}

		if _, err := strconv.Atoi(after); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		if _, err := strconv.Atoi(before); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		if limit == "" {
			limit = defaultLimit
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := types.QueryTxListParamsV2{
			Address: address,
			TxType:  typeInt,
			After:   after,
			Before:  before,
			Limit:   limitInt,
		}
		req := cliCtx.Codec.MustMarshalJSON(params)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryTxListV2), req)
		common.HandleResponseV2(w, res, err)
	}
}

func dealsHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		product := r.URL.Query().Get("instrument_id")
		side := strings.ToUpper(r.URL.Query().Get("side"))
		after := r.URL.Query().Get("after")
		before := r.URL.Query().Get("before")
		limit := r.URL.Query().Get("limit")

		// validate request
		if address == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidAddress)
			return
		}
		if _, err := strconv.Atoi(after); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		if _, err := strconv.Atoi(before); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		// default limit 100
		if limit == "" {
			limit = defaultLimit
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := types.QueryDealsParamsV2{
			Address: address,
			Product: product,
			Side:    side,
			After:   after,
			Before:  before,
			Limit:   limitInt,
		}
		req := cliCtx.Codec.MustMarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryDealListV2), req)
		common.HandleResponseV2(w, res, err)
	}
}

func feesHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		after := r.URL.Query().Get("after")
		before := r.URL.Query().Get("before")
		limit := r.URL.Query().Get("limit")

		// validate request
		if address == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidAddress)
			return
		}
		if _, err := strconv.Atoi(after); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		if _, err := strconv.Atoi(before); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		// default limit 100
		if limit == "" {
			limit = defaultLimit
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := types.QueryFeeDetailsParamsV2{
			Address: address,
			After:   after,
			Before:  before,
			Limit:   limitInt,
		}

		req := cliCtx.Codec.MustMarshalJSON(params)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryFeeDetailsV2), req)
		common.HandleResponseV2(w, res, err)
	}
}

func matchHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		product := vars["instrument_id"]

		after := r.URL.Query().Get("after")
		before := r.URL.Query().Get("before")
		limit := r.URL.Query().Get("limit")

		// validate request
		if product == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}
		if _, err := strconv.Atoi(after); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		if _, err := strconv.Atoi(before); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		// default limit 100
		if limit == "" {
			limit = defaultLimit
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := types.QueryMatchParamsV2{
			Product: product,
			After:   after,
			Before:  before,
			Limit:   limitInt,
		}
		req := cliCtx.Codec.MustMarshalJSON(params)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryMatchResultsV2), req)
		common.HandleResponseV2(w, res, err)
	}
}

func candleHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		product := vars["instrument_id"]

		strGranularity := r.URL.Query().Get("granularity")
		strSize := r.URL.Query().Get("size")

		if len(strSize) == 0 {
			strSize = defaultLimit
		}

		if len(strGranularity) == 0 {
			strGranularity = "60"
		}

		size, err := strconv.Atoi(strSize)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		granularity, err := strconv.Atoi(strGranularity)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := types.NewQueryKlinesParams(product, granularity, size)

		req := cliCtx.Codec.MustMarshalJSON(params)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryCandleListV2), req)
		common.HandleResponseV2(w, res, err)
	}
}

func tickerListHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryTickerListV2), nil)

		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorServerException)
			return
		}

		common.HandleSuccessResponseV2(w, res)

	}
}

func tickerHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		product := vars["instrument_id"]

		params := types.QueryTickerParams{
			Product: product,
		}

		req := cliCtx.Codec.MustMarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryTickerV2), req)

		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorServerException)
			return
		}

		common.HandleSuccessResponseV2(w, res)

	}
}

func instrumentsHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryInstrumentsV2), nil)

		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorServerException)
			return
		}

		common.HandleSuccessResponseV2(w, res)
	}
}

func orderOpenListHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := r.URL.Query().Get("instrument_id")
		address := r.URL.Query().Get("address")
		side := r.URL.Query().Get("side")
		after := r.URL.Query().Get("after")
		before := r.URL.Query().Get("before")
		limit := r.URL.Query().Get("limit")

		// validate request
		if product == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}
		// default limit 100
		if limit == "" {
			limit = defaultLimit
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := types.QueryOrderParamsV2{
			Product: product,
			Address: address,
			Side:    side,
			After:   after,
			Before:  before,
			Limit:   limitInt,
			IsOpen:  true,
		}

		req := cliCtx.Codec.MustMarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryOrderListV2), req)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorServerException)
			return
		}

		common.HandleSuccessResponseV2(w, res)
	}
}

func orderClosedListHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := r.URL.Query().Get("instrument_id")
		address := r.URL.Query().Get("address")
		side := r.URL.Query().Get("side")
		after := r.URL.Query().Get("after")
		before := r.URL.Query().Get("before")
		limit := r.URL.Query().Get("limit")

		// validate request
		if product == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}
		// default limit 100
		if limit == "" {
			limit = defaultLimit
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		params := types.QueryOrderParamsV2{
			Product: product,
			Address: address,
			Side:    side,
			After:   after,
			Before:  before,
			Limit:   limitInt,
			IsOpen:  false,
		}

		req := cliCtx.Codec.MustMarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryOrderListV2), req)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorServerException)
			return
		}

		common.HandleSuccessResponseV2(w, res)
	}
}

func orderHandlerV2(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["order_id"]

		// validate request
		if orderID == "" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorMissingRequiredParam)
			return
		}

		params := types.QueryOrderParamsV2{
			OrderID: orderID,
		}

		req := cliCtx.Codec.MustMarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryOrderV2), req)

		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorServerException)
			return
		}

		common.HandleSuccessResponseV2(w, res)
	}
}
