package client

import (
	govclient "github.com/okex/exchain/libs/cosmos-sdk/x/gov/client"
	"github.com/okex/exchain/libs/cosmos-sdk/x/upgrade/client/cli"
	"github.com/okex/exchain/libs/cosmos-sdk/x/upgrade/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
