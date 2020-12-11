package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = uint32

const (
	DefaultCodespace string = ModuleName

	CodeInvalidInput           uint32 = 66000
	CodePoolAlreadyExist       uint32 = 66001
	CodeNoFarmPoolFound        uint32 = 66002
	CodePoolNotInWhiteList     uint32 = 66003
	CodeInvalidLockInfo        uint32 = 66004
	CodeTokenNotExist          uint32 = 66005
	CodePoolNotFinished        uint32 = 66006
	CodeUnexpectedProposalType uint32 = 66007
	CodeInvalidAddress         uint32 = 66008
	CodeUnknownRequest         uint32 = 66009
	CodeInternal			   uint32 = 66010
	CodeGetEarningsFailed	   uint32 = 66011
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
	errInternal			  	  = sdkerrors.Register(DefaultCodespace, CodeInternal, "error occur internal")
	errGetEarningsFailed	  = sdkerrors.Register(DefaultCodespace, CodeGetEarningsFailed, "get earning failed")
)

// ErrInvalidInput returns an error when an input parameter is invalid
func ErrInvalidInput(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, "failed. invalid input: %s", msg)}
}

// ErrNoFarmPoolFound returns an error when a farm pool doesn't exist
func ErrNoFarmPoolFound(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errNoFarmPoolFound, "failed. farm pool %s does not exist", poolName)}
}

// ErrPoolAlreadyExist returns an error when a pool exist
func ErrPoolAlreadyExist(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errPoolAlreadyExist, "failed. farm pool %s already exists", poolName)}
}

// ErrTokenNotExist returns an error when a token not exists
func ErrTokenNotExist(tokenName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errTokenNotExist, "failed. token %s does not exist", tokenName)}
}

// ErrNoLockInfoFound returns an error when an address doesn't have any lock infos
func ErrNoLockInfoFound(addr string, pool string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidLockInfo, "failed. %s hasn't locked in pool %s", addr, pool)}
}

// ErrRemainingAmountNotZero returns an error when the remaining amount in yieldedTokenInfo is not zero
func ErrRemainingAmountNotZero(amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,
		"failed. the remaining amount is %s, so it's not enable to provide token repeatedly util amount become zero",
		amount)}
}

// ErrInvalidPoolOwner returns an error when an input address is not the owner of pool
func ErrInvalidPoolOwner(address string, poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, "failed. %s isn't the owner of pool %s", address, poolName)}
}

// ErrInvalidDenom returns an error when it provides an unmatched token name
func ErrInvalidDenom(symbolLocked string, token string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,
		"failed. the coin name should be %s, not %s", symbolLocked, token)}
}

// ErrInvalidInputAmount returns an error when an input amount is invaild
func ErrInvalidInputAmount(amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, "failed. the input amount %s is invalid", amount)}
}

// ErrInsufficientAmount returns an error when there is no enough tokens
func ErrInsufficientAmount(amount string, inputAmount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput,
		"failed. the actual amount %s is less than %s", amount, inputAmount)}
}

// ErrInvalidStartHeight returns an error when the start_height_to_yield parameter is invalid
func ErrInvalidStartHeight() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, "failed. the start height to yield is less than current height")}
}

// ErrNilAddress returns an error when an empty address appears
func ErrNilAddress() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidAddress, "failed. address is nil")}
}

// ErrPoolNotFinished returns an error when the pool is not finished and can not be destroyed
func ErrPoolNotFinished(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errPoolNotFinished,
		"failed. the pool %s that is with unclaimed rewards or locked coins can not be destroyed", poolName)}
}

// ErrPoolNameNotExistedInWhiteList returns an error when the pool name is not existed in the white list
func ErrPoolNameNotExistedInWhiteList(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errPoolNotInWhiteList,
		"failed. the pool name %s not exists in the white list", poolName)}
}

// ErrUnexpectedProposalType returns an error when the proposal type is not supported in farm module
func ErrUnexpectedProposalType(proposalType string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errUnexpectedProposalType,
		"failed. the proposal type %s is not supported in farm module", proposalType)}
}

// ErrPoolNameLength returns an error when length of pool name is invalid
func ErrPoolNameLength(poolName string, got, max int) sdk.EnvelopedErr {
	msg := fmt.Sprintf("invalid pool name length for %v, got length %v, max is %v", poolName, got, max)
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, msg)}
}

// ErrLockAmountBelowMinimum returns an error when lock amount belows minimum
func ErrLockAmountBelowMinimum(minLockAmount, amount sdk.Dec) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInvalidInput, "lock amount %s must be greater than the pool`s min lock amount %s",
		amount.String(), minLockAmount.String())}
}

func ErrUnknownRequest(content string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errUnknownRequest, content)}
}

func ErrInternal(content string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errInternal, content)}
}

func ErrGetEarningsFailed(content string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.Wrapf(errGetEarningsFailed, content)}
}