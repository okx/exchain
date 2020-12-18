package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const uint32
const (
	CodeInvalidAddress							uint32 = 63000
	CodeProductIsEmpty 							uint32 = 63001
	CodeSizeIsInvalid  							uint32 = 63002
	CodeGetTokenPairFailed						uint32 = 63003
	CodeSendCoinsFromAccountToAccountFaile		uint32 = 63004
	CodeTradingPairIsdelisting					uint32 = 63005
	CodeRoundedPriceEqual						uint32 = 63006
	CodeRoundedQuantityEqual					uint32 = 63007
	CodeMsgQuantityLessThan						uint32 = 63008
	CodeGetOrderFailed							uint32 = 63013
	CodeTokenPairNotFound						uint32 = 63014
	CodeCheckTokenPairUnderDexDelistFailed		uint32 = 63015
	CodeIsProductLocked							uint32 = 63016
	CodeNoOrdersIsCanceled						uint32 = 63017
	CodeOrderStatusIsNotOpen					uint32 = 63018
	CodeOrderIsNotExist							uint32 = 63020
	CodeUnknownOrderQueryType					uint32 = 63021
	CodeOrderItemCountsBiggerThanLimit			uint32 = 63022
	CodeOrderItemCountsIsEmpty					uint32 = 63023
	CodeOrderItemProductCountsIsEmpty			uint32 = 63024
	CodeOrderItemProductSymbolError				uint32 = 63025
	CodeOrderItemProductSymbolIsEqual			uint32 = 63026
	CodeOrderItemSideIsNotBuyAndSell			uint32 = 63027
	CodeOrderItemPriceOrQuantityIsNotPositive	uint32 = 63028
	CodeOrderIDsIsEmpty							uint32 = 63029
	CodeOrderIDCountsBiggerThanMultiCancelOrderItemLimit	uint32 = 63030
	CodeOrderIDsHasDuplicatedID					uint32 = 63031
	CodeUserinputOrderIDIsEmpty					uint32 = 63032
	CodeInputSenderNotEqualOrderSender			uint32 = 63033
)

func ErrInvalidAddress(address string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("invalid address: %s", address))
}
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

func ErrGetOrderFailed(order string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetOrderFailed, fmt.Sprintf("order(%v) does not exist", order))
}

func ErrCheckTokenPairUnderDexDelistFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeCheckTokenPairUnderDexDelistFailed, "check token pair under dex delist failed")
}

func ErrIsProductLocked() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsProductLocked, "current product is locked")
}

func ErrNoOrdersIsCanceled() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNoOrdersIsCanceled, "no order is cancled")
}

func ErrOrderStatusIsNotOpen() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderStatusIsNotOpen, "order status is not open")
}

func ErrOrderIsNotExist() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIsNotExist, "order is not exist")
}

func ErrUnknowOrdernQueryType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownOrderQueryType, "unknown order query type")
}

func ErrOrderItemCountsBiggerThanLimit() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemCountsBiggerThanLimit, "order item counts bigger than limit")
}

func ErrOrderItemCountsIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemCountsIsEmpty, "order item counts is empty")
}

func ErrOrderItemProductCountsIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemProductCountsIsEmpty, "order item's product counts is empty")
}

func ErrOrderItemProductSymbolError() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemProductSymbolError, "order item's product symbol is error")
}

func ErrOrderItemProductSymbolIsEqual() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemProductSymbolIsEqual, "order item's product two symbols is equal")
}

func ErrOrderItemSideIsNotBuyAndSell() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemSideIsNotBuyAndSell, "order item's side is not buy and sell")
}

func ErrOrderItemPriceOrQuantityIsNotPositive() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemPriceOrQuantityIsNotPositive, "order item's price or quantity is not positive")
}

func ErrOrderIDsIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIDsIsEmpty, "order Ids is empty")
}

func ErrOrderIDCountsBiggerThanMultiCancelOrderItemLimit() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIDCountsBiggerThanMultiCancelOrderItemLimit, "order Id counts bigger than multi cancel order item limits")
}

func ErrOrderIDsHasDuplicatedID() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIDsHasDuplicatedID, "order Ids has duplicated ID")
}

func ErrUserinputOrderIDIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUserinputOrderIDIsEmpty, "user input order id is empty")
}

func ErrInputSenderNotEqualOrderSender() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInputSenderNotEqualOrderSender, "user input sender address is not equal order's sender")
}