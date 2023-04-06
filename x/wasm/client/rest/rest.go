package rest

import (
	"github.com/gorilla/mux"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

// RegisterRoutes registers wasm-related REST handlers to a router
func RegisterRoutes(cliCtx clientCtx.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
	registerNewTxRoutes(cliCtx, r)
}

// RegisterRoutesV2 registers wasm-related v2 REST handlers to a router
func RegisterRoutesV2(cliCtx clientCtx.CLIContext, r *mux.Router) {
	registerQueryRoutesV2(cliCtx, r)
}
