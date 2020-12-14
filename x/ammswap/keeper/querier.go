package keeper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/okex/okexchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/ammswap/types"
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
		case types.QuerySwapTokens:
			res, err = querySwapTokens(ctx, req, k)
		case types.QuerySwapQuoteInfo:
			res, err = querySwapQuoteInfo(ctx, req, k)
		case types.QuerySwapLiquidityHistories:
			res, err = querySwapLiquidityHistories(ctx, req, k)
		case types.QuerySwapAddLiquidityQuote:
			res, err = querySwapAddLiquidityQuote(ctx, req, k)

		default:
			return nil, types.ErrUnknownRequest()
		}

		if err != nil {
			response := common.GetErrorResponse(types.CodeInternalError, "", err.Error())
			res, errJSON := json.Marshal(response)
			if errJSON != nil {
				return nil, types.ErrInternal()
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
		return nil, types.ErrInternal()
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
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}
	errToken := types.ValidateSwapAmountName(queryParams.TokenToBuy)
	if errToken != nil {
		return nil, types.ErrUnknownRequest()
	}
	errToken = types.ValidateSwapAmountName(queryParams.SoldToken.Denom)
	if errToken != nil {
		return nil, types.ErrUnknownRequest()
	}
	params := keeper.GetParams(ctx)
	var buyAmount sdk.Dec
	swapTokenPair := types.GetSwapTokenPairName(queryParams.SoldToken.Denom, queryParams.TokenToBuy)
	tokenPair, errTokenPair := keeper.GetSwapTokenPair(ctx, swapTokenPair)
	if errTokenPair == nil {
		if tokenPair.BasePooledCoin.IsZero() || tokenPair.QuotePooledCoin.IsZero() {
			return nil, types.ErrInternal()
		}
		buyAmount = CalculateTokenToBuy(tokenPair, queryParams.SoldToken, queryParams.TokenToBuy, params).Amount
	} else {
		tokenPairName1 := types.GetSwapTokenPairName(queryParams.SoldToken.Denom, sdk.DefaultBondDenom)
		tokenPair1, err := keeper.GetSwapTokenPair(ctx, tokenPairName1)
		if err != nil {
			return nil, types.ErrUnknownRequest()
		}
		if tokenPair1.BasePooledCoin.IsZero() || tokenPair1.QuotePooledCoin.IsZero() {
			return nil, types.ErrInternal()
		}
		tokenPairName2 := types.GetSwapTokenPairName(queryParams.TokenToBuy, sdk.DefaultBondDenom)
		tokenPair2, err := keeper.GetSwapTokenPair(ctx, tokenPairName2)
		if err != nil {
			return nil, types.ErrUnknownRequest()
		}
		if tokenPair2.BasePooledCoin.IsZero() || tokenPair2.QuotePooledCoin.IsZero() {
			return nil, types.ErrInternal()
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
		return nil, types.ErrUnknownRequest()
	}

	liquidity, errDec := sdk.NewDecFromStr(path[1])
	if errDec != nil {
		return nil, errDec

	}
	var tokenList sdk.SysCoins
	baseToken, quoteToken, err := keeper.GetRedeemableAssets(ctx, swapTokenPair.BasePooledCoin.Denom, swapTokenPair.QuotePooledCoin.Denom, liquidity)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}
	tokenList = append(tokenList, baseToken, quoteToken)
	bz := keeper.cdc.MustMarshalJSON(tokenList)
	return bz, nil
}

// querySwapTokens returns tokens which are supported to swap in ammswap module
func querySwapTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapTokensParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}

	if queryParams.BusinessType == "" {
		return nil, types.ErrUnknownRequest()
	}

	// coins in account
	var accountCoins sdk.SysCoins
	if queryParams.Address != "" {
		addr, err := sdk.AccAddressFromBech32(queryParams.Address)
		if err != nil {
			return nil, types.ErrInvalidAddress(fmt.Sprintf("invalid address：%s", queryParams.Address))
		}
		accountCoins = keeper.tokenKeeper.GetCoins(ctx, addr)
	}

	var tokens []string
	switch queryParams.BusinessType {
	case types.SwapBusinessTypeCreate:
		tokens = getSwapCreateLiquidityTokens(ctx, keeper)
	case types.SwapBusinessTypeAdd:
		tokens = getSwapAddLiquidityTokens(ctx, keeper, queryParams.BaseTokenName)
	case types.SwapBusinessTypeSwap:
		tokens = getSwapTokens(ctx, keeper, queryParams.BaseTokenName)
	}

	swapTokensMap := make(map[string]sdk.Dec, len(tokens))
	for _, token := range tokens {
		swapTokensMap[token] = sdk.ZeroDec()
	}

	// update amount by coins in account
	for _, coin := range accountCoins {
		if _, ok := swapTokensMap[coin.Denom]; ok {
			swapTokensMap[coin.Denom] = coin.Amount
		}
	}

	// sort token list by account balance
	var swapTokens types.SwapTokens
	for symbol, available := range swapTokensMap {
		swapTokens = append(swapTokens, types.NewSwapToken(symbol, available))
	}
	sort.Sort(swapTokens)

	swapTokensResp := types.SwapTokensResponse{
		NativeToken: common.NativeToken,
		Tokens:      swapTokens,
	}

	response := common.GetBaseResponse(swapTokensResp)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, types.ErrInternal()
	}
	return bz, nil
}

func getSwapCreateLiquidityTokens(ctx sdk.Context, keeper Keeper) []string {
	var tokens []string
	allTokens := keeper.tokenKeeper.GetTokensInfo(ctx)
	for _, token := range allTokens {
		if !strings.HasPrefix(token.Symbol, types.PoolTokenPrefix) {
			tokens = append(tokens, token.Symbol)
		}
	}
	return tokens
}

func getSwapAddLiquidityTokens(ctx sdk.Context, keeper Keeper, baseTokenName string) []string {
	var tokens []string

	// all swap token pairs
	swapTokenPairs := keeper.GetSwapTokenPairs(ctx)
	for _, swapTokenPair := range swapTokenPairs {
		if baseTokenName == "" {
			tokens = append(tokens, swapTokenPair.BasePooledCoin.Denom)
			tokens = append(tokens, swapTokenPair.QuotePooledCoin.Denom)
		} else if baseTokenName == swapTokenPair.BasePooledCoin.Denom {
			tokens = append(tokens, swapTokenPair.QuotePooledCoin.Denom)
		} else if baseTokenName == swapTokenPair.QuotePooledCoin.Denom {
			tokens = append(tokens, swapTokenPair.BasePooledCoin.Denom)
		}
	}
	return tokens
}

func getSwapTokens(ctx sdk.Context, keeper Keeper, baseTokenName string) []string {
	var tokens []string
	// swap by route
	hasSwapRoute := false
	if baseTokenName != "" {
		nativePairName := types.GetSwapTokenPairName(baseTokenName, common.NativeToken)
		if _, err := keeper.GetSwapTokenPair(ctx, nativePairName); err == nil {
			hasSwapRoute = true
		}
	}
	// all swap token pairs
	swapTokenPairs := keeper.GetSwapTokenPairs(ctx)
	for _, swapTokenPair := range swapTokenPairs {
		if baseTokenName == "" {
			tokens = append(tokens, swapTokenPair.BasePooledCoin.Denom)
			tokens = append(tokens, swapTokenPair.QuotePooledCoin.Denom)
		} else if baseTokenName == swapTokenPair.BasePooledCoin.Denom {
			tokens = append(tokens, swapTokenPair.QuotePooledCoin.Denom)
		} else if baseTokenName == swapTokenPair.QuotePooledCoin.Denom {
			tokens = append(tokens, swapTokenPair.BasePooledCoin.Denom)
		}

		// swap by route
		if hasSwapRoute {
			if swapTokenPair.BasePooledCoin.Denom == common.NativeToken && swapTokenPair.QuotePooledCoin.Denom != baseTokenName {
				tokens = append(tokens, swapTokenPair.QuotePooledCoin.Denom)
			} else if swapTokenPair.QuotePooledCoin.Denom == common.NativeToken && swapTokenPair.BasePooledCoin.Denom != baseTokenName {
				tokens = append(tokens, swapTokenPair.BasePooledCoin.Denom)
			}
		}
	}
	return tokens
}

// querySwapQuoteInfo returns infos when swap token
func querySwapQuoteInfo(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapBuyInfoParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}
	if queryParams.SellTokenAmount == "" || queryParams.BuyToken == "" {
		return nil, types.ErrUnknownRequest()
	}

	sellAmount, err := sdk.ParseDecCoin(queryParams.SellTokenAmount)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}

	if sellAmount.Denom == queryParams.BuyToken {
		return nil, types.ErrUnknownRequest()
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
			return nil, types.ErrUnknownRequest()
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
			return nil, types.ErrUnknownRequest()
		}
		if tokenPair1.BasePooledCoin.Amount.IsZero() || tokenPair1.QuotePooledCoin.IsZero() {
			return nil, types.ErrUnknownRequest()
		}
		tokenPairName2 := types.GetSwapTokenPairName(queryParams.BuyToken, common.NativeToken)
		tokenPair2, err := keeper.GetSwapTokenPair(ctx, tokenPairName2)
		if err != nil {
			return nil, types.ErrUnknownRequest()
		}
		if tokenPair2.BasePooledCoin.Amount.IsZero() || tokenPair2.QuotePooledCoin.IsZero() {
			return nil, types.ErrUnknownRequest()
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
		return nil, types.ErrInternal()
	}
	return bz, nil

}

// querySwapLiquidityHistories returns liquidity info of the address
func querySwapLiquidityHistories(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapLiquidityInfoParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}
	// check params
	if queryParams.Address == "" {
		return nil, types.ErrUnknownRequest()
	}

	// coins in account
	addr, err := sdk.AccAddressFromBech32(queryParams.Address)
	if err != nil {
		return nil, types.ErrInvalidAddress(queryParams.Address)
	}

	var liquidityInfoList []types.SwapLiquidityInfo
	// coins in account
	accountCoins := keeper.tokenKeeper.GetCoins(ctx, addr)
	for _, coin := range accountCoins {
		// check if the token is pool token
		if !strings.HasPrefix(coin.Denom, types.PoolTokenPrefix) {
			continue
		}
		// check token pair name
		tokenPairName := coin.Denom[len(types.PoolTokenPrefix):]
		if queryParams.TokenPairName != "" && queryParams.TokenPairName != tokenPairName {
			continue
		}
		// get swap token pair
		swapTokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
		if err != nil {
			continue
		}
		poolTokenAmount := keeper.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
		baseDec := common.MulAndQuo(swapTokenPair.BasePooledCoin.Amount, coin.Amount, poolTokenAmount)
		quoteDec := common.MulAndQuo(swapTokenPair.QuotePooledCoin.Amount, coin.Amount, poolTokenAmount)
		baseAmount := sdk.NewDecCoinFromDec(swapTokenPair.BasePooledCoin.Denom, baseDec)
		quoteAmount := sdk.NewDecCoinFromDec(swapTokenPair.QuotePooledCoin.Denom, quoteDec)

		liquidityInfo := types.SwapLiquidityInfo{
			BasePooledCoin:  baseAmount,
			QuotePooledCoin: quoteAmount,
			PoolTokenCoin:   coin,
			PoolTokenRatio:  coin.Amount.Quo(poolTokenAmount),
		}
		liquidityInfoList = append(liquidityInfoList, liquidityInfo)
	}

	response := common.GetBaseResponse(liquidityInfoList)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, types.ErrInternal()
	}
	return bz, nil

}

// querySwapAddLiquidityQuote returns swap information of adding liquidity
func querySwapAddLiquidityQuote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapAddInfoParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}
	// check params
	if queryParams.QuoteTokenAmount == "" {
		return nil, types.ErrUnknownRequest()
	}
	if queryParams.BaseToken == "" {
		return nil, sdk.ErrUnknownRequest("invalid params: base_token is required")
	}
	queryTokenAmount, err := sdk.ParseDecCoin(queryParams.QuoteTokenAmount)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}

	tokenPairName := types.GetSwapTokenPairName(queryParams.BaseToken, queryTokenAmount.Denom)
	swapTokenPair, err := keeper.GetSwapTokenPair(ctx, tokenPairName)
	if err != nil {
		return nil, types.ErrUnknownRequest()
	}
	if swapTokenPair.BasePooledCoin.Amount.IsZero() && swapTokenPair.QuotePooledCoin.Amount.IsZero() {
		return nil, types.ErrUnknownRequest()
	}

	totalSupply := keeper.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
	if totalSupply.IsZero() {
		return nil, types.ErrUnknownRequest()
	}

	var addAmount sdk.Dec
	var liquidity sdk.Dec
	if swapTokenPair.BasePooledCoin.Denom == queryTokenAmount.Denom {
		addAmount = common.MulAndQuo(queryTokenAmount.Amount, swapTokenPair.QuotePooledCoin.Amount, swapTokenPair.BasePooledCoin.Amount)
		liquidity = common.MulAndQuo(queryTokenAmount.Amount, totalSupply, swapTokenPair.BasePooledCoin.Amount)
	} else {
		addAmount = common.MulAndQuo(queryTokenAmount.Amount, swapTokenPair.BasePooledCoin.Amount, swapTokenPair.QuotePooledCoin.Amount)
		liquidity = common.MulAndQuo(queryTokenAmount.Amount, totalSupply, swapTokenPair.QuotePooledCoin.Amount)
	}
	addInfo := types.SwapAddInfo{
		BaseTokenAmount: addAmount,
		PoolShare:       liquidity.Quo(totalSupply.Add(liquidity)),
	}
	response := common.GetBaseResponse(addInfo)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, types.ErrInternal()
	}
	return bz, nil

}
