package rest

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	govRest "github.com/okex/exchain/x/gov/client/rest"
)

// ProposeValidatorProposalRESTHandler defines propose validator proposal handler
func ProposeValidatorProposalRESTHandler(context.CLIContext) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{}
}
