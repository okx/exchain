package client

import (
	"github.com/okex/exchain/x/erc20/client/cli"
	"github.com/okex/exchain/x/erc20/client/rest"
	govcli "github.com/okex/exchain/x/gov/client"
)

var (
	// TokenMappingProposalHandler alias gov NewProposalHandler
	TokenMappingProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdTokenMappingProposal,
		rest.TokenMappingProposalRESTHandler,
	)
)
