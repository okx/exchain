package rest

import (
	"fmt"
	"net/http"

	"github.com/okex/exchain/x/token/types"

	"encoding/json"
	"strings"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"

	"github.com/okex/exchain/dependence/cosmos-sdk/client/context"
	"github.com/okex/exchain/dependence/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/okex/exchain/x/common"
)

// RegisterRoutes, a central function to define routes
// which is called by the rest module in main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/token/{symbol}"), tokenHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/tokens"), tokensHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/currency/describe"), currencyDescribeHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/accounts/{address}"), spotAccountsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/upload"), uploadAccountsHandler(cliCtx, storeName)).Methods("GET")
}

func tokenHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenName := vars["symbol"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/info/%s", storeName, tokenName), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, err.Error())
			return
		}
		result := common.GetBaseResponse("hello")
		result2, err2 := json.Marshal(result)
		if err2 != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err2.Error())
			return
		}
		result2 = []byte(strings.Replace(string(result2), "\"hello\"", string(res), 1))
		rest.PostProcessResponse(w, cliCtx, result2)
	}
}

func tokensHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerAddress := r.URL.Query().Get("address")
		if _, err := sdk.AccAddressFromBech32(ownerAddress); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/tokens/%s", storeName, ownerAddress), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, err.Error())
			return
		}

		result := common.GetBaseResponse("hello")
		result2, err2 := json.Marshal(result)
		if err2 != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err2.Error())
			return
		}
		result2 = []byte(strings.Replace(string(result2), "\"hello\"", string(res), 1))
		rest.PostProcessResponse(w, cliCtx, result2)
	}
}

func currencyDescribeHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/currency/describe", storeName), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, err.Error())
			return
		}

		result := common.GetBaseResponse("hello")
		result2, err2 := json.Marshal(result)
		if err2 != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err2.Error())
			return
		}
		result2 = []byte(strings.Replace(string(result2), "\"hello\"", string(res), 1))
		rest.PostProcessResponse(w, cliCtx, result2)
	}
}

func spotAccountsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]

		symbol := r.URL.Query().Get("symbol")
		show := r.URL.Query().Get("show")

		if show == "" {
			show = "partial"
		}
		if show != "partial" && show != "all" {
			result := common.GetErrorResponseJSON(1, "", "param show not valid")
			rest.PostProcessResponse(w, cliCtx, result)
			return
		}

		accountParam := types.AccountParam{
			Symbol: symbol,
			Show:   show,
			//QueryPage: token.QueryPage{
			//	Page:    pageInt,
			//	PerPage: perPageInt,
			//},
		}

		bz, err := cliCtx.Codec.MarshalJSON(accountParam)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/accounts/%s", storeName, address), bz)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, err.Error())
			return
		}

		result := common.GetBaseResponse("hello")
		result2, err2 := json.Marshal(result)
		if err2 != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err2.Error())
			return
		}
		result2 = []byte(strings.Replace(string(result2), "\"hello\"", string(res), 1))
		rest.PostProcessResponse(w, cliCtx, result2)
	}
}

func uploadAccountsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/upload", storeName), nil)
		if err != nil {
			sdkErr := common.ParseSDKError(err.Error())
			common.HandleErrorMsg(w, cliCtx, sdkErr.Code, err.Error())
			return
		}
		result := common.GetBaseResponse("hello")
		result2, err2 := json.Marshal(result)
		if err2 != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err2.Error())
			return
		}
		result2 = []byte(strings.Replace(string(result2), "\"hello\"", string(res), 1))
		rest.PostProcessResponse(w, cliCtx, result2)
	}
}
