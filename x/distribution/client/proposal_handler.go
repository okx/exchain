package client

import (
	"github.com/okex/exchain/x/distribution/client/cli"
	"github.com/okex/exchain/x/distribution/client/rest"
	govclient "github.com/okex/exchain/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
