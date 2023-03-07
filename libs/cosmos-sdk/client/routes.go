package client

import (
	"github.com/gorilla/mux"

	"github.com/okx/exchain/libs/cosmos-sdk/client/context"
	"github.com/okx/exchain/libs/cosmos-sdk/client/rpc"
)

// Register routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	rpc.RegisterRPCRoutes(cliCtx, r)
}
