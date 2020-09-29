package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
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
		case types.QueryParameters:
			return queryParams(ctx, k)
		case types.QueryWhitelist:
			return queryWhitelist(ctx, k)
		case types.QueryAccount:
			return queryAccount(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("failed. unknown farm query endpoint")
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
		return nil, types.ErrNoFarmPoolFound(types.DefaultCodespace, params.PoolName)
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
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
	if !(params.Page == 1 && params.Limit == 0) {
		start, end := client.Paginate(len(pools), params.Page, params.Limit, defaultPoolsDisplayedNum)
		if start < 0 || end < 0 {
			start, end = 0, 0
		}
		pools = pools[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pools)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryEarnings(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryEarningsParams
	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, defaultQueryErrParseParams(err)
	}

	earnings, sdkErr := k.GetEarnings(ctx, params.PoolName, params.AccAddress)
	if sdkErr != nil {
		return nil, sdkErr
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, earnings)
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

	poolNames := k.GetFarmPoolNamesForAccount(ctx, params.AccAddress)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, poolNames)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func defaultQueryErrJSONMarshal(err error) sdk.Error {
	return sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", err.Error()))
}

func defaultQueryErrParseParams(err error) sdk.Error {
	return sdk.ErrInternal(fmt.Sprintf("failed to parse params. %s", err))
}
