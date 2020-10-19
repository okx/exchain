package client

import (
	"github.com/okex/okexchain/x/farm/client/cli"
	"github.com/okex/okexchain/x/farm/client/rest"
	govcli "github.com/okex/okexchain/x/gov/client"
)

var (
	// ManageWhiteListProposalHandler alias gov NewProposalHandler
	ManageWhiteListProposalHandler = govcli.NewProposalHandler(cli.GetCmdManageWhiteListProposal, rest.ManageWhiteListProposalRESTHandler)
)
