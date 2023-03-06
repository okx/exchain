package rest

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	govRest "github.com/okx/okbchain/x/gov/client/rest"
)

// ProposeValidatorProposalRESTHandler defines propose validator proposal handler
func ProposeValidatorProposalRESTHandler(context.CLIContext) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{}
}
