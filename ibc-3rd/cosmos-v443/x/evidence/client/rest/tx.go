package rest

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/client"

	"github.com/gorilla/mux"
)

func registerTxRoutes(clientCtx client.Context, r *mux.Router, handlers []EvidenceRESTHandler) {
	// TODO: Register tx handlers.
}
