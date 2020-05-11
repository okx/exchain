package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	AssetCodespace sdk.CodespaceType = "asset"

	CodeTotalSupplyExceeds        CodeType = 61001
	CodeBadSymbolGeneration       CodeType = 61002
	CodeBadCoinsMintage           CodeType = 61003
	CodeBadCoinsSendingToModule   CodeType = 61004
	CodeBadCoinsSendingFromModule CodeType = 61005
	CodeInsufficientFees          CodeType = 61006
	CodeUnauthorizedIdentity      CodeType = 61007
	CodeBadCoinsBurning           CodeType = 61008
	CodeCoinsNotMintable          CodeType = 61009
	CodeInsufficientBalance       CodeType = 61010
	CodeInvalidModificationSet    CodeType = 61011
	CodeEmptyOriginalSymbol       CodeType = 61012
	CodeInvalidOriginalSymbol     CodeType = 61013
	CodeInvalidWholeName          CodeType = 61014
	CodeTokenDescriptionExceeds   CodeType = 61015
	CodeTransfersLengthExceeds    CodeType = 61016
	CodeInvalidMultisignCheck     CodeType = 61017
	CodeInvalidCoins              CodeType = 61018
	CodeMintageAmountExceeds      CodeType = 61019
	CodeEmptySymbol               CodeType = 61020
	CodeInvalidSymbol             CodeType = 61021
	CodeInsufficientCoins         CodeType = 61022
)

// ErrInvalidMultisignCheck returns an error with an invalid check result of multi-sign
func ErrInvalidMultisignCheck(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMultisignCheck,
		"failed. invalid check result of multi-sign")
}

// ErrInvalidSymbol returns an error with an invalid symbol of token
func ErrInvalidSymbol(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSymbol,
		"failed. invalid symbol of token: %s", symbol)
}

// ErrEmptySymbol returns an error with an empty symbol of token
func ErrEmptySymbol(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptySymbol,
		"failed. empty symbol of token")
}

// ErrTransfersLengthExceeds returns an error when the transfers' length of multi-send exceeds the limit
func ErrTransfersLengthExceeds(codespace sdk.CodespaceType, limit int64) sdk.Error {
	return sdk.NewError(codespace, CodeTransfersLengthExceeds,
		"failed. the length of transfers in multi-send exceeds the limit: %d", limit)
}

// ErrMintageAmountExceeds returns an error when the amount of mintage exceeds the limit
func ErrMintageAmountExceeds(codespace sdk.CodespaceType, limit int64) sdk.Error {
	return sdk.NewError(codespace, CodeMintageAmountExceeds,
		"failed. the amount of mintage exceeds the limit: %d", limit)
}

// ErrInvalidCoins returns an error with invalid coins
func ErrInvalidCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCoins,
		"failed. invalid coins with negative amount or illegal denomination")
}

// ErrTokenDescriptionExceeds returns an error when the token description exceeds the length limit
func ErrTokenDescriptionExceeds(codespace sdk.CodespaceType, lenLimit int) sdk.Error {
	return sdk.NewError(codespace, CodeTokenDescriptionExceeds,
		"failed. token description exceeds the length limit: %d", lenLimit)
}

// ErrInvalidWholeName returns an error with an invalid whole name of token
func ErrInvalidWholeName(codespace sdk.CodespaceType, wholeName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidWholeName,
		"failed. invalid whole name of token: %s", wholeName)
}

// ErrInvalidOriginalSymbol returns an error with an invalid original symbol of token
func ErrInvalidOriginalSymbol(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOriginalSymbol,
		"failed. invalid original symbol of token: %s", symbol)
}

// ErrEmptyOriginalSymbol returns an error with an empty original symbol of token
func ErrEmptyOriginalSymbol(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyOriginalSymbol,
		"failed. empty original symbol of token")
}

// ErrInvalidModification returns an error when neither "IsWholeNameModified" nor "IsDescriptionModified" is true
func ErrInvalidModification(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidModificationSet,
		`failed. invalid set for token's modification: neither "IsWholeNameModified" nor "IsDescriptionModified" is true`)
}

// ErrInsufficientBalance returns an error when the balance of an account is insufficient
func ErrInsufficientBalance(codespace sdk.CodespaceType, expectedAmount string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientBalance,
		"failed. insufficient balance, needs %s", expectedAmount)
}

// ErrCoinsNotMintable returns an error with the mintage of coins which are not mintable
func ErrCoinsNotMintable(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeCoinsNotMintable,
		"failed. token %s is not mintable", symbol)
}

// ErrBadCoinsBurning returns an error with the bad coins burning
func ErrBadCoinsBurning(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadCoinsBurning,
		"failed. bad coins burning: %s", msg)
}

// ErrUnauthorizedIdentity returns an error with the unauthorized identity of the owner
func ErrUnauthorizedIdentity(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeUnauthorizedIdentity,
		"failed. not the owner of token: %s", symbol)
}

// ErrTotalSupplyExceeds returns an error when the token's total supply exceeds the upper limit
func ErrTotalSupplyExceeds(codespace sdk.CodespaceType, totalSupply string, upperLimit int64) sdk.Error {
	return sdk.NewError(codespace, CodeTotalSupplyExceeds,
		"failed. total supply %s exceeds the upper limit %d", totalSupply, upperLimit)
}

// ErrBadSymbolGeneration returns an error with the bad unique symbol generation
func ErrBadSymbolGeneration(codespace sdk.CodespaceType, originalSymbol string) sdk.Error {
	return sdk.NewError(codespace, CodeBadSymbolGeneration,
		"failed. bad unique symbol generation for token %s", originalSymbol)
}

// ErrBadCoinsMintage returns an error with the bad coin mintage
func ErrBadCoinsMintage(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadCoinsMintage,
		"failed. bad coins mintage: %s", msg)
}

// ErrBadCoinsSendingToModule returns an error with the bad coins sending to module account
func ErrBadCoinsSendingToModule(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadCoinsSendingToModule,
		"failed. bad coins sending to module account: %s", msg)
}

// ErrBadCoinsSendingFromModule returns an error with the bad coins sending from module account
func ErrBadCoinsSendingFromModule(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeBadCoinsSendingFromModule,
		"failed. bad coins sending from module account: %s", msg)
}

// ErrInsufficientFee returns an error with the insufficient fees
func ErrInsufficientFees(codespace sdk.CodespaceType, fees string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientFees,
		"failed. insufficient fees, needs %s", fees)
}

func ErrInsufficientCoins(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientCoins,
		"failed. insufficient coins: %s", msg)
}