package client

import (
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/distribution/client/cli"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/distribution/client/rest"
	govclient "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/gov/client"
)

// ProposalHandler is the community spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
