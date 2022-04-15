package rest

import (
	"github.com/gorilla/mux"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

func RegisterOriginRPCRoutersForGRPC(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(
		"/cosmos/auth/v1beta1/accounts/{address}",
		QueryAccountRequestHandlerFn(storeName, cliCtx),
	).Methods("GET")
}
