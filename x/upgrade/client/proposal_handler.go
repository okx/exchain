package client

import (
	govclient "github.com/okex/okexchain/x/gov/client"
	"github.com/okex/okexchain/x/upgrade/client/cli"
	"github.com/okex/okexchain/x/upgrade/client/rest"
)

// ProposalHandler is param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
