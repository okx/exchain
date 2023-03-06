package client

import (
	"github.com/okx/okbchain/x/evm/client/cli"
	"github.com/okx/okbchain/x/evm/client/rest"
	govcli "github.com/okx/okbchain/x/gov/client"
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
	ManageSysContractAddressProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageSysContractAddressProposal,
		rest.ManageSysContractAddressProposalRESTHandler,
	)
	ManageContractMethodGuFactorProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageContractMethodGuFactorProposal,
		rest.ManageContractMethodBlockedListProposalRESTHandler,
	)

	ManageContractByteCodeProposalHandler = govcli.NewProposalHandler(
		cli.GetCmdManageContractByteCodeProposal,
		rest.ManageContractBytecodeProposalRESTHandler)
)
