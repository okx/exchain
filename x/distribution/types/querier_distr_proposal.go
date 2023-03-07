package types

import sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"

// querier keys
const (
	QueryValidatorOutstandingRewards = "validator_outstanding_rewards"
	QueryDelegatorTotalRewards       = "delegator_total_rewards"
	QueryDelegatorValidators         = "delegator_validators"
	QueryDelegationRewards           = "delegation_rewards"

	ParamDistributionType        = "distribution_type"
	ParamWithdrawRewardEnabled   = "withdraw_reward_enabled"
	ParamRewardTruncatePrecision = "reward_truncate_precision"
)

// params for query 'custom/distr/delegator_total_rewards' and 'custom/distr/delegator_validators'
type QueryDelegatorParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

// creates a new instance of QueryDelegationRewardsParams
func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddress: delegatorAddr,
	}
}

// params for query 'custom/distr/delegation_rewards'
type QueryDelegationRewardsParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

// creates a new instance of QueryDelegationRewardsParams
func NewQueryDelegationRewardsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) QueryDelegationRewardsParams {
	return QueryDelegationRewardsParams{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
	}
}

// params for query 'custom/distr/validator_outstanding_rewards'
type QueryValidatorOutstandingRewardsParams struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

// creates a new instance of QueryValidatorOutstandingRewardsParams
func NewQueryValidatorOutstandingRewardsParams(validatorAddr sdk.ValAddress) QueryValidatorOutstandingRewardsParams {
	return QueryValidatorOutstandingRewardsParams{
		ValidatorAddress: validatorAddr,
	}
}
