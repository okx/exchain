package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const CodeType
const (
	CodeInvalidProduct      = 1
	CodeTokenPairNotFound   = 2
	CodeDelistOwnerNotMatch = 3

	CodeInvalidBalanceNotEnough sdk.CodeType = 4
	CodeInvalidHeight           sdk.CodeType = 5
	CodeInvalidAsset            sdk.CodeType = 6
	CodeInvalidCommon           sdk.CodeType = 7
)

// CodeType to Message
func CodeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidProduct:
		return "invalid product"
	case CodeTokenPairNotFound:
		return "tokenpair not found"
	case CodeDelistOwnerNotMatch:
		return "tokenpair delistor should be it's owner "
	default:
		return fmt.Sprintf("unknown code %d", code)
	}
}

// SDK Errors Functor
// All error raised in this module is kept here.
// Global errors which can be seen outside.
func ErrInvalidProduct(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidProduct, CodeToDefaultMsg(CodeInvalidProduct)+": %s", msg)
}

func ErrTokenPairNotFound(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTokenPairNotFound, CodeToDefaultMsg(CodeTokenPairNotFound)+": %s", msg)
}

func ErrDelistOwnerNotMatch(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDelistOwnerNotMatch, CodeToDefaultMsg(CodeDelistOwnerNotMatch)+": %s", msg)
}

func ErrInvalidBalanceNotEnough(codespace sdk.CodespaceType, message string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBalanceNotEnough, message)
}

func ErrInvalidHeight(codespace sdk.CodespaceType, h, ch, max int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHeight, fmt.Sprintf("Height %d must be greater than current block height %d and less than %d + %d.", h, ch, ch, max))
}

func ErrInvalidCommon(codespace sdk.CodespaceType, message string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCommon, message)
}

func ErrFailToDeleteTokenPair(codespace sdk.CodespaceType, tokenPair string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCommon, "Failed to delete token pair: %s", tokenPair)
}
