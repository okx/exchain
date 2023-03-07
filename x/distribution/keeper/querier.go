package keeper

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	comm "github.com/okx/okbchain/x/common"

	"github.com/okx/okbchain/x/distribution/types"
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
	case types.ParamWithdrawRewardEnabled:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetWithdrawRewardEnabled(ctx))
		if err != nil {
			return nil, comm.ErrMarshalJSONFailed(err.Error())
		}
		return bz, nil
	case types.ParamRewardTruncatePrecision:
		bz, err := codec.MarshalJSONIndent(k.cdc, k.GetRewardTruncatePrecision(ctx))
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
