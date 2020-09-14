package keeper

import (
	"github.com/okex/okexchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/ammswap/types"
)

// NewQuerier creates a new querier for swap clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QuerySwapTokenPair:
			return querySwapTokenPair(ctx, path[1:], req, k)
		case types.QueryParams:
			return queryParams(ctx, path[1:], req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown swap query endpoint")
		}
	}
}

// nolint
func querySwapTokenPair(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte,
	err sdk.Error) {
	tokenPairName := path[0] + "_" + common.NativeToken
	tokenPair, error := keeper.GetSwapTokenPair(ctx, tokenPairName)
	if error != nil {
		return nil, sdk.ErrUnknownRequest(error.Error())
	}
	bz := keeper.cdc.MustMarshalJSON(tokenPair)
	return bz, nil
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	return keeper.cdc.MustMarshalJSON(keeper.GetParams(ctx)), nil
}
