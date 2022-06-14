package client

import (
	"github.com/okex/exchain/x/distribution/client/cli"
	"github.com/okex/exchain/x/distribution/client/rest"
	govclient "github.com/okex/exchain/x/gov/client"
)

// param change proposal handler
var (
	CommunityPoolSpendProposalHandler      = govclient.NewProposalHandler(cli.GetCmdCommunityPoolSpendProposal, rest.ProposalRESTHandler)
	ChangeDistributionModelProposalHandler = govclient.NewProposalHandler(cli.GetChangeDistributionModelProposal, rest.ProposalRESTHandler) //TODO zhujianguo
)
