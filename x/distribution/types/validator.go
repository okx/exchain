package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// ValidatorAccumulatedCommission is the accumulated commission for a validator
// kept as a running counter, can be withdrawn at any time
type ValidatorAccumulatedCommission = sdk.SysCoins

// InitialValidatorAccumulatedCommission returns the initial accumulated commission (zero)
func InitialValidatorAccumulatedCommission() ValidatorAccumulatedCommission {
	return ValidatorAccumulatedCommission{}
}

type ValidatorHistoricalRewards struct {
	CumulativeRewardRatio sdk.SysCoins `json:"cumulative_reward_ratio" yaml:"cumulative_reward_ratio"`
	ReferenceCount        uint16       `json:"reference_count" yaml:"reference_count"`
}

// create a new ValidatorHistoricalRewards
func NewValidatorHistoricalRewards(cumulativeRewardRatio sdk.SysCoins, referenceCount uint16) ValidatorHistoricalRewards {
	return ValidatorHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
	}
}

// current rewards and current period for a validator
// kept as a running counter and incremented each block
// as long as the validator's tokens remain constant
type ValidatorCurrentRewards struct {
	Rewards sdk.SysCoins `json:"rewards" yaml:"rewards"` // current rewards
	Period  uint64       `json:"period" yaml:"period"`   // current period
}

// create a new ValidatorCurrentRewards
func NewValidatorCurrentRewards(rewards sdk.SysCoins, period uint64) ValidatorCurrentRewards {
	return ValidatorCurrentRewards{
		Rewards: rewards,
		Period:  period,
	}
}

// outstanding (un-withdrawn) rewards for a validator
// inexpensive to track, allows simple sanity checks
type ValidatorOutstandingRewards = sdk.SysCoins
