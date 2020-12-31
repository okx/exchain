package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const CodeType
const (
	CodeProductIsEmpty 							sdk.CodeType = 63001
	CodeSizeIsInvalid  							sdk.CodeType = 63002
	CodeGetTokenPairFailed						sdk.CodeType = 63003
	CodeSendCoinsFromAccountToAccountFaile		sdk.CodeType = 63004
	CodeTradingPairIsdelisting					sdk.CodeType = 63005
	CodeRoundedPriceEqual						sdk.CodeType = 63006
	CodeRoundedQuantityEqual					sdk.CodeType = 63007
	CodeMsgQuantityLessThan						sdk.CodeType = 63008
	CodeUnknownRequest							sdk.CodeType = 63009
	CodeInternal								sdk.CodeType = 63010
	CodeInsufficientCoins						sdk.CodeType = 63011
	CodeUnauthorized							sdk.CodeType = 63012
)

// invalid size
func ErrInvalidSizeParam(size uint) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSizeIsInvalid, fmt.Sprintf("invalid param: size= %d", size))
}

func ErrGetTokenPairFailed(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeGetTokenPairFailed, message)
}

func ErrSendCoinsFromAccountToAccountFaile(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSendCoinsFromAccountToAccountFaile, message)
}

func ErrTradingPairIsdelisting(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTradingPairIsdelisting, message)
}

func ErrRoundedPriceEqual(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRoundedPriceEqual, message)
}

func ErrRoundedQuantityEqual(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeRoundedQuantityEqual, message)
}

func ErrMsgQuantityLessThan(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeMsgQuantityLessThan, message)
}

func ErrInternal(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInternal, message)
}

func ErrInsufficientCoins(message string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInsufficientCoins, message)
}
