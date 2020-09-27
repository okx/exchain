package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidFarmPool  CodeType = 101
	CodeInvalidLockInfo  CodeType = 102
	CodeInvalidInput     CodeType = 103
	CodePoolAlreadyExist CodeType = 104
	CodeTokenNotExist    CodeType = 105
	CodePoolNotFinished  CodeType = 106
	CodeInvalidAddress            = sdk.CodeInvalidAddress
	CodeUnknownRequest            = sdk.CodeUnknownRequest
)

// ErrNoFarmPoolFound returns an error when a farm pool doesn't exist
func ErrNoFarmPoolFound(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFarmPool, "failed. farm pool %s does not exist", poolName)
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

// ErrInvalidTokenOwner returns an error when an input address is not the owner of token
func ErrInvalidTokenOwner(codespace sdk.CodespaceType, addr string, token string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. %s isn't the owner of token %s", addr, token)
}

// ErrInvalidPoolOwner returns an error when an input address is not the owner of pool
func ErrInvalidPoolOwner(codespace sdk.CodespaceType, address string, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. %s isn't the owner of pool %s", address, poolName)
}

// ErrInvalidAmount returns an error when an input amount is invalid
func ErrInvalidAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. the input amount %s is invalid", amount)
}

// ErrInsufficientAmount returns an error when the total amount is not enough to yield in one block
func ErrInsufficientAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. the provided amount %s is not enough", amount)
}

// ErrInvalidStartHeight returns an error when the start_height_to_yield parameter is invalid
func ErrInvalidStartHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. the start height to yield is less than current height")
}

// ErrInvalidInput returns an error when an input parameter is invalid
func ErrInvalidInput(codespace sdk.CodespaceType, input string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "failed. the input parameter %s is invalid", input)
}

// ErrNilAddress returns an error when an empty address appears
func ErrNilAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, "failed. address is nil")
}

// ErrPoolNotFinished returns an error when the pool is not finished and can not be destroyed
func ErrPoolNotFinished(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodePoolNotFinished, "failed. the pool %s is not finished and can not be destroyed", poolName)
}
