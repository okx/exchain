package keeper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/okex/okexchain/x/ammswap"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// querySwapWatchlist returns watchlist of swap
func querySwapWatchlist(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapWatchlistParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(sdk.AppendMsgToErr("incorrectly formatted request data", err.Error()))
	}

	offset, limit := common.GetPage(queryParams.Page, queryParams.PerPage)
	if offset < 0 || limit < 0 {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("invalid params: page=%d or per_page=%d", queryParams.Page, queryParams.PerPage))
	}

	//check sort column param
	switch queryParams.SortColumn {
	case "":
	case types.SwapWatchlistLiquidity:
	case types.SwapWatchlistVolume24h:
	case types.SwapWatchlistFeeApy:
	case types.SwapWatchlistLastPrice:
	case types.SwapWatchlistChange24h:
	default:
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("invalid sort_column: %s", queryParams.SortColumn))
	}

	// whitelist map
	whitelistMap := getSwapWhitelistMap(keeper)

	// all swap token pairs
	swapTokenPairs := keeper.swapKeeper.GetSwapTokenPairs(ctx)
	startTime := ctx.BlockTime().Add(-24 * time.Hour).Unix()
	// query last 24 hours swap infos in orm db
	swapInfos := keeper.Orm.GetSwapInfo(startTime)
	swapVolumePriceMap := make(map[string]types.SwapVolumePriceInfo, len(swapTokenPairs))
	for _, swapInfo := range swapInfos {
		var err error
		price, err := sdk.NewDecFromStr(swapInfo.Price)
		if err != nil {
			continue
		}

		// check if in whitelist
		if _, found := whitelistMap[swapInfo.TokenPairName]; !found {
			continue
		}

		// calculate volume in dollar
		sellAmount, err := sdk.ParseDecCoin(swapInfo.SellAmount)
		if err != nil {
			continue
		}
		buyAmount, err := sdk.ParseDecCoin(swapInfo.BuysAmount)
		if err != nil {
			continue
		}
		volume := calculateDollarAmount(ctx, keeper, sellAmount, buyAmount)

		volumePriceInfo, ok := swapVolumePriceMap[swapInfo.TokenPairName]
		// not exist
		if !ok {
			swapVolumePriceMap[swapInfo.TokenPairName] = types.SwapVolumePriceInfo{
				Volume:    volume,
				Price24h:  price,
				Timestamp: swapInfo.Timestamp,
			}
			continue
		}

		// update swapVolumePriceMap
		if volumePriceInfo.Timestamp > swapInfo.Timestamp {
			volumePriceInfo.Price24h = price
		}
		volumePriceInfo.Volume = volumePriceInfo.Volume.Add(volume)
		swapVolumePriceMap[swapInfo.TokenPairName] = volumePriceInfo
	}

	// total watchlist
	var totalWatchlist []types.SwapWatchlist
	swapParams := keeper.swapKeeper.GetParams(ctx)
	for _, swapTokenPair := range swapTokenPairs {
		tokenPairName := swapTokenPair.TokenPairName()
		// check if in whitelist
		if _, found := whitelistMap[tokenPairName]; !found {
			continue
		}
		// calculate liquidity in dollar
		liquidity := calculateDollarAmount(ctx, keeper, swapTokenPair.BasePooledCoin, swapTokenPair.QuotePooledCoin)

		// calculate last price
		lastPrice := sdk.ZeroDec()
		if swapTokenPair.QuotePooledCoin.Amount.IsPositive() {
			lastPrice = swapTokenPair.BasePooledCoin.Amount.Quo(swapTokenPair.QuotePooledCoin.Amount)
		}

		// 24h volume and price
		volume24h := sdk.ZeroDec()
		price24h := sdk.ZeroDec()
		if volumePriceInfo, ok := swapVolumePriceMap[tokenPairName]; ok {
			volume24h = volumePriceInfo.Volume
			price24h = volumePriceInfo.Price24h
		}

		// calculate fee apy
		feeApy := sdk.ZeroDec()
		if liquidity.IsPositive() && liquidity.IsPositive() {
			feeApy = volume24h.Mul(swapParams.FeeRate).Quo(liquidity).Mul(sdk.NewDec(365))
		}

		// calculate price change
		change24h := sdk.ZeroDec()
		if price24h.IsPositive() {
			change24h = lastPrice.Sub(price24h).Quo(price24h)
		}

		totalWatchlist = append(totalWatchlist, types.SwapWatchlist{
			SwapPair:  tokenPairName,
			Liquidity: liquidity,
			Volume24h: volume24h,
			FeeApy:    feeApy,
			LastPrice: lastPrice,
			Change24h: change24h,
		})
	}

	// sort watchlist
	if queryParams.SortColumn != "" && len(totalWatchlist) != 0 {
		watchlistSorter := &types.SwapWatchlistSorter{
			Watchlist:     totalWatchlist,
			SortField:     queryParams.SortColumn,
			SortDirectory: queryParams.SortDirection,
		}
		sort.Sort(watchlistSorter)
		totalWatchlist = watchlistSorter.Watchlist
	}

	total := len(totalWatchlist)
	switch {
	case total < offset:
		totalWatchlist = totalWatchlist[0:0]
	case total < offset+limit:
		totalWatchlist = totalWatchlist[offset:]
	default:
		totalWatchlist = totalWatchlist[offset : offset+limit]
	}
	var response *common.ListResponse
	if len(totalWatchlist) > 0 {
		response = common.GetListResponse(total, queryParams.Page, queryParams.PerPage, totalWatchlist)
	} else {
		response = common.GetEmptyListResponse(total, queryParams.Page, queryParams.PerPage)
	}

	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// calculate baseAmount and quoteAmount in dollar by usdk
func calculateDollarAmount(ctx sdk.Context, keeper Keeper, baseAmount sdk.SysCoin, quoteAmount sdk.SysCoin) sdk.Dec {
	dollarAmount := sdk.ZeroDec()
	baseTokenDollar := sdk.ZeroDec()
	quoteTokenDollar := sdk.ZeroDec()

	if baseAmount.Denom == types.DollarQuoteToken {
		baseTokenDollar = baseAmount.Amount
	} else {
		baseTokenPairName := ammswap.GetSwapTokenPairName(baseAmount.Denom, types.DollarQuoteToken)
		if baseTokenPair, err := keeper.swapKeeper.GetSwapTokenPair(ctx, baseTokenPairName); err == nil {
			if baseTokenPair.BasePooledCoin.Denom == types.DollarQuoteToken && baseTokenPair.QuotePooledCoin.Amount.IsPositive() {
				baseTokenDollar = common.MulAndQuo(baseTokenPair.BasePooledCoin.Amount, baseAmount.Amount, baseTokenPair.QuotePooledCoin.Amount)
			} else if baseTokenPair.BasePooledCoin.Amount.IsPositive() {
				baseTokenDollar = common.MulAndQuo(baseTokenPair.QuotePooledCoin.Amount, baseAmount.Amount, baseTokenPair.BasePooledCoin.Amount)
			}
		}
	}

	if quoteAmount.Denom == types.DollarQuoteToken {
		quoteTokenDollar = quoteAmount.Amount
	} else {
		quoteTokenPairName := ammswap.GetSwapTokenPairName(quoteAmount.Denom, types.DollarQuoteToken)
		if quoteTokenPair, err := keeper.swapKeeper.GetSwapTokenPair(ctx, quoteTokenPairName); err == nil {
			if quoteTokenPair.BasePooledCoin.Denom == types.DollarQuoteToken && quoteTokenPair.QuotePooledCoin.Amount.IsPositive() {
				quoteTokenDollar = common.MulAndQuo(quoteTokenPair.BasePooledCoin.Amount, quoteAmount.Amount, quoteTokenPair.QuotePooledCoin.Amount)
			} else if quoteTokenPair.BasePooledCoin.Amount.IsPositive() {
				quoteTokenDollar = common.MulAndQuo(quoteTokenPair.QuotePooledCoin.Amount, quoteAmount.Amount, quoteTokenPair.BasePooledCoin.Amount)
			}
		}
	}

	if baseTokenDollar.IsZero() && quoteTokenDollar.IsPositive() && baseAmount.Amount.IsPositive() {
		baseTokenDollar = quoteTokenDollar
	}

	if quoteTokenDollar.IsZero() && baseTokenDollar.IsPositive() && quoteAmount.Amount.IsPositive() {
		quoteTokenDollar = baseTokenDollar
	}

	dollarAmount = baseTokenDollar.Add(quoteTokenDollar)
	return dollarAmount
}

// querySwapTokens returns tokens which are supported to swap in ammswap module
func querySwapTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapTokensParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}

	if queryParams.BusinessType == "" {
		return nil, swaptypes.ErrIsZeroValue("input business type param")
	}

	// coins in account
	var accountCoins sdk.SysCoins
	if queryParams.Address != "" {
		addr, err := sdk.AccAddressFromBech32(queryParams.Address)
		if err != nil {
			return nil, common.ErrCreateAddrFromBech32Failed(queryParams.Address, err.Error())
		}
		accountCoins = keeper.TokenKeeper.GetCoins(ctx, addr)
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
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func getSwapCreateLiquidityTokens(ctx sdk.Context, keeper Keeper) []string {
	var tokens []string
	allTokens := keeper.TokenKeeper.GetTokensInfo(ctx)
	for _, token := range allTokens {
		if !strings.HasPrefix(token.Symbol, swaptypes.PoolTokenPrefix) {
			tokens = append(tokens, token.Symbol)
		}
	}
	return tokens
}

func getSwapAddLiquidityTokens(ctx sdk.Context, keeper Keeper, baseTokenName string) []string {
	var tokens []string

	// whitelist map
	whitelistMap := getSwapWhitelistMap(keeper)
	// all swap token pairs
	swapTokenPairs := keeper.swapKeeper.GetSwapTokenPairs(ctx)
	for _, swapTokenPair := range swapTokenPairs {
		// check if in whitelist
		if _, found := whitelistMap[swapTokenPair.TokenPairName()]; !found {
			continue
		}
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

	// whitelist map
	whitelistMap := getSwapWhitelistMap(keeper)

	// swap by route
	hasSwapRoute := false
	if baseTokenName != "" {
		nativePairName := swaptypes.GetSwapTokenPairName(baseTokenName, common.NativeToken)
		// check if in whitelist
		if _, found := whitelistMap[nativePairName]; found {
			if _, err := keeper.swapKeeper.GetSwapTokenPair(ctx, nativePairName); err == nil {
				hasSwapRoute = true
			}
		}
	}
	// all swap token pairs
	swapTokenPairs := keeper.swapKeeper.GetSwapTokenPairs(ctx)
	for _, swapTokenPair := range swapTokenPairs {
		// check if in whitelist
		if _, found := whitelistMap[swapTokenPair.TokenPairName()]; !found {
			continue
		}
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

// nolint
func querySwapTokenPairs(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte,
	err sdk.Error) {
	// whitelist map
	whitelistMap := getSwapWhitelistMap(keeper)

	swapTokenPairs := keeper.swapKeeper.GetSwapTokenPairs(ctx)
	var whitelistTokenPair []swaptypes.SwapTokenPair
	for _, swapTokenPair := range swapTokenPairs {
		// check if in whitelist
		if _, found := whitelistMap[swapTokenPair.TokenPairName()]; !found {
			continue
		}
		whitelistTokenPair = append(whitelistTokenPair, swapTokenPair)
	}
	response := common.GetBaseResponse(whitelistTokenPair)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// querySwapLiquidityHistories returns liquidity info of the address
func querySwapLiquidityHistories(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QuerySwapLiquidityInfoParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	// check params
	if queryParams.Address == "" {
		return nil, swaptypes.ErrQueryParamsAddressIsEmpty()
	}

	// coins in account
	addr, err := sdk.AccAddressFromBech32(queryParams.Address)
	if err != nil {
		return nil, common.ErrCreateAddrFromBech32Failed(queryParams.Address, err.Error())
	}

	// whitelist map
	whitelistMap := getSwapWhitelistMap(keeper)

	var liquidityInfoList []types.SwapLiquidityInfo
	// coins in account
	accountCoins := keeper.TokenKeeper.GetCoins(ctx, addr)
	for _, coin := range accountCoins {
		// check if the token is pool token
		if !strings.HasPrefix(coin.Denom, swaptypes.PoolTokenPrefix) {
			continue
		}
		// check token pair name
		tokenPairName := coin.Denom[len(swaptypes.PoolTokenPrefix):]
		if queryParams.TokenPairName != "" && queryParams.TokenPairName != tokenPairName {
			continue
		}
		// check if in whitelist
		if _, found := whitelistMap[tokenPairName]; !found {
			continue
		}
		// get swap token pair
		swapTokenPair, err := keeper.swapKeeper.GetSwapTokenPair(ctx, tokenPairName)
		if err != nil {
			continue
		}
		poolTokenAmount := keeper.swapKeeper.GetPoolTokenAmount(ctx, swapTokenPair.PoolTokenName)
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
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil

}

func getSwapWhitelistMap(keeper Keeper) map[string]struct{} {
	swapWhitelist := keeper.Orm.GetSwapWhitelist()
	whitelistMap := make(map[string]struct{}, len(swapWhitelist))
	for _, whitelist := range swapWhitelist {
		whitelistMap[whitelist.TokenPairName] = struct{}{}
	}
	return whitelistMap
}
