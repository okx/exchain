package rest

import (
	"github.com/gorilla/mux"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	govRest "github.com/okex/exchain/x/gov/client/rest"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {

}

// TokenMappingProposalRESTHandler defines erc20 proposal handler
func TokenMappingProposalRESTHandler(context.CLIContext) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{}
}
