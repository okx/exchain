package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	commonCodespace sdk.CodespaceType = "common"

	codeMissingRequiredParam CodeType = 60001
	codeInvalidRequestParam  CodeType = 60002
)

// ErrMissingRequiredParam returns an error when the required param is missing
func ErrMissingRequiredParam() sdk.Error {
	return sdk.NewError(commonCodespace, codeMissingRequiredParam, "failed. the required param is missing")
}

// ErrInvalidRequestParam returns an error when the the request param is invalid
func ErrInvalidRequestParam(msg string) sdk.Error {
	return sdk.NewError(commonCodespace, codeInvalidRequestParam, "failed. invalid request param: %s", msg)
}
