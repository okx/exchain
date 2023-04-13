package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
)

func CM45QueryBalancesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryBalanceParams(addr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/bank/balances", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var coins sdk.Coins
		cliCtx.Codec.MustUnmarshalJSON(res, &coins)
		wrappedCoins := types.NewWrappedBalances(coins)
		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, wrappedCoins)
	}
}
