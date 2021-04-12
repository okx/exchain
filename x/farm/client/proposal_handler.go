package client

import (
	"github.com/okex/exchain/x/farm/client/cli"
	"github.com/okex/exchain/x/farm/client/rest"
	govcli "github.com/okex/exchain/x/gov/client"
)

var (
	// ManageWhiteListProposalHandler alias gov NewProposalHandler
	ManageWhiteListProposalHandler = govcli.NewProposalHandler(cli.GetCmdManageWhiteListProposal, rest.ManageWhiteListProposalRESTHandler)
)
