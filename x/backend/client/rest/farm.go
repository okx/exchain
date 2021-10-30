package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/okex/exchain/dependence/cosmos-sdk/client/context"
	"github.com/okex/exchain/dependence/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/okex/exchain/x/backend/types"
	"github.com/okex/exchain/x/common"
)

func registerFarmQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r = r.PathPrefix("/farm").Subrouter()
	r.HandleFunc("/pools/{whitelistOrNormal}", farmPoolsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/dashboard/{address}", farmDashboardHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/whitelist/max_apy", farmWhitelistMaxApyHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/pools/{poolName}/staked_info", farmStakedInfoHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/first_pool", firstPoolHandler(cliCtx)).Methods("GET")
}

func farmPoolsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		whitelistOrNormal := vars["whitelistOrNormal"]
		if whitelistOrNormal != "whitelist" && whitelistOrNormal != "normal" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sortColumn := r.URL.Query().Get("sort_column")
		sortDirection := r.URL.Query().Get("sort_direction")
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}
		params := types.NewQueryFarmPoolsParams(whitelistOrNormal, sortColumn, sortDirection, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryFarmPools), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func farmDashboardHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]
		pageStr := r.URL.Query().Get("page")
		perPageStr := r.URL.Query().Get("per_page")

		page, perPage, err := common.Paginate(pageStr, perPageStr)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeInvalidPaginateParam, err.Error())
			return
		}
		// validate request
		if address == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressIsRequired, "bad request: address is required")
			return
		}
		params := types.NewQueryFarmDashboardParams(address, page, perPage)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryFarmDashboard), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func farmWhitelistMaxApyHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryFarmMaxApy), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func farmStakedInfoHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		poolName := vars["poolName"]
		address := r.URL.Query().Get("address")
		// validate request
		if address == "" {
			common.HandleErrorMsg(w, cliCtx, types.CodeAddressIsRequired, "bad request: address is required")
			return
		}
		params := types.NewQueryFarmStakedInfoParams(poolName, address)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryFarmStakedInfo), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func firstPoolHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		poolName := r.URL.Query().Get("pool_name")
		stakeAtStr := r.URL.Query().Get("stake_at")
		claimHeightStr := r.URL.Query().Get("claim_height")

		stakeAt, err := strconv.ParseInt(stakeAtStr, 10, 64)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
			return
		}

		claimHeight, err := strconv.ParseInt(claimHeightStr, 10, 64)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeStrconvFailed, err.Error())
			return
		}

		params := types.NewQueryFarmFirstPoolParams(poolName, address, stakeAt, claimHeight)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/backend/%s", types.QueryFarmFirstPool), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
