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
)

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
