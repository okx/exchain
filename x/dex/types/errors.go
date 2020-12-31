package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const CodeType
const (
	CodeAddrAndProductAllRequired           sdk.CodeType = 64000
	codeInvalidTokenPair                    sdk.CodeType = 64001
	codeTokenPairNotFound                   sdk.CodeType = 64002
	codeDelistOwnerNotMatch                 sdk.CodeType = 64003
	codeInvalidBalanceNotEnough             sdk.CodeType = 64004
	codeInvalidAsset                        sdk.CodeType = 64005
	codeUnknownOperator                     sdk.CodeType = 64006
	codeExistOperator                       sdk.CodeType = 64007
	codeInvalidWebsiteLength                sdk.CodeType = 64008
	codeInvalidWebsiteURL                   sdk.CodeType = 64009
	CodeTokenPairIsInvalid                  sdk.CodeType = 64010
	CodeTokenPairSaveFailed                 sdk.CodeType = 64011
	CodeInsufficientFeeCoins                sdk.CodeType = 64012
	CodeTokenPairAlreadyExist               sdk.CodeType = 64013
	CodeMustTokenPairOwner                  sdk.CodeType = 64014
	CodeDepositOnlySupportDefaultBondDenom  sdk.CodeType = 64015
	CodeInsufficientDepositCoins            sdk.CodeType = 64016
	CodeWithdrawOnlySupportDefaultBondDenom sdk.CodeType = 64017
	CodeInsufficientWithdrawCoins           sdk.CodeType = 64018
	CodeWithdrawDepositsError               sdk.CodeType = 64019
	CodeMustOperatorOwner                   sdk.CodeType = 64020
)

// CodeType to Message
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case codeInvalidTokenPair:
		return "invalid tokenpair"
	case codeTokenPairNotFound:
		return "tokenpair not found"
	case codeDelistOwnerNotMatch:
		return "tokenpair delistor should be it's owner "
	default:
		return fmt.Sprintf("unknown code %d", code)
	}
}

// Addr and Product All Required
func ErrAddrAndProductAllRequired() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeAddrAndProductAllRequired, "bad request: address、base_asset and quote_asset could not be empty at the same time")
}

// invalid tokenpair
func ErrTokenPairIsInvalid() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTokenPairIsInvalid, "the nil pointer is not expected")
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

// ErrTokenPairExisted returns an error when the token pair is existing during the process of listing
func ErrTokenPairExisted(baseAsset string, quoteAsset string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTokenPairAlreadyExist,
		fmt.Sprintf("failed. the token pair exists with %s and %s", baseAsset, quoteAsset))
}

// ErrInvalidTokenPair returns invalid product error
func ErrInvalidTokenPair(tokenPair string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidTokenPair, codeToDefaultMsg(codeInvalidTokenPair)+": %s", tokenPair)
}
func ErrTokenPairSaveFailed(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTokenPairSaveFailed, fmt.Sprintf("failed to SaveTokenPair: %s", msg))
}
func ErrInsufficientFeeCoins(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInsufficientFeeCoins, fmt.Sprintf("insufficient fee coins(need %s)", msg))
}
func ErrMustTokenPairOwner(addr string, tokenPair string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeMustTokenPairOwner, fmt.Sprintf("failed because %s is not the owner of product:%s", addr, tokenPair))
}
func ErrDepositOnlySupportDefaultBondDenom(defaultBondDenom string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDepositOnlySupportDefaultBondDenom, fmt.Sprintf("failed to deposit because deposits only support %s token", defaultBondDenom))
}
func ErrInsufficientDepositCoins(msg string, depositCoins string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInsufficientDepositCoins, fmt.Sprintf("failed: %s, because insufficient deposit coins(need %s)", msg, depositCoins))
}
func ErrWithdrawOnlySupportDefaultBondDenom(defaultBondDenom string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeWithdrawOnlySupportDefaultBondDenom, fmt.Sprintf("failed to withdraws because deposits only support %s token", defaultBondDenom))
}
func ErrInsufficientWithdrawCoins(depositCoins string, amount string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInsufficientWithdrawCoins, fmt.Sprintf("failed to withdraws because deposits:%s is less than withdraw:%s", depositCoins, amount))
}
func ErrWithdrawDepositsError(depositCoins string, msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeWithdrawDepositsError, fmt.Sprintf("withdraw deposits:%s error:%s", depositCoins, msg))
}
func ErrMustOperatorOwnerOwner(addr string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeMustOperatorOwner, fmt.Sprintf("failed because %s is not the owner of operator", addr))
}
