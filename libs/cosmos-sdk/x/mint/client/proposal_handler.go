package client

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint/client/cli"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint/client/rest"
	govcli "github.com/okx/okbchain/x/gov/client"
)

var (
	ManageTreasuresProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageTreasuresProposal,
		rest.ManageTreasuresProposalRESTHandler,
	)

	ExtraProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdExtraProposal,
		rest.ExtraProposalRESTHandler,
	)
)
