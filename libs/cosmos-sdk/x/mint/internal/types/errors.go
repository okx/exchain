package types

import (
	"fmt"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
)

// NOTE: We can't use 1 since that error code is reserved for internal errors.
const (
	DefaultCodespace string = ModuleName
)

var (
	// ErrInvalidState returns an error resulting from an invalid Storage State.
	ErrEmptyTreasures = sdkerrors.Register(ModuleName, 2, "treasures is not empty")

	ErrDuplicatedTreasure         = sdkerrors.Register(ModuleName, 3, "treasures can not be duplicate")
	ErrUnexpectedProposalType     = sdkerrors.Register(ModuleName, 4, "unsupported proposal type of mint module")
	ErrProposerMustBeValidator    = sdkerrors.Register(ModuleName, 5, "the proposal of proposer must be validator")
	ErrNextBlockUpdateTooLate     = sdkerrors.Register(ModuleName, 7, "the next block to update is too late")
	ErrCodeInvalidHeight          = sdkerrors.Register(ModuleName, 8, "height must be greater than current block")
	ErrHandleExtraProposal        = sdkerrors.Register(ModuleName, 9, "handle extra proposal error")
	ErrUnknownExtraProposalAction = sdkerrors.Register(ModuleName, 10, "extra proposal's action unknown")
)

// ErrTreasuresInternal returns an error when the length of address list in the proposal is larger than the max limitation
func ErrTreasuresInternal(err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			DefaultParamspace,
			11,
			fmt.Sprintf("treasures error:%s", err.Error()))}
}

func ErrExtraProposalParams(desc string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, 6, fmt.Sprintf("mint extra proposal error:%s", desc))
}
