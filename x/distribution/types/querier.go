package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

// querier keys
const (
	QueryParams                  = "params"
	QueryValidatorCommission     = "validator_commission"
	QueryCM45ValidatorCommission = "commission"
	QueryWithdrawAddr            = "withdraw_addr"
	QueryCommunityPool           = "community_pool"

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

// QueryValidatorCommissionRequest is the request type for the
// Query/ValidatorCommission RPC method
type QueryValidatorCommissionRequest struct {
	// validator_address defines the validator address to query for.
	ValidatorAddress string `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
}

// NewQueryValidatorCommissionRequest creates a new instance of NewQueryValidatorCommissionRequest
func NewQueryValidatorCommissionRequest(validatorAddr string) QueryValidatorCommissionRequest {
	return QueryValidatorCommissionRequest{
		ValidatorAddress: validatorAddr,
	}
}

// QueryValidatorCommissionResponse is the response type for the
// Query/ValidatorCommission RPC method
type QueryValidatorCommissionResponse struct {
	// commission defines the commision the validator received.
	Commission ValidatorAccumulatedCommission `protobuf:"bytes,1,opt,name=commission,proto3" json:"commission"`
}

type WrappedCommission struct {
	// commission defines the commision the validator received.
	Response QueryValidatorCommissionResponse `json:"commission"`
}

func NewWrappedCommission(r QueryValidatorCommissionResponse) WrappedCommission {
	return WrappedCommission{
		Response: r,
	}
}
