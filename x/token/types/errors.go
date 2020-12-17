package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CodeParamTokenUnknown						uint32 = 61000
	CodeInvalidDexList          				uint32 = 61001
	CodeInvalidBalanceNotEnough 				uint32 = 61002
	CodeInvalidHeight          					uint32 = 61003
	CodeInvalidAsset           					uint32 = 61004
	CodeInvalidCommon           				uint32 = 61005
	CodeBlockedRecipient        				uint32 = 61006
	CodeSendDisabled            				uint32 = 61007
	CodeSendCoinsFromAccountToModuleFailed		uint32 = 61008
	CodeUnrecognizedLockCoinsType				uint32 = 61009
	CodeFailedToUnlockAddress					uint32 = 61010
	CodeUnknownRequest							uint32 = 61011
	CodeInternal								uint32 = 61012
	CodeInvalidCoins							uint32 = 61013
	CodeInsufficientCoins						uint32 = 61014
	CodeUnauthorized							uint32 = 61015
	CodeInvalidPriceDigit   	    			uint32 = 61016
	CodeInvalidMinTradeSize     				uint32 = 61017
	CodeInvalidAddress							uint32 = 61018
	CodeGetConfirmOwnership						uint32 = 61019
	CodeUpdateLockedCoins						uint32 = 61020
	CodeUnknownQueryType						uint32 = 61021
	CodeUserInputSymbolIsEmpty					uint32 = 61022
	CodeInvalidOriginalSymbol					uint32 = 61023
	CodeNotAllowedOriginalSymbol				uint32 = 61023
	CodeWholeNameIsNotValid						uint32 = 61024
	CodeDescLenBiggerThanLimit					uint32 = 61025
	CodeNewDecFromStrFailed						uint32 = 61026
	CodeTotalSupplyOutOfRange					uint32 = 61027
	CodeAmountBiggerThanTotalSupplyUpperbound	uint32 = 61028
	CodeAmountIsNotValid						uint32 = 61029
	CodeMsgSymbolIsEmpty						uint32 = 61030
	CodeMintCoinsFailed							uint32 = 61031
	CodeSendCoinsFromModuleToAccountFailed		uint32 = 61032
	CodeBurnCoinsFailed							uint32 = 61033
	CodeConfirmOwnershipNotExistOrBlockTimeAfter	uint32 = 61034
	CodeWholeNameAndDescriptionIsNotModified 	uint32 = 61035
)

var (
	errInvalidDexList          = sdkerrors.Register(DefaultCodespace, CodeInvalidDexList, "invalid dex list")
	errInvalidBalanceNotEnough = sdkerrors.Register(DefaultCodespace, CodeInvalidBalanceNotEnough, "invalid balance not enough")
	errInvalidHeight           = sdkerrors.Register(DefaultCodespace, CodeInvalidHeight, "invalid height")
	errInvalidAsset            = sdkerrors.Register(DefaultCodespace, CodeInvalidAsset, "invalid asset")
	errInvalidCommon           = sdkerrors.Register(DefaultCodespace, CodeInvalidCommon, "invalid common")
	errBlockedRecipient        = sdkerrors.Register(DefaultCodespace, CodeBlockedRecipient, "blocked recipient")
	errSendDisabled            = sdkerrors.Register(DefaultCodespace, CodeSendDisabled, "send disabled")
	errCodeSendCoinsFromAccountToModuleFailed	= sdkerrors.Register(DefaultCodespace, CodeSendCoinsFromAccountToModuleFailed, "send to module account failed")
	errCodeUnrecognizedLockCoinsType			= sdkerrors.Register(DefaultCodespace, CodeUnrecognizedLockCoinsType, "unrecognized lock coins")
	errCodeFailedToUnlockAddress				= sdkerrors.Register(DefaultCodespace, CodeFailedToUnlockAddress, "unlock address failed")
	errCodeUnknownRequest						= sdkerrors.Register(DefaultCodespace, CodeUnknownRequest, "unlock address failed")
	errCodeInternal							= sdkerrors.Register(DefaultCodespace, CodeInternal, "err occur internal")
	errCodeInvalidCoins						 	= sdkerrors.Register(DefaultCodespace, CodeInvalidCoins, "invalid coins")
	errCodeInsufficientCoins					= sdkerrors.Register(DefaultCodespace, CodeInsufficientCoins, "insufficient coins")
	errCodeUnauthorized							= sdkerrors.Register(DefaultCodespace, CodeUnauthorized	, "code unauthorized")
	errCodeInvalidPriceDigit       				= sdkerrors.Register(DefaultCodespace, CodeInvalidPriceDigit, "invalid price digit")
	errCodeInvalidMinTradeSize     				= sdkerrors.Register(DefaultCodespace, CodeInvalidMinTradeSize, "invalid min trade size")
	errCodeInvalidAddress						= sdkerrors.Register(DefaultCodespace, CodeInvalidAddress, "invalid address")
	errCodeGetConfirmOwnership					= sdkerrors.Register(DefaultCodespace, CodeGetConfirmOwnership, "get confirm ownership failed")
	errCodeUpdateLockedCoins					= sdkerrors.Register(DefaultCodespace, CodeUpdateLockedCoins, "update locked coins failed")
	errCodeUnknownQueryType						= sdkerrors.Register(DefaultCodespace, CodeUnknownQueryType, "unknown query type")
	errCodeOriginalSymbolIsEmpty				= sdkerrors.Register(DefaultCodespace, CodeUserInputSymbolIsEmpty, "user input symbol is empty")
	errCodeNotAllowedOriginalSymbol				= sdkerrors.Register(DefaultCodespace, CodeNotAllowedOriginalSymbol, "not allowed original symbol")
	errCodeWholeNameIsNotValid					= sdkerrors.Register(DefaultCodespace, CodeWholeNameIsNotValid, "whole name is not valid")
	errCodeDescLenBiggerThanLimit				= sdkerrors.Register(DefaultCodespace, CodeDescLenBiggerThanLimit, "description len bigger than limit")
	errCodeNewDecFromStrFailed					= sdkerrors.Register(DefaultCodespace, CodeNewDecFromStrFailed, "new dec from string")
	errCodeTotalSupplyOutOfRange				= sdkerrors.Register(DefaultCodespace, CodeTotalSupplyOutOfRange, "total supply out of range")
	errCodeAmountBiggerThanTotalSupplyUpperbound= sdkerrors.Register(DefaultCodespace, CodeAmountBiggerThanTotalSupplyUpperbound, "amount bigger than total supply upperbound")
	errCodeAmountIsNotValid						= sdkerrors.Register(DefaultCodespace, CodeAmountIsNotValid, "amount is not valid")
	errCodeMsgSymbolIsEmpty						= sdkerrors.Register(DefaultCodespace, CodeMsgSymbolIsEmpty, "msg symbol is empty")
	errCodeMintCoinsFailed						= sdkerrors.Register(DefaultCodespace, CodeMintCoinsFailed, "mint coins failed")
	errCodeSendCoinsFromModuleToAccountFailed	= sdkerrors.Register(DefaultCodespace, CodeSendCoinsFromModuleToAccountFailed, "send coins from module to account failed")
	errCodeBurnCoinsFailed						= sdkerrors.Register(DefaultCodespace, CodeBurnCoinsFailed	, "burn coins failed")
	errCodeConfirmOwnershipNotExistOrBlockTimeAfter	= sdkerrors.Register(DefaultCodespace, CodeConfirmOwnershipNotExistOrBlockTimeAfter	, "confirm ownership not exist or blocktime after")
	errCodeWholeNameAndDescriptionIsNotModified = sdkerrors.Register(DefaultCodespace, CodeWholeNameAndDescriptionIsNotModified	, "whole name and description is not modified")
)

// ErrBlockedRecipient returns an error when a transfer is tried on a blocked recipient
func ErrBlockedRecipient(blockedAddr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errBlockedRecipient, "failed. %s is not allowed to receive transactions", blockedAddr)}
}

// ErrSendDisabled returns an error when the transaction sending is disabled in bank module
func ErrSendDisabled() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errSendDisabled, "failed. send transactions are currently disabled")}
}

func ErrInvalidDexList(message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidDexList, message)}
}

func ErrInvalidBalanceNotEnough(message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidBalanceNotEnough, message)}
}

func ErrInvalidHeight(h, ch, max int64) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidHeight, fmt.Sprintf("Height %d must be greater than current block height %d and less than %d + %d.", h, ch, ch, max))}
}

func ErrInvalidCommon(message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidCommon, message)}
}

func ErrSendCoinsFromAccountToModuleFailed(message string) sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeSendCoinsFromAccountToModuleFailed, message)}
}

func ErrUnrecognizedLockCoinsType(lockCoinsType int) sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUnrecognizedLockCoinsType, fmt.Sprintf("unrecognized lock coins type: %d", lockCoinsType))}
}

func ErrFailedToUnlockAddress(coins string, addr string) sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeFailedToUnlockAddress, fmt.Sprintf("failed to unlock <%s>. Address <%s>, coins locked <0>", coins, addr))}
}

func ErrUnknownRequest() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUnknownRequest, "unknown request: %s")}
}

func ErrInternal() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeInternal, "occur error internal")}
}

func ErrInvalidCoins() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeInvalidCoins, "unknown token")}
}

func ErrInsufficientCoins(coins string) sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeInsufficientCoins, fmt.Sprintf("insufficient coins(need %s)", coins))}
}

func ErrUnauthorized() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUnauthorized, "input address is not ownership address")}
}

func ErrInvalidAddress() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeInvalidAddress, "invalid address")}
}

func ErrGetConfirmOwnership() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeGetConfirmOwnership, "get confirm ownership info failed")}
}

func ErrUpdateLockedCoins() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUpdateLockedCoins, "get confirm ownership failed")}
}

func ErrUnknownQueryType() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUnknownQueryType, "unknown query type")}
}

func ErrUserInputSymbolIsEmpty() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeOriginalSymbolIsEmpty, "user intput symbol is empty")}
}

func ErrNotAllowedOriginalSymbol() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeNotAllowedOriginalSymbol, "not allowed original symbol")}
}

func ErrWholeNameIsNotValidl() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeWholeNameIsNotValid, "whole name is not valid")}
}

func ErrDescLenBiggerThanLimit() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeDescLenBiggerThanLimit, "description len bigger than limit")}
}

func ErrNewDecFromStrFailed() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeNewDecFromStrFailed, "total supply convert to decimal failed")}
}

func ErrTotalSupplyOutOfRange() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeTotalSupplyOutOfRange, "new dec from string")}
}

func ErrAmountBiggerThanTotalSupplyUpperbound() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeAmountBiggerThanTotalSupplyUpperbound, "amount bigger than total supply upperbound")}
}

func ErrAmountIsNotValid(amount string) sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeAmountIsNotValid, fmt.Sprintf("amount is not valid amount is %s", amount))}
}

func ErrMsgSymbolIsEmpty() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeMsgSymbolIsEmpty, "msg symbol is empty")}
}

func ErrMintCoinsFailed() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeMintCoinsFailed, "mint coins failed")}
}

func ErrSendCoinsFromModuleToAccountFailed() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeSendCoinsFromModuleToAccountFailed, "send coins from module to account failed")}
}

func ErrBurnCoinsFailed() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeBurnCoinsFailed, "burn coins failed")}
}

func ErrConfirmOwnershipNotExistOrBlockTimeAfter() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeConfirmOwnershipNotExistOrBlockTimeAfter, "confirm ownership not exist or blocktime after")}
}

func ErrWholeNameAndDescriptionIsNotModified() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeWholeNameAndDescriptionIsNotModified, "whole name and description is not modified")}
}