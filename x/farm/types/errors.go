package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidInput           CodeType = 101
	CodePoolAlreadyExist       CodeType = 102
	CodeNoFarmPoolFound        CodeType = 103
	CodePoolNotInWhiteList     CodeType = 104
	CodeInvalidLockInfo        CodeType = 105
	CodeTokenNotExist          CodeType = 106
	CodePoolNotFinished        CodeType = 107
	CodeUnexpectedProposalType CodeType = 108
	CodeInvalidAddress                  = sdk.CodeInvalidAddress
	CodeUnknownRequest                  = sdk.CodeUnknownRequest
)

// ErrInvalidInput returns an error when an input parameter is invalid
func ErrInvalidInput(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. invalid input: %s", msg)
}

// ErrNoFarmPoolFound returns an error when a farm pool doesn't exist
func ErrNoFarmPoolFound(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodeNoFarmPoolFound, "failed. farm pool %s does not exist", poolName)
}

// ErrPoolAlreadyExist returns an error when a pool exist
func ErrPoolAlreadyExist(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodePoolAlreadyExist, "failed. farm pool %s already exists", poolName)
}

// ErrTokenNotExist returns an error when a token not exists
func ErrTokenNotExist(codespace sdk.CodespaceType, tokenName string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotExist, "failed. token %s does not exist", tokenName)
}

// ErrNoLockInfoFound returns an error when an address doesn't have any lock infos
func ErrNoLockInfoFound(codespace sdk.CodespaceType, addr string, pool string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLockInfo, "failed. %s haven't locked in pool %s", addr, pool)
}

// ErrRemainingAmountNotZero returns an error when the remaining amount in yieldedTokenInfo is not zero
func ErrRemainingAmountNotZero(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput,
		"failed. the remaining amount is %s, so it's not enable to provide token repeatedly util amount become zero",
		amount)
}

// ErrInvalidPoolOwner returns an error when an input address is not the owner of pool
func ErrInvalidPoolOwner(codespace sdk.CodespaceType, address string, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. %s isn't the owner of pool %s", address, poolName)
}

// ErrInvalidDenom returns an error when it provides an unmatched token name
func ErrInvalidDenom(codespace sdk.CodespaceType, symbolLocked string, token string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput,
		"failed. the coin name should be %s, not %s", symbolLocked, token)
}

// ErrInvalidInputAmount returns an error when an input amount is invaild
func ErrInvalidInputAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. the input amount %s is invalid", amount)
}

// ErrInsufficientAmount returns an error when there is no enough tokens
func ErrInsufficientAmount(codespace sdk.CodespaceType, amount string, inputAmount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput,
		"failed. the actual amount %s is less than %s", amount, inputAmount)
}

// ErrInvalidStartHeight returns an error when the start_height_to_yield parameter is invalid
func ErrInvalidStartHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. the start height to yield is less than current height")
}

// ErrNilAddress returns an error when an empty address appears
func ErrNilAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, "failed. address is nil")
}

// ErrPoolNotFinished returns an error when the pool is not finished and can not be destroyed
func ErrPoolNotFinished(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodePoolNotFinished,
		"failed. the pool %s that is with unclaimed rewards or locked coins can not be destroyed", poolName)
}

// ErrPoolNameNotExistedInWhiteList returns an error when the pool name is not existed in the white list
func ErrPoolNameNotExistedInWhiteList(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodePoolNotInWhiteList,
		"failed. the pool name %s not exists in the white list", poolName)
}

// ErrUnexpectedProposalType returns an error when the proposal type is not supported in farm module
func ErrUnexpectedProposalType(codespace sdk.CodespaceType, proposalType string) sdk.Error {
	return sdk.NewError(codespace, CodeUnexpectedProposalType,
		"failed. the proposal type %s is not supported in farm module", proposalType)
}

// ErrPoolNameLength returns an error when length of pool name is invalid
func ErrPoolNameLength(codespace sdk.CodespaceType, poolName string, got, max int) sdk.Error {
	msg := fmt.Sprintf("invalid pool name length for %v, got length %v, max is %v", poolName, got, max)
	return sdk.NewError(codespace, CodeInvalidInput, msg)
}

// ErrLockAmountBelowMinimum returns an error when lock amount belows minimum
func ErrLockAmountBelowMinimum(codespace sdk.CodespaceType, minLockAmount, amount sdk.Dec) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "lock amount %s must be greater than the pool`s min lock amount %s",
		amount.String(), minLockAmount.String())
}
