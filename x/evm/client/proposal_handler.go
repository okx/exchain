package client

import (
	"github.com/okex/okexchain/x/evm/client/cli"
	"github.com/okex/okexchain/x/evm/client/rest"
	govcli "github.com/okex/okexchain/x/gov/client"
)

var (
	// ManageContractDeploymentWhitelistProposalHandler alias gov NewProposalHandler
	ManageContractDeploymentWhitelistProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageContractDeploymentWhitelistProposal,
		rest.ManageContractDeploymentWhitelistProposalRESTHandler,
	)
)
