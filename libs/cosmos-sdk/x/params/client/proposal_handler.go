package client

import (
	govclient "github.com/okx/okbchain/libs/cosmos-sdk/x/gov/client"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/params/client/cli"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/params/client/rest"
)

// ProposalHandler handles param change proposals
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
