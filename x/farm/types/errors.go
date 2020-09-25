package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidFarmPool CodeType = 101
	CodeInvalidLockInfo CodeType = 102
	CodeInvalidInput    CodeType = 103
	CodeInvalidAddress           = sdk.CodeInvalidAddress
	CodeUnknownRequest           = sdk.CodeUnknownRequest
)

// ErrNoFarmPoolFound returns an error when a farm pool doesn't exist
func ErrNoFarmPoolFound(codespace sdk.CodespaceType, poolName string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidFarmPool, "farm pool %s does not exist", poolName)
}

// ErrNoLockInfoFound returns an error when an address doesn't have any lock infos
func ErrNoLockInfoFound(codespace sdk.CodespaceType, addr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLockInfo, "%s doesn't have any lock infos", addr)
}

// ErrInvalidTokenOwner returns an error when an input address is not the owner of token
func ErrInvalidTokenOwner(codespace sdk.CodespaceType, addr string, token string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "%s isn't the owner of token %s", addr, token)
}

// ErrInvalidAmount returns an error when an input amount is invaild
func ErrInvalidAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "the input amount %s is invaild", amount)
}

// ErrInsufficientAmount returns an error when the total amount is not enough to yield in one block
func ErrInsufficientAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "the provided amount %s is not enough", amount)
}

// ErrInvalidStartHeight returns an error when the start_height_to_yield parameter is invaild
func ErrInvalidStartHeight(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "the start height to yield is less than current height")
}

// ErrInvalidInput returns an error when an input parameter is invaild
func ErrInvalidInput(codespace sdk.CodespaceType, input string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "the input parameter %s is invaild", input)
}

// ErrNilAddress returns an error when an empty address appears
func ErrNilAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "address is nil")
}