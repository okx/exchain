package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// DelegatorWithdrawInfo is the address for where distributions rewards are withdrawn to by default
// this struct is only used at genesis to feed in default withdraw addresses
type DelegatorWithdrawInfo struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	WithdrawAddress  sdk.AccAddress `json:"withdraw_address" yaml:"withdraw_address"`
}

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

// NewGenesisState creates a new object of GenesisState
func NewGenesisState(params Params) GenesisState {

	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState returns default genesis
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the genesis state of distribution genesis input
func ValidateGenesis(gs GenesisState) error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}
	return nil
}
