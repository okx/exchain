package keeper

import ()

//func queryUndelegation2(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
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
//	k.GetUndelegating()
//	unbond, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
//	if !found {
//		return nil, status.Errorf(
//			codes.NotFound,
//			"unbonding delegation with delegator %s not found for validator %s",
//			req.DelegatorAddr, req.ValidatorAddr)
//	}
//
//	return &types.QueryUnbondingDelegationResponse{Unbond: unbond}, nil
//
//	var params types.QueryDelegatorParams
//	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
//	if err != nil {
//		return nil, common.ErrUnMarshalJSONFailed(err.Error())
//	}
//
//	undelegation, found := k.GetUndelegating(ctx, params.DelegatorAddr)
//	if !found {
//		return nil, types.ErrNoUnbondingDelegation()
//	}
//
//	res, err := codec.MarshalJSONIndent(types.ModuleCdc, undelegation)
//	if err != nil {
//		return nil, common.ErrMarshalJSONFailed(err.Error())
//	}
//
//	return res, nil
//}
