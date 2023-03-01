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
	validatorProposal, ok := proposal.Content.(types.ProposeValidatorProposal)
	if !ok {
		return types.ErrUnexpectedProposalType
	}

	validator := validatorProposal.Validator
	var valKey [sdk.AddrLen]byte
	copy(valKey[:], validator.ValidatorAddress[:])
	// verify proposed validator with validator set
	valSetCount, err := verifyProposalWithValSet(ctx, k, valKey, validatorProposal.IsAdd)
	if err != nil {
		return err
	}

	// proposed validators
	proposedValidators, found := k.GetProposeValidators(ctx)
	if found {
		proposedValidators[valKey] = validatorProposal.IsAdd
	} else {
		proposedValidators = map[[sdk.AddrLen]byte]bool{
			valKey: validatorProposal.IsAdd,
		}
	}

	// verify validator count
	if err := verifyValidatorCount(ctx, k, valSetCount, proposedValidators); err != nil {
		return err
	}

	if _, found := k.GetValidator(ctx, validator.ValidatorAddress); !found {
		// create validator
		createValMsg := NewMsgCreateValidator(validator.ValidatorAddress, validator.PubKey,
			validator.Description, validator.MinSelfDelegation)
		_, err := handleMsgCreateValidator(ctx, createValMsg, *k)
		return err
	}

	k.SetProposeValidators(ctx, proposedValidators)
	return nil
}

func verifyProposalWithValSet(ctx sdk.Context, k *Keeper, valKey [sdk.AddrLen]byte, isAdd bool) (int, sdk.Error) {
	lastValSet := k.GetLastValidatorsByAddr(ctx)
	_, inValSet := lastValSet[valKey] // exist in validator set
	if isAdd && inValSet {
		return 0, types.ErrProposedInValSet
	}
	if !isAdd && !inValSet {
		return 0, types.ErrProposedNotInValSet
	}
	return len(lastValSet), nil
}

func verifyValidatorCount(ctx sdk.Context, k *Keeper, valSetCount int, proposedValidators map[[sdk.AddrLen]byte]bool) sdk.Error {
	maxCount := k.GetParams(ctx).MaxValidators
	proposeCount := 0
	for _, isAdd := range proposedValidators {
		if isAdd {
			proposeCount++
		} else {
			proposeCount--
		}
	}
	if valSetCount+proposeCount > int(maxCount) {
		return types.ErrProposedExceedMax
	}
	return nil
}
