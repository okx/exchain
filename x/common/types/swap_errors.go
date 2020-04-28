package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const
const (
	SwapCodespace sdk.CodespaceType = "swap"

)

// nolint
func ErrQuoteOnlySupportsNativeToken(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. quote token only supports okt")
}

// nolint
func ErrSwapTokenPairIsExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. swapTokenPair already exists")
}

// nolint
func ErrBlockTimeExceededDeadline(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. blockTime exceeded deadline")
}

// nolint
func ErrGetPoolTokenFailed(codespace sdk.CodespaceType, poolToken, msg string) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. get poolToken %s failed: %s", poolToken, msg)
}

// nolint
func ErrUnexpectedTotalSupply(codespace sdk.CodespaceType, poolToken string) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. unexpected totalSupply in poolToken %s", poolToken)
}

// nolint
func ErrInvalidSwapTokenPair(codespace sdk.CodespaceType, tokenPair string) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. invalid swapTokenPair %s", tokenPair)
}

// nolint
func ErrMaxBaseAmountIsTooLow(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. MaxBaseAmount is too low")
}

// nolint
func ErrMinLiquidityIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. MinLiquidity is too high")
}

// nolint
func ErrFailToMintPoolCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. fail to mint poolCoins")
}

// nolint
func ErrInsufficientPoolToken(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. insufficient poolToken")
}

// nolint
func ErrMinBaseAmountIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. MinBaseAmount is too high")
}
// nolint
func ErrMinQuoteAmountIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. MinQuoteAmount is too high")
}

// nolint
func ErrFailToBurnPoolCoins(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. fail to burn poolCoins")
}

// nolint
func ErrMinBountTokenAmountIsTooHigh(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, 63001, "failed. MinBoughtTokenAmount is too high")
}