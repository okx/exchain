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

func ErrGetTokenPairFailed(product string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetTokenPairFailed, fmt.Sprintf("failed. token pair %s doesn't exist", product))
}

func ErrSendCoinsFromAccountToAccountFaile(coins string, to string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromAccountToAccountFaile, fmt.Sprintf("send fee(%s) to address(%s) failed\n", coins, to))
}

func ErrTradingPairIsdelisting(product string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTradingPairIsdelisting, fmt.Sprintf("trading pair '%s' is delisting", product))
}

func ErrRoundedPriceEqual(priceDigit int64) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeRoundedPriceEqual, fmt.Sprintf("price over accuracy(%d)", priceDigit))
}

func ErrRoundedQuantityEqual(quantityDigit int64) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeRoundedQuantityEqual, fmt.Sprintf("quantity over accuracy(%d)", quantityDigit))
}

func ErrMsgQuantityLessThan(minQuantity string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMsgQuantityLessThan, fmt.Sprintf("quantity should be greater than %s", minQuantity))
}

func ErrInternal() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternal, "occur error internal")
}

func ErrInsufficientCoins() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientCoins, "insufficient coins")
}

func ErrUnknownRequest() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, "unknown request")
}

func ErrGetOrderFailed(order string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetOrderFailed, fmt.Sprintf("order(%v) does not exist", order))
}

func ErrTokenPairNotFound(product string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenPairNotFound, fmt.Sprintf("token pair not found %s", product))
}