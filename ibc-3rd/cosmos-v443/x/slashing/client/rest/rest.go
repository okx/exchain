package rest

import (
	"github.com/gorilla/mux"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/client"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/client/rest"
)

func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := rest.WithHTTPDeprecationHeaders(rtr)

	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}
