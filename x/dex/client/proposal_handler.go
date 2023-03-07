package client

import (
	"github.com/okx/exchain/x/dex/client/cli"
	"github.com/okx/exchain/x/dex/client/rest"
	govclient "github.com/okx/exchain/x/gov/client"
)

// param change proposal handler
var (
	// DelistProposalHandler alias gov NewProposalHandler
	DelistProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitDelistProposal, rest.DelistProposalRESTHandler)
)
