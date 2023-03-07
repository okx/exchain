package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
)

// FeeSplit defines an instance that organizes fee distribution conditions for
// the owner of a given smart contract
type FeeSplit struct {
	// hex address of registered contract
	ContractAddress common.Address `json:"contract_address,omitempty"`
	// bech32 address of contract deployer
	DeployerAddress sdk.AccAddress `json:"deployer_address,omitempty"`
	// bech32 address of account receiving the transaction fees it defaults to
	// deployer_address
	WithdrawerAddress sdk.AccAddress `json:"withdrawer_address,omitempty"`
}

// NewFeeSplit returns an instance of FeeSplit. If the provided withdrawer
// address is empty, it sets the value to an empty string.
func NewFeeSplit(contract common.Address, deployer, withdrawer sdk.AccAddress) FeeSplit {
	if withdrawer.Empty() {
		withdrawer = deployer
	}

	return FeeSplit{
		ContractAddress:   contract,
		DeployerAddress:   deployer,
		WithdrawerAddress: withdrawer,
	}
}

// Validate performs a stateless validation of a FeeSplit
func (fs FeeSplit) Validate() error {
	if bytes.Equal(fs.ContractAddress.Bytes(), common.Address{}.Bytes()) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress, "address '%s' is not a valid ethereum hex address",
			fs.ContractAddress.String(),
		)
	}

	if fs.DeployerAddress.Empty() {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress, "empty address string is not allowed",
		)
	}

	return nil
}
