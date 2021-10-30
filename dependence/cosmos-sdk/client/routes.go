package client

import (
	"github.com/gorilla/mux"

	"github.com/okex/exchain/dependence/cosmos-sdk/client/context"
	"github.com/okex/exchain/dependence/cosmos-sdk/client/rpc"
)

// Register routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	rpc.RegisterRPCRoutes(cliCtx, r)
}
