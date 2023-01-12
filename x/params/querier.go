package params

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	"github.com/okex/exchain/x/params/types"
)

// NewQuerier returns all query handlers
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, req, keeper)
		case types.QueryUpgrade:
			if len(path) != 2 {
				keeper.Logger(ctx).Error("invalid query path", "path", path)
			}
			return queryUpgrade(ctx, path[1], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown params query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, _ abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, fmt.Sprintf("could not marshal result to JSON %s", err.Error()))
	}
	return bz, nil
}

func queryUpgrade(ctx sdk.Context, name string, keeper Keeper) ([]byte, sdk.Error) {
	info, err := keeper.readUpgradeInfo(ctx, name)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, info)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, fmt.Sprintf("could not marshal result to JSON %s", err.Error()))
	}
	return bz, nil
}
