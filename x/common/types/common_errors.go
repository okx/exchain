package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

// const
const (
	CommonCodespace sdk.CodespaceType = "common"

	codeMissingRequiredParam CodeType = 60001
	codeInvalidRequestParam  CodeType = 60002
	codeInternalServer       CodeType = 60003
	codeDataNotExists        CodeType = 60004
	codeInvalidAddress       CodeType = 60005
	codeUnknownQueryEndpoint CodeType = 60006
	codeUnknownMsgType       CodeType = 60007
	codeBadJSONMarshaling    CodeType = 60008
)

// ErrMissingRequiredParam returns an error when the required param is missing
func ErrMissingRequiredParam(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeMissingRequiredParam, "failed. the required param is missing")
}

// ErrInvalidRequestParam returns an error when the the request param is invalid
func ErrInvalidRequestParam(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, codeInvalidRequestParam, "failed. invalid request param: %s", msg)
}

// ErrInternalServer returns an error when error occurs in internal server
func ErrInternalServer(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeInternalServer, "failed. error occurs in internal server")
}

// ErrDataNotExists returns an error when the target data for request doesn't exist
func ErrDataNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeDataNotExists, "failed. target data for request doesn't exist")
}

// ErrInvalidAddress returns an error when the address is invalid
func ErrInvalidAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeInvalidAddress, "failed. invalid address")
}

// ErrUnknownQueryEndpoint returns an error when the the query endpoint is unknown
func ErrUnknownQueryEndpoint(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeUnknownQueryEndpoint, "failed. unknown query endpoint")
}

// ErrUnknownMsgType returns an error when the the msg type is unknown
func ErrUnknownMsgType(codespace sdk.CodespaceType, msgType string) sdk.Error {
	return sdk.NewError(codespace, codeUnknownMsgType, "failed. unknown msg type: %s", msgType)
}

// ErrBadJSONMarshaling returns an error with the bad encoding of JSON
func ErrBadJSONMarshaling(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, codeBadJSONMarshaling, "failed. marshal JSON unsuccessfully: %s", msg)
}