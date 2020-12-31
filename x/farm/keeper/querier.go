package keeper

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okexchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

const (
	defaultPoolsDisplayedNum = 20
)

// NewQuerier creates a new querier for farm clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryPool:
			return queryPool(ctx, req, k)
		case types.QueryPools:
			return queryPools(ctx, req, k)
		case types.QueryEarnings:
			return queryEarnings(ctx, req, k)
		case types.QueryLockInfo:
			return queryLockInfo(ctx, req, k)
		case types.QueryParameters:
			return queryParams(ctx, k)
		case types.QueryWhitelist:
			return queryWhitelist(ctx, k)
		case types.QueryAccount:
			return queryAccount(ctx, req, k)
		case types.QueryAccountsLockedTo:
			return queryAccountsLockedTo(ctx, req, k)
		case types.QueryPoolNum:
			return queryPoolNum(ctx, k)
		default:
			return nil, types.ErrUnknownFarmQueryType("failed. unknown farm query endpoint")
		}
	}
}

func queryPool(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryPoolParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, defaultQueryErrParseParams(err)
	}

	pool, found := k.GetFarmPool(ctx, params.PoolName)
	if !found {
		return nil, types.ErrNoFarmPoolFound(params.PoolName)
	}

	updatedPool, _ := k.CalculateAmountYieldedBetween(ctx, pool)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, updatedPool)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

// support query by page && limit
func queryPools(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryPoolsParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, defaultQueryErrParseParams(err)
	}

	pools := k.GetFarmPools(ctx)
	var updatedPools types.FarmPools
	for _, pool := range pools {
		updatedPool, _ := k.CalculateAmountYieldedBetween(ctx, pool)
		updatedPools = append(updatedPools, updatedPool)
	}

	if !(params.Page == 1 && params.Limit == 0) {
		start, end := client.Paginate(len(updatedPools), params.Page, params.Limit, defaultPoolsDisplayedNum)
		if start < 0 || end < 0 {
			start, end = 0, 0
		}
		updatedPools = updatedPools[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, updatedPools)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryEarnings(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryPoolAccountParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, defaultQueryErrParseParams(err)
	}

	earnings, sdkErr := k.getEarnings(ctx, params.PoolName, params.AccAddress)
	if sdkErr != nil {
		return nil, sdkErr
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, earnings)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryLockInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryPoolAccountParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, defaultQueryErrParseParams(err)
	}

	lockInfo, found := k.GetLockInfo(ctx, params.AccAddress, params.PoolName)
	if !found {
		return nil, types.ErrNoLockInfoFound(params.AccAddress.String(), params.PoolName)
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, lockInfo)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryWhitelist(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	whitelist := k.GetWhitelist(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, whitelist)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryAccount(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAccountParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, defaultQueryErrParseParams(err)
	}

	poolNames := k.getFarmPoolNamesForAccount(ctx, params.AccAddress)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, poolNames)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryAccountsLockedTo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryPoolParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, defaultQueryErrParseParams(err)
	}

	accAddrList := k.getAccountsLockedTo(ctx, params.PoolName)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, accAddrList)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryPoolNum(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	poolNum := k.getPoolNum(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, poolNum)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func defaultQueryErrJSONMarshal(err error) sdk.Error {
	return common.ErrMarshalJSONFailed(err.Error())
}

func defaultQueryErrParseParams(err error) sdk.Error {
	return common.ErrUnMarshalJSONFailed(err.Error())
}
