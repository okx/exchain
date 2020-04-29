package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	SwapCodespace sdk.CodespaceType = "swap"

	CodeQuoteOnlySupportsNativeToken = 63001
	CodeSwapTokenPairIsExist         = 63002
	CodeBlockTimeExceededDeadline    = 63003
	CodeGetPoolTokenFailed           = 63004
	CodeUnexpectedTotalSupply        = 63005
	CodeInvalidSwapTokenPair         = 63006
	CodeMaxBaseAmountIsTooLow        = 63007
	CodeMinLiquidityIsTooHigh        = 63008
	CodeFailToMintPoolCoins          = 63009
	CodeInsufficientPoolToken        = 63010
	CodeMinBaseAmountIsTooHigh       = 63011
	CodeMinQuoteAmountIsTooHigh      = 63012
	CodeFailToBurnPoolCoins          = 63013
	CodeMinBountTokenAmountIsTooHigh = 63014
)

// nolint
func ErrQuoteOnlySupportsNativeToken(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeQuoteOnlySupportsNativeToken, "failed. quote token only supports okt")
}

// nolint
func ErrSwapTokenPairIsExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSwapTokenPairIsExist, "failed. swapTokenPair already exists")
}

// nolint
func ErrBlockTimeExceededDeadline(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeBlockTimeExceededDeadline, "failed. blockTime exceeded deadline")
}

// nolint
func ErrGetPoolTokenFailed(codespace sdk.CodespaceType, poolToken, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeGetPoolTokenFailed, "failed. get poolToken %s failed: %s", poolToken, msg)
}

// nolint
func ErrUnexpectedTotalSupply(codespace sdk.CodespaceType, poolToken string) sdk.Error {
	return sdk.NewError(codespace, CodeUnexpectedTotalSupply, "failed. unexpected totalSupply in poolToken %s", poolToken)
}

// nolint
func ErrInvalidSwapTokenPair(codespace sdk.CodespaceType, tokenPair string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSwapTokenPair, "failed. invalid swapTokenPair %s", tokenPair)
}

// nolint
func ErrMaxBaseAmountIsTooLow(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMaxBaseAmountIsTooLow, "failed. MaxBaseAmount is too low")
}

// nolint
func ErrMinLiquidityIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinLiquidityIsTooHigh, "failed. MinLiquidity is too high")
}

// nolint
func ErrFailToMintPoolCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeFailToMintPoolCoins, "failed. fail to mint poolCoins")
}

// nolint
func ErrInsufficientPoolToken(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientPoolToken, "failed. insufficient poolToken")
}

// nolint
func ErrMinBaseAmountIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinBaseAmountIsTooHigh, "failed. MinBaseAmount is too high")
}

// nolint
func ErrMinQuoteAmountIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinQuoteAmountIsTooHigh, "failed. MinQuoteAmount is too high")
}

// nolint
func ErrFailToBurnPoolCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeFailToBurnPoolCoins, "failed. fail to burn poolCoins")
}

// nolint
func ErrMinBountTokenAmountIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMinBountTokenAmountIsTooHigh, "failed. MinBoughtTokenAmount is too high")
}
