package types

import (
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	DefaultCodespace string = ModuleName
)

var (
	// ErrChainConfigNotFound returns an error if the chain config cannot be found on the store.
	ErrChainConfigNotFound = sdkerrors.Register(ModuleName, 1, "chain configuration not found")
	// ErrKeyNotFound returns an error if the target key not found in database.
	ErrKeyNotFound = sdkerrors.Register(ModuleName, 2, "Key not found in database")
	// ErrUnexpectedProposalType returns an error when the proposal type is not supported in erc20 module
	ErrUnexpectedProposalType = sdkerrors.Register(ModuleName, 3, "Unsupported proposal type of erc20 module")
	// ErrEmptyAddressList returns an error if the address list is empty
	ErrEmptyAddressList = sdkerrors.Register(ModuleName, 4, "Empty account address list")
	ErrIbcDenomInvalid  = sdkerrors.Register(ModuleName, 5, "ibc denom is invalid")
)

func ErrRegisteredContract(contract string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{
		Err: sdkerrors.New(
			ModuleName,
			21,
			fmt.Sprintf("the contract is already registered: %s", contract),
		),
	}
}
