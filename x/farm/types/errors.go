package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidFarmPool CodeType = 101
	CodeInvalidLockInfo CodeType = 102
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