package keeper

import (
	"context"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/okx/okbchain/x/staking/typesadapter"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	k Keeper
}

func NewGrpcQuerier(k Keeper) *Querier {
	return &Querier{k: k}
}

var _ typesadapter.QueryServer = (*Querier)(nil)

// Validators queries all validators that match the given status
func (k Querier) Validators(c context.Context, req *typesadapter.QueryValidatorsRequest) (*typesadapter.QueryValidatorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) Validator(ctx context.Context, request *typesadapter.QueryValidatorRequest) (*typesadapter.QueryValidatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) ValidatorDelegations(ctx context.Context, request *typesadapter.QueryValidatorDelegationsRequest) (*typesadapter.QueryValidatorDelegationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) ValidatorUnbondingDelegations(ctx context.Context, request *typesadapter.QueryValidatorUnbondingDelegationsRequest) (*typesadapter.QueryValidatorUnbondingDelegationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) Delegation(ctx context.Context, request *typesadapter.QueryDelegationRequest) (*typesadapter.QueryDelegationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) UnbondingDelegation(ctx context.Context, request *typesadapter.QueryUnbondingDelegationRequest) (*typesadapter.QueryUnbondingDelegationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) DelegatorDelegations(ctx context.Context, request *typesadapter.QueryDelegatorDelegationsRequest) (*typesadapter.QueryDelegatorDelegationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) DelegatorUnbondingDelegations(ctx context.Context, request *typesadapter.QueryDelegatorUnbondingDelegationsRequest) (*typesadapter.QueryDelegatorUnbondingDelegationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) Redelegations(ctx context.Context, request *typesadapter.QueryRedelegationsRequest) (*typesadapter.QueryRedelegationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) DelegatorValidators(ctx context.Context, request *typesadapter.QueryDelegatorValidatorsRequest) (*typesadapter.QueryDelegatorValidatorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) DelegatorValidator(ctx context.Context, request *typesadapter.QueryDelegatorValidatorRequest) (*typesadapter.QueryDelegatorValidatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) HistoricalInfo(ctx context.Context, request *typesadapter.QueryHistoricalInfoRequest) (*typesadapter.QueryHistoricalInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) Pool(ctx context.Context, request *typesadapter.QueryPoolRequest) (*typesadapter.QueryPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validators not implemented")
}

func (q *Querier) Params(goCtx context.Context, request *typesadapter.QueryParamsRequest) (*typesadapter.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.k.GetParams(ctx)

	ret := typesadapter.Params{}
	ret.From(params)
	return &typesadapter.QueryParamsResponse{Params: ret}, nil

}
