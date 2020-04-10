package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegatorWithdrawInfo is the address for where distributions rewards are withdrawn to by default
// this struct is only used at genesis to feed in default withdraw addresses
type DelegatorWithdrawInfo struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	WithdrawAddress  sdk.AccAddress `json:"withdraw_address" yaml:"withdraw_address"`
}

// ValidatorAccumulatedCommissionRecord is used for import / export via genesis json
type ValidatorAccumulatedCommissionRecord struct {
	ValidatorAddress sdk.ValAddress                 `json:"validator_address" yaml:"validator_address"`
	Accumulated      ValidatorAccumulatedCommission `json:"accumulated" yaml:"accumulated"`
}

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	WithdrawAddrEnabled             bool                                   `json:"withdraw_addr_enabled" yaml:"withdraw_addr_enabled"`
	DelegatorWithdrawInfos          []DelegatorWithdrawInfo                `json:"delegator_withdraw_infos" yaml:"delegator_withdraw_infos"`
	PreviousProposer                sdk.ConsAddress                        `json:"previous_proposer" yaml:"previous_proposer"`
	ValidatorAccumulatedCommissions []ValidatorAccumulatedCommissionRecord `json:"validator_accumulated_commissions" yaml:"validator_accumulated_commissions"`
}

// NewGenesisState creates a new object of GenesisState
func NewGenesisState(withdrawAddrEnabled bool, dwis []DelegatorWithdrawInfo, pp sdk.ConsAddress,
	acc []ValidatorAccumulatedCommissionRecord) GenesisState {
	return GenesisState{
		WithdrawAddrEnabled:             withdrawAddrEnabled,
		DelegatorWithdrawInfos:          dwis,
		PreviousProposer:                pp,
		ValidatorAccumulatedCommissions: acc,
	}
}

// DefaultGenesisState returns default genesis
func DefaultGenesisState() GenesisState {
	return GenesisState{
		WithdrawAddrEnabled:             true,
		DelegatorWithdrawInfos:          []DelegatorWithdrawInfo{},
		PreviousProposer:                nil,
		ValidatorAccumulatedCommissions: []ValidatorAccumulatedCommissionRecord{},
	}
}

// ValidateGenesis validates the genesis state of distribution genesis input
func ValidateGenesis(data GenesisState) error {
	return nil
}
