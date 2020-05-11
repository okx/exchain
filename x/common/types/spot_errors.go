package types

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	SpotCodespace sdk.CodespaceType = "spot"

	CodeEmptyProduct         CodeType = 62001
	CodeExistingProduct      CodeType = 62002
	CodeNonexistentProduct   CodeType = 62003
	CodeInvalidProduct       CodeType = 62004
	CodeInvalidProductOwner  CodeType = 62005
	CodeInvalidWithdraw      CodeType = 62006
	CodeInvalidToken         CodeType = 62007
	CodeSaveProductFailed    CodeType = 62008
	CodeOverAccuracyPrice    CodeType = 62009
	CodeOverAccuracyQuantity CodeType = 62010
	CodeOverAccuracy         CodeType = 62011
	CodeInvaildQuantity      CodeType = 62012
	CodeLockedProduct        CodeType = 62013
	CodeEmptyOrders          CodeType = 62014
	CodeOverLimitedOrders    CodeType = 62015
	CodeInvaildFormatProduct CodeType = 62016
	CodeInvaildSideParam     CodeType = 62017
	CodeNegativeParam        CodeType = 62018
	CodeEmptyOrderId         CodeType = 62019
	CodeOverLimitedCancelOrders CodeType = 62020
	CodeDuplicatedOrderId    CodeType = 62021
	CodeNonexistentOrder     CodeType = 62022
)

// ErrEmptySymbol returns an error with an empty symbol of token
func ErrEmptyProduct(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyProduct,
		"failed. product cannot be empty")
}

// ErrExistingProduct returns an error with an existing product
func ErrExistingProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeExistingProduct,
		fmt.Sprintf("failed. product %s already exists", product))
}

// ErrExistingProduct returns an error with the nonexistent product
func ErrNonexistentProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeNonexistentProduct,
		fmt.Sprintf("failed. product %s does not exist", product))
}

// ErrInvalidProduct returns an error with an invalid product
func ErrInvalidProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProduct,
		fmt.Sprintf("failed. product %s is invalid", product))
}

// ErrProductUnauthorizedIdentity returns an error with the unauthorized identity of the owner
func ErrProductUnauthorizedIdentity(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProductOwner,
		fmt.Sprintf("failed. not the owner of product: %s", product))
}

// ErrInvalidToken returns an error with the invalid token
func ErrInvalidToken(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidToken,
		fmt.Sprintf("failed. deposits and withdraw only support %s token", symbol))
}

// ErrInvalidWithdrawAmount returns an error with the invalid withdraw info
func ErrInvalidWithdrawAmount(codespace sdk.CodespaceType, depositsAmount, withdrawAmount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidWithdraw,
		fmt.Sprintf("failed. deposits:%s is less than withdraw:%s", depositsAmount, withdrawAmount))
}

// ErrSaveProduct returns an error when save product to db failed
func ErrSaveProduct(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeSaveProductFailed,
		fmt.Sprintf("failed. save product error: %s", msg))
}

func ErrOverAccuracyPrice(codespace sdk.CodespaceType, price sdk.Dec, digit int64) sdk.Error {
	return sdk.NewError(codespace, CodeOverAccuracyPrice,
		fmt.Sprintf("failed. price(%v) is over accuracy(%d)", price, digit))
}

func ErrOverAccuracyQuantity(codespace sdk.CodespaceType, quantity sdk.Dec, digit int64) sdk.Error {
	return sdk.NewError(codespace, CodeOverAccuracyQuantity,
		fmt.Sprintf("failed. quantity(%v) is over accuracy(%d)", quantity, digit))
}

func ErrOverAccuracy(codespace sdk.CodespaceType, price sdk.Dec, quantity sdk.Dec, digit int64) sdk.Error {
	return sdk.NewError(codespace, CodeOverAccuracy,
		fmt.Sprintf("failed. price(%v) * quantity(%v) is over accuracy(%d)", price, quantity, digit))
}

func ErrInvaildQuantity(codespace sdk.CodespaceType, quantity sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvaildQuantity,
		fmt.Sprintf("failed. quantity should be greater than %s", quantity))
}

func ErrInDexlistProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeInvaildQuantity,
		fmt.Sprintf("failed. product %s is delisting", product))
}

func ErrLockedProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeLockedProduct,
		fmt.Sprintf("failed. the product (%s) is locked, please retry later", product))
}

func ErrEmptyOrders(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyOrders,
		"failed. orderItems is empty")
}

func ErrOverLimitedOrders(codespace sdk.CodespaceType, limit int) sdk.Error {
	return sdk.NewError(codespace, CodeOverLimitedOrders,
		fmt.Sprintf("failed. the number of NewOrderItem should not be more than "+strconv.Itoa(limit)))
}

func ErrInvaildFormatProduct(codespace sdk.CodespaceType, product string) sdk.Error {
	return sdk.NewError(codespace, CodeInvaildFormatProduct,
		fmt.Sprintf("failed. product %s must be in the format of \"base_quote\"", product))
}

func ErrInvaildSideParam(codespace sdk.CodespaceType, side string) sdk.Error {
	return sdk.NewError(codespace, CodeInvaildSideParam,
		fmt.Sprintf("failed. side is expected to be \"BUY\" or \"SELL\", but got \"%s\"", side))
}

func ErrNegativeParam(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNegativeParam,
		"failed. price and quantity must be positive")
}

func ErrEmptyOrderId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyOrderId,
		"failed. order id is empty")
}

func ErrOverLimitedCancelOrders(codespace sdk.CodespaceType, limit int) sdk.Error {
	return sdk.NewError(codespace, CodeOverLimitedCancelOrders,
		fmt.Sprintf("failed. the number of CancelOrderItem should not be more than "+strconv.Itoa(limit)))
}

func ErrDuplicatedOrderId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicatedOrderId,
		"failed. duplicated order ids detected")
}

func ErrNonexistentOrder(codespace sdk.CodespaceType, id string) sdk.Error {
	return sdk.NewError(codespace, CodeNonexistentOrder,
		fmt.Sprintf("failed. order %s does not exist", id))
}