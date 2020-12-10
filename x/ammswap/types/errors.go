package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace 						string = "ammswap"

	CodeUnexistwapTokenPair       			uint32 = 65000
	CodeUnexistPoolToken					uint32 = 65001
	CodeMintCoinsFailed						uint32 = 65002
	CodeSendCoinsFromAccountToModule 		uint32 = 65003
	CodeBaseAmountNameBigerQuoteAmountName	uint32 = 65004
	CodeValidateSwapAmountName				uint32 = 65005
	CodeBaseAmountNameEqualQuoteAmountName	uint32 = 65006
	CodeValidateDenom						uint32 = 65007
	CodeNotAllowedOriginSymbol				uint32 = 65008
	CodeInsufficientPoolToken				uint32 = 65009
	CodeUnknownRequest						uint32 = 65010
	CodeInsufficientCoins					uint32 = 65011
	CodeTokenNotExist						uint32 = 65012
	CodeInvalidCoins						uint32 = 65013
	CodeInternalError						uint32 = 65014
)

func ErrUnexistswapTokenPair(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeUnexistwapTokenPair, "validator address is nil")
}

func ErrUnexistPoolToken(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeUnexistPoolToken, message)
}

func ErrCodeMinCoinsFailed(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeMintCoinsFailed, message)
}

func ErrSendCoinsFromAccountToModule(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeSendCoinsFromAccountToModule, message)
}

func ErrBaseAmountNameBigerQuoteAmountName(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeBaseAmountNameBigerQuoteAmountName, message)
}

func ErrValidateSwapAmountName(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeValidateSwapAmountName, message)
}

func ErrBaseAmountNameEqualQuoteAmountName(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeBaseAmountNameEqualQuoteAmountName, message)
}

func ErrValidateDenom(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeValidateDenom, message)
}

func ErrNotAllowedOriginSymbol(codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeNotAllowedOriginSymbol, message)
}

func ErrInsufficientPoolToken (codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeInsufficientPoolToken, message)
}

func ErrUnknownRequest (codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeUnknownRequest, message)
}

func ErrTokenNotExist (codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeTokenNotExist, message)
}

func ErrInvalidCoins (codespace string, message string) sdk.Error {
	return sdkerrors.New(codespace, CodeInvalidCoins, message)
}

func ErrInternalError(codespace string, message string) (sdk.Error) {
	return sdkerrors.New(codespace, CodeInternalError, message)
}

func ErrInsufficientCoins(codespace string, message string) (sdk.Error) {
	return sdkerrors.New(codespace, CodeInsufficientCoins, message)
}