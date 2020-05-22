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
	MarginCodespace sdk.CodespaceType = "margin"

	CodeInvalidLeverage     CodeType = 62100
	CodeInvalidTradePair    CodeType = 62101
	CodeEmptyAccountDeposit CodeType = 62102
	CodeInvalidBorrowAmount CodeType = 62103
)

func ErrInvalidLeverage(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLeverage,
		"failed. invalid leverage : %s", msg)
}

func ErrInvalidTradePair(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTradePair,
		"failed. invalid trade pair : %s", msg)
}

func ErrEmptyAccountDeposit(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyAccountDeposit,
		"failed. invalid request : %s", msg)
}

func ErrInvalidBorrowAmount(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBorrowAmount,
		"failed. invalid borrow amount : %s", msg)
}
