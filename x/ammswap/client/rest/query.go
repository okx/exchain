package rest

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/common"
	"net/http"
	"strings"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/ammswap/swap_token_pair", querySwapTokenPairHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/ammswap/swap_token_pairs", querySwapTokenPairsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/ammswap/params", queryParamsHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/ammswap/buy_amount", queryBuyAmountHandler(cliCtx)).Methods("GET")
	r.HandleFunc("/ammswap/redeemable_assets", queryRedeemableAssetsHandler(cliCtx)).Methods("GET")
}

func querySwapTokenPairHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		baseToken := r.URL.Query().Get("base_token")
		quoteToken := r.URL.Query().Get("quote_token")

		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QuerySwapTokenPair, baseToken, quoteToken), nil)
		if err != nil {
			common.HandleErrorMsg(w, cliContext, err.Error())
			return
		}

		formatAndReturnResult(w, cliContext, res)
	}

}

func querySwapTokenPairsHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySwapTokenPairs), nil)
		if err != nil {
			common.HandleErrorMsg(w, cliContext, err.Error())
			return
		}

		formatAndReturnResult(w, cliContext, res)
	}

}

func queryParamsHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/params", types.QuerierRoute), nil)
		if err != nil {
			common.HandleErrorMsg(w, cliContext, err.Error())
			return
		}

		formatAndReturnResult(w, cliContext, res)
	}

}

func queryBuyAmountHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		soldTokenStr := r.URL.Query().Get("sold_token")
		tokenToBuyStr := r.URL.Query().Get("token_to_buy")

		sellToken, err := sdk.ParseDecCoin(soldTokenStr)
		if err != nil {
			common.HandleErrorMsg(w, cliContext, err.Error())
			return
		}
		params := types.QueryBuyAmountParams{
			SoldToken:  sellToken,
			TokenToBuy: tokenToBuyStr,
		}
		bz, err := codec.Cdc.MarshalJSON(params)
		if err != nil {
			common.HandleErrorMsg(w, cliContext, err.Error())
			return
		}
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBuyAmount), bz)
		if err != nil {
			common.HandleErrorMsg(w, cliContext, err.Error())
			return
		}

		formatAndReturnResult(w, cliContext, res)
	}

}

func queryRedeemableAssetsHandler(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		baseTokenName := r.URL.Query().Get("base_token_name")
		quoteTokenName := r.URL.Query().Get("quote_token_name")
		liquidity := r.URL.Query().Get("liquidity")
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s/%s", types.QuerierRoute, types.QueryRedeemableAssets, baseTokenName, quoteTokenName, liquidity), nil)
		if err != nil {
			common.HandleErrorMsg(w, cliContext, err.Error())
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
		common.HandleErrorMsg(w, cliContext, err.Error())
		return
	}
	resultJson = []byte(strings.Replace(string(resultJson), "\"" + replaceStr + "\"", string(data), 1))

	rest.PostProcessResponse(w, cliContext, resultJson)
}