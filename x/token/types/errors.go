package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"time"
)

const (
	CodeInvalidAsset                               uint32 = 61000
	CodeBlockedRecipient                           uint32 = 61001
	CodeSendDisabled                               uint32 = 61002
	CodeSendCoinsFromAccountToModuleFailed         uint32 = 61003
	CodeUnrecognizedLockCoinsType                  uint32 = 61004
	CodeFailedToUnlockAddress                      uint32 = 61005
	CodeInvalidCoins                               uint32 = 61006
	CodeInvalidPriceDigit                          uint32 = 61007
	CodeInvalidMinTradeSize                        uint32 = 61008
	CodeAddressIsRequired                          uint32 = 61009
	CodeGetConfirmOwnership                        uint32 = 61010
	CodeUpdateLockedCoins                          uint32 = 61011
	CodeUnknownTokenQueryType                      uint32 = 61012
	CodeUserInputSymbolIsEmpty                     uint32 = 61013
	CodeNotAllowedOriginalSymbol                   uint32 = 61014
	CodeWholeNameIsNotValid                        uint32 = 61015
	CodeDescLenBiggerThanLimit                     uint32 = 61016
	CodeTotalSupplyOutOfRange                      uint32 = 61017
	CodeAmountBiggerThanTotalSupplyUpperbound      uint32 = 61018
	CodeAmountIsNotValid                           uint32 = 61019
	CodeMsgSymbolIsEmpty                           uint32 = 61020
	CodeMintCoinsFailed                            uint32 = 61021
	CodeSendCoinsFromModuleToAccountFailed         uint32 = 61022
	CodeBurnCoinsFailed                            uint32 = 61023
	CodeConfirmOwnershipNotExistOrBlockTimeAfter   uint32 = 61024
	CodeWholeNameAndDescriptionIsNotModified       uint32 = 61025
	CodeTokenIsNotMintable                         uint32 = 61026
	CodeMsgTransfersAmountBiggerThanSendLimit      uint32 = 61027
	CodeInputOwnerIsNotEqualTokenOwner             uint32 = 61028
	CodeinputFromAddressIsNotEqualTokenInfoOwner   uint32 = 61029
	CodeConfirmOwnershipAddressNotEqualsMsgAddress uint32 = 61030
	CodeGetDecimalFromDecimalStringFailed          uint32 = 61031
	CodeTotalsupplyExceedsTheUpperLimit            uint32 = 61032
)

var (
	errInvalidAsset                                   = sdkerrors.Register(DefaultCodespace, CodeInvalidAsset, "invalid asset")
	errBlockedRecipient                               = sdkerrors.Register(DefaultCodespace, CodeBlockedRecipient, "blocked recipient")
	errSendDisabled                                   = sdkerrors.Register(DefaultCodespace, CodeSendDisabled, "send disabled")
	errCodeSendCoinsFromAccountToModuleFailed         = sdkerrors.Register(DefaultCodespace, CodeSendCoinsFromAccountToModuleFailed, "send to module account failed")
	errCodeUnrecognizedLockCoinsType                  = sdkerrors.Register(DefaultCodespace, CodeUnrecognizedLockCoinsType, "unrecognized lock coins")
	errCodeFailedToUnlockAddress                      = sdkerrors.Register(DefaultCodespace, CodeFailedToUnlockAddress, "unlock address failed")
	errCodeInvalidCoins                               = sdkerrors.Register(DefaultCodespace, CodeInvalidCoins, "invalid coins")
	errCodeInvalidPriceDigit                          = sdkerrors.Register(DefaultCodespace, CodeInvalidPriceDigit, "invalid price digit")
	errCodeInvalidMinTradeSize                        = sdkerrors.Register(DefaultCodespace, CodeInvalidMinTradeSize, "invalid min trade size")
	errCodeAddressIsRequired                          = sdkerrors.Register(DefaultCodespace, CodeAddressIsRequired, "address is required")
	errCodeGetConfirmOwnership                        = sdkerrors.Register(DefaultCodespace, CodeGetConfirmOwnership, "get confirm ownership failed")
	errCodeUpdateLockedCoins                          = sdkerrors.Register(DefaultCodespace, CodeUpdateLockedCoins, "update locked coins failed")
	errCodeUnknownTokenQueryType                      = sdkerrors.Register(DefaultCodespace, CodeUnknownTokenQueryType, "unknown token query type")
	errCodeOriginalSymbolIsEmpty                      = sdkerrors.Register(DefaultCodespace, CodeUserInputSymbolIsEmpty, "user input symbol is empty")
	errCodeNotAllowedOriginalSymbol                   = sdkerrors.Register(DefaultCodespace, CodeNotAllowedOriginalSymbol, "not allowed original symbol")
	errCodeWholeNameIsNotValid                        = sdkerrors.Register(DefaultCodespace, CodeWholeNameIsNotValid, "whole name is not valid")
	errCodeDescLenBiggerThanLimit                     = sdkerrors.Register(DefaultCodespace, CodeDescLenBiggerThanLimit, "description len bigger than limit")
	errCodeTotalSupplyOutOfRange                      = sdkerrors.Register(DefaultCodespace, CodeTotalSupplyOutOfRange, "total supply out of range")
	errCodeAmountBiggerThanTotalSupplyUpperbound      = sdkerrors.Register(DefaultCodespace, CodeAmountBiggerThanTotalSupplyUpperbound, "amount bigger than total supply upperbound")
	errCodeAmountIsNotValid                           = sdkerrors.Register(DefaultCodespace, CodeAmountIsNotValid, "amount is not valid")
	errCodeMsgSymbolIsEmpty                           = sdkerrors.Register(DefaultCodespace, CodeMsgSymbolIsEmpty, "msg symbol is empty")
	errCodeMintCoinsFailed                            = sdkerrors.Register(DefaultCodespace, CodeMintCoinsFailed, "mint coins failed")
	errCodeSendCoinsFromModuleToAccountFailed         = sdkerrors.Register(DefaultCodespace, CodeSendCoinsFromModuleToAccountFailed, "send coins from module to account failed")
	errCodeBurnCoinsFailed                            = sdkerrors.Register(DefaultCodespace, CodeBurnCoinsFailed, "burn coins failed")
	errCodeConfirmOwnershipNotExistOrBlockTimeAfter   = sdkerrors.Register(DefaultCodespace, CodeConfirmOwnershipNotExistOrBlockTimeAfter, "confirm ownership not exist or blocktime after")
	errCodeWholeNameAndDescriptionIsNotModified       = sdkerrors.Register(DefaultCodespace, CodeWholeNameAndDescriptionIsNotModified, "whole name and description is not modified")
	errCodeTokenIsNotMintable                         = sdkerrors.Register(DefaultCodespace, CodeTokenIsNotMintable, "token is not mintable")
	errCodeMsgTransfersAmountBiggerThanSendLimit      = sdkerrors.Register(DefaultCodespace, CodeMsgTransfersAmountBiggerThanSendLimit, "use transfer amount bigger than send limit")
	errCodeInputOwnerIsNotEqualTokenOwner             = sdkerrors.Register(DefaultCodespace, CodeInputOwnerIsNotEqualTokenOwner, "input owner is not equal token owner")
	errCodeinputFromAddressIsNotEqualTokenInfoOwner   = sdkerrors.Register(DefaultCodespace, CodeinputFromAddressIsNotEqualTokenInfoOwner, "input from address is not equal token owner")
	errCodeConfirmOwnershipAddressNotEqualsMsgAddress = sdkerrors.Register(DefaultCodespace, CodeConfirmOwnershipAddressNotEqualsMsgAddress, "input address is not equal confirm ownership address")
	errCodeGetDecimalFromDecimalStringFailed          = sdkerrors.Register(DefaultCodespace, CodeGetDecimalFromDecimalStringFailed, "create a decimal from an input decimal string failed")
	errCodeTotalsupplyExceedsTheUpperLimit            = sdkerrors.Register(DefaultCodespace, CodeTotalsupplyExceedsTheUpperLimit, "total-supply exceeds the upper limit")
)

// ErrBlockedRecipient returns an error when a transfer is tried on a blocked recipient
func ErrBlockedRecipient(blockedAddr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errBlockedRecipient, "failed. %s is not allowed to receive transactions", blockedAddr)}
}

// ErrSendDisabled returns an error when the transaction sending is disabled in bank module
func ErrSendDisabled() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errSendDisabled, "failed. send transactions are currently disabled")}
}

func ErrSendCoinsFromAccountToModuleFailed(message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeSendCoinsFromAccountToModuleFailed, "failed to send to module account: %s", message)}
}

func ErrUnrecognizedLockCoinsType(lockCoinsType int) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUnrecognizedLockCoinsType, fmt.Sprintf("unrecognized lock coins type: %d", lockCoinsType))}
}

func ErrFailedToUnlockAddress(coins string, addr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeFailedToUnlockAddress, fmt.Sprintf("failed to unlock <%s>. Address <%s>, coins locked <0>", coins, addr))}
}

func ErrInvalidCoins() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeInvalidCoins, "invalid coins")}
}

func ErrGetConfirmOwnership() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeGetConfirmOwnership, "get confirm ownership info failed")}
}

func ErrUpdateLockedCoins() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUpdateLockedCoins, "update locked coins failed")}
}

func ErrUnknownTokenQueryType() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeUnknownTokenQueryType, "unknown token query type")}
}

func ErrUserInputSymbolIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeOriginalSymbolIsEmpty, "user intput symbol is empty")}
}

func ErrNotAllowedOriginalSymbol() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeNotAllowedOriginalSymbol, "not allowed original symbol")}
}

func ErrWholeNameIsNotValidl() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeWholeNameIsNotValid, "whole name is not valid")}
}

func ErrDescLenBiggerThanLimit() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeDescLenBiggerThanLimit, "description len bigger than limit")}
}

func ErrTotalSupplyOutOfRange() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeTotalSupplyOutOfRange, "total supply out of range")}
}

func ErrAmountBiggerThanTotalSupplyUpperbound() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeAmountBiggerThanTotalSupplyUpperbound, "amount bigger than total supply upperbound")}
}

func ErrAmountIsNotValid(amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeAmountIsNotValid, fmt.Sprintf("amount %s is not valid", amount))}
}

func ErrMsgSymbolIsEmpty() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeMsgSymbolIsEmpty, "msg symbol is empty")}
}

func ErrMintCoinsFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeMintCoinsFailed, fmt.Sprintf("mint coins failed: %s", msg))}
}

func ErrSendCoinsFromModuleToAccountFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeSendCoinsFromModuleToAccountFailed, fmt.Sprintf("send coins from module to account failed: %s", msg))}
}

func ErrBurnCoinsFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeBurnCoinsFailed, fmt.Sprintf("burn coins failed: %s", msg))}
}

func ErrConfirmOwnershipNotExistOrBlockTimeAfter(expire time.Time) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeConfirmOwnershipNotExistOrBlockTimeAfter, "confirm ownership not exist or blocktime after %s", expire)}
}

func ErrWholeNameAndDescriptionIsNotModified() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeWholeNameAndDescriptionIsNotModified, "whole name and description is not modified")}
}

func ErrTokenIsNotMintable() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeTokenIsNotMintable, "token is not mintable")}
}

func ErrMsgTransfersAmountBiggerThanSendLimit() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeMsgTransfersAmountBiggerThanSendLimit, "use transfer amount bigger than send limit")}
}

func ErrInputOwnerIsNotEqualTokenOwner(address sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeInputOwnerIsNotEqualTokenOwner, fmt.Sprintf("input owner(%s) is not equal token owner", address))}
}

func ErrCodeinputFromAddressIsNotEqualTokenInfoOwner(address sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeinputFromAddressIsNotEqualTokenInfoOwner, fmt.Sprintf("input from address is not equal token owner: %s", address))}
}

func ErrCodeConfirmOwnershipAddressNotEqualsMsgAddress(address sdk.AccAddress) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeConfirmOwnershipAddressNotEqualsMsgAddress, fmt.Sprintf("input address (%s) is not equal confirm ownership address", address))}
}

func ErrAddressIsRequired() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeAddressIsRequired, "address is required")}
}

func ErrGetDecimalFromDecimalStringFailed(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeGetDecimalFromDecimalStringFailed, fmt.Sprintf("create a decimal from an input decimal string failed: %", msg))}
}

func ErrCodeTotalsupplyExceedsTheUpperLimit(totalSupplyAfterMint sdk.Dec, TotalSupplyUpperbound int64) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errCodeTotalsupplyExceedsTheUpperLimit, fmt.Sprintf("total-supply(%s) exceeds the upper limit(%d)", totalSupplyAfterMint, TotalSupplyUpperbound))}
}
