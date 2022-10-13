package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
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
	var withdrawerAddr sdk.AccAddress
	if !withdrawer.Empty() {
		withdrawerAddr = withdrawer
	}

	return FeeSplit{
		ContractAddress:   contract,
		DeployerAddress:   deployer,
		WithdrawerAddress: withdrawerAddr,
	}
}

// GetContractAddr returns the contract address
//func (fs FeeSplit) GetContractAddr() common.Address {
//	return common.HexToAddress(fs.ContractAddress)
//}

// GetDeployerAddr returns the contract deployer address
//func (fs FeeSplit) GetDeployerAddr() sdk.AccAddress {
//	return sdk.MustAccAddressFromBech32(fs.DeployerAddress)
//}

// GetWithdrawerAddr returns the account address to where the funds proceeding
// from the fees will be received. If the withdraw address is not defined, it
// defaults to the deployer address.
//func (fs FeeSplit) GetWithdrawerAddr() sdk.AccAddress {
//	if fs.WithdrawerAddress == "" {
//		return nil
//	}
//
//	return sdk.MustAccAddressFromBech32(fs.WithdrawerAddress)
//}

// Validate performs a stateless validation of a FeeSplit
func (fs FeeSplit) Validate() error {
	if bytes.Equal(fs.ContractAddress.Bytes(), common.Address{}.Bytes()) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress, "address '%s' is not a valid ethereum hex address",
			fs.ContractAddress.String(),
		)
	}

	//if _, err := sdk.AccAddressFromBech32(fs.DeployerAddress); err != nil {
	//	return err
	//}

	//if fs.WithdrawerAddress != "" {
	//	if _, err := sdk.AccAddressFromBech32(fs.WithdrawerAddress); err != nil {
	//		return err
	//	}
	//}

	return nil
}
