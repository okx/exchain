package client

import (
	"github.com/okx/exchain/libs/cosmos-sdk/x/mint/client/cli"
	"github.com/okx/exchain/libs/cosmos-sdk/x/mint/client/rest"
	govcli "github.com/okx/exchain/x/gov/client"
)

var (
	ManageTreasuresProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageTreasuresProposal,
		rest.ManageTreasuresProposalRESTHandler,
	)
	ModifyNextBlockUpdateProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdModifyNextBlockUpdateProposal,
		rest.ModifyNextBlockUpdateProposalRESTHandler,
	)
)
