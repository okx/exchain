package types

import (
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace = ModuleName

	CodeNonExistSwapTokenPair                   uint32 = 65000
	CodeNonExistPoolToken                       uint32 = 65001
	CodeMintCoinsFailed                         uint32 = 65002
	CodeSendCoinsFromAccountToModule            uint32 = 65003
	CodeBaseAmountNameBiggerThanQuoteAmountName uint32 = 65004
	CodeValidateSwapAmountName                  uint32 = 65005
	CodeBaseNameEqualQuoteName                  uint32 = 65006
	CodeValidateDenom                           uint32 = 65007
	CodeNotAllowedOriginSymbol                  uint32 = 65008
	CodeInsufficientPoolToken                   uint32 = 65009
	CodeTokenNotExist                           uint32 = 65010
	CodeInvalidCoins                            uint32 = 65011
	CodeInvalidTokenPair                        uint32 = 65012
	CodeAddressIsRequire                        uint32 = 65013
	CodeIsZeroValue                             uint32 = 65014
	CodeBlockTimeBigThanDeadline                uint32 = 65015
	CodeLessThan                                uint32 = 65016
	CodeMintPoolTokenFailed                     uint32 = 65017
	CodeSendCoinsFromPoolToAccountFailed        uint32 = 65018
	CodeBurnPoolTokenFailed                     uint32 = 65019
	CodeSendCoinsToPoolFailed                   uint32 = 65020
	CodeSwapUnknownMsgType                      uint32 = 65021
	CodeSwapUnknownQueryTypes                   uint32 = 65022
	CodeSellAmountOrBuyTokenIsEmpty             uint32 = 65023
	CodeSellAmountEqualBuyToken                 uint32 = 65024
	CodeQueryParamsAddressIsEmpty               uint32 = 65025
	CodeQueryParamsQuoteTokenAmountIsEmpty      uint32 = 65026
	CodeQueryParamsBaseTokenIsEmpty             uint32 = 65027
	CodeMinLiquidityIsNegative                  uint32 = 65028
	CodeMaxBaseAmountOrQuoteAmountIsNegative    uint32 = 65029
	CodeMaxBaseAmount                           uint32 = 65030
	CodeQuoteAmount                             uint32 = 65031
	CodeMinBaseAmount                           uint32 = 65032
	CodeMinQuoteAmount                          uint32 = 65033
	CodeSoldTokenAmountIsNegative               uint32 = 65034
	CodeToken0NameEqualToken1Name               uint32 = 65035
	CodeSoldTokenAmount                         uint32 = 65036
	CodeMinBoughtTokenAmount                    uint32 = 65037
	CodeConvertSellTokenAmount                  uint32 = 65038
	CodeConvertQuoteTokenAmount                 uint32 = 65039
	CodeSendCoinsFailed                         uint32 = 65040
	CodeMsgDeadlineLessThanBlockTime            uint32 = 65041
	CodeBaseTokensAmountBiggerThanMax           uint32 = 65042
	CodeIsSwapTokenPairExist                    uint32 = 65043
	CodeIsPoolTokenPairExist                    uint32 = 65044
	CodeInternalError                           uint32 = 65045
)

func ErrNonExistSwapTokenPair(tokenPairName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeNonExistSwapTokenPair, fmt.Sprintf("swap token pair is not exist: %s", tokenPairName))}
}

func ErrNonExistPoolToken(token string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeNonExistPoolToken, fmt.Sprintf("pool token %s does not exist", token))}
}

func ErrCodeMinCoinsFailed(err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMintCoinsFailed, fmt.Sprintf("mint coins failed: %s", err))}
}

func ErrSendCoinsFromAccountToModule(err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSendCoinsFromAccountToModule, fmt.Sprintf("send coins from account to module failed: %s", err))}
}

func ErrBaseAmountNameBiggerThanQuoteAmountName() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeBaseAmountNameBiggerThanQuoteAmountName, "the lexicographic order of BaseTokenName must be less than QuoteTokenName")}
}

func ErrValidateSwapAmountName() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeValidateSwapAmountName, "validate swap amount name failed")}
}

func ErrBaseNameEqualQuoteName() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeBaseNameEqualQuoteName, "base token name equal token name")}
}

func ErrValidateDenom(tokenName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeValidateDenom, fmt.Sprintf("invalid token name: %s", tokenName))}
}

func ErrNotAllowedOriginSymbol() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeNotAllowedOriginSymbol, fmt.Sprintf("liquidity-pool-token(with prefix %s is not allowed to be a base or quote token", PoolTokenPrefix))}
}

func ErrInsufficientPoolToken() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInsufficientPoolToken, "insufficient pool token")}
}

func ErrTokenNotExist() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeTokenNotExist, "token does not exist")}
}

func ErrInvalidCoins() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidCoins, "failed to create exchange with equal token name")}
}

func ErrInvalidTokenPair(swapTokenPair string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeInvalidTokenPair, fmt.Sprintf("invalid token pair %s", swapTokenPair))}
}

func ErrAddressIsRequire(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeAddressIsRequire, fmt.Sprintf("%s address is require", msg))}
}

func ErrIsZeroValue(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeIsZeroValue, fmt.Sprintf("%s is zero value", msg))}
}

func ErrBlockTimeBigThanDeadline() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeBlockTimeBigThanDeadline, fmt.Sprintf("block time big than deadline"))}
}

func ErrLessThan(param1 string, param2 string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeLessThan, fmt.Sprintf("%s value less than %s value", param1, param2))}
}

func ErrMintPoolTokenFailed(err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMintPoolTokenFailed, fmt.Sprintf("mint pool token failed: %s", err))}
}

func ErrSendCoinsFromPoolToAccountFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSendCoinsFromPoolToAccountFailed, fmt.Sprintf("send coins from pool to account failed: %s", msg))}
}

func ErrBurnPoolTokenFailed(err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeBurnPoolTokenFailed, fmt.Sprintf("burn pool token failed: %s", err))}
}

func ErrSendCoinsToPoolFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSendCoinsToPoolFailed, fmt.Sprintf("send coins to pool failed: %s", msg))}
}

func ErrSwapUnknownMsgType() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSwapUnknownMsgType, "swap unknown msg type")}
}

func ErrSwapUnknownQueryType() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSwapUnknownQueryTypes, "unknown swap query endpoint")}
}

func ErrSellAmountOrBuyTokenIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSellAmountOrBuyTokenIsEmpty, "sell token amount or buy token is empty")}
}

func ErrSellAmountEqualBuyToken() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSellAmountEqualBuyToken, "sell token name should not be equal to buy token name")}
}

func ErrQueryParamsAddressIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeQueryParamsAddressIsEmpty, "query param address is empty")}
}

func ErrQueryParamsQuoteTokenAmountIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeQueryParamsQuoteTokenAmountIsEmpty, "query param quote token amount is empty")}
}

func ErrQueryParamsBaseTokenIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeQueryParamsBaseTokenIsEmpty, "query param base token is empty")}
}

func ErrMinLiquidityIsNegative() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMinLiquidityIsNegative, "min liquidity is negative")}
}

func ErrMaxBaseAmountOrQuoteAmountIsNegative() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMaxBaseAmountOrQuoteAmountIsNegative, "max base amount or quote amount is negative")}
}

func ErrMaxBaseAmount() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMaxBaseAmount, "max base amount is negative or not validate denom")}
}

func ErrQuoteAmount() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeQuoteAmount, "quote amount is negative or not validate denom")}
}

func ErrMinBaseAmount() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMinBaseAmount, "min base amount is negative or not validate denom")}
}

func ErrMinQuoteAmount() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMinQuoteAmount, "min quote amount is negative or not validate denom")}
}

func ErrSoldTokenAmountIsNegative() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSoldTokenAmountIsNegative, "sold token amount is negative")}
}

func ErrToken0NameEqualToken1Name() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeToken0NameEqualToken1Name, "token0 name is equal token1 name")}
}

func ErrSoldTokenAmount() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSoldTokenAmount, "sold token amount is negative or not validate denom")}
}

func ErrMinBoughtTokenAmount() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMinBoughtTokenAmount, "min bought token amount is negative or not validate denom")}
}

func ErrConvertSellTokenAmount(sellTokenAmount string, err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeConvertSellTokenAmount, fmt.Sprintf("invalid params, parse sell_token_amount:%s error:%s",
		sellTokenAmount, err))}
}

func ErrConvertQuoteTokenAmount(quoteTokenAmount string, err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeConvertQuoteTokenAmount, fmt.Sprintf("invalid params, parse quote_token_amount:%s error:%s",
		quoteTokenAmount, err))}
}

func ErrSendCoinsFailed(err error) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeSendCoinsFailed, fmt.Sprintf("send coin failed: %s", err))}
}

func ErrMsgDeadlineLessThanBlockTime() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMsgDeadlineLessThanBlockTime, "input deadline less than block time")}
}

func ErrBaseTokensAmountBiggerThanMax() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeBaseTokensAmountBiggerThanMax, "base token amount bigger than max base amount")}
}

func ErrSwapTokenPairExist() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeIsSwapTokenPairExist, "the swap token pair already exists")}
}

func ErrPoolTokenPairExist() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeIsPoolTokenPairExist, "the pool token pair already exists")}
}
