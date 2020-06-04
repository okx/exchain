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
	CodeInvalidToken            sdk.CodeType = 8
)

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
