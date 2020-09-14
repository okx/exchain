package client

import (
	"github.com/okex/okexchain/x/dex/client/cli"
	"github.com/okex/okexchain/x/dex/client/rest"
	govclient "github.com/okex/okexchain/x/gov/client"
)

// param change proposal handler
var (
	// DelistProposalHandler alias gov NewProposalHandler
	DelistProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitDelistProposal, rest.DelistProposalRESTHandler)
)
