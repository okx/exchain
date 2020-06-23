package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const CodeType
const (
	codeInvalidProduct      sdk.CodeType = 1
	codeTokenPairNotFound   sdk.CodeType = 2
	codeDelistOwnerNotMatch sdk.CodeType = 3

	codeInvalidBalanceNotEnough sdk.CodeType = 4
	codeInvalidAsset            sdk.CodeType = 5
)

// CodeType to Message
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case codeInvalidProduct:
		return "invalid product"
	case codeTokenPairNotFound:
		return "tokenpair not found"
	case codeDelistOwnerNotMatch:
		return "tokenpair delistor should be it's owner "
	default:
		return fmt.Sprintf("unknown code %d", code)
	}
}

// ErrInvalidProduct returns invalid product error
func ErrInvalidProduct(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidProduct, codeToDefaultMsg(codeInvalidProduct)+": %s", msg)
}

// ErrTokenPairNotFound returns token pair not found error
func ErrTokenPairNotFound(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeTokenPairNotFound, codeToDefaultMsg(codeTokenPairNotFound)+": %s", msg)
}

// ErrDelistOwnerNotMatch returns delist owner not match error
func ErrDelistOwnerNotMatch(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeDelistOwnerNotMatch, codeToDefaultMsg(codeDelistOwnerNotMatch)+": %s", msg)
}

// ErrInvalidBalanceNotEnough returns invalid balance not enough error
func ErrInvalidBalanceNotEnough(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidBalanceNotEnough, message)
}

// ErrInvalidAsset returns invalid asset error
func ErrInvalidAsset(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidAsset, message)
}

// ErrTokenPairExisted returns an error when the token pair is existing during the process of listing
func ErrTokenPairExisted(baseAsset, quoteAsset string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidAsset,
		fmt.Sprintf("failed. the token pair exists with %s and %s", baseAsset, quoteAsset))
}
