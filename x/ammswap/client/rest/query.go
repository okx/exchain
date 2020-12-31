package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/common"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r = r.PathPrefix("/swap").Subrouter()
	r.HandleFunc("/tokens", swapTokensHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/token_pair/{name}", querySwapTokenPairHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/token_pairs", querySwapTokenPairsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/params", queryParamsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/liquidity/add_quote/{token}", swapAddQuoteHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/liquidity/histories", swapLiquidityHistoriesHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/liquidity/remove_quote/{token_pair}", queryRedeemableAssetsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/quote/{token}", swapQuoteHandler(cliCtx)).Methods("GET")
}

func querySwapTokenPairHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenPairName := vars["name"]
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySwapTokenPair, tokenPairName), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliContext, sdkErr.Code, sdkErr.Message)
			return
		}
		rest.PostProcessResponse(w, cliContext, res)
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

		formatAndReturnResult(w, cliContext, res)
	}

}

func queryParamsHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/params", types.QuerierRoute), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliContext, sdkErr.Code, sdkErr.Message)
			return
		}

		formatAndReturnResult(w, cliContext, res)
	}

}

func queryRedeemableAssetsHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenPair := vars["token_pair"]
		liquidity := r.URL.Query().Get("liquidity")
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryRedeemableAssets, tokenPair, liquidity), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliContext, sdkErr.Code, sdkErr.Message)
			return
		}
		formatAndReturnResult(w, cliContext, res)
	}

}

func formatAndReturnResult(w http.ResponseWriter, cliContext context.CLIContext, data []byte) {
	replaceStr := "replaceHere"
	result := common.GetBaseResponse(replaceStr)
	resultJson, err := json.Marshal(result)
	if err != nil {
		common.HandleErrorMsg(w, cliContext, common.CodeMarshalJSONFailed, err.Error())
		return
	}
	resultJson = []byte(strings.Replace(string(resultJson), "\""+replaceStr+"\"", string(data), 1))

	rest.PostProcessResponse(w, cliContext, resultJson)
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

func swapQuoteHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		buyToken := vars["token"]
		sellTokenAmount := r.URL.Query().Get("sell_token_amount")

		params := types.NewQuerySwapBuyInfoParams(sellTokenAmount, buyToken)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySwapQuoteInfo), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
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

func swapAddQuoteHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		baseToken := vars["token"]
		quoteTokenAmount := r.URL.Query().Get("quote_token_amount")

		params := types.NewQuerySwapAddInfoParams(quoteTokenAmount, baseToken)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySwapAddLiquidityQuote), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, sdkErr.Message)
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
