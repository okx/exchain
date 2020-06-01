package types

// TODO: Fill out some custom errors for the module
// You can see how they are constructed below:
// var (
//	ErrInvalid = sdkerrors.Register(ModuleName, 1, "custom error message")
// )
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	Codespace sdk.CodespaceType = ModuleName

	CodeInvalidLeverage     CodeType = 62100
	CodeInvalidTradePair    CodeType = 62101
	CodeAccountNotExist     CodeType = 62102
	CodeInvalidBorrowAmount CodeType = 62103
	CodeNotAllowed          CodeType = 62104
)

func ErrInvalidLeverage(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLeverage,
		"failed. invalid leverage : %s", msg)
}

func ErrInvalidTradePair(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTradePair,
		"failed. invalid trade pair : %s", msg)
}

func ErrAccountNotExist(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeAccountNotExist,
		"failed. invalid request : %s", msg)
}

func ErrInvalidBorrowAmount(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBorrowAmount,
		"failed. invalid borrow amount : %s", msg)
}

func ErrNotAllowed(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeNotAllowed, "not allowed: %s", msg)
}
