package types

import (
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

// const uint32
const (
	CodeInvalidAddress                        uint32 = 63000
	CodeSizeIsInvalid                         uint32 = 63001
	CodeTokenPairNotExist                     uint32 = 63002
	CodeSendCoinsFailed                        uint32 = 63003
	CodeTradingPairIsDelisting                uint32 = 63004
	CodePriceOverAccuracy                     uint32 = 63005
	CodeQuantityOverAccuracy                  uint32 = 63006
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
	CodeUserInputOrderIDIsEmpty               uint32 = 63025
	CodeNotOrderOwner                         uint32 = 63026
	CodeProductIsEmpty                        uint32 = 63027
	CodeAllOrderFailedToExecute               uint32 = 63028
)

func ErrInvalidAddress(address string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("invalid address: %s", address))}
}

// invalid size
func ErrInvalidSizeParam(size uint) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSizeIsInvalid, fmt.Sprintf("invalid param: size= %d", size))}
}

func ErrTokenPairNotExist(product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTokenPairNotExist, fmt.Sprintf("token pair %s doesn't exist", product))}
}

func ErrSendCoinsFailed(coins string, to string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSendCoinsFailed, fmt.Sprintf("send fee(%s) to address(%s) failed\n", coins, to))}
}

func ErrTradingPairIsDelisting(product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTradingPairIsDelisting, fmt.Sprintf("trading pair '%s' is delisting", product))}
}

func ErrPriceOverAccuracy(price sdk.Dec, priceDigit int64) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodePriceOverAccuracy, fmt.Sprintf("price(%v) over accuracy(%d)", price, priceDigit))}
}

func ErrQuantityOverAccuracy(quantity sdk.Dec, quantityDigit int64) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeQuantityOverAccuracy, fmt.Sprintf("quantity(%v) over accuracy(%d)", quantity, quantityDigit))}
}

func ErrMsgQuantityLessThan(minQuantity string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMsgQuantityLessThan, fmt.Sprintf("quantity should be greater than %s", minQuantity))}
}

func ErrOrderIsNotExist(order string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderIsNotExist, fmt.Sprintf("order(%v) does not exist", order))}
}

func ErrCheckTokenPairUnderDexDelistFailed() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeCheckTokenPairUnderDexDelistFailed, "check token pair under dex delist failed")}
}

func ErrIsProductLocked(product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeIsProductLocked, fmt.Sprintf("the trading pair (%s) is locked, please retry later", product))}
}

func ErrNoOrdersIsCanceled() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeNoOrdersIsCanceled, "no order is cancled")}
}

func ErrOrderStatusIsNotOpen() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderStatusIsNotOpen, "order status is not open")}
}

func ErrOrderIsNotExistOrClosed(orderID string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderIsNotExistOrClosed, fmt.Sprintf("order(%s) does not exist or already closed", orderID))}
}

func ErrUnknownOrderQueryType() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeUnknownOrderQueryType, "unknown order query endpoint")}
}

func ErrOrderItemCountsBiggerThanLimit(limit int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderItemCountsBiggerThanLimit, fmt.Sprintf("order item counts bigger than limit %d", limit))}
}

func ErrOrderItemCountsIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderItemCountsIsEmpty, "order item counts is empty")}
}

func ErrOrderItemProductCountsIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderItemProductCountsIsEmpty, "order item's product counts is empty")}
}

func ErrOrderItemProductFormat() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderItemProductSymbolError, "order item's product format is error")}
}

func ErrOrderItemProductSymbolIsEqual() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderItemProductSymbolIsEqual, "order item's product two symbols is equal")}
}

func ErrOrderItemSideIsNotBuyAndSell() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderItemSideIsNotBuyAndSell, "order item's side is not \"BUY\" or \"SELL\"")}
}

func ErrOrderItemPriceOrQuantityIsNotPositive() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderItemPriceOrQuantityIsNotPositive, "order item's price or quantity is not positive")}
}

func ErrOrderIDsIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderIDsIsEmpty, "order IDs is empty")}
}

func ErrCancelOrderBiggerThanLimit(limit int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeCancelOrderBiggerThanLimit, fmt.Sprintf("Numbers of CancelOrderItem should not be bigger than limit %d", limit))}
}

func ErrOrderIDsHasDuplicatedID() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderIDsHasDuplicatedID, "order IDs has duplicated ID")}
}

func ErrUserInputOrderIDIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeUserInputOrderIDIsEmpty, "user input order id is empty")}
}

func ErrNotOrderOwner(orderID string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeNotOrderOwner, fmt.Sprintf("not the owner of order(%v)", orderID))}
}

func ErrAllOrderFailedToExecute() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeAllOrderFailedToExecute, "all order items failed to execute")}
}
