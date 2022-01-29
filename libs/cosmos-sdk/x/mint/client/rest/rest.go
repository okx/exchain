package rest

import (
	"github.com/gorilla/mux"
	govRest "github.com/okex/exchain/x/gov/client/rest"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}

// ManageContractMethodBlockedListProposalRESTHandler defines evm proposal handler
func ManageTreasuresProposalRESTHandler(context.CLIContext) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{}
}
