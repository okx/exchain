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
		case types.QueryBuyAmount:
			return queryBuyAmount(ctx, path[1:], req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown swap query endpoint")
		}
	}
}

// nolint
func querySwapTokenPair(
	ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper,
) (res []byte, err sdk.Error) {
	tokenPairName := path[0] + "_" + common.NativeToken
	tokenPair, error := keeper.GetSwapTokenPair(ctx, tokenPairName)
	if error != nil {
		return nil, sdk.ErrUnknownRequest(error.Error())
	}
	bz := keeper.cdc.MustMarshalJSON(tokenPair)
	return bz, nil
}

// nolint
func queryBuyAmount(
	ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper,
) ([]byte, sdk.Error) {
	var queryParams types.QueryBuyAmountParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	params := keeper.GetParams(ctx)
	var buyAmount sdk.Dec
	if (queryParams.SellToken.Denom == sdk.DefaultBondDenom) {
		tokenPairName := queryParams.BuyTokenName + "_" + queryParams.SellToken.Denom
		tokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
		if err != nil {
			return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
		}
		buyAmount = CalculateTokenToBuy(tokenPair, queryParams.SellToken, queryParams.BuyTokenName, params).Amount
	} else if (queryParams.BuyTokenName == sdk.DefaultBondDenom) {
		tokenPairName := queryParams.SellToken.Denom + "_" + queryParams.BuyTokenName
		tokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
		if err != nil {
			return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
		}
		buyAmount = CalculateTokenToBuy(tokenPair, queryParams.SellToken, queryParams.BuyTokenName, params).Amount
	} else {
		tokenPairName1 := queryParams.SellToken.Denom + "_" + sdk.DefaultBondDenom
		tokenPair1, err := keeper.GetSwapTokenPair(ctx, tokenPairName1)
		if err != nil {
			return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
		}

		tokenPairName2 := queryParams.BuyTokenName + "_" + sdk.DefaultBondDenom
		tokenPair2, err := keeper.GetSwapTokenPair(ctx, tokenPairName2)
		if err != nil {
			return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
		}

		nativeToken := CalculateTokenToBuy(tokenPair1, queryParams.SellToken, sdk.DefaultBondDenom, params)
		buyAmount = CalculateTokenToBuy(tokenPair2, nativeToken, queryParams.BuyTokenName, params).Amount
	}

	bz := keeper.cdc.MustMarshalJSON(buyAmount)

	return bz, nil
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	return keeper.cdc.MustMarshalJSON(keeper.GetParams(ctx)), nil
}
