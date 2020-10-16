package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidFarmPool           CodeType = 101
	CodeInvalidPoolCurrentRewards CodeType = 102
	CodeInvalidLockInfo           CodeType = 103
	CodeInvalidInput              CodeType = 104
	CodePoolAlreadyExist          CodeType = 105
	CodeTokenNotExist             CodeType = 106
	CodePoolNotFinished           CodeType = 107
	CodeInvalidAddress                     = sdk.CodeInvalidAddress
	CodeUnknownRequest                     = sdk.CodeUnknownRequest
)

// ErrNoFarmPoolFound returns an error when a farm pool doesn't exist
func ErrNoFarmPoolFound(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFarmPool, "failed. Farm pool %s does not exist", poolName)
}

// ErrPoolAlreadyExist returns an error when a pool exist
func ErrPoolAlreadyExist(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodePoolAlreadyExist, "failed. farm pool %s already exists", poolName)
}

// ErrTokenNotExist returns an error when a token not exists
func ErrTokenNotExist(codespace sdk.CodespaceType, tokenName string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotExist, "failed. lock token %s not exists", tokenName)
}

// ErrNoLockInfoFound returns an error when an address doesn't have any lock infos
func ErrNoLockInfoFound(codespace sdk.CodespaceType, addr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLockInfo, "failed. %s doesn't have any lock infos", addr)
}

// ErrInvalidDenom returns an error when the remaining amount in yieldedTokenInfo is not zero
func ErrRemainingAmountNotZero(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput,
		"failed. The remaining amount is %s, so it's not enable to provide token repeatedly util amount become zero",
		amount)
}

// ErrInvalidTokenOwner returns an error when an input address is not the owner of token
func ErrInvalidTokenOwner(codespace sdk.CodespaceType, addr string, token string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. %s isn't the owner of token %s", addr, token)
}

// ErrInvalidPoolOwner returns an error when an input address is not the owner of pool
func ErrInvalidPoolOwner(codespace sdk.CodespaceType, address string, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. %s isn't the owner of pool %s", address, poolName)
}

// ErrInvalidDenom returns an error when it provides an unmatched token name
func ErrInvalidDenom(codespace sdk.CodespaceType, symbolLocked string, token string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. The coin name should be %s, not %s", symbolLocked, token)
}

// ErrInvalidInputAmount returns an error when an input amount is invaild
func ErrInvalidInputAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. The input amount %s is invaild", amount)
}

// ErrinsufficientAmount returns an error when there is no enough tokens
func ErrinsufficientAmount(codespace sdk.CodespaceType, amount string, inputAmount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. The actual amount %s is less than %s", amount, inputAmount)
}

// ErrInvalidStartHeight returns an error when the start_height_to_yield parameter is invalid
func ErrInvalidStartHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. The start height to yield is less than current height")
}

// ErrInvalidInput returns an error when an input parameter is invalid
func ErrInvalidInput(codespace sdk.CodespaceType, input string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. The input parameter %s is invaild", input)
}

// ErrNilAddress returns an error when an empty address appears
func ErrNilAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, "failed. Address is nil")
}

// ErrPoolNotFinished returns an error when the pool is not finished and can not be destroyed
func ErrPoolNotFinished(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodePoolNotFinished, "failed. the pool %s is not finished and can not be destroyed", poolName)
}
