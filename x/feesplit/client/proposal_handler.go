package client

import (
	"github.com/okx/exchain/x/feesplit/client/cli"
	"github.com/okx/exchain/x/feesplit/client/rest"
	govcli "github.com/okx/exchain/x/gov/client"
)

var (
	// FeeSplitSharesProposalHandler alias gov NewProposalHandler
	FeeSplitSharesProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdFeeSplitSharesProposal,
		rest.FeeSplitSharesProposalRESTHandler,
	)
)
