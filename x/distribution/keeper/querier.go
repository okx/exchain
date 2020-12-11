package keeper

import (
	comm "github.com/okex/okexchain/x/common"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okexchain/x/distribution/types"
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

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
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

	default:
		return nil, types.ErrUnknownRequest()
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
