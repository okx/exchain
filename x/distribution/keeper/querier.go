package keeper

import (
	"encoding/json"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	comm "github.com/okex/exchain/x/common"

	"github.com/okex/exchain/x/distribution/types"
)

// NewQuerier creates a querier for distribution REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, path[1:], req, k)

		case types.QueryValidatorCommission:
			return queryValidatorCommission(ctx, path[1:], req, k)

		case types.QueryWithdrawAddr:
			return queryDelegatorWithdrawAddress(ctx, path[1:], req, k)

		case types.QueryCommunityPool:
			return queryCommunityPool(ctx, path[1:], req, k)

		case types.QueryDelegatorValidators:
			return queryDelegatorValidators(ctx, path[1:], req, k)

		case types.QueryDelegationRewards:
			return queryDelegationRewards(ctx, path[1:], req, k)

		case types.QueryDelegatorTotalRewards:
			return queryDelegatorTotalRewards(ctx, path[1:], req, k)

		case types.QueryValidatorOutstandingRewards:
			return queryValidatorOutstandingRewards(ctx, path[1:], req, k)

		default:
			return nil, types.ErrUnknownDistributionQueryType()
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	switch path[0] {
	case types.ParamCommunityTax:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetCommunityTax(ctx))
		if err != nil {
			return nil, comm.ErrMarshalJSONFailed(err.Error())
		}
		return bz, nil
	case types.ParamWithdrawAddrEnabled:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetWithdrawAddrEnabled(ctx))
		if err != nil {
			return nil, comm.ErrMarshalJSONFailed(err.Error())
		}
		return bz, nil
	case types.ParamDistributionType:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetDistributionType(ctx))
		if err != nil {
			return nil, comm.ErrMarshalJSONFailed(err.Error())
		}
		return bz, nil
	default:
		return nil, types.ErrUnknownDistributionParamType()
	}
}

func queryValidatorCommission(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryValidatorCommissionParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, comm.ErrUnMarshalJSONFailed(err.Error())
	}

	commission := k.GetValidatorAccumulatedCommission(ctx, params.ValidatorAddress)
	if commission == nil { //TODO
		commission = types.ValidatorAccumulatedCommission{}
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, commission)
	if err != nil {
		return nil, comm.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

func queryDelegatorWithdrawAddress(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorWithdrawAddrParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, comm.ErrUnMarshalJSONFailed(err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, params.DelegatorAddress)

	bz, err := codec.MarshalJSONIndent(k.cdc, withdrawAddr)
	if err != nil {
		return nil, comm.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

func queryCommunityPool(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	pool := k.GetFeePoolCommunityCoins(ctx)
	if pool == nil {
		pool = sdk.SysCoins{}
	}

	bz, err := k.cdc.MarshalJSON(pool)
	if err != nil {
		return nil, comm.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

func queryDelegatorValidators(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	delegator := k.stakingKeeper.Delegator(ctx, params.DelegatorAddress)

	bz, err := codec.MarshalJSONIndent(k.cdc, delegator.GetShareAddedValidatorAddresses())
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDelegationRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	if !k.CheckDistributionProposalValid(ctx) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidVersion, "not support")
	}

	var params types.QueryDelegationRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	val := k.stakingKeeper.Validator(ctx, params.ValidatorAddress)
	if val == nil {
		return nil, sdkerrors.Wrap(types.ErrCodeEmptyValidatorDistInfo(), params.ValidatorAddress.String())
	}

	del := k.stakingKeeper.Delegator(ctx, params.DelegatorAddress)
	if del == nil {
		return nil, types.ErrCodeEmptyDelegationDistInfo()
	}

	found := false
	for _, valAddr := range del.GetShareAddedValidatorAddresses() {
		if valAddr.Equals(params.ValidatorAddress) {
			found = true
		}
	}
	if !found {
		return nil, sdkerrors.Wrap(types.ErrCodeEmptyValidatorDistInfo(), params.ValidatorAddress.String())
	}

	logger := k.Logger(ctx)
	if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), params.DelegatorAddress) && !del.GetLastAddedShares().IsZero() {
		k.initExistedDelegationStartInfo(ctx, val, del)
	}

	endingPeriod := k.incrementValidatorPeriod(ctx, val)
	rewards := k.calculateDelegationRewards(ctx, val, params.DelegatorAddress, endingPeriod)
	if rewards == nil {
		rewards = sdk.DecCoins{}
	}

	logger.Debug("queryDelegationRewards", "Validator", val.GetOperator(),
		"Delegator", params.DelegatorAddress, "Reward", rewards)

	bz, err := codec.MarshalJSONIndent(k.cdc, rewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDelegatorTotalRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	if !k.CheckDistributionProposalValid(ctx) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidVersion, "not support")
	}

	var params types.QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	del := k.stakingKeeper.Delegator(ctx, params.DelegatorAddress)
	if del == nil {
		return nil, types.ErrCodeEmptyDelegationDistInfo()
	}

	total := sdk.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	for _, valAddr := range del.GetShareAddedValidatorAddresses() {
		val := k.stakingKeeper.Validator(ctx, valAddr)
		if val == nil {
			continue
		}

		logger := k.Logger(ctx)
		if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), params.DelegatorAddress) && !del.GetLastAddedShares().IsZero() {
			k.initExistedDelegationStartInfo(ctx, val, del)
		}

		endingPeriod := k.incrementValidatorPeriod(ctx, val)
		delReward := k.calculateDelegationRewards(ctx, val, params.DelegatorAddress, endingPeriod)
		if delReward == nil {
			delReward = sdk.DecCoins{}
		}
		delRewards = append(delRewards, types.NewDelegationDelegatorReward(valAddr, delReward))
		total = total.Add(delReward...)
		logger.Debug("queryDelegatorTotalRewards", "Validator", val.GetOperator(),
			"Delegator", params.DelegatorAddress, "Reward", delReward)
	}

	totalRewards := types.NewQueryDelegatorTotalRewardsResponse(delRewards, total)

	bz, err := json.Marshal(totalRewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryValidatorOutstandingRewards(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	if !k.CheckDistributionProposalValid(ctx) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidVersion, "not support")
	}

	var params types.QueryValidatorOutstandingRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	rewards := k.GetValidatorOutstandingRewards(ctx, params.ValidatorAddress)
	if rewards == nil {
		rewards = sdk.DecCoins{}
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, rewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
