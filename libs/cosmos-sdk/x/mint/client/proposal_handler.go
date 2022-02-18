package client

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/client/cli"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/client/rest"
	govcli "github.com/okex/exchain/x/gov/client"
)

var (
	ManageTreasuresProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageTreasuresProposal,
		rest.ManageTreasuresProposalRESTHandler,
	)
)
