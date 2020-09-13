package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/okex/okchain/x/ammswap/types"
	"github.com/okex/okchain/x/common"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/ammswap/exchange", swapExchangeHandler(cliCtx)).Methods("GET")
}

func swapExchangeHandler(cliCtx context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tokenName := vars["token"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/ammswap/swapTokenPair/%s", tokenName), nil)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, err.Error())
			return
		}

		exchange := types.SwapTokenPair{}
		codec.Cdc.MustUnmarshalJSON(res, exchange)
		response := common.GetBaseResponse(exchange)
		resBytes, err := json.Marshal(response)
		if err != nil {
			common.HandleErrorMsg(w, cliCtx, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, resBytes)
	}
}
