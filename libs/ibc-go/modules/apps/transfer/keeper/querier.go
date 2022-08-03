package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

// NewQuerier creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryDenomTraces:
			return queryDenomTraces(ctx, req, k)
		default:
			return nil, types.ErrUnexpectedEndOfGroupTransfer
		}
	}
}

func queryDenomTraces(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDenomTracesRequest

	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	//delegationResps, err := k.DelegatorDelegations(ctx, &params)

	denomTracesResp, err := k.DenomTraces(sdk.WrapSDKContext(ctx), &params)
	if denomTracesResp == nil {
		denomTracesResp = &types.QueryDenomTracesResponse{}
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, denomTracesResp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}
