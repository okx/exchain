package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
)

// note: it should act like protobuf ,but for now ,we only need the `pure struct`

// QueryValidatorDelegationsResponse is response type for the
// Query/ValidatorDelegations RPC method
type QueryValidatorDelegationsResponse struct {
	DelegationResponses types.DelegationResponses `json:"delegation_responses"`
	Pagination          *query.PageResponse       `json:"pagination"`
}

type QueryUnbondingDelegationResponse struct {
	// unbond defines the unbonding information of a delegation.
	Unbond types.UnbondingDelegation `protobuf:"bytes,1,opt,name=unbond,proto3" json:"unbond"`
}

type QueryDelegatorUnbondingDelegationsResponse struct {
	UnbondingResponses []types.UnbondingDelegation `protobuf:"bytes,1,rep,name=unbonding_responses,json=unbondingResponses,proto3" json:"unbonding_responses"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

type QueryDelegatorDelegationsResponse struct {
	// delegation_responses defines all the delegations' info of a delegator.
	DelegationResponses []types.DelegationResponse `protobuf:"bytes,1,rep,name=delegation_responses,json=delegationResponses,proto3" json:"delegation_responses"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
