package rest

import (
	"github.com/gorilla/mux"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx clientCtx.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
	registerNewTxRoutes(cliCtx, r)
}
