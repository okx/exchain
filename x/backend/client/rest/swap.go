package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/okex/exchain/x/backend/types"
	"github.com/okex/exchain/x/common"
)

func registerSwapQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r = r.PathPrefix("/swap").Subrouter()
	r.HandleFunc("/watchlist", swapWatchlistHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/tokens", swapTokensHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/token_pairs", querySwapTokenPairsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/liquidity/histories", swapLiquidityHistoriesHandler(cliCtx)).Methods("GET")
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

func swapTokensHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		businessType := r.URL.Query().Get("business_type")
		address := r.URL.Query().Get("address")
		baseTokenName := r.URL.Query().Get("base_token")

		params := types.NewQuerySwapTokensParams(businessType, address, baseTokenName)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySwapTokens), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func querySwapTokenPairsHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySwapTokenPairs), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliContext, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}

}

func swapLiquidityHistoriesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		tokenPairName := r.URL.Query().Get("token_pair_name")

		params := types.NewQuerySwapLiquidityInfoParams(address, tokenPairName)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySwapLiquidityHistories), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
