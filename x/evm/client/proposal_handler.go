package client

import (
	"github.com/okex/exchain/x/evm/client/cli"
	"github.com/okex/exchain/x/evm/client/rest"
	govcli "github.com/okex/exchain/x/gov/client"
)

var (
	// ManageContractDeploymentWhitelistProposalHandler alias gov NewProposalHandler
	ManageContractDeploymentWhitelistProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageContractDeploymentWhitelistProposal,
		rest.ManageContractDeploymentWhitelistProposalRESTHandler,
	)

	// ManageContractBlockedListProposalHandler alias gov NewProposalHandler
	ManageContractBlockedListProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageContractBlockedListProposal,
		rest.ManageContractBlockedListProposalRESTHandler,
	)
	ManageContractMethodBlockedListProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageContractMethodBlockedListProposal,
		rest.ManageContractMethodBlockedListProposalRESTHandler,
	)
)
