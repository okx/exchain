package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	commonCodespace sdk.CodespaceType = "common"

	codeMissingRequiredParam CodeType = 60001
	codeInvalidRequestParam  CodeType = 60002
	codeInternalServer       CodeType = 60003
	codeDataNotExsits        CodeType = 60004
	codeInvalidAddress       CodeType = 60005
	codeUnknownQueryEndpoint CodeType = 60006
	codeUnknownMsgType       CodeType = 60007
)

// ErrMissingRequiredParam returns an error when the required param is missing
func ErrMissingRequiredParam() sdk.Error {
	return sdk.NewError(commonCodespace, codeMissingRequiredParam, "failed. the required param is missing")
}

// ErrInvalidRequestParam returns an error when the the request param is invalid
func ErrInvalidRequestParam(msg string) sdk.Error {
	return sdk.NewError(commonCodespace, codeInvalidRequestParam, "failed. invalid request param: %s", msg)
}

// ErrInternalServer returns an error when error occurs in internal server
func ErrInternalServer() sdk.Error {
	return sdk.NewError(commonCodespace, codeInternalServer, "failed. error occurs in internal server")
}

// ErrDataNotExists returns an error when the target data for request doesn't exist
func ErrDataNotExists() sdk.Error {
	return sdk.NewError(commonCodespace, codeDataNotExsits, "failed. target data for request doesn't exist")
}

// ErrInvalidAddress returns an error when the address is invalid
func ErrInvalidAddress() sdk.Error {
	return sdk.NewError(commonCodespace, codeInvalidAddress, "failed. invalid address")
}

// ErrUnknownQueryEndpoint returns an error when the the query endpoint is unknown
func ErrUnknownQueryEndpoint(msg string) sdk.Error {
	return sdk.NewError(commonCodespace, codeUnknownQueryEndpoint, "failed. unknown query endpoint: %s", msg)
}

// ErrUnknownMsgType returns an error when the the msg type is unknown
func ErrUnknownMsgType(msg string) sdk.Error {
	return sdk.NewError(commonCodespace, codeUnknownMsgType, "failed. unknown msg type: %s", msg)
}
