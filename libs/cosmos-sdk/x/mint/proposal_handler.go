package mint

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	"github.com/okex/exchain/x/common"

	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewManageContractDeploymentWhitelistProposalHandler handles "gov" type message in "evm"
func NewManageTreasuresProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.ManageTreasuresProposal:
			return handleManageContractDeploymentWhitelistProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleManageContractDeploymentWhitelistProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	// check
	manageTreasuresProposal, ok := proposal.Content.(types.ManageTreasuresProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}

	if manageTreasuresProposal.IsAdded {
		// add deployer addresses into whitelist
		if err := k.UpdateTreasures(ctx, manageTreasuresProposal.Treasures); err != nil {
			return types.ErrTreasuresInternal(err)
		}
		return nil
	}

	// remove deployer addresses from whitelist
	if err := k.DeleteTreasures(ctx, manageTreasuresProposal.Treasures); err != nil {
		return types.ErrTreasuresInternal(err)
	}
	return nil
}
