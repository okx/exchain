package mint

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/common"
	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewManageTreasuresProposalHandler handles "gov" type message in "mint"
func NewManageTreasuresProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.ManageTreasuresProposal:
			return handleManageTreasuresProposal(ctx, k, proposal)
		case types.ModifyNextBlockUpdateProposal:
			return handleModifyNextBlockUpdateProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleManageTreasuresProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	// check
	manageTreasuresProposal, ok := proposal.Content.(types.ManageTreasuresProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}

	if manageTreasuresProposal.IsAdded {
		// add/update treasures into state
		if err := k.UpdateTreasures(ctx, manageTreasuresProposal.Treasures); err != nil {
			return types.ErrTreasuresInternal(err)
		}
		return nil
	}

	// delete treasures into state
	if err := k.DeleteTreasures(ctx, manageTreasuresProposal.Treasures); err != nil {
		return types.ErrTreasuresInternal(err)
	}
	return nil
}

func handleModifyNextBlockUpdateProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	modifyProposal, ok := proposal.Content.(types.ModifyNextBlockUpdateProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}
	if modifyProposal.BlockNum <= uint64(global.GetGlobalHeight())+1 {
		return types.ErrNextBlockUpdateTooLate
	}

	minter := k.GetMinterCustom(ctx)
	minter.NextBlockToUpdate = modifyProposal.BlockNum
	k.SetMinterCustom(ctx, minter)
	return nil
}
