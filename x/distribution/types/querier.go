package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

// querier keys
const (
	QueryParams              = "params"
	QueryValidatorCommission = "validator_commission"
	QueryWithdrawAddr        = "withdraw_addr"
	QueryCommunityPool       = "community_pool"

	ParamCommunityTax        = "community_tax"
	ParamWithdrawAddrEnabled = "withdraw_addr_enabled"
)

// QueryValidatorCommissionParams is the struct of params for query 'custom/distr/validator_commission'
type QueryValidatorCommissionParams struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
}

// NewQueryValidatorCommissionParams creates a new instance of QueryValidatorCommissionParams
func NewQueryValidatorCommissionParams(validatorAddr sdk.ValAddress) QueryValidatorCommissionParams {
	return QueryValidatorCommissionParams{
		ValidatorAddress: validatorAddr,
	}
}

// QueryDelegatorWithdrawAddrParams is the struct of params for query 'custom/distr/withdraw_addr'
type QueryDelegatorWithdrawAddrParams struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

// NewQueryDelegatorWithdrawAddrParams creates a new instance of QueryDelegatorWithdrawAddrParams.
func NewQueryDelegatorWithdrawAddrParams(delegatorAddr sdk.AccAddress) QueryDelegatorWithdrawAddrParams {
	return QueryDelegatorWithdrawAddrParams{DelegatorAddress: delegatorAddr}
}
