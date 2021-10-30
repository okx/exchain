package types

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
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
	Params                          Params                                 `json:"params" yaml:"params"`
	FeePool                         FeePool                                `json:"fee_pool" yaml:"fee_pool"`
	DelegatorWithdrawInfos          []DelegatorWithdrawInfo                `json:"delegator_withdraw_infos" yaml:"delegator_withdraw_infos"`
	PreviousProposer                sdk.ConsAddress                        `json:"previous_proposer" yaml:"previous_proposer"`
	ValidatorAccumulatedCommissions []ValidatorAccumulatedCommissionRecord `json:"validator_accumulated_commissions" yaml:"validator_accumulated_commissions"`
}

// NewGenesisState creates a new object of GenesisState
func NewGenesisState( params Params, feePool FeePool,
	dwis []DelegatorWithdrawInfo, pp sdk.ConsAddress, acc []ValidatorAccumulatedCommissionRecord) GenesisState {

	return GenesisState{
		Params:                          params,
		FeePool:                         feePool,
		DelegatorWithdrawInfos:          dwis,
		PreviousProposer:                pp,
		ValidatorAccumulatedCommissions: acc,
	}
}

// DefaultGenesisState returns default genesis
func DefaultGenesisState() GenesisState {
	return GenesisState{
		FeePool:                         InitialFeePool(),
		Params:                          DefaultParams(),
		DelegatorWithdrawInfos:          []DelegatorWithdrawInfo{},
		PreviousProposer:                nil,
		ValidatorAccumulatedCommissions: []ValidatorAccumulatedCommissionRecord{},
	}
}

// ValidateGenesis validates the genesis state of distribution genesis input
func ValidateGenesis(gs GenesisState) error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}
	return gs.FeePool.ValidateGenesis()
}
