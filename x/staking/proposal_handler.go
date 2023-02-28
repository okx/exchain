package staking

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/staking/types"

	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewProposalHandler handles "gov" type message in "staking"
func NewProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch content := proposal.Content.(type) {
		case types.ProposeValidatorProposal:
			return handleProposeValidatorProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(types.DefaultCodespace, content.ProposalType())
		}
	}
}

func handleProposeValidatorProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	// check
	proposeValidatorProposal, ok := proposal.Content.(types.ProposeValidatorProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}
	validator := proposeValidatorProposal.Validator
	if proposeValidatorProposal.IsAdd {
		createValMsg := NewMsgCreateValidator(validator.ValidatorAddress, validator.PubKey,
			validator.Description, validator.MinSelfDelegation)
		_, err := handleMsgCreateValidator(ctx, createValMsg, *k)
		return err
	} else {
		delValMsg := NewMsgDestroyValidator(sdk.AccAddress(validator.ValidatorAddress))
		_, err := handleMsgDestroyValidator(ctx, delValMsg, *k)
		return err
	}
	return nil
}
