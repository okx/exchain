package client

import (
	"github.com/okex/okchain/x/dex/client/cli"
	"github.com/okex/okchain/x/dex/client/rest"
	govclient "github.com/okex/okchain/x/gov/client"
)

// param change proposal handler
var (
	DelistProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitDelistProposal, rest.DelistProposalRESTHandler)
)
