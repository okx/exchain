package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/evm/types"
	govTypes "github.com/okex/okexchain/x/gov/types"
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

	//if sdkErr := k.CheckMsgManageContractBlockedListProposal(ctx, manageContractBlockedListProposal); sdkErr != nil {
	//	return sdkErr
	//}

	// TODO
	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	if manageContractBlockedListProposal.IsAdded {
		// add contract address into blocked list
		csdb.SetContractBlockedListMember(manageContractBlockedListProposal.ContractAddrs[0])
		return nil
	}

	// remove contract address from blocked list
	csdb.DeleteContractBlockedListMember(manageContractBlockedListProposal.ContractAddrs[0])
	return nil
}
