package types

import sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

var (
	// ErrChainConfigNotFound returns an error if the chain config cannot be found on the store.
	ErrChainConfigNotFound = sdkerrors.Register(ModuleName, 1, "chain configuration not found")

	ErrCallerOfEvmEmpty = sdkerrors.Register(ModuleName, 2, "the caller of evm can not be nil")

	ErrCannotCreate = sdkerrors.Register(ModuleName, 3, "create is not supprot for vmbridge")

	ErrIsNotWasmAddr = sdkerrors.Register(ModuleName, 4, "call wasm contract must use wasmaddress")
	ErrIsNotEvmAddr  = sdkerrors.Register(ModuleName, 5, "call evm contract must use evmaddress")

	ErrAmountNegative = sdkerrors.Register(ModuleName, 6, "the amount can not negative")
)
