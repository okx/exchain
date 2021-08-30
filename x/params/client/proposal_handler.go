package client

import (
	govclient "github.com/okex/exchain/x/gov/client"
	"github.com/okex/exchain/x/params/client/cli"
	"github.com/okex/exchain/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler in cmsdk
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
