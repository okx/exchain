package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeInvalidPriceDigit       sdk.CodeType = 1
	CodeInvalidMinTradeSize     sdk.CodeType = 2
	CodeInvalidDexList          sdk.CodeType = 3
	CodeInvalidBalanceNotEnough sdk.CodeType = 4
	CodeInvalidHeight           sdk.CodeType = 5
	CodeInvalidAsset            sdk.CodeType = 6
	CodeInvalidCommon           sdk.CodeType = 7
	CodeBlockedRecipient        sdk.CodeType = 8
	CodeSendDisabled            sdk.CodeType = 9
)

// ErrBlockedRecipient returns an error when a transfer is tried on a blocked recipient
func ErrBlockedRecipient(codespace sdk.CodespaceType, blockedAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeBlockedRecipient, "failed. %s is not allowed to receive transactions", blockedAddr)
}

// ErrSendDisabled returns an error when the transaction sending is disabled in bank module
func ErrSendDisabled(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSendDisabled, "failed. send transactions are currently disabled")
}

func ErrInvalidDexList(codespace sdk.CodespaceType, message string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDexList, message)
}

func ErrInvalidBalanceNotEnough(codespace sdk.CodespaceType, message string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBalanceNotEnough, message)
}

func ErrInvalidHeight(codespace sdk.CodespaceType, h, ch, max int64) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHeight, fmt.Sprintf("Height %d must be greater than current block height %d and less than %d + %d.", h, ch, ch, max))
}

func ErrInvalidCommon(codespace sdk.CodespaceType, message string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCommon, message)
}
