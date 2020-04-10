package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	CodeInvalidMaxProposalNum sdk.CodeType = 4
)

// ErrInvalidMaxProposalNum returns error when the number of params to change are out of limit
func ErrInvalidMaxProposalNum(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMaxProposalNum, msg)
}
