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
		FeePool                         v034distr.FeePool                                `json:"fee_pool" yaml:"fee_pool"`
		CommunityTax                    sdk.Dec                                          `json:"community_tax" yaml:"community_tax"`
		WithdrawAddrEnabled             bool                                             `json:"withdraw_addr_enabled"`
		DelegatorWithdrawInfos          []v034distr.DelegatorWithdrawInfo                `json:"delegator_withdraw_infos"`
		PreviousProposer                sdk.ConsAddress                                  `json:"previous_proposer"`
		ValidatorAccumulatedCommissions []v034distr.ValidatorAccumulatedCommissionRecord `json:"validator_accumulated_commissions"`
	}
)

func NewGenesisState(feePool v034distr.FeePool, withdrawAddrEnabled bool, dwis []v034distr.DelegatorWithdrawInfo,
	pp sdk.ConsAddress, acc []v034distr.ValidatorAccumulatedCommissionRecord) GenesisState {

	return GenesisState{
		FeePool:				         feePool,
		CommunityTax:                    sdk.NewDecWithPrec(2,2),
		WithdrawAddrEnabled:             withdrawAddrEnabled,
		DelegatorWithdrawInfos:          dwis,
		PreviousProposer:                pp,
		ValidatorAccumulatedCommissions: acc,
	}
}
