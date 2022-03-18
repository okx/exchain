package rest

import (
	"fmt"
	"net/http"

	"github.com/okex/exchain/x/token/types"

	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
)

// RegisterRoutesV2, a central function to define routes
// which is called by the rest module in main application
func RegisterRoutesV2(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/tokens/{currency}"), tokenHandlerV2(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/tokens"), tokensHandlerV2(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/accounts/{address}"), accountsHandlerV2(cliCtx, storeName)).Methods("GET")
}

func tokenHandlerV2(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenName := vars["currency"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", storeName, types.QueryTokenV2, tokenName), nil)
		common.HandleResponseV2(w, res, err)
	}
}

func tokensHandlerV2(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", storeName, types.QueryTokensV2), nil)
		common.HandleResponseV2(w, res, err)
	}
}

func accountsHandlerV2(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		address := vars["address"]

		currency := r.URL.Query().Get("currency")
		hideZero := r.URL.Query().Get("hide_zero")

		if hideZero == "" {
			hideZero = "yes"
		}
		if hideZero != "yes" && hideZero != "no" {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidParam)
			return
		}

		// valid address
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			common.HandleErrorResponseV2(w, http.StatusBadRequest, common.ErrorInvalidAddress)
			return
		}

		accountParam := types.AccountParamV2{
			Currency: currency,
			HideZero: hideZero,
		}

		req, err := cliCtx.Codec.MarshalJSON(accountParam)

		if err != nil {
			common.HandleErrorMsg(w, cliCtx, common.CodeMarshalJSONFailed, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", storeName, types.QueryAccountV2, address), req)
		common.HandleResponseV2(w, res, err)
	}
}
