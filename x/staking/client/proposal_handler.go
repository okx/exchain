package client

import (
	"github.com/okx/okbchain/x/staking/client/cli"
	"github.com/okx/okbchain/x/staking/client/rest"
	govcli "github.com/okx/okbchain/x/gov/client"
)

var (
	ProposeValidatorProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdProposeValidatorProposal,
		rest.ProposeValidatorProposalRESTHandler,
	)
)

