package keeper

import (
	"encoding/json"

	"github.com/okex/exchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/ammswap/types"
)

// NewQuerier creates a new querier for swap clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryBuyAmount:
			res, err = queryBuyAmount(ctx, path[1:], req, k)
		case types.QuerySwapTokenPair:
			res, err = querySwapTokenPair(ctx, path[1:], req, k)
		case types.QueryParams:
			res, err = queryParams(ctx, path[1:], req, k)
		case types.QuerySwapTokenPairs:
			res, err = querySwapTokenPairs(ctx, path[1:], req, k)
		case types.QueryRedeemableAssets:
			res, err = queryRedeemableAssets(ctx, path[1:], req, k)
		case types.QuerySwapQuoteInfo:
			res, err = querySwapQuoteInfo(ctx, req, k)
		case types.QuerySwapAddLiquidityQuote:
			res, err = querySwapAddLiquidityQuote(ctx, req, k)

		default:
			return nil, types.ErrSwapUnknownQueryType()
		}

		if err != nil {
			response := common.GetErrorResponse(types.CodeInternalError, "", err.Error())
			res, errJSON := json.Marshal(response)
			if errJSON != nil {
				return nil, common.ErrMarshalJSONFailed(errJSON.Error())
			}
			return res, err
		}

		return res, nil
	}
}

// nolint
func querySwapTokenPair(
	ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper,
) ([]byte, sdk.Error) {
	tokenPairName := path[0]
	var response *common.BaseResponse
	tokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
	// return nil when token pair not exists
	if err != nil {
		response = common.GetBaseResponse(nil)
	} else {
		response = common.GetBaseResponse(tokenPair)
	}

	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// nolint
func queryBuyAmount(
	ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper,
) ([]byte, sdk.Error) {
	var queryParams types.QueryBuyAmountParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	errToken := types.ValidateSwapAmountName(queryParams.TokenToBuy)
	if errToken != nil {
		return nil, errToken
	}
	errToken = types.ValidateSwapAmountName(queryParams.SoldToken.Denom)
	if errToken != nil {
		return nil, errToken
	}
	params := keeper.GetParams(ctx)
	var buyAmount sdk.Dec
	swapTokenPair := types.GetSwapTokenPairName(queryParams.SoldToken.Denom, queryParams.TokenToBuy)
	tokenPair, errTokenPair := keeper.GetSwapTokenPair(ctx, swapTokenPair)
	if errTokenPair == nil {
		if tokenPair.BasePooledCoin.IsZero() || tokenPair.QuotePooledCoin.IsZero() {
			return nil, types.ErrIsZeroValue("base pooled coin or quote pooled coin")
		}
		buyAmount = CalculateTokenToBuy(tokenPair, queryParams.SoldToken, queryParams.TokenToBuy, params).Amount
	} else {
		tokenPairName1 := types.GetSwapTokenPairName(queryParams.SoldToken.Denom, sdk.DefaultBondDenom)
		tokenPair1, err := keeper.GetSwapTokenPair(ctx, tokenPairName1)
		if err != nil {
			return nil, err
		}
		if tokenPair1.BasePooledCoin.IsZero() || tokenPair1.QuotePooledCoin.IsZero() {
			return nil, types.ErrIsZeroValue("base pooled coin or quote pooled coin")
		}
		tokenPairName2 := types.GetSwapTokenPairName(queryParams.TokenToBuy, sdk.DefaultBondDenom)
		tokenPair2, err := keeper.GetSwapTokenPair(ctx, tokenPairName2)
		if err != nil {
			return nil, err
		}
		if tokenPair2.BasePooledCoin.IsZero() || tokenPair2.QuotePooledCoin.IsZero() {
			return nil, types.ErrIsZeroValue("base pooled coin or quote pooled coin")
		}
		nativeToken := CalculateTokenToBuy(tokenPair1, queryParams.SoldToken, sdk.DefaultBondDenom, params)
		buyAmount = CalculateTokenToBuy(tokenPair2, nativeToken, queryParams.TokenToBuy, params).Amount
	}

	bz := keeper.cdc.MustMarshalJSON(buyAmount)

	return bz, nil
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	return keeper.cdc.MustMarshalJSON(keeper.GetParams(ctx)), nil
}

// nolint
func querySwapTokenPairs(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte,
	err sdk.Error) {
	return keeper.cdc.MustMarshalJSON(keeper.GetSwapTokenPairs(ctx)), nil
}

// nolinte
func queryRedeemableAssets(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	tokenPairName := path[0]
	swapTokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
	if err != nil {
		return nil, err
	}

	liquidity, errDec := sdk.NewDecFromStr(path[1])
	if errDec != nil {
		return nil, errDec

	}
	var tokenList sdk.SysCoins
	baseToken, quoteToken, err := keeper.GetRedeemableAssets(ctx, swapTokenPair.BasePooledCoin.Denom, swapTokenPair.QuotePooledCoin.Denom, liquidity)
	if err != nil {
		return nil, err
	}
	tokenList = append(tokenList, baseToken, quoteToken)
	bz := keeper.cdc.MustMarshalJSON(tokenList)
	return bz, nil
}

// querySwapQuoteInfo returns infos when swap token
func querySwapQuoteInfo(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapBuyInfoParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	if queryParams.SellTokenAmount == "" || queryParams.BuyToken == "" {
		return nil, types.ErrSellAmountOrBuyTokenIsEmpty()
	}

	sellAmount, err := sdk.ParseDecCoin(queryParams.SellTokenAmount)
	if err != nil {
		return nil, types.ErrConvertSellTokenAmount(queryParams.SellTokenAmount, err)
	}

	if sellAmount.Denom == queryParams.BuyToken {
		return nil, types.ErrSellAmountEqualBuyToken()
	}

	var route string
	var fee sdk.SysCoin
	buyAmount := sdk.ZeroDec()
	marketPrice := sdk.ZeroDec()

	swapParams := keeper.GetParams(ctx)
	tokenPairName := types.GetSwapTokenPairName(sellAmount.Denom, queryParams.BuyToken)
	tokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
	if err == nil {
		if tokenPair.BasePooledCoin.Amount.IsZero() || tokenPair.QuotePooledCoin.IsZero() {
			return nil, types.ErrIsZeroValue("base pooled coin or quote pooled coin")
		}
		buyAmount = CalculateTokenToBuy(tokenPair, sellAmount, queryParams.BuyToken, swapParams).Amount
		// calculate market price
		if tokenPair.BasePooledCoin.Denom == sellAmount.Denom {
			marketPrice = tokenPair.QuotePooledCoin.Amount.Quo(tokenPair.BasePooledCoin.Amount)
		} else {
			marketPrice = tokenPair.BasePooledCoin.Amount.Quo(tokenPair.QuotePooledCoin.Amount)
		}
		// calculate fee
		fee = sdk.NewDecCoinFromDec(sellAmount.Denom, sellAmount.Amount.Mul(swapParams.FeeRate))
	} else {
		tokenPairName1 := types.GetSwapTokenPairName(sellAmount.Denom, common.NativeToken)
		tokenPair1, err := keeper.GetSwapTokenPair(ctx, tokenPairName1)
		if err != nil {
			return nil, err
		}
		if tokenPair1.BasePooledCoin.Amount.IsZero() || tokenPair1.QuotePooledCoin.IsZero() {
			return nil, types.ErrIsZeroValue("base pooled coin or quote pooled coin")
		}
		tokenPairName2 := types.GetSwapTokenPairName(queryParams.BuyToken, common.NativeToken)
		tokenPair2, err := keeper.GetSwapTokenPair(ctx, tokenPairName2)
		if err != nil {
			return nil, err
		}
		if tokenPair2.BasePooledCoin.Amount.IsZero() || tokenPair2.QuotePooledCoin.IsZero() {
			return nil, types.ErrIsZeroValue("base pooled coin or quote pooled coin")
		}
		nativeToken := CalculateTokenToBuy(tokenPair1, sellAmount, common.NativeToken, swapParams)
		buyAmount = CalculateTokenToBuy(tokenPair2, nativeToken, queryParams.BuyToken, swapParams).Amount

		// calculate market price
		var sellTokenMarketPrice sdk.Dec
		var routeTokenMarketPrice sdk.Dec
		if tokenPair1.BasePooledCoin.Denom == sellAmount.Denom {
			sellTokenMarketPrice = tokenPair1.QuotePooledCoin.Amount.Quo(tokenPair1.BasePooledCoin.Amount)
		} else {
			sellTokenMarketPrice = tokenPair1.BasePooledCoin.Amount.Quo(tokenPair1.QuotePooledCoin.Amount)
		}
		if tokenPair2.BasePooledCoin.Denom == common.NativeToken {
			routeTokenMarketPrice = tokenPair2.QuotePooledCoin.Amount.Quo(tokenPair2.BasePooledCoin.Amount)
		} else {
			routeTokenMarketPrice = tokenPair2.BasePooledCoin.Amount.Quo(tokenPair2.QuotePooledCoin.Amount)
		}
		if routeTokenMarketPrice.IsPositive() && sellTokenMarketPrice.IsPositive() {
			marketPrice = sellTokenMarketPrice.Mul(routeTokenMarketPrice)
		}

		// calculate fee
		fee1 := sdk.NewDecCoinFromDec(sellAmount.Denom, sellAmount.Amount.Mul(swapParams.FeeRate))
		routeTokenFee := sdk.NewDecCoinFromDec(common.NativeToken, nativeToken.Amount.Mul(swapParams.FeeRate))
		fee2 := CalculateTokenToBuy(tokenPair1, routeTokenFee, sellAmount.Denom, swapParams)
		fee = fee1.Add(fee2)

		// swap by route
		route = common.NativeToken
	}

	// calculate price
	price := sdk.ZeroDec()
	if sellAmount.Amount.IsPositive() {
		price = buyAmount.Quo(sellAmount.Amount)
	}

	// calculate price impact
	var priceImpact sdk.Dec
	if marketPrice.IsPositive() {
		if marketPrice.GT(price) {
			priceImpact = marketPrice.Sub(price).Quo(marketPrice)
		} else {
			priceImpact = price.Sub(marketPrice).Quo(marketPrice)
		}
	}

	swapBuyInfo := types.SwapBuyInfo{
		BuyAmount:   buyAmount,
		Price:       price,
		PriceImpact: priceImpact,
		Fee:         fee.String(),
		Route:       route,
	}

	response := common.GetBaseResponse(swapBuyInfo)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil

}

// querySwapAddLiquidityQuote returns swap information of adding liquidity
func querySwapAddLiquidityQuote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapAddInfoParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	// check params
	if queryParams.QuoteTokenAmount == "" {
		return nil, types.ErrQueryParamsQuoteTokenAmountIsEmpty()
	}
	if queryParams.BaseToken == "" {
		return nil, types.ErrQueryParamsBaseTokenIsEmpty()
	}
	queryTokenAmount, err := sdk.ParseDecCoin(queryParams.QuoteTokenAmount)
	if err != nil {
		return nil, types.ErrConvertQuoteTokenAmount(queryParams.QuoteTokenAmount, err)
	}

	tokenPairName := types.GetSwapTokenPairName(queryParams.BaseToken, queryTokenAmount.Denom)
	swapTokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
	if err != nil {
		return nil, err
	}
	if swapTokenPair.BasePooledCoin.Amount.IsZero() && swapTokenPair.QuotePooledCoin.Amount.IsZero() {
		return nil, types.ErrIsZeroValue("base pooled coin or quote pooled coin")
	}

	totalSupply := keeper.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
	if totalSupply.IsZero() {
		return nil, types.ErrIsZeroValue("total supply")
	}

	var addAmount sdk.Dec
	var liquidity sdk.Dec
	if swapTokenPair.BasePooledCoin.Denom == queryTokenAmount.Denom {
		addAmount = common.MulAndQuo(queryTokenAmount.Amount, swapTokenPair.QuotePooledCoin.Amount, swapTokenPair.BasePooledCoin.Amount)
		liquidity = common.MulAndQuo(addAmount, totalSupply, swapTokenPair.QuotePooledCoin.Amount)
	} else {
		addAmount = common.MulAndQuo(queryTokenAmount.Amount, swapTokenPair.BasePooledCoin.Amount, swapTokenPair.QuotePooledCoin.Amount)
		liquidity = common.MulAndQuo(queryTokenAmount.Amount, totalSupply, swapTokenPair.QuotePooledCoin.Amount)
	}
	addInfo := types.SwapAddInfo{
		BaseTokenAmount: addAmount,
		PoolShare:       liquidity.Quo(totalSupply.Add(liquidity)),
		Liquidity:       liquidity,
	}
	response := common.GetBaseResponse(addInfo)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil

}
