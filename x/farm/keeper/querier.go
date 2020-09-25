package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// NewQuerier creates a new querier for farm clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryPool:
			return queryPool(ctx, req, k)
		case types.QueryPools:
			return queryPools(ctx, k)
		case types.QueryEarnings:
			return queryEarnings(ctx, req, k)
		case types.QueryParams:
			return queryParams(ctx, k)
		case types.QueryWhitelist:
			return queryWhitelist(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("failed. unknown farm query endpoint")
		}
	}
}

// TODO: queriers
func queryPool(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	return nil, nil
}

func queryPools(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	// TODO: get pools from ctx with keeper
	pools := types.NewTestStruct("test pools")
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pools)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryEarnings(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	return nil, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	// TODO: get params from ctx with keeper
	params := types.NewTestStruct("test params")
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func queryWhitelist(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	// TODO: get whitelist from ctx with keeper
	whitelist := types.NewTestStruct("test whitelist")
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, whitelist)
	if err != nil {
		return nil, defaultQueryErrJSONMarshal(err)
	}

	return res, nil
}

func defaultQueryErrJSONMarshal(err error) sdk.Error {
	return sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", err.Error()))
}
