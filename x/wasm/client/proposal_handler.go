package client

import (
	govclient "github.com/okex/exchain/x/gov/client"
	"github.com/okex/exchain/x/wasm/client/cli"
	"github.com/okex/exchain/x/wasm/client/rest"
)

// ProposalHandlers define the wasm cli proposal types and rest handler.
var ProposalHandlers = []govclient.ProposalHandler{
	//govclient.NewProposalHandler(cli.ProposalStoreCodeCmd, rest.StoreCodeProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalInstantiateContractCmd, rest.InstantiateProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalMigrateContractCmd, rest.MigrateProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalExecuteContractCmd, rest.ExecuteProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalSudoContractCmd, rest.SudoProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalUpdateContractAdminCmd, rest.UpdateContractAdminProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalClearContractAdminCmd, rest.ClearContractAdminProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalPinCodesCmd, rest.PinCodeProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalUnpinCodesCmd, rest.UnpinCodeProposalHandler),
	//govclient.NewProposalHandler(cli.ProposalUpdateInstantiateConfigCmd, rest.UpdateInstantiateConfigProposalHandler),
}

// UpdateContractAdminProposalHandler is a proposal handler which can update admin of a contract.
var UpdateContractAdminProposalHandler = govclient.NewProposalHandler(cli.ProposalUpdateContractAdminCmd, rest.UpdateContractAdminProposalHandler)

// ClearContractAdminProposalHandler is a proposal handler which can clear admin of a contract.
var ClearContractAdminProposalHandler = govclient.NewProposalHandler(cli.ProposalClearContractAdminCmd, rest.ClearContractAdminProposalHandler)

// MigrateContractProposalHandler is a proposal handler which can migrate contract to disable some methods of the contract.
var MigrateContractProposalHandler = govclient.NewProposalHandler(cli.ProposalMigrateContractCmd, rest.MigrateProposalHandler)

// UpdateDeploymentWhitelistProposalHandler is a custom proposal handler which defines whitelist to deploy contracts.
var UpdateDeploymentWhitelistProposalHandler = govclient.NewProposalHandler(cli.ProposalUpdateDeploymentWhitelistCmd, rest.EmptyProposalRestHandler)

// UpdateWASMContractMethodBlockedListProposalHandler is a custom proposal handler which defines methods blacklist of a contract.
var UpdateWASMContractMethodBlockedListProposalHandler = govclient.NewProposalHandler(cli.ProposalUpdateWASMContractMethodBlockedListCmd, rest.EmptyProposalRestHandler)
