package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/order/keeper"
	"github.com/okex/okexchain/x/order/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/order/depthbook", orderBookHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/order/{orderID}", orderDetailHandler(cliCtx)).Methods("GET")
}

func orderDetailHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderID := vars["orderID"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/order/detail/%s", orderID), nil)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, err.Error())
			return
		}

		order2 := &types.Order{}
		codec.Cdc.MustUnmarshalJSON(res, order2)
		response := common.GetBaseResponse(order2)
		resBytes, err2 := json.Marshal(response)
		if err2 != nil {
			common.HandleErrorMsg(w, cliCtx, err2.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, resBytes)
	}
}

func orderBookHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		product := r.URL.Query().Get("product")
		sizeStr := r.URL.Query().Get("size")
		// validate request
		if product == "" {
			common.HandleErrorMsg(w, cliCtx, "Bad request: product is empty")
			return
		}
		var size int
		var err error
		if sizeStr != "" {
			size, err = strconv.Atoi(sizeStr)
			if err != nil {
				common.HandleErrorMsg(w, cliCtx, err.Error())
				return
			}
		}
		if size < 0 {
			common.HandleErrorMsg(w, cliCtx, "Bad request: size is invalid")
			return
		}
		params := keeper.NewQueryDepthBookParams(product, uint(size))
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/order/depthbook", bz)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, err.Error())
			return
		}

		bookRes := &keeper.BookRes{}
		codec.Cdc.MustUnmarshalJSON(res, bookRes)
		response := common.GetBaseResponse(bookRes)
		resBytes, err2 := json.Marshal(response)
		if err2 != nil {
			common.HandleErrorMsg(w, cliCtx, err2.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, resBytes)
	}
}
