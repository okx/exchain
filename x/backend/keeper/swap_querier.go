package keeper

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/okex/okexchain/x/ammswap"

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
	totalWatchlist := make([]types.SwapWatchlist, len(swapTokenPairs))
	swapParams := keeper.swapKeeper.GetParams(ctx)
	for i, swapTokenPair := range swapTokenPairs {
		tokenPairName := swapTokenPair.TokenPairName()
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

		totalWatchlist[i] = types.SwapWatchlist{
			SwapPair:  tokenPairName,
			Liquidity: liquidity,
			Volume24h: volume24h,
			FeeApy:    feeApy,
			LastPrice: lastPrice,
			Change24h: change24h,
		}
	}

	// sort watchlist
	if queryParams.SortColumn != "" {
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
		return nil, sdk.ErrInternal(err.Error())
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
