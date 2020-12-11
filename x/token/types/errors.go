package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CodeParamTokenUnknown					uint32 = 61000
	CodeInvalidDexList          			uint32 = 61001
	CodeInvalidBalanceNotEnough 			uint32 = 61002
	CodeInvalidHeight          				uint32 = 61003
	CodeInvalidAsset           				uint32 = 61004
	CodeInvalidCommon           			uint32 = 61005
	CodeBlockedRecipient        			uint32 = 61006
	CodeSendDisabled            			uint32 = 61007
	CodeSendCoinsFromAccountToModuleFailed	uint32 = 61008
	CodeUnrecognizedLockCoinsType			uint32 = 61009
	CodeFailedToUnlockAddress				uint32 = 61010
	CodeUnknownRequest						uint32 = 61011
	CodeInternal							uint32 = 61012
	CodeInvalidCoins						uint32 = 61013
	CodeInsufficientCoins					uint32 = 61014
	CodeUnauthorized						uint32 = 61015
	CodeInvalidPriceDigit       			uint32 = 61016
	CodeInvalidMinTradeSize     			uint32 = 61017
	CodeInvalidAddress						uint32 = 61018
)

var (
	errInvalidDexList          = sdkerrors.Register(DefaultCodespace, CodeInvalidDexList, "invalid dex list")
	errInvalidBalanceNotEnough = sdkerrors.Register(DefaultCodespace, CodeInvalidBalanceNotEnough, "invalid balance not enough")
	errInvalidHeight           = sdkerrors.Register(DefaultCodespace, CodeInvalidHeight, "invalid height")
	errInvalidAsset            = sdkerrors.Register(DefaultCodespace, CodeInvalidAsset, "invalid asset")
	errInvalidCommon           = sdkerrors.Register(DefaultCodespace, CodeInvalidCommon, "invalid common")
	errBlockedRecipient        = sdkerrors.Register(DefaultCodespace, CodeBlockedRecipient, "blocked recipient")
	errSendDisabled            = sdkerrors.Register(DefaultCodespace, CodeSendDisabled, "send disabled")
	errCodeSendCoinsFromAccountToModuleFailed	= sdkerrors.Register(DefaultCodespace, CodeSendCoinsFromAccountToModuleFailed, "send failed")
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
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUnauthorized, "code unauthorized")}
}

func ErrInvalidAddress() sdk.Error {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeInvalidAddress, "invalid address")}
}