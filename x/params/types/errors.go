package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	CodeInvalidMaxProposalNum sdk.CodeType = 4
	CodeInvalidRate           sdk.CodeType = 5
	CodeErrCodec              sdk.CodeType = 6
)

// ErrInvalidMaxProposalNum returns an error when the number of params to change are out of limit
func ErrInvalidMaxProposalNum(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMaxProposalNum, msg)
}

// ErrInvalidRateNum returns an error when the specific number of rate is invalid
func ErrInvalidRate(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRate,
		"failed. params rate is only allowed within the scope of [0,1]")
}

// ErrUnmarshalJSON returns an error when it fails to unmarshal JSON for params value
func ErrUnmarshalJSON(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeErrCodec,
		"failed. unmarshal JSON error for params value")
}
