package client

import (
	"github.com/okex/okexchain/x/distribution/client/cli"
	"github.com/okex/okexchain/x/distribution/client/rest"
	govclient "github.com/okex/okexchain/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
