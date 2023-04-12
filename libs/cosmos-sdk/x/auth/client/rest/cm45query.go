package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	genutilrest "github.com/okex/exchain/libs/cosmos-sdk/x/genutil/client/rest"
)

func CM45QueryTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(
				w, http.StatusBadRequest,
				fmt.Sprintf("failed to parse query parameters: %s", err),
			)
			return
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponseBare(w, cliCtx, txs)
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				genutilrest.QueryGenesisTxs(cliCtx, w)
				return
			}
		}

		pr, err := rest.ParseCM45PageRequest(r)
		if err != nil {
			rest.WriteErrorResponse(
				w, http.StatusBadRequest,
				fmt.Sprintf("failed to parse page request: %s", err),
			)
			return
		}
		page, limit, err = query.ParsePagination(pr)

		// parse Orderby is not supported for now

		events, err = rest.ParseEvents(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		for _, event := range events {
			if !strings.Contains(event, "=") || strings.Count(event, "=") > 1 {
				rest.WriteErrorResponse(w, http.StatusBadRequest,
					fmt.Sprintf("invalid event; event %s should be of the format: %s", event, "{eventType}.{eventAttribute}={value}"))
				return
			}
		}

		searchResult, err := utils.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		txsList := make([]sdk.Tx, len(searchResult.Txs))
		for i, tx := range searchResult.Txs {
			ourTx := tx.Tx
			if !ok {
				return
			}

			txsList[i] = ourTx
		}
		res := types.CustomGetTxsEventResponse{
			Txs:         txsList,
			TxResponses: searchResult.Txs,
			Pagination: &query.PageResponse{
				Total: uint64(searchResult.TotalCount),
			},
		}
		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}
