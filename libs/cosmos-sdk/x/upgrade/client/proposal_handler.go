package client

import (
	govclient "github.com/okx/okbchain/libs/cosmos-sdk/x/gov/client"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/upgrade/client/cli"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/upgrade/client/rest"
)

var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitUpgradeProposal, rest.ProposalRESTHandler)
