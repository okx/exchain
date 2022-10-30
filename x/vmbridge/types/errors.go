package types

import (
	"fmt"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

var (
	// ErrChainConfigNotFound returns an error if the chain config cannot be found on the store.
	ErrChainConfigNotFound = sdkerrors.Register(ModuleName, 1, "chain configuration not found")

	ErrCallerOfEvmEmpty = sdkerrors.Register(ModuleName, 2, "the caller of evm can not be nil")

	ErrCannotCreate = sdkerrors.Register(ModuleName, 3, "create is not supprot for vmbridge")

	ErrIsNotWasmAddr = sdkerrors.Register(ModuleName, 4, "call wasm contract must use wasmaddress")
	ErrIsNotEvmAddr  = sdkerrors.Register(ModuleName, 5, "call evm contract must use evmaddress")

	ErrAmountNegative   = sdkerrors.Register(ModuleName, 6, "the amount can not negative")
	ErrEvmExecuteFailed = sdkerrors.Register(ModuleName, 7, "the evm execute")

	ErrVMBridgeEnable = sdkerrors.Register(ModuleName, 8, "the vmbridge is disable")
	ErrIsNotOKCAddr   = sdkerrors.Register(ModuleName, 9, "the address prefix must be ex")
	ErrIsNotETHAddr   = sdkerrors.Register(ModuleName, 10, "the address prefix must be 0x")
)

func ErrMsgSendToEvm(str string) *sdkerrors.Error {
	return sdkerrors.Register(ModuleName, 11, fmt.Sprintf("MsgSendToEvm ValidateBasic: %s", str))
}
