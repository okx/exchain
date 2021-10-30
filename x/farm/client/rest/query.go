package rest

import (
	"fmt"
	"net/http"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/dependence/cosmos-sdk/client/context"
	"github.com/okex/exchain/dependence/cosmos-sdk/types/rest"
	"github.com/okex/exchain/x/farm/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// get the current state of the all farm pools
	r.HandleFunc(
		"/farm/pools",
		queryPoolsHandlerFn(cliCtx),
	).Methods("GET")

	// get a single pool info by the farm pool's name
	r.HandleFunc(
		"/farm/pool/{poolName}",
		queryPoolHandlerFn(cliCtx),
	).Methods("GET")

	// get the current earnings of an account in a farm pool
	r.HandleFunc(
		"/farm/earnings/{poolName}/{accAddr}",
		queryEarningsHandlerFn(cliCtx),
	).Methods("GET")

	// get the white list info
	r.HandleFunc(
		"/farm/whitelist",
		queryWhitelistHandlerFn(cliCtx),
	).Methods("GET")

	// get the current farm parameter values
	r.HandleFunc(
		"/farm/parameters",
		queryParamsHandlerFn(cliCtx),
	).Methods("GET")

	// get the names of all farm pools that the account has locked coins in
	r.HandleFunc(
		"/farm/account/{accAddr}",
		queryAccountHandlerFn(cliCtx),
	).Methods("GET")
}

func queryAccountHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accAddr, err := sdk.AccAddressFromBech32(mux.Vars(r)["accAddr"])
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeCreateAddrFromBech32Failed, err.Error())
		}

		jsonBytes, err := cliCtx.Codec.MarshalJSON(types.NewQueryAccountParams(accAddr))
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorCodecFails)
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccount)
		res, height, err := cliCtx.QueryWithData(route, jsonBytes)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryEarningsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		varsMap := mux.Vars(r)
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		accAddr, err := sdk.AccAddressFromBech32(varsMap["accAddr"])
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeCreateAddrFromBech32Failed, err.Error())
		}

		jsonBytes, err := cliCtx.Codec.MarshalJSON(types.NewQueryPoolAccountParams(varsMap["poolName"], accAddr))
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorCodecFails)
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryEarnings)
		res, height, err := cliCtx.QueryWithData(route, jsonBytes)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryPoolHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		poolName := mux.Vars(r)["poolName"]
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryPoolParams(poolName)

		jsonBytes, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorCodecFails)
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPool)
		res, height, err := cliCtx.QueryWithData(route, jsonBytes)
		if err != nil {
			common.HandleErrorResponseV2(w, http.StatusInternalServerError, common.ErrorABCIQueryFails)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
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

func queryWhitelistHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryWhitelist)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
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

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
