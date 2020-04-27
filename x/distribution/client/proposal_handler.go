package client

import (
	"github.com/okex/okchain/x/distribution/client/cli"
	"github.com/okex/okchain/x/distribution/client/rest"
	govclient "github.com/okex/okchain/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
