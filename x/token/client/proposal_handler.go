package client

import (
	govclient "github.com/okex/okchain/x/gov/client"
	"github.com/okex/okchain/x/token/client/cli"
	"github.com/okex/okchain/x/token/client/rest"
)

// ProposalHandler is the param change proposal handler in cmsdk
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
