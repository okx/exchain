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

	if sdkErr := k.CheckMsgManageContractDeploymentWhitelistProposal(ctx, manageContractDeploymentWhitelistProposal); sdkErr != nil {
		return sdkErr
	}

	if manageContractDeploymentWhitelistProposal.IsAdded {
		// add deployer address into whitelist
		k.SetContractDeploymentWhitelistMember(ctx, manageContractDeploymentWhitelistProposal.DeployerAddr)
		return nil
	}

	// remove deployer address from whitelist
	k.DeleteContractDeploymentWhitelistMember(ctx, manageContractDeploymentWhitelistProposal.DeployerAddr)
	return nil
}
