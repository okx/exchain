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
	CodeTokenNotExist						uint32 = 65012
	CodeInvalidCoins						uint32 = 65013
	CodeInternalError						uint32 = 65014
	CodeGetSwapTokenPairFailed				uint32 = 65015
	CodeInvalidAddress						uint32 = 65016
	CodeIsZeroValue							uint32 = 65017
	CodeGetRedeemableAssetsFailed			uint32 = 65018
	CodeBlockTimeBigThanDeadline			uint32 = 65019
	CodeGetPoolTokenInfoFailed				uint32 = 65020
	CodeTokenGreaterThanBaseAccount			uint32 = 65021
	CodeIsTokenExist						uint32 = 65022
	CodeMintPoolCoinsToUserFailed			uint32 = 65023
	CodeSendCoinsFromPoolToAccountFailed	uint32 = 65024
	CodeBurnPoolCoinsFromUserFailed			uint32 = 65025
	CodeCalculateTokenToBuyFailed			uint32 = 65026
	CodeSendCoinsToPoolFailed				uint32 = 65027
)

func ErrUnexistswapTokenPair() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnexistwapTokenPair, "token pair is not exist")
}

func ErrUnexistPoolToken() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnexistPoolToken, "pool token does not exist")
}

func ErrCodeMinCoinsFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMintCoinsFailed, "mint coins failed")
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

func ErrInsufficientPoolToken() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInsufficientPoolToken, "insufficient pool token")
}

func ErrUnknownRequest() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnknownRequest, "unknown request")
}

func ErrTokenNotExist() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenNotExist, "token does not exist")
}

func ErrInvalidCoins() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoins, "failed to create exchange with pool token")
}

func ErrInternalError() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternalError, "internal error")
}

func ErrGetSwapTokenPair() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetSwapTokenPairFailed, "get swap token pair failed")
}

func ErrInvalidAddress(address string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidAddress, fmt.Sprintf("invalid address: %s", address))
}

func ErrIsZeroValue() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsZeroValue, fmt.Sprintf("is zero value"))
}

func ErrGetRedeemableAssetsFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetRedeemableAssetsFailed, fmt.Sprintf("get redeemable assets failed"))
}

func ErrBlockTimeBigThanDeadline() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBlockTimeBigThanDeadline, fmt.Sprintf("block time big than deadline"))
}

func ErrGetPoolTokenInfoFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetPoolTokenInfoFailed, fmt.Sprintf("get pool token info failed"))
}

func ErrTokenGreaterThanBaseAccount() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenGreaterThanBaseAccount, fmt.Sprintf("token greater than base account"))
}

func ErrLiquidityLessThanMsg() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenGreaterThanBaseAccount, fmt.Sprintf("value less than msg"))
}

func ErrIsTokenExist(tokenName string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsTokenExist, fmt.Sprintf("%s is token exist", tokenName))
}

func ErrMintPoolCoinsToUserFailed(tokenName string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMintPoolCoinsToUserFailed, "mint pool coins to user failed")
}

func ErrSendCoinsFromPoolToAccountFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromPoolToAccountFailed, "send coins from pool to account failed")
}

func ErrBurnPoolCoinsFromUserFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBurnPoolCoinsFromUserFailed, "burn pool coins fro user failed")
}

func ErrCalculateTokenToBuyFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeCalculateTokenToBuyFailed, "calculate token to buy failed")
}

func ErrSendCoinsToPoolFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsToPoolFailed, "send coins to pool failed")
}
