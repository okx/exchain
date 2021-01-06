package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = uint32

const (
	DefaultCodespace string = ModuleName

	CodeInvalidInput                       uint32 = 66000
	CodePoolAlreadyExist                   uint32 = 66001
	CodeNoFarmPoolFound                    uint32 = 66002
	CodePoolNotInWhiteList                 uint32 = 66003
	CodeInvalidLockInfo                    uint32 = 66004
	CodeTokenNotExist                      uint32 = 66005
	CodePoolNotFinished                    uint32 = 66006
	CodeUnexpectedProposalType             uint32 = 66007
	CodeInvalidAddress                     uint32 = 66008
	CodeRemainingAmountNotZero             uint32 = 66009
	CodeInvalidPoolOwner                   uint32 = 66010
	CodeInvalidDenom                       uint32 = 66011
	CodeSendCoinsFromAccountToModuleFailed uint32 = 66012
	CodeUnknownFarmMsgType                 uint32 = 66013
	CodeUnknownFarmQueryType               uint32 = 66014
	CodeInvalidInputAmount                 uint32 = 66015
	CodeInsufficientAmount                 uint32 = 66016
	CodeInvalidStartHeight                 uint32 = 66017
	CodePoolNameLength                     uint32 = 66018
	CodeLockAmountBelowMinimum             uint32 = 66019
	CodeSendCoinsFromModuleToAccountFailed uint32 = 66020
	CodeSwapTokenPairNotExist              uint32 = 66021
)

// ErrInvalidInput returns an error when an input parameter is invalid
func ErrInvalidInput(msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInvalidInput, fmt.Sprintf("failed. invalid input: %s", msg))}
}

// ErrPoolAlreadyExist returns an error when a pool exist
func ErrPoolAlreadyExist(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodePoolAlreadyExist, fmt.Sprintf("failed. farm pool %s already exists", poolName))}
}

// ErrNoFarmPoolFound returns an error when a farm pool doesn't exist
func ErrNoFarmPoolFound(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeNoFarmPoolFound, fmt.Sprintf("failed. farm pool %s does not exist", poolName))}
}

// ErrPoolNameNotExistedInWhiteList returns an error when the pool name is not existed in the white list
func ErrPoolNameNotExistedInWhiteList(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodePoolNotInWhiteList,
		fmt.Sprintf("failed. the pool name %s not exists in the white list", poolName))}
}

// ErrNoLockInfoFound returns an error when an address doesn't have any lock infos
func ErrNoLockInfoFound(addr string, pool string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInvalidLockInfo, fmt.Sprintf("failed. %s hasn't locked in pool %s", addr, pool))}
}

// ErrTokenNotExist returns an error when a token not exists
func ErrTokenNotExist(tokenName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeTokenNotExist, fmt.Sprintf("failed. token %s does not exist", tokenName))}
}

// ErrPoolNotFinished returns an error when the pool is not finished and can not be destroyed
func ErrPoolNotFinished(poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodePoolNotFinished,
		fmt.Sprintf("failed. the pool %s that is with unclaimed rewards or locked coins can not be destroyed", poolName))}
}

// ErrUnexpectedProposalType returns an error when the proposal type is not supported in farm module
func ErrUnexpectedProposalType(proposalType string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeUnexpectedProposalType,
		fmt.Sprintf("failed. the proposal type %s is not supported in farm module", proposalType))}
}

// ErrNilAddress returns an error when an empty address appears
func ErrNilAddress() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInvalidAddress, fmt.Sprintf("failed. address is required"))}
}

// ErrRemainingAmountNotZero returns an error when the remaining amount in yieldedTokenInfo is not zero
func ErrRemainingAmountNotZero(amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeRemainingAmountNotZero,
		fmt.Sprintf("failed. the remaining amount is %s, so it's not enable to provide token repeatedly util amount become zero", amount))}
}

// ErrInvalidPoolOwner returns an error when an input address is not the owner of pool
func ErrInvalidPoolOwner(address string, poolName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInvalidPoolOwner, fmt.Sprintf("failed. %s isn't the owner of pool %s", address, poolName))}
}

// ErrInvalidDenom returns an error when it provides an unmatched token name
func ErrInvalidDenom(symbolLocked string, token string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInvalidDenom,
		fmt.Sprintf("failed. the coin name should be %s, not %s", symbolLocked, token))}
}

func ErrSendCoinsFromAccountToModuleFailed(content string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeSendCoinsFromAccountToModuleFailed, fmt.Sprintf("failed. send coins from account to module failed %s", content))}
}

func ErrUnknownFarmMsgType(content string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeUnknownFarmMsgType, fmt.Sprintf("failed. unknown farm msg type %s", content))}
}

func ErrUnknownFarmQueryType(content string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeUnknownFarmQueryType, fmt.Sprintf("failed. unknown farm query type %s", content))}
}

// ErrInvalidInputAmount returns an error when an input amount is invaild
func ErrInvalidInputAmount(amount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInvalidInputAmount, fmt.Sprintf("failed. the input amount %s is invalid", amount))}
}

// ErrInsufficientAmount returns an error when there is no enough tokens
func ErrInsufficientAmount(amount string, inputAmount string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInsufficientAmount,
		fmt.Sprintf("failed. the actual amount %s is less than %s", amount, inputAmount))}
}

// ErrInvalidStartHeight returns an error when the start_height_to_yield parameter is invalid
func ErrInvalidStartHeight() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeInvalidStartHeight, "failed. the start height to yield is less than current height")}
}

// ErrPoolNameLength returns an error when length of pool name is invalid
func ErrPoolNameLength(poolName string, got, max int) sdk.EnvelopedErr {
	msg := fmt.Sprintf("invalid pool name length for %v, got length %v, max is %v", poolName, got, max)
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodePoolNameLength, msg)}
}

// ErrLockAmountBelowMinimum returns an error when lock amount belows minimum
func ErrLockAmountBelowMinimum(minLockAmount, amount sdk.Dec) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeLockAmountBelowMinimum, fmt.Sprintf("lock amount %s must be greater than the pool's min lock amount %s", amount.String(), minLockAmount.String()))}
}

func ErrSendCoinsFromModuleToAccountFailed(content string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeSendCoinsFromModuleToAccountFailed, fmt.Sprintf("failed. send coins from module to account failed %s", content))}
}

// ErrSwapTokenPairNotExist returns an error when a swap token pair not exists
func ErrSwapTokenPairNotExist(tokenName string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultParamspace, CodeSwapTokenPairNotExist, fmt.Sprintf("failed. swap token pair %s does not exist", tokenName))}
}