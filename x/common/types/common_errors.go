package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

// const
const (
	CommonCodespace sdk.CodespaceType = "common"

	CodeMissingRequiredParam CodeType = 60001
	CodeInvalidRequestParam  CodeType = 60002
	CodeInternalServer       CodeType = 60003
	CodeDataNotExists        CodeType = 60004
	CodeInvalidAddress       CodeType = 60005
	CodeUnknownQueryEndpoint CodeType = 60006
	CodeUnknownMsgType       CodeType = 60007
	CodeBadJSONMarshaling    CodeType = 60008
	CodeBadJSONUnmarshaling  CodeType = 60009
)

// ErrMissingRequiredParam returns an error when the required param is missing
func ErrMissingRequiredParam(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMissingRequiredParam, "failed. the required param is missing")
}

// ErrInvalidRequestParam returns an error when the the request param is invalid
func ErrInvalidRequestParam(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequestParam, "failed. invalid request param: %s", msg)
}

// ErrInternalServer returns an error when error occurs in internal server
func ErrInternalServer(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInternalServer, "failed. error occurs in internal server")
}

// ErrDataNotExists returns an error when the target data for request doesn't exist
func ErrDataNotExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDataNotExists, "failed. target data for request doesn't exist")
}

// ErrInvalidAddress returns an error when the address is invalid
func ErrInvalidAddress(codespace sdk.CodespaceType, kind string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, "failed. invalid %s address", kind)
}

// ErrUnknownQueryEndpoint returns an error when the the query endpoint is unknown
func ErrUnknownQueryEndpoint(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownQueryEndpoint, "failed. unknown query endpoint")
}

// ErrUnknownMsgType returns an error when the the msg type is unknown
func ErrUnknownMsgType(codespace sdk.CodespaceType, msgType string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownMsgType, "failed. unknown msg type: %s", msgType)
}

// ErrBadJSONMarshaling returns an error with the bad encoding of JSON
func ErrBadJSONMarshaling(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadJSONMarshaling, "failed. marshal JSON unsuccessfully: %s", msg)
}

// ErrBadJSONUnmarshaling returns an error with the bad decoding of JSON
func ErrBadJSONUnmarshaling(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadJSONUnmarshaling, "failed. unmarshal JSON unsuccessfully: %s", msg)
}
