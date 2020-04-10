// DONTCOVER
// nolint
package v0_36

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v034distr "github.com/okex/okchain/x/distribution/legacy/v0_34"
)

// ----------------------------------------------------------------------------
// Types and Constants
// ----------------------------------------------------------------------------

const (
	ModuleName = "distribution"
)

type (
	ValidatorAccumulatedCommission = sdk.DecCoins

	GenesisState struct {
		WithdrawAddrEnabled             bool                                             `json:"withdraw_addr_enabled"`
		DelegatorWithdrawInfos          []v034distr.DelegatorWithdrawInfo                `json:"delegator_withdraw_infos"`
		PreviousProposer                sdk.ConsAddress                                  `json:"previous_proposer"`
		ValidatorAccumulatedCommissions []v034distr.ValidatorAccumulatedCommissionRecord `json:"validator_accumulated_commissions"`
	}
)

func NewGenesisState(withdrawAddrEnabled bool, dwis []v034distr.DelegatorWithdrawInfo, pp sdk.ConsAddress, acc []v034distr.ValidatorAccumulatedCommissionRecord) GenesisState {

	return GenesisState{
		WithdrawAddrEnabled:             withdrawAddrEnabled,
		DelegatorWithdrawInfos:          dwis,
		PreviousProposer:                pp,
		ValidatorAccumulatedCommissions: acc,
	}
}
