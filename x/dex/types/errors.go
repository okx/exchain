package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const CodeType
const (
	codeInvalidProduct          sdk.CodeType = 1
	codeTokenPairNotFound       sdk.CodeType = 2
	codeDelistOwnerNotMatch     sdk.CodeType = 3
	codeInvalidBalanceNotEnough sdk.CodeType = 4
	codeInvalidAsset            sdk.CodeType = 5
	codeUnknownOperator         sdk.CodeType = 6
	codeExistOperator           sdk.CodeType = 7
	codeInvalidWebsiteLength    sdk.CodeType = 8
	codeInvalidWebsiteURL       sdk.CodeType = 9
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

func ErrUnknownOperator(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeUnknownOperator, fmt.Sprintf("unknown dex operator with address %s", addr.String()))
}

func ErrExistOperator(addr sdk.AccAddress) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeExistOperator, fmt.Sprintf("dex operator already exists with address %s", addr.String()))
}

func ErrInvalidWebsiteLength(got, max int) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidWebsiteLength, fmt.Sprintf("invalid website length, got length %v, max is %v", got, max))
}

func ErrInvalidWebsiteURL(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidWebsiteURL, fmt.Sprintf("invalid website URL: %s", msg))
}
