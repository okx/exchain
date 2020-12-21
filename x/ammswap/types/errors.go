package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace = ModuleName

	CodeUnexistSwapTokenPair                             uint32 = 65000
	CodeUnexistPoolToken                                 uint32 = 65001
	CodeMintCoinsFailed                                  uint32 = 65002
	CodeSendCoinsFromAccountToModule                     uint32 = 65003
	CodeBaseAmountNameBigerQuoteAmountName               uint32 = 65004
	CodeValidateSwapAmountName                           uint32 = 65005
	CodeBaseAmountNameEqualQuoteAmountName               uint32 = 65006
	CodeValidateDenom                                    uint32 = 65007
	CodeNotAllowedOriginSymbol                           uint32 = 65008
	CodeInsufficientPoolToken                            uint32 = 65009
	CodeTokenNotExist                                    uint32 = 65012
	CodeInvalidCoins                                     uint32 = 65013
	CodeInternalError                                    uint32 = 65014
	CodeGetSwapTokenPairFailed                           uint32 = 65015
	AddressIsRequire                                     uint32 = 65016
	CodeIsZeroValue                                      uint32 = 65017
	CodeGetRedeemableAssetsFailed                        uint32 = 65018
	CodeBlockTimeBigThanDeadline                         uint32 = 65019
	CodeGetPoolTokenInfoFailed                           uint32 = 65020
	CodeTokenGreaterThanBaseAccount                      uint32 = 65021
	CodeLessThan                                         uint32 = 65022
	CodeMintPoolCoinsToUserFailed                        uint32 = 65024
	CodeSendCoinsFromPoolToAccountFailed                 uint32 = 65025
	CodeBurnPoolCoinsFromUserFailed                      uint32 = 65026
	CodeCalculateTokenToBuyFailed                        uint32 = 65027
	CodeSendCoinsToPoolFailed                            uint32 = 65028
	CodeSwapUnknownMsgType                               uint32 = 65029
	CodeSwapUnknowQueryTypes                             uint32 = 65030
	CodeSellAmountOrBuyTokenIsEmpty                      uint32 = 65031
	CodeSellAmountEqualBuyToken                          uint32 = 65032
	CodeQueryParamsAddressIsEmpty                        uint32 = 65033
	CodeQueryParamsQuoteTokenAmountIsEmpty               uint32 = 65034
	CodeQueryParamsBaseTokenIsEmpty                      uint32 = 65035
	CodeMinLiquidityIsNegative                           uint32 = 65036
	CodeMaxBaseAmountOrMsgQuoteAmountIsNegative          uint32 = 65037
	CodeMaxBaseAmountIsNegativeOrNotValidateDenom        uint32 = 65038
	CodeQuoteAmountIsNegativeOrNotValidateDenom          uint32 = 65039
	CodeMinBaseAmountIsNegativeOrNotValidateDenom        uint32 = 65040
	CodeMinQuoteAmountIsNegativeOrNotValidateDenom       uint32 = 65041
	CodeSoldTokenAmountIsNegative                        uint32 = 65042
	CodeToken0NameEqualToken1Name                        uint32 = 65043
	CodeSoldTokenAmountIsNegativeOrNotValidateDenom      uint32 = 65044
	CodeMinBoughtTokenAmountIsNegativeOrNotValidateDenom uint32 = 65045
	CodeConvertSellTokenAmountToDecimal                  uint32 = 65046
	CodeConvertQuoteTokenAmountToDecimal                 uint32 = 65047
	CodeSendCoinsFromAccountToModuleFailed               uint32 = 65048
	CodeMsgDeadlineLessThanBlockTime                     uint32 = 65049
	CodeBaseTokensAmountBiggerThanMaxBaseAmount          uint32 = 65050
	CodeIsSwapTokenPairExist                             uint32 = 65051
	CodeIsPoolTokenPairExist                             uint32 = 65052
)

func ErrSwapTokenPairExist() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsSwapTokenPairExist, "the swap token pair already exists")
}

func ErrPoolTokenPairExist() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsPoolTokenPairExist, "the pool token pair already exists")
}

func ErrUnexistSwapTokenPair(tokenPairName string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeUnexistSwapTokenPair, fmt.Sprintf("swap token pair is not exist: %s", tokenPairName))
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

func ErrTokenNotExist() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeTokenNotExist, "token does not exist")
}

func ErrInvalidCoins() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCoins, "failed to create exchange with equal token name")
}

func ErrInternalError() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInternalError, "internal error")
}

func ErrGetSwapTokenPair() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeGetSwapTokenPairFailed, "get swap token pair failed")
}

func ErrAddressIsRequire(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, AddressIsRequire, fmt.Sprintf("%s address Is require", msg))
}

func ErrIsZeroValue(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeIsZeroValue, fmt.Sprintf("%s is zero value", msg))
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

func ErrLessThan(param1 string, param2 string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeLessThan, fmt.Sprintf("%s value less than %s value", param1, param2))
}

func ErrMintPoolCoinsToUserFailed(tokenName string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMintPoolCoinsToUserFailed, "mint pool coins to user failed")
}

func ErrSendCoinsFromPoolToAccountFailed(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromPoolToAccountFailed, fmt.Sprintf("send coins from pool to account failed: %s", msg))
}

func ErrBurnPoolCoinsFromUserFailed(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBurnPoolCoinsFromUserFailed, fmt.Sprintf("burn pool coins fro user failed: %s", msg))
}

func ErrSendCoinsToPoolFailed(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsToPoolFailed, fmt.Sprintf("send coins to pool failed: %s", msg))
}

func ErrSwapUnknownMsgType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSwapUnknownMsgType, "swap unknown msg type")
}

func ErrSwapUnknownQueryType() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSwapUnknowQueryTypes, "swap unknown query type")
}

func ErrSellAmountOrBuyTokenIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSellAmountOrBuyTokenIsEmpty, "sell amount or buy amount token is empty")
}

func ErrSellAmountEqualBuyToken() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSellAmountEqualBuyToken, "sell amount equal buy token")
}

func ErrQueryParamsAddressIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeQueryParamsAddressIsEmpty, "query param address is empty")
}

func ErrQueryParamsQuoteTokenAmountIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeQueryParamsQuoteTokenAmountIsEmpty, "query param quote token amount is empty")
}

func ErrQueryParamsBaseTokenIsEmpty() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeQueryParamsBaseTokenIsEmpty, "query param base token is empty")
}

func ErrMinLiquidityIsNegative() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMinLiquidityIsNegative, "min liquidity is negative")
}

func ErrMaxBaseAmountOrMsgQuoteAmountIsNegative() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMaxBaseAmountOrMsgQuoteAmountIsNegative, "max base amount or msg quote amount is negative")
}

func ErrMaxBaseAmountIsNegativeOrNotValidateDenom() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMaxBaseAmountIsNegativeOrNotValidateDenom, "max base amount is negative or not validate denom")
}

func ErrQuoteAmountIsNegativeOrNotValidateDenom() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeQuoteAmountIsNegativeOrNotValidateDenom, "quote amount is negative or not validate denom")
}

func ErrMinBaseAmountIsNegativeOrNotValidateDenom() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMinBaseAmountIsNegativeOrNotValidateDenom, "min base amount is negative or not validate denom")
}

func ErrMinQuoteAmountIsNegativeOrNotValidateDenom() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMinQuoteAmountIsNegativeOrNotValidateDenom, "min quote amount is negative or not validate denom")
}

func ErrSoldTokenAmountIsNegative() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSoldTokenAmountIsNegative, "sold token amount is negative")
}

func ErrToken0NameEqualToken1Name() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeToken0NameEqualToken1Name, "token 0 name is equal token 1 name")
}

func ErrSoldTokenAmountIsNegativeOrNotValidateDenom() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSoldTokenAmountIsNegativeOrNotValidateDenom, "sold token amount is negative or not validate denom")
}

func ErrMinBoughtTokenAmountIsNegativeOrNotValidateDenom() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMinBoughtTokenAmountIsNegativeOrNotValidateDenom, "min bought token amount is negative or not validate denom")
}

func ErrConvertSellTokenAmountToDecimal() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeConvertSellTokenAmountToDecimal, "parse dec coin query params' sell token amount")
}

func ErrConvertQuoteTokenAmountToDecimal(msg string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeConvertQuoteTokenAmountToDecimal, "parse dec coin query params' quote token amount")
}

func ErrSendCoinsFromAccountToModuleFailed() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeSendCoinsFromAccountToModuleFailed, "send coin from input account to module account Failed")
}

func ErrMsgDeadlineLessThanBlockTime() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeMsgDeadlineLessThanBlockTime, "input deadline less than block time")
}

func ErrBaseTokensAmountBiggerThanMaxBaseAmountAmount() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeBaseTokensAmountBiggerThanMaxBaseAmount, "base token amount bigger than max base amount")
}
