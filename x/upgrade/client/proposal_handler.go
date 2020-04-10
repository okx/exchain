package client

import (
	govclient "github.com/okex/okchain/x/gov/client"
	"github.com/okex/okchain/x/upgrade/client/cli"
	"github.com/okex/okchain/x/upgrade/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
