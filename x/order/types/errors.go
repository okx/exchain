package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const uint32
const (
	CodeProductIsEmpty 							uint32 = 63001
	CodeSizeIsInvalid  							uint32 = 63002
	CodeGetTokenPairFailed						uint32 = 63003
	CodeSendCoinsFromAccountToAccountFaile		uint32 = 63004
	CodeTradingPairIsdelisting					uint32 = 63005
	CodeRoundedPriceEqual						uint32 = 63006
	CodeRoundedQuantityEqual					uint32 = 63007
	CodeMsgQuantityLessThan						uint32 = 63008
	CodeUnknownRequest							uint32 = 63009
	CodeInternal								uint32 = 63010
	CodeInsufficientCoins						uint32 = 63011
	CodeUnauthorized							uint32 = 63012
	CodeGetOrderFailed							uint32 = 63013
	CodeTokenPairNotFound						uint32 = 63014
)

// invalid size
func ErrInvalidSizeParam(size uint) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSizeIsInvalid, fmt.Sprintf("invalid param: size= %d", size))
}

func ErrGetTokenPairFailed(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetTokenPairFailed, message)
}

func ErrSendCoinsFromAccountToAccountFaile(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromAccountToAccountFaile, message)
}

func ErrTradingPairIsdelisting(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTradingPairIsdelisting, message)
}

func ErrRoundedPriceEqual(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeRoundedPriceEqual, message)
}

func ErrRoundedQuantityEqual(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeRoundedQuantityEqual, message)
}

func ErrMsgQuantityLessThan(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMsgQuantityLessThan, message)
}

func ErrInternal(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternal, message)
}

func ErrInsufficientCoins(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientCoins, message)
}

func ErrUnknownRequest(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, message)
}

func ErrGetOrderFailed(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetOrderFailed, message)
}

func ErrTokenPairNotFound(message string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenPairNotFound, message)
}