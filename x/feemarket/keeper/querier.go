package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/feemarket/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		if len(path) < 1 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
				"Insufficient parameters, at least 1 parameter is required")
		}

		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, keeper)
		case types.QueryBaseFee:
			return queryBaseFee(ctx, keeper)
		case types.QueryBlockGas:
			return queryBlockGas(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}
func queryParams(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	params := keeper.GetParams(ctx)
	res := &types.QueryParamsResponse{
		Params: params,
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, res)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
func queryBaseFee(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	bnRes := types.QueryBaseFeeResponse{}
	baseFee := keeper.GetBaseFee(ctx)
	if baseFee != nil {
		aux := sdk.NewIntFromBigInt(baseFee)
		bnRes.BaseFee = aux.BigInt()
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, bnRes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
func queryBlockGas(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	gas := keeper.GetBlockGasUsed(ctx)

	bnRes := &types.QueryBlockGasResponse{
		Gas: int64(gas),
	}
	bz, err := codec.MarshalJSONIndent(keeper.cdc, bnRes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
