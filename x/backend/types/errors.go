package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// const uint32
const (
	DefaultCodespace = ModuleName

	CodeProductIsRequired             uint32 = 62000
	CodeAddressIsRequired             uint32 = 62001
	CodeOrderStatusMustBeOpenOrClosed uint32 = 62002
	CodeAddressAndProductRequired     uint32 = 62003
	CodeGetChainHeightFailed          uint32 = 62004
	CodeGetBlockTxHashesFailed        uint32 = 62005
	CodeOrderSideMustBuyOrSell        uint32 = 62006
	CodeProductDoesNotExist           uint32 = 62007
	CodeBackendPluginNotEnabled       uint32 = 62008
	CodeGoroutinePanic                uint32 = 62009
	CodeBackendModuleUnknownQueryType uint32 = 62010
	CodeGetCandlesFailed              uint32 = 62011
	CodeGetCandlesByMarketFailed      uint32 = 62012
	CodeGetTickerByProductsFailed     uint32 = 62013
	CodeParamNotCorrect               uint32 = 62014
	CodeNoKlinesFunctionFound         uint32 = 62015
	CodeMarketkeeperNotInitialized    uint32 = 62016
	CodeGetInvalidateGranularity      uint32 = 62017
	CodeGetInvalidTickerByProducts    uint32 = 62018
)

// invalid param side, must be buy or sell
func ErrOrderSideParamMustBuyOrSell(side string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeOrderSideMustBuyOrSell, fmt.Sprintf("Side should not be %s", side))}
}

// product is required
func ErrProductIsRequired() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeProductIsRequired, "invalid params: product is required")}
}

// address is required
func ErrAddressIsRequired() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeAddressIsRequired, "bad request: address is required")}
}

// product does not exist
func ErrProductDoesNotExist(product string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeProductDoesNotExist, fmt.Sprintf("product %s does not exist", product))}
}

func ErrBackendPluginNotEnabled() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeBackendPluginNotEnabled, "backend is not enabled")}
}

func ErrParamNotCorrect(size int, granularity int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeParamNotCorrect, fmt.Sprintf("parameter is not correct, size: %d, granularity: %d", size, granularity))}
}

func ErrNoKlinesFunctionFound() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeNoKlinesFunctionFound, "no klines constructor function found")}
}

func ErrMarketkeeperNotInitialized() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeMarketkeeperNotInitialized, "market keeper is not initialized properly")}
}

func ErrBackendModuleUnknownQueryType() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeBackendModuleUnknownQueryType, "backend module unknown query type")}
}

func ErrGetInvalidateGranularity(msg string, key string, field string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeGetInvalidateGranularity, fmt.Sprintf("server error: %s, key=%s, can not convert timestamp %s", msg, key, field))}
}

func ErrGetInvalidTickerByProducts(key string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeGetInvalidTickerByProducts, fmt.Sprintf("No value found for key: %s", key))}
}
