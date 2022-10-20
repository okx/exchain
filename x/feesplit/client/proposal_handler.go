package client

import (
	"github.com/okex/exchain/x/feesplit/client/cli"
	"github.com/okex/exchain/x/feesplit/client/rest"
	govcli "github.com/okex/exchain/x/gov/client"
)

var (
	// FeeSplitSharesProposalHandler alias gov NewProposalHandler
	FeeSplitSharesProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdFeeSplitSharesProposal,
		rest.FeeSplitSharesProposalRESTHandler,
	)
)
