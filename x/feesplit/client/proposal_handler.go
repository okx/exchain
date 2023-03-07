package client

import (
	"github.com/okx/okbchain/x/feesplit/client/cli"
	"github.com/okx/okbchain/x/feesplit/client/rest"
	govcli "github.com/okx/okbchain/x/gov/client"
)

var (
	// FeeSplitSharesProposalHandler alias gov NewProposalHandler
	FeeSplitSharesProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdFeeSplitSharesProposal,
		rest.FeeSplitSharesProposalRESTHandler,
	)
)
