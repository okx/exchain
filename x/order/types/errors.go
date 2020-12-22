package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// const uint32
const (
	CodeInvalidAddress                        uint32 = 63000
	CodeSizeIsInvalid                         uint32 = 63001
	CodeTokenPairNotExist                     uint32 = 63002
	CodeSendCoinsFaile                        uint32 = 63003
	CodeTradingPairIsDelisting                uint32 = 63004
	CodeRoundedPriceEqual                     uint32 = 63005
	CodeRoundedQuantityEqual                  uint32 = 63006
	CodeMsgQuantityLessThan                   uint32 = 63007
	CodeOrderIsNotExist                       uint32 = 63008
	CodeCheckTokenPairUnderDexDelistFailed    uint32 = 63009
	CodeIsProductLocked                       uint32 = 63010
	CodeNoOrdersIsCanceled                    uint32 = 63011
	CodeOrderStatusIsNotOpen                  uint32 = 63012
	CodeOrderIsNotExistOrClosed               uint32 = 63013
	CodeUnknownOrderQueryType                 uint32 = 63014
	CodeOrderItemCountsBiggerThanLimit        uint32 = 63015
	CodeOrderItemCountsIsEmpty                uint32 = 63016
	CodeOrderItemProductCountsIsEmpty         uint32 = 63017
	CodeOrderItemProductSymbolError           uint32 = 63018
	CodeOrderItemProductSymbolIsEqual         uint32 = 63019
	CodeOrderItemSideIsNotBuyAndSell          uint32 = 63020
	CodeOrderItemPriceOrQuantityIsNotPositive uint32 = 63021
	CodeOrderIDsIsEmpty                       uint32 = 63022
	CodeCancelOrderBiggerThanLimit            uint32 = 63023
	CodeOrderIDsHasDuplicatedID               uint32 = 63024
	CodeUserinputOrderIDIsEmpty               uint32 = 63025
	CodeNotOrderOwner                         uint32 = 63026
	CodeProductIsEmpty                        uint32 = 63027
)

func ErrInvalidAddress(address string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("invalid address: %s", address))
}

// invalid size
func ErrInvalidSizeParam(size uint) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSizeIsInvalid, fmt.Sprintf("invalid param: size= %d", size))
}

func ErrTokenPairNotExist(product string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenPairNotExist, fmt.Sprintf("token pair %s doesn't exist", product))
}

func ErrSendCoinsFaile(coins string, to string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFaile, fmt.Sprintf("send fee(%s) to address(%s) failed\n", coins, to))
}

func ErrTradingPairIsDelisting(product string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTradingPairIsDelisting, fmt.Sprintf("trading pair '%s' is delisting", product))
}

func ErrPriceOverAccuracy(price sdk.Dec, priceDigit int64) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeRoundedPriceEqual, fmt.Sprintf("price(%v) over accuracy(%d)", price, priceDigit))
}

func ErrQuantityOverAccuracy(quantity sdk.Dec, quantityDigit int64) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeRoundedQuantityEqual, fmt.Sprintf("quantity(%v) over accuracy(%d)", quantity, quantityDigit))
}

func ErrMsgQuantityLessThan(minQuantity string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMsgQuantityLessThan, fmt.Sprintf("quantity should be greater than %s", minQuantity))
}

func ErrOrderIsNotExist(order string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIsNotExist, fmt.Sprintf("order(%v) does not exist", order))
}

func ErrCheckTokenPairUnderDexDelistFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeCheckTokenPairUnderDexDelistFailed, "check token pair under dex delist failed")
}

func ErrIsProductLocked(product string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsProductLocked, fmt.Sprintf("the trading pair (%s) is locked", product))
}

func ErrNoOrdersIsCanceled() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNoOrdersIsCanceled, "no order is cancled")
}

func ErrOrderStatusIsNotOpen() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderStatusIsNotOpen, "order status is not open")
}

func ErrOrderIsNotExistOrClosed(orderID string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIsNotExistOrClosed, fmt.Sprintf("order(%s) does not exist or already closed", orderID))
}

func ErrUnknowOrdernQueryType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownOrderQueryType, "unknown order query endpoint")
}

func ErrOrderItemCountsBiggerThanLimit(limit int) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemCountsBiggerThanLimit, fmt.Sprintf("order item counts bigger than limit %d", limit))
}

func ErrOrderItemCountsIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemCountsIsEmpty, "order item counts is empty")
}

func ErrOrderItemProductCountsIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemProductCountsIsEmpty, "order item's product counts is empty")
}

func ErrOrderItemProductFormat() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemProductSymbolError, "order item's product format is error")
}

func ErrOrderItemProductSymbolIsEqual() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemProductSymbolIsEqual, "order item's product two symbols is equal")
}

func ErrOrderItemSideIsNotBuyAndSell() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemSideIsNotBuyAndSell, "order item's side is not \"BUY\" or \"SELL\"")
}

func ErrOrderItemPriceOrQuantityIsNotPositive() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderItemPriceOrQuantityIsNotPositive, "order item's price or quantity is not positive")
}

func ErrOrderIDsIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIDsIsEmpty, "order IDs is empty")
}

func ErrCancelOrderBiggerThanLimit(limit int) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeCancelOrderBiggerThanLimit, fmt.Sprintf("Numbers of CancelOrderItem should not be bigger than limit %d", limit))
}

func ErrOrderIDsHasDuplicatedID() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeOrderIDsHasDuplicatedID, "order IDs has duplicated ID")
}

func ErrUserinputOrderIDIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUserinputOrderIDIsEmpty, "user input order id is empty")
}

func ErrNotOrderOwner(orderID string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNotOrderOwner, fmt.Sprintf("not the owner of order(%v)", orderID))
}
