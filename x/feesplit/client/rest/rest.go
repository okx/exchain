package rest

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	govRest "github.com/okex/exchain/x/gov/client/rest"
)

// FeeSplitSharesProposalRESTHandler defines feesplit proposal handler
func FeeSplitSharesProposalRESTHandler(context.CLIContext) govRest.ProposalRESTHandler {
	return govRest.ProposalRESTHandler{}
}
