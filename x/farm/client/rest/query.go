package rest

import (
	"fmt"
	"github.com/okex/okexchain/x/common"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/okex/okexchain/x/farm/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// TODO: Define your GET REST endpoints
	r.HandleFunc(
		"/farm/pools",
		queryPoolsHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/farm/parameters",
		queryParamsHandlerFn(cliCtx),
	).Methods("GET")
}

func queryPoolsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorArgsWithLimit)
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryPoolsParams(page, limit)
		jsonBytes, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorCodecFails)
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPools)
		res, height, err := cliCtx.QueryWithData(route, jsonBytes)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		fmt.Println(1)
		route := fmt.Sprintf("custom/%s/parameters", types.QuerierRoute)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println(2)

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
