package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// const CodeType
const (
	codeInvalidProduct          uint32 = 1
	codeTokenPairNotFound       uint32 = 2
	codeDelistOwnerNotMatch     uint32 = 3
	codeInvalidBalanceNotEnough uint32 = 4
	codeInvalidAsset            uint32 = 5
	codeUnknownOperator         uint32 = 6
	codeExistOperator           uint32 = 7
	codeInvalidWebsiteLength    uint32 = 8
	codeInvalidWebsiteURL       uint32 = 9
)

var (
	errInvalidProduct 			= sdkerrors.Register(DefaultCodespace, codeInvalidProduct, "invalid product")
	errTokenPairNotFound 		= sdkerrors.Register(DefaultCodespace, codeTokenPairNotFound, "token pair not found")
	errDelistOwnerNotMatch 		= sdkerrors.Register(DefaultCodespace, codeDelistOwnerNotMatch, "delist owner not match")
	errInvalidBalanceNotEnough 	= sdkerrors.Register(DefaultCodespace, codeInvalidBalanceNotEnough, "invalid balance not enough")
	errInvalidAsset 			= sdkerrors.Register(DefaultCodespace, codeInvalidAsset, "invalid asset")
	errUnknownOperator 			= sdkerrors.Register(DefaultCodespace, codeUnknownOperator, "unknown operator")
	errExistOperator 			= sdkerrors.Register(DefaultCodespace, codeExistOperator, "exist operator")
	errInvalidWebsiteLength 	= sdkerrors.Register(DefaultCodespace, codeInvalidWebsiteLength, "invalid website length")
	errInvalidWebsiteURL 		= sdkerrors.Register(DefaultCodespace, codeInvalidWebsiteURL, "invalid website URL")
)

// CodeType to Message
func codeToDefaultMsg(code uint32) string {
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
func ErrInvalidProduct(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidProduct, codeToDefaultMsg(codeInvalidProduct)+": %s", msg)}
}

// ErrTokenPairNotFound returns token pair not found error
func ErrTokenPairNotFound(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errTokenPairNotFound, codeToDefaultMsg(codeTokenPairNotFound)+": %s", msg)}
}

// ErrDelistOwnerNotMatch returns delist owner not match error
func ErrDelistOwnerNotMatch(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errDelistOwnerNotMatch, codeToDefaultMsg(codeDelistOwnerNotMatch)+": %s", msg)}
}

// ErrInvalidBalanceNotEnough returns invalid balance not enough error
func ErrInvalidBalanceNotEnough(message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidBalanceNotEnough, message)}
}

// ErrInvalidAsset returns invalid asset error
func ErrInvalidAsset(message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidAsset, message)}
}

func ErrUnknownOperator(addr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errUnknownOperator, fmt.Sprintf("unknown dex operator with address %s", addr.String()))}
}

func ErrExistOperator(addr sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errExistOperator, fmt.Sprintf("dex operator already exists with address %s", addr.String()))}
}

func ErrInvalidWebsiteLength(got, max int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidWebsiteLength, fmt.Sprintf("invalid website length, got length %v, max is %v", got, max))}
}

func ErrInvalidWebsiteURL(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidWebsiteURL, fmt.Sprintf("invalid website URL: %s", msg))}
}

// ErrTokenPairExisted returns an error when the token pair is existed during the process of listing
// ErrTokenPairExisted returns an error when the token pair is existing during the process of listing
func ErrTokenPairExisted(baseAsset, quoteAsset string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidAsset, fmt.Sprintf("failed. the token pair exists with %s and %s", baseAsset, quoteAsset))}
}
