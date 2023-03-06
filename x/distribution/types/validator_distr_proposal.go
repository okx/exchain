package types

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

// historical rewards for a validator
// height is implicit within the store key
// cumulative reward ratio is the sum from the zeroeth period
// until this period of rewards / tokens, per the spec
// The reference count indicates the number of objects
// which might need to reference this historical entry
// at any point.
// ReferenceCount =
//    number of outstanding delegations which ended the associated period (and might need to read that record)
//  + number of slashes which ended the associated period (and might need to read that record)
//  + one per validator for the zeroeth period, set on initialization
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
