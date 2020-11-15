package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = uint32

const (
	DefaultCodespace string = ModuleName

	CodeInvalidInput           CodeType = 101
	CodePoolAlreadyExist       CodeType = 102
	CodeNoFarmPoolFound        CodeType = 103
	CodePoolNotInWhiteList     CodeType = 104
	CodeInvalidLockInfo        CodeType = 105
	CodeTokenNotExist          CodeType = 106
	CodePoolNotFinished        CodeType = 107
	CodeUnexpectedProposalType CodeType = 108
	CodeInvalidAddress         CodeType = 109
	CodeUnknownRequest         CodeType = 110
)

var (
	errInvalidInput           = sdkerrors.Register(DefaultCodespace, CodeInvalidInput, "invalid input")
	errPoolAlreadyExist       = sdkerrors.Register(DefaultCodespace, CodePoolAlreadyExist, "pool already exist")
	errNoFarmPoolFound        = sdkerrors.Register(DefaultCodespace, CodeNoFarmPoolFound, "no farm pool found")
	errPoolNotInWhiteList     = sdkerrors.Register(DefaultCodespace, CodePoolNotInWhiteList, "pool not in white list")
	errInvalidLockInfo        = sdkerrors.Register(DefaultCodespace, CodeInvalidLockInfo, "invalid lock info")
	errTokenNotExist          = sdkerrors.Register(DefaultCodespace, CodeTokenNotExist, "token not exist")
	errPoolNotFinished        = sdkerrors.Register(DefaultCodespace, CodePoolNotFinished, "pool not finished")
	errUnexpectedProposalType = sdkerrors.Register(DefaultCodespace, CodeUnexpectedProposalType, "unexpected proposal type")
	errInvalidAddress         = sdkerrors.Register(DefaultCodespace, CodeInvalidAddress, "invalid address")
	errUnknownRequest         = sdkerrors.Register(DefaultCodespace, CodeUnknownRequest, "unknown request")
)

// ErrInvalidInput returns an error when an input parameter is invalid
func ErrInvalidInput(codespace string, msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,  "failed. invalid input: %s", msg)}
}

// ErrNoFarmPoolFound returns an error when a farm pool doesn't exist
func ErrNoFarmPoolFound(codespace string, poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errNoFarmPoolFound,  "failed. farm pool %s does not exist", poolName)}
}

// ErrPoolAlreadyExist returns an error when a pool exist
func ErrPoolAlreadyExist(codespace string, poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errPoolAlreadyExist,  "failed. farm pool %s already exists", poolName)}
}

// ErrTokenNotExist returns an error when a token not exists
func ErrTokenNotExist(codespace string, tokenName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errTokenNotExist,  "failed. token %s does not exist", tokenName)}
}

// ErrNoLockInfoFound returns an error when an address doesn't have any lock infos
func ErrNoLockInfoFound(codespace string, addr string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidLockInfo,  "failed. %s doesn't have any lock infos", addr)}
}

// ErrRemainingAmountNotZero returns an error when the remaining amount in yieldedTokenInfo is not zero
func ErrRemainingAmountNotZero(codespace string, amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, 
		"failed. the remaining amount is %s, so it's not enable to provide token repeatedly util amount become zero",
		amount)}
}

// ErrInvalidPoolOwner returns an error when an input address is not the owner of pool
func ErrInvalidPoolOwner(codespace string, address string, poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,  "failed. %s isn't the owner of pool %s", address, poolName)}
}

// ErrInvalidDenom returns an error when it provides an unmatched token name
func ErrInvalidDenom(codespace string, symbolLocked string, token string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, 
		"failed. the coin name should be %s, not %s", symbolLocked, token)}
}

// ErrInvalidInputAmount returns an error when an input amount is invaild
func ErrInvalidInputAmount(codespace string, amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,  "failed. the input amount %s is invalid", amount)}
}

// ErrInsufficientAmount returns an error when there is no enough tokens
func ErrInsufficientAmount(codespace string, amount string, inputAmount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, 
		"failed. the actual amount %s is less than %s", amount, inputAmount)}
}

// ErrInvalidStartHeight returns an error when the start_height_to_yield parameter is invalid
func ErrInvalidStartHeight(codespace string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,  "failed. the start height to yield is less than current height")}
}

// ErrNilAddress returns an error when an empty address appears
func ErrNilAddress(codespace string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidAddress,  "failed. address is nil")}
}

// ErrPoolNotFinished returns an error when the pool is not finished and can not be destroyed
func ErrPoolNotFinished(codespace string, poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errPoolNotFinished, 
		"failed. the pool %s is not finished and can not be destroyed", poolName)}
}

// ErrPoolNameNotExistedInWhiteList returns an error when the pool name is not existed in the white list
func ErrPoolNameNotExistedInWhiteList(codespace string, poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errPoolNotInWhiteList, 
		"failed. the pool name %s not exists in the white list", poolName)}
}

// ErrUnexpectedProposalType returns an error when the proposal type is not supported in farm module
func ErrUnexpectedProposalType(codespace string, proposalType string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errUnexpectedProposalType, 
		"failed. the proposal type %s is not supported in farm module", proposalType)}
}

// ErrPoolNameLength returns an error when length of pool name is invalid
func ErrPoolNameLength(codespace string, poolName string, got, max int) sdk.EnvelopedErr {
	msg := fmt.Sprintf("invalid pool name length for %v, got length %v, max is %v", poolName, got, max)
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,  msg)}
}

// ErrLockAmountBelowMinimum returns an error when lock amount belows minimum
func ErrLockAmountBelowMinimum(codespace string, minLockAmount, amount sdk.Dec) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,  "lock amount %s must be greater than the pool`s min lock amount %s",
		amount.String(), minLockAmount.String())}
}
