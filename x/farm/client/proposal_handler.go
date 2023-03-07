package client

import (
	"github.com/okx/exchain/x/farm/client/cli"
	"github.com/okx/exchain/x/farm/client/rest"
	govcli "github.com/okx/exchain/x/gov/client"
)

var (
	// ManageWhiteListProposalHandler alias gov NewProposalHandler
	ManageWhiteListProposalHandler = govcli.NewProposalHandler(cli.GetCmdManageWhiteListProposal, rest.ManageWhiteListProposalRESTHandler)
)
