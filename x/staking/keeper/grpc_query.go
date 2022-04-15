package keeper

import ()

//func (k Keeper) ValidatorDelegations(c context.Context, adar string, pageReq *query.PageRequest) (*types.QueryValidatorDelegationsResponse, error) {
//
//	if adar == "" {
//		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
//	}
//	var delegations []types.Delegation
//	ctx := sdk.UnwrapSDKContext(c)
//
//	store := ctx.KVStore(k.storeKey)
//	valStore := prefix.NewStore(store, types.DelegationKey)
//	pageRes, err := query.FilteredPaginate(valStore, pageReq, func(key []byte, value []byte, accumulate bool) (bool, error) {
//		delegation, err := types.UnmarshalDelegation(k.cdcMarshl.GetCdc(), value)
//		if err != nil {
//			return false, err
//		}
//
//		valAddr, err := sdk.ValAddressFromBech32(adar)
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
//	delResponses, err := DelegationsToDelegationResponses(ctx, k, delegations)
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &types.QueryValidatorDelegationsResponse{
//		DelegationResponses: delResponses, Pagination: pageRes}, nil
//}
