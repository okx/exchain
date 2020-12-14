package types

import (
	"fmt"
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
	CodeGetSwapTokenPairFailed				uint32 = 65015
	CodeInvalidAddress						uint32 = 65016
)

func ErrUnexistswapTokenPair () sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnexistwapTokenPair, "validator address is nil")
}

func ErrUnexistPoolToken() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnexistPoolToken, "pool token does not exist")
}

func ErrCodeMinCoinsFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMintCoinsFailed, "min coins failed")
}

func ErrSendCoinsFromAccountToModule() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromAccountToModule, "send coins from account to module failed")
}

func ErrBaseAmountNameBigerQuoteAmountName() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBaseAmountNameBigerQuoteAmountName, "base amount name biger quote amount name")
}

func ErrValidateSwapAmountName() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidateSwapAmountName, "validate swap amount name failed")
}

func ErrBaseAmountNameEqualQuoteAmountName() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBaseAmountNameEqualQuoteAmountName, "base amount name equal quote amount name")
}

func ErrValidateDenom() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeValidateDenom, "validate denom failed")
}

func ErrNotAllowedOriginSymbol() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNotAllowedOriginSymbol, "liquidity-pool-token is not allowed to be a base or quote token")
}

func ErrInsufficientPoolToken () sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientPoolToken, "insufficient pool token")
}

func ErrUnknownRequest () sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, "unknown request")
}

func ErrTokenNotExist () sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenNotExist, "token does not exist")
}

func ErrInvalidCoins () sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoins, "failed to create exchange with pool token")
}

func ErrInternal() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternalError, "internal error")
}

func ErrInsufficientCoins() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientCoins, "insufficient coins")
}

func ErrGetSwapTokenPair() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetSwapTokenPairFailed, "get swap token pair failed")
}

func ErrInvalidAddress(address string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("invalid address: %s", address))
}