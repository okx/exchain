package rest

import (
	"net/http"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
)

func queryDenomTraces(ctx context.CLIContext, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params types.QueryDenomTracesRequest
		pr, err := rest.ParseCM45PageRequest(r)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		params = types.QueryDenomTracesRequest{
			Pagination: pr,
		}
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}

		bz, err := clientCtx.CodecProy.GetCdc().MarshalJSON(params)
		res, height, err := clientCtx.QueryWithData(endpoint, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
