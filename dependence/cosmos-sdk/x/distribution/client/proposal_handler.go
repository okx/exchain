package client

import (
	"github.com/okex/exchain/dependence/cosmos-sdk/x/distribution/client/cli"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/distribution/client/rest"
	govclient "github.com/okex/exchain/dependence/cosmos-sdk/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
