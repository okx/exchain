package client

import (
	govclient "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/gov/client"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/upgrade/client/cli"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/upgrade/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
var CancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal, rest.ProposalCancelRESTHandler)
