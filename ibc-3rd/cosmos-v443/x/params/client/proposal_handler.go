package client

import (
	govclient "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/gov/client"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/params/client/cli"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, rest.ProposalRESTHandler)
