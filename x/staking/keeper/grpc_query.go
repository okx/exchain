package keeper

import (
	"context"

	"github.com/okex/exchain/x/staking/typesadapter"
)

//import (
//	"context"
//	"strings"
//
//	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
//
//	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
//	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
//	"github.com/okex/exchain/x/staking/typesadapter"
//
//	"google.golang.org/grpc/codes"
//	"google.golang.org/grpc/status"
//)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	k Keeper
}

func NewGrpcQuerier(k Keeper) *Querier {
	return &Querier{k: k}
}

var _ typesadapter.QueryServer = (*Querier)(nil)

func (q *Querier) Validators(ctx context.Context, request *typesadapter.QueryValidatorsRequest) (*typesadapter.QueryValidatorsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) Validator(ctx context.Context, request *typesadapter.QueryValidatorRequest) (*typesadapter.QueryValidatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) ValidatorDelegations(ctx context.Context, request *typesadapter.QueryValidatorDelegationsRequest) (*typesadapter.QueryValidatorDelegationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) ValidatorUnbondingDelegations(ctx context.Context, request *typesadapter.QueryValidatorUnbondingDelegationsRequest) (*typesadapter.QueryValidatorUnbondingDelegationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) Delegation(ctx context.Context, request *typesadapter.QueryDelegationRequest) (*typesadapter.QueryDelegationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) UnbondingDelegation(ctx context.Context, request *typesadapter.QueryUnbondingDelegationRequest) (*typesadapter.QueryUnbondingDelegationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) DelegatorDelegations(ctx context.Context, request *typesadapter.QueryDelegatorDelegationsRequest) (*typesadapter.QueryDelegatorDelegationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) DelegatorUnbondingDelegations(ctx context.Context, request *typesadapter.QueryDelegatorUnbondingDelegationsRequest) (*typesadapter.QueryDelegatorUnbondingDelegationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) Redelegations(ctx context.Context, request *typesadapter.QueryRedelegationsRequest) (*typesadapter.QueryRedelegationsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) DelegatorValidators(ctx context.Context, request *typesadapter.QueryDelegatorValidatorsRequest) (*typesadapter.QueryDelegatorValidatorsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) DelegatorValidator(ctx context.Context, request *typesadapter.QueryDelegatorValidatorRequest) (*typesadapter.QueryDelegatorValidatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) HistoricalInfo(ctx context.Context, request *typesadapter.QueryHistoricalInfoRequest) (*typesadapter.QueryHistoricalInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) Pool(ctx context.Context, request *typesadapter.QueryPoolRequest) (*typesadapter.QueryPoolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q *Querier) Params(ctx context.Context, request *typesadapter.QueryParamsRequest) (*typesadapter.QueryParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

//
//// Validators queries all validators that match the given status
//func (k Querier) Validators(c context.Context, req *typesadapter.QueryValidatorsRequest) (*typesadapter.QueryValidatorsResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	// validate the provided status, return all the validators if the status is empty
//	if req.Status != "" && !(req.Status == typesadapter.Bonded.String() || req.Status == typesadapter.Unbonded.String() || req.Status == typesadapter.Unbonding.String()) {
//		return nil, status.Errorf(codes.InvalidArgument, "invalid validator status %s", req.Status)
//	}
//
//	var validators typesadapter.Validators
//	ctx := sdk.UnwrapSDKContext(c)
//
//	store := ctx.KVStore(k.storeKey)
//	valStore := prefix.NewStore(store, typesadapter.ValidatorsKey)
//
//	pageRes, err := query.FilteredPaginate(valStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
//		val, err := typesadapter.UnmarshalValidator(k.cdc, value)
//		if err != nil {
//			return false, err
//		}
//
//		if req.Status != "" && !strings.EqualFold(val.GetStatus().String(), req.Status) {
//			return false, nil
//		}
//
//		if accumulate {
//			validators = append(validators, val)
//		}
//
//		return true, nil
//	})
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryValidatorsResponse{Validators: validators, Pagination: pageRes}, nil
//}
//
//// Validator queries validator info for given validator address
//func (k Querier) Validator(c context.Context, req *typesadapter.QueryValidatorRequest) (*typesadapter.QueryValidatorResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.ValidatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
//	}
//
//	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//	validator, found := k.k.GetValidator(ctx, valAddr)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "validator %s not found", req.ValidatorAddr)
//	}
//
//	return &typesadapter.QueryValidatorResponse{Validator: validator}, nil
//}
//
//// ValidatorDelegations queries delegate info for given validator
//func (k Querier) ValidatorDelegations(c context.Context, req *typesadapter.QueryValidatorDelegationsRequest) (*typesadapter.QueryValidatorDelegationsResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.ValidatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
//	}
//	var delegations []typesadapter.Delegation
//	ctx := sdk.UnwrapSDKContext(c)
//
//	store := ctx.KVStore(k.storeKey)
//	valStore := prefix.NewStore(store, typesadapter.DelegationKey)
//	pageRes, err := query.FilteredPaginate(valStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
//		delegation, err := typesadapter.UnmarshalDelegation(k.cdc, value)
//		if err != nil {
//			return false, err
//		}
//
//		valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
//		if err != nil {
//			return false, err
//		}
//
//		if !delegation.GetValidatorAddr().Equals(valAddr) {
//			return false, nil
//		}
//
//		if accumulate {
//			delegations = append(delegations, delegation)
//		}
//		return true, nil
//	})
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	delResponses, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryValidatorDelegationsResponse{
//		DelegationResponses: delResponses, Pagination: pageRes,
//	}, nil
//}
//
//// ValidatorUnbondingDelegations queries unbonding delegations of a validator
//func (k Querier) ValidatorUnbondingDelegations(c context.Context, req *typesadapter.QueryValidatorUnbondingDelegationsRequest) (*typesadapter.QueryValidatorUnbondingDelegationsResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.ValidatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
//	}
//	var ubds typesadapter.UnbondingDelegations
//	ctx := sdk.UnwrapSDKContext(c)
//
//	store := ctx.KVStore(k.storeKey)
//
//	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	srcValPrefix := typesadapter.GetUBDsByValIndexKey(valAddr)
//	ubdStore := prefix.NewStore(store, srcValPrefix)
//	pageRes, err := query.Paginate(ubdStore, req.Pagination, func(key []byte, value []byte) error {
//		storeKey := typesadapter.GetUBDKeyFromValIndexKey(append(srcValPrefix, key...))
//		storeValue := store.Get(storeKey)
//
//		ubd, err := typesadapter.UnmarshalUBD(k.cdc, storeValue)
//		if err != nil {
//			return err
//		}
//		ubds = append(ubds, ubd)
//		return nil
//	})
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryValidatorUnbondingDelegationsResponse{
//		UnbondingResponses: ubds,
//		Pagination:         pageRes,
//	}, nil
//}
//
//// Delegation queries delegate info for given validator delegator pair
//func (k Querier) Delegation(c context.Context, req *typesadapter.QueryDelegationRequest) (*typesadapter.QueryDelegationResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.DelegatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
//	}
//	if req.ValidatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
//	if !found {
//		return nil, status.Errorf(
//			codes.NotFound,
//			"delegation with delegator %s not found for validator %s",
//			req.DelegatorAddr, req.ValidatorAddr)
//	}
//
//	delResponse, err := DelegationToDelegationResponse(ctx, k.Keeper, delegation)
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryDelegationResponse{DelegationResponse: &delResponse}, nil
//}
//
//// UnbondingDelegation queries unbonding info for give validator delegator pair
//func (k Querier) UnbondingDelegation(c context.Context, req *typesadapter.QueryUnbondingDelegationRequest) (*typesadapter.QueryUnbondingDelegationResponse, error) {
//	if req == nil {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	if req.DelegatorAddr == "" {
//		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
//	}
//	if req.ValidatorAddr == "" {
//		return nil, status.Errorf(codes.InvalidArgument, "validator address cannot be empty")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	unbond, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
//	if !found {
//		return nil, status.Errorf(
//			codes.NotFound,
//			"unbonding delegation with delegator %s not found for validator %s",
//			req.DelegatorAddr, req.ValidatorAddr)
//	}
//
//	return &typesadapter.QueryUnbondingDelegationResponse{Unbond: unbond}, nil
//}
//
//// DelegatorDelegations queries all delegations of a give delegator address
//func (k Querier) DelegatorDelegations(c context.Context, req *typesadapter.QueryDelegatorDelegationsRequest) (*typesadapter.QueryDelegatorDelegationsResponse, error) {
//	if req == nil {
//		return nil, status.Errorf(codes.InvalidArgument, "empty request")
//	}
//
//	if req.DelegatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
//	}
//	var delegations typesadapter.Delegations
//	ctx := sdk.UnwrapSDKContext(c)
//
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	store := ctx.KVStore(k.storeKey)
//	delStore := prefix.NewStore(store, typesadapter.GetDelegationsKey(delAddr))
//	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
//		delegation, err := typesadapter.UnmarshalDelegation(k.cdc, value)
//		if err != nil {
//			return err
//		}
//		delegations = append(delegations, delegation)
//		return nil
//	})
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	delegationResps, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryDelegatorDelegationsResponse{DelegationResponses: delegationResps, Pagination: pageRes}, nil
//}
//
//// DelegatorValidator queries validator info for given delegator validator pair
//func (k Querier) DelegatorValidator(c context.Context, req *typesadapter.QueryDelegatorValidatorRequest) (*typesadapter.QueryDelegatorValidatorResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.DelegatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
//	}
//	if req.ValidatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
//	}
//
//	ctx := sdk.UnwrapSDKContext(c)
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	validator, err := k.GetDelegatorValidator(ctx, delAddr, valAddr)
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryDelegatorValidatorResponse{Validator: validator}, nil
//}
//
//// DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address
//func (k Querier) DelegatorUnbondingDelegations(c context.Context, req *typesadapter.QueryDelegatorUnbondingDelegationsRequest) (*typesadapter.QueryDelegatorUnbondingDelegationsResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.DelegatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
//	}
//	var unbondingDelegations typesadapter.UnbondingDelegations
//	ctx := sdk.UnwrapSDKContext(c)
//
//	store := ctx.KVStore(k.storeKey)
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	unbStore := prefix.NewStore(store, typesadapter.GetUBDsKey(delAddr))
//	pageRes, err := query.Paginate(unbStore, req.Pagination, func(key []byte, value []byte) error {
//		unbond, err := typesadapter.UnmarshalUBD(k.cdc, value)
//		if err != nil {
//			return err
//		}
//		unbondingDelegations = append(unbondingDelegations, unbond)
//		return nil
//	})
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryDelegatorUnbondingDelegationsResponse{
//		UnbondingResponses: unbondingDelegations, Pagination: pageRes,
//	}, nil
//}
//
//// HistoricalInfo queries the historical info for given height
//func (k Querier) HistoricalInfo(c context.Context, req *typesadapter.QueryHistoricalInfoRequest) (*typesadapter.QueryHistoricalInfoResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.Height < 0 {
//		return nil, status.Error(codes.InvalidArgument, "height cannot be negative")
//	}
//	ctx := sdk.UnwrapSDKContext(c)
//	hi, found := k.GetHistoricalInfo(ctx, req.Height)
//	if !found {
//		return nil, status.Errorf(codes.NotFound, "historical info for height %d not found", req.Height)
//	}
//
//	return &typesadapter.QueryHistoricalInfoResponse{Hist: &hi}, nil
//}
//
//// Redelegations queries redelegations of given address
//func (k Querier) Redelegations(c context.Context, req *typesadapter.QueryRedelegationsRequest) (*typesadapter.QueryRedelegationsResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	var redels typesadapter.Redelegations
//	var pageRes *query.PageResponse
//	var err error
//
//	ctx := sdk.UnwrapSDKContext(c)
//	store := ctx.KVStore(k.storeKey)
//	switch {
//	case req.DelegatorAddr != "" && req.SrcValidatorAddr != "" && req.DstValidatorAddr != "":
//		redels, err = queryRedelegation(ctx, k, req)
//	case req.DelegatorAddr == "" && req.SrcValidatorAddr != "" && req.DstValidatorAddr == "":
//		redels, pageRes, err = queryRedelegationsFromSrcValidator(store, k, req)
//	default:
//		redels, pageRes, err = queryAllRedelegations(store, k, req)
//	}
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//	redelResponses, err := RedelegationsToRedelegationResponses(ctx, k.Keeper, redels)
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryRedelegationsResponse{RedelegationResponses: redelResponses, Pagination: pageRes}, nil
//}
//
//// DelegatorValidators queries all validators info for given delegator address
//func (k Querier) DelegatorValidators(c context.Context, req *typesadapter.QueryDelegatorValidatorsRequest) (*typesadapter.QueryDelegatorValidatorsResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//
//	if req.DelegatorAddr == "" {
//		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
//	}
//	var validators typesadapter.Validators
//	ctx := sdk.UnwrapSDKContext(c)
//
//	store := ctx.KVStore(k.storeKey)
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	delStore := prefix.NewStore(store, typesadapter.GetDelegationsKey(delAddr))
//	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
//		delegation, err := typesadapter.UnmarshalDelegation(k.cdc, value)
//		if err != nil {
//			return err
//		}
//
//		validator, found := k.GetValidator(ctx, delegation.GetValidatorAddr())
//		if !found {
//			return typesadapter.ErrNoValidatorFound
//		}
//
//		validators = append(validators, validator)
//		return nil
//	})
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &typesadapter.QueryDelegatorValidatorsResponse{Validators: validators, Pagination: pageRes}, nil
//}
//
//// Pool queries the pool info
//func (k Querier) Pool(c context.Context, _ *typesadapter.QueryPoolRequest) (*typesadapter.QueryPoolResponse, error) {
//	ctx := sdk.UnwrapSDKContext(c)
//	bondDenom := k.BondDenom(ctx)
//	bondedPool := k.GetBondedPool(ctx)
//	notBondedPool := k.GetNotBondedPool(ctx)
//
//	pool := typesadapter.NewPool(
//		k.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount,
//		k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount,
//	)
//
//	return &typesadapter.QueryPoolResponse{Pool: pool}, nil
//}
//
//// Params queries the staking parameters
//func (k Querier) Params(c context.Context, _ *typesadapter.QueryParamsRequest) (*typesadapter.QueryParamsResponse, error) {
//	ctx := sdk.UnwrapSDKContext(c)
//	params := k.GetParams(ctx)
//
//	return &typesadapter.QueryParamsResponse{Params: params}, nil
//}
//
//func queryRedelegation(ctx sdk.Context, k Querier, req *typesadapter.QueryRedelegationsRequest) (redels typesadapter.Redelegations, err error) {
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	srcValAddr, err := sdk.ValAddressFromBech32(req.SrcValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	dstValAddr, err := sdk.ValAddressFromBech32(req.DstValidatorAddr)
//	if err != nil {
//		return nil, err
//	}
//
//	redel, found := k.GetRedelegation(ctx, delAddr, srcValAddr, dstValAddr)
//	if !found {
//		return nil, status.Errorf(
//			codes.NotFound,
//			"redelegation not found for delegator address %s from validator address %s",
//			req.DelegatorAddr, req.SrcValidatorAddr)
//	}
//	redels = []typesadapter.Redelegation{redel}
//
//	return redels, err
//}
//
//func queryRedelegationsFromSrcValidator(store sdk.KVStore, k Querier, req *typesadapter.QueryRedelegationsRequest) (redels typesadapter.Redelegations, res *query.PageResponse, err error) {
//	valAddr, err := sdk.ValAddressFromBech32(req.SrcValidatorAddr)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	srcValPrefix := typesadapter.GetREDsFromValSrcIndexKey(valAddr)
//	redStore := prefix.NewStore(store, srcValPrefix)
//	res, err = query.Paginate(redStore, req.Pagination, func(key []byte, value []byte) error {
//		storeKey := typesadapter.GetREDKeyFromValSrcIndexKey(append(srcValPrefix, key...))
//		storeValue := store.Get(storeKey)
//		red, err := typesadapter.UnmarshalRED(k.cdc, storeValue)
//		if err != nil {
//			return err
//		}
//		redels = append(redels, red)
//		return nil
//	})
//
//	return redels, res, err
//}
//
//func queryAllRedelegations(store sdk.KVStore, k Querier, req *typesadapter.QueryRedelegationsRequest) (redels typesadapter.Redelegations, res *query.PageResponse, err error) {
//	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	redStore := prefix.NewStore(store, typesadapter.GetREDsKey(delAddr))
//	res, err = query.Paginate(redStore, req.Pagination, func(key []byte, value []byte) error {
//		redelegation, err := typesadapter.UnmarshalRED(k.cdc, value)
//		if err != nil {
//			return err
//		}
//		redels = append(redels, redelegation)
//		return nil
//	})
//
//	return redels, res, err
//}
