package client

import (
	"github.com/okex/exchain/x/staking/client/cli"
	"github.com/okex/exchain/x/staking/client/rest"
	govcli "github.com/okex/exchain/x/gov/client"
)

var (
	ProposeValidatorProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdProposeValidatorProposal,
		rest.ProposeValidatorProposalRESTHandler,
	)
)

