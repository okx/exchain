package client

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/x/distribution/client/cli"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/distribution/client/rest"
	govclient "github.com/okx/okbchain/libs/cosmos-sdk/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
