package evm

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/evm/types"
	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewManageContractDeploymentWhitelistProposalHandler handles "gov" type message in "evm"
func NewManageContractDeploymentWhitelistProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.ManageContractDeploymentWhitelistProposal:
			return handleManageContractDeploymentWhitelistProposal(ctx, k, proposal)
		case types.ManageContractBlockedListProposal:
			return handleManageContractBlockedlListProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleManageContractDeploymentWhitelistProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	// check
	manageContractDeploymentWhitelistProposal, ok := proposal.Content.(types.ManageContractDeploymentWhitelistProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}

	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	if manageContractDeploymentWhitelistProposal.IsAdded {
		// add deployer addresses into whitelist
		csdb.SetContractDeploymentWhitelist(manageContractDeploymentWhitelistProposal.DistributorAddrs)
		return nil
	}

	// remove deployer addresses from whitelist
	csdb.DeleteContractDeploymentWhitelist(manageContractDeploymentWhitelistProposal.DistributorAddrs)
	return nil
}

func handleManageContractBlockedlListProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	// check
	manageContractBlockedListProposal, ok := proposal.Content.(types.ManageContractBlockedListProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}

	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	csdb.SetCache(k.ConfigCache)
	if manageContractBlockedListProposal.IsAdded {
		// add contract addresses into blocked list
		csdb.SetContractBlockedList(manageContractBlockedListProposal.ContractAddrs)
		return nil
	}

	// remove contract addresses from blocked list
	csdb.DeleteContractBlockedList(manageContractBlockedListProposal.ContractAddrs)
	return nil
}
