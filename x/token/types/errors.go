package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CodeInvalidPriceDigit       uint32 = 1
	CodeInvalidMinTradeSize     uint32 = 2
	CodeInvalidDexList          uint32 = 3
	CodeInvalidBalanceNotEnough uint32 = 4
	CodeInvalidHeight           uint32 = 5
	CodeInvalidAsset            uint32 = 6
	CodeInvalidCommon           uint32 = 7
	CodeBlockedRecipient        uint32 = 8
	CodeSendDisabled            uint32 = 9
)

var (
	errInvalidDexList          = sdkerrors.Register(DefaultCodespace, CodeInvalidDexList, "invalid dex list")
	errInvalidBalanceNotEnough = sdkerrors.Register(DefaultCodespace, CodeInvalidBalanceNotEnough, "invalid balance not enough")
	errInvalidHeight           = sdkerrors.Register(DefaultCodespace, CodeInvalidHeight, "invalid height")
	errInvalidAsset            = sdkerrors.Register(DefaultCodespace, CodeInvalidAsset, "invalid asset")
	errInvalidCommon           = sdkerrors.Register(DefaultCodespace, CodeInvalidCommon, "invalid common")
	errBlockedRecipient        = sdkerrors.Register(DefaultCodespace, CodeBlockedRecipient, "blocked recipient")
	errSendDisabled            = sdkerrors.Register(DefaultCodespace, CodeSendDisabled, "send disabled")
)

// ErrBlockedRecipient returns an error when a transfer is tried on a blocked recipient
func ErrBlockedRecipient(codespace string, blockedAddr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errBlockedRecipient, "failed. %s is not allowed to receive transactions", blockedAddr)}
}

// ErrSendDisabled returns an error when the transaction sending is disabled in bank module
func ErrSendDisabled(codespace string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errSendDisabled, "failed. send transactions are currently disabled")}
}

func ErrInvalidDexList(codespace string, message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidDexList, message)}
}

func ErrInvalidBalanceNotEnough(codespace string, message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidBalanceNotEnough, message)}
}

func ErrInvalidHeight(codespace string, h, ch, max int64) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidHeight, fmt.Sprintf("Height %d must be greater than current block height %d and less than %d + %d.", h, ch, ch, max))}
}

func ErrInvalidCommon(codespace string, message string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidCommon, message)}
}
