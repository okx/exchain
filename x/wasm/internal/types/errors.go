package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// const CodeType
const (
	codeCreateFailed      sdk.CodeType = 1
	codeAccountExists     sdk.CodeType = 2
	codeInstantiateFailed sdk.CodeType = 3
	codeExecuteFailed     sdk.CodeType = 4
	codeGasLimit          sdk.CodeType = 5
	codeInvalidGenesis    sdk.CodeType = 6
	codeNotFound          sdk.CodeType = 7
	codeQueryFailed       sdk.CodeType = 8
	codeInvalidMsg        sdk.CodeType = 9
)

// CodeType to Message
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case codeCreateFailed:
		return "create wasm contract failed"
	case codeAccountExists:
		return "contract account already exists"
	case codeInstantiateFailed:
		return "instantiate wasm contract failed"
	case codeExecuteFailed:
		return "execute wasm contract failed"
	case codeGasLimit:
		return "insufficient gas"
	case codeInvalidGenesis:
		return "invalid genesis"
	case codeNotFound:
		return "not found"
	case codeQueryFailed:
		return "query wasm contract failed"
	case codeInvalidMsg:
		return "invalid CosmosMsg from the contract"
	default:
		return fmt.Sprintf("unknown code %d", code)
	}
}

// ErrCreateFailed error for wasm code that has already been uploaded or failed
func ErrCreateFailed(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeCreateFailed, codeToDefaultMsg(codeCreateFailed)+": %s", msg)
}

// ErrAccountExists error for a contract account that already exists
func ErrAccountExists(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeAccountExists, codeToDefaultMsg(codeAccountExists)+": %s", msg)
}

// ErrInstantiateFailed error for rust instantiate contract failure
func ErrInstantiateFailed(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInstantiateFailed, codeToDefaultMsg(codeInstantiateFailed)+": %s", msg)
}

// ErrExecuteFailed error for rust execution contract failure
func ErrExecuteFailed(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeExecuteFailed, codeToDefaultMsg(codeExecuteFailed)+": %s", msg)
}

// ErrGasLimit error for out of gas
func ErrGasLimit(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeGasLimit, codeToDefaultMsg(codeGasLimit)+": %s", msg)
}

//ErrInvalidGenesis error for invalid genesis file syntax
func ErrInvalidGenesis(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidGenesis, codeToDefaultMsg(codeInvalidGenesis)+": %s", msg)
}

// ErrNotFound error for an entry not found in the store
func ErrNotFound(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeNotFound, codeToDefaultMsg(codeNotFound)+": %s", msg)
}

// ErrQueryFailed error for rust smart query contract failure
func ErrQueryFailed(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeQueryFailed, codeToDefaultMsg(codeQueryFailed)+": %s", msg)
}

// ErrInvalidMsg error when we cannot process the error returned from the contract
func ErrInvalidMsg(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, codeInvalidMsg, codeToDefaultMsg(codeInvalidMsg)+": %s", msg)
}
