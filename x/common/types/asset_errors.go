package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	AssetCodespace sdk.CodespaceType = "asset"

	codeTotalSupplyExceeds        CodeType = 61001
	codeBadSymbolGeneration       CodeType = 61002
	codeBadCoinsMintage           CodeType = 61003
	codeBadCoinsSendingToModule   CodeType = 61004
	codeBadCoinsSendingFromModule CodeType = 61005
	codeInsufficientFees          CodeType = 61006
	codeUnauthorizedIdentity      CodeType = 61007
	codeBadCoinsBurning           CodeType = 61008
	codeCoinsNotMintable          CodeType = 61009
	codeInsufficientBalance       CodeType = 61010
	codeInvalidModificationSet    CodeType = 61011
	codeEmptyOriginalSymbol       CodeType = 61012
	codeInvalidOriginalSymbol     CodeType = 61013
	codeInvalidWholeName          CodeType = 61014
	codeTokenDescriptionExceeds   CodeType = 61015
	codeTransfersLengthExceeds    CodeType = 61016
	codeInvalidMultisignCheck     CodeType = 61017
	codeInvalidCoins              CodeType = 61018
	codeMintageAmountExceeds      CodeType = 61019
	codeEmptySymbol               CodeType = 61020
	codeInvalidSymbol             CodeType = 61021
)

// ErrInvalidMultisignCheck returns an error with an invalid check result of multi-sign
func ErrInvalidMultisignCheck(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeInvalidMultisignCheck,
		"failed. invalid check result of multi-sign")
}

// ErrInvalidSymbol returns an error with an invalid symbol of token
func ErrInvalidSymbol(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, codeInvalidSymbol,
		"failed. invalid symbol of token: %s", symbol)
}

// ErrEmptySymbol returns an error with an empty symbol of token
func ErrEmptySymbol(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeEmptySymbol,
		"failed. empty symbol of token")
}

// ErrTransfersLengthExceeds returns an error when the transfers' length of multi-send exceeds the limit
func ErrTransfersLengthExceeds(codespace sdk.CodespaceType, limit int64) sdk.Error {
	return sdk.NewError(codespace, codeTransfersLengthExceeds,
		"failed. the length of transfers in multi-send exceeds the limit: %d", limit)
}

// ErrMintageAmountExceeds returns an error when the amount of mintage exceeds the limit
func ErrMintageAmountExceeds(codespace sdk.CodespaceType, limit int64) sdk.Error {
	return sdk.NewError(codespace, codeMintageAmountExceeds,
		"failed. the amount of mintage exceeds the limit: %d", limit)
}

// ErrInvalidCoins returns an error with invalid coins
func ErrInvalidCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeInvalidCoins,
		"failed. invalid coins with negative amount or illegal denomination")
}

// ErrTokenDescriptionExceeds returns an error when the token description exceeds the length limit
func ErrTokenDescriptionExceeds(codespace sdk.CodespaceType, lenLimit int) sdk.Error {
	return sdk.NewError(codespace, codeTokenDescriptionExceeds,
		"failed. token description exceeds the length limit: %d", lenLimit)
}

// ErrInvalidWholeName returns an error with an invalid whole name of token
func ErrInvalidWholeName(codespace sdk.CodespaceType, wholeName string) sdk.Error {
	return sdk.NewError(codespace, codeInvalidWholeName,
		"failed. invalid whole name of token: %s", wholeName)
}

// ErrInvalidOriginalSymbol returns an error with an invalid original symbol of token
func ErrInvalidOriginalSymbol(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, codeInvalidOriginalSymbol,
		"failed. invalid original symbol of token: %s", symbol)
}

// ErrEmptyOriginalSymbol returns an error with an empty original symbol of token
func ErrEmptyOriginalSymbol(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeEmptyOriginalSymbol,
		"failed. empty original symbol of token")
}

// ErrInvalidModification returns an error when neither "IsWholeNameModified" nor "IsDescriptionModified" is true
func ErrInvalidModification(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, codeInvalidModificationSet,
		`failed. invalid set for token's modification: neither "IsWholeNameModified" nor "IsDescriptionModified" is true`)
}

// ErrInsufficientBalance returns an error when the balance of an account is insufficient
func ErrInsufficientBalance(codespace sdk.CodespaceType, expectedAmount string) sdk.Error {
	return sdk.NewError(codespace, codeInsufficientBalance,
		"failed. insufficient balance, needs %s", expectedAmount)
}

// ErrCoinsNotMintable returns an error with the mintage of coins which are not mintable
func ErrCoinsNotMintable(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, codeCoinsNotMintable,
		"failed. token %s is not mintable", symbol)
}

// ErrBadCoinsBurning returns an error with the bad coins burning
func ErrBadCoinsBurning(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, codeBadCoinsBurning,
		"failed. bad coins burning: %s", msg)
}

// ErrUnauthorizedIdentity returns an error with the unauthorized identity of the owner
func ErrUnauthorizedIdentity(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, codeUnauthorizedIdentity,
		"failed. not the owner of token: %s", symbol)
}

// ErrTotalSupplyExceeds returns an error when the token's total supply exceeds the upper limit
func ErrTotalSupplyExceeds(codespace sdk.CodespaceType, totalSupply string, upperLimit int64) sdk.Error {
	return sdk.NewError(codespace, codeTotalSupplyExceeds,
		"failed. total supply %s exceeds the upper limit %d", totalSupply, upperLimit)
}

// ErrBadSymbolGeneration returns an error with the bad unique symbol generation
func ErrBadSymbolGeneration(codespace sdk.CodespaceType, originalSymbol string) sdk.Error {
	return sdk.NewError(codespace, codeBadSymbolGeneration,
		"failed. bad unique symbol generation for token %s", originalSymbol)
}

// ErrBadCoinsMintage returns an error with the bad coin mintage
func ErrBadCoinsMintage(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, codeBadCoinsMintage,
		"failed. bad coins mintage: %s", msg)
}

// ErrBadCoinsSendingToModule returns an error with the bad coins sending to module account
func ErrBadCoinsSendingToModule(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, codeBadCoinsSendingToModule,
		"failed. bad coins sending to module account: %s", msg)
}

// ErrBadCoinsSendingFromModule returns an error with the bad coins sending from module account
func ErrBadCoinsSendingFromModule(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, codeBadCoinsSendingFromModule,
		"failed. bad coins sending from module account: %s", msg)
}

// ErrInsufficientFee returns an error with the insufficient fees
func ErrInsufficientFees(codespace sdk.CodespaceType, fees string) sdk.Error {
	return sdk.NewError(codespace, codeInsufficientFees,
		"failed. insufficient fees, needs %s", fees)
}
