package periodicauction

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
)

func preMatchProcessing(book *types.DepthBook) (buyAmountSum, sellAmountSum []sdk.Dec) {
	bookLength := len(book.Items)
	if bookLength == 0 {
		return
	}

	buyAmountSum = make([]sdk.Dec, bookLength)
	sellAmountSum = make([]sdk.Dec, bookLength)

	buyAmountSum[0] = book.Items[0].BuyQuantity
	for i := 1; i < bookLength; i++ {
		buyAmountSum[i] = buyAmountSum[i-1].Add(book.Items[i].BuyQuantity)
	}

	sellAmountSum[bookLength-1] = book.Items[bookLength-1].SellQuantity
	for i := bookLength - 2; i >= 0; i-- {
		sellAmountSum[i] = sellAmountSum[i+1].Add(book.Items[i].SellQuantity)
	}

	return
}

func execRule0(buyAmountSum, sellAmountSum []sdk.Dec) (maxExecution sdk.Dec, execution []sdk.Dec) {
	maxExecution = sdk.ZeroDec()
	bookLength := len(buyAmountSum)
	execution = make([]sdk.Dec, bookLength)
	for i := 0; i < bookLength; i++ {
		execution[i] = sdk.MinDec(buyAmountSum[i], sellAmountSum[i])
		maxExecution = sdk.MaxDec(execution[i], maxExecution)
	}

	return
}

func execRule1(maxExecution sdk.Dec, execution []sdk.Dec) (indexesRule1 []int) {
	bookLength := len(execution)
	for i := 0; i < bookLength; i++ {
		if execution[i].Equal(maxExecution) {
			indexesRule1 = append(indexesRule1, i)
		}
	}

	return
}

func execRule2(buyAmountSum, sellAmountSum []sdk.Dec, indexesRule1 []int) (indexesRule2 []int, imbalance []sdk.Dec) {
	indexLen1 := len(indexesRule1)
	imbalance = make([]sdk.Dec, indexLen1)
	for i := 0; i < indexLen1; i++ {
		imbalance[i] = buyAmountSum[indexesRule1[i]].Sub(sellAmountSum[indexesRule1[i]])
	}
	minAbsImbalance := imbalance[0].Abs()
	for i := 1; i < indexLen1; i++ {
		minAbsImbalance = sdk.MinDec(minAbsImbalance, imbalance[i].Abs())
	}
	for i := 0; i < indexLen1; i++ {
		if imbalance[i].Abs().Equal(minAbsImbalance) {
			indexesRule2 = append(indexesRule2, indexesRule1[i])
		}
	}

	return
}

func execRule3(book *types.DepthBook, offset int, refPrice sdk.Dec, pricePrecision int64,
	indexesRule2 []int, imbalance []sdk.Dec) (bestPrice sdk.Dec) {
	indexLen2 := len(indexesRule2)
	if imbalance[indexesRule2[0]-offset].GT(sdk.ZeroDec()) {
		// rule3a: all imbalances are positive, buy side pressure
		newRefPrice := refPrice.Mul(sdk.MustNewDecFromStr("1.05"))
		newRefPrice = newRefPrice.RoundDecimal(pricePrecision)
		bestPrice = bestPriceFromRefPrice(book.Items[indexesRule2[0]].Price,
			book.Items[indexesRule2[indexLen2-1]].Price, newRefPrice)
	} else if imbalance[indexesRule2[indexLen2-1]-offset].LT(sdk.ZeroDec()) {
		// rule3b: all imbalances are negative, sell side pressure
		newRefPrice := refPrice.Mul(sdk.MustNewDecFromStr("0.95"))
		newRefPrice = newRefPrice.RoundDecimal(pricePrecision)
		bestPrice = bestPriceFromRefPrice(book.Items[indexesRule2[0]].Price,
			book.Items[indexesRule2[indexLen2-1]].Price, newRefPrice)
	} else {
		// rule3c: some imbalance > 0, and some imbalance < 0, no buyer pressure or seller pressure
		newRefPrice := refPrice.RoundDecimal(pricePrecision)
		bestPrice = bestPriceFromRefPrice(book.Items[indexesRule2[0]].Price,
			book.Items[indexesRule2[indexLen2-1]].Price, newRefPrice)
	}

	return
}

// Calculate periodic auction match price, return the best price and execution amount
// The best price is found according following rules:
// rule0: No match, bestPrice = 0, maxExecution=0
// rule1: Maximum execution volume.
//        If there are more than one price with the same max execution, following rule2
// rule2: Minimum imbalance. We should select the price with minimum absolute value of imbalance.
//        If more than one price satisfy rule2, following rule3
// rule3: Market Pressure. There are 3 cases:
// rule3a: All imbalances are positive. It indicates buy side pressure. Set reference price with
//         last execute price plus a upper limit percentage(e.g. 5%). Then choose the price
//         which is closest to reference price.
// rule3b: All imbalances are negative. It indicates sell side pressure. Set reference price with
//         last execute price minus a lower limit percentage(e.g. 5%). Then choose the price
//         which is closest to reference price.
// rule3c: Otherwise, it indicates no one side pressure. Set reference price with last execute
//         price. Then choose the price which is closest to reference price.
func periodicAuctionMatchPrice(book *types.DepthBook, pricePrecision int64,
	refPrice sdk.Dec) (bestPrice sdk.Dec, maxExecution sdk.Dec) {

	buyAmountSum, sellAmountSum := preMatchProcessing(book)
	if len(buyAmountSum) == 0 {
		return sdk.ZeroDec(), sdk.ZeroDec()
	}

	maxExecution, execution := execRule0(buyAmountSum, sellAmountSum)
	if maxExecution.IsZero() {
		return refPrice, maxExecution
	}

	indexesRule1 := execRule1(maxExecution, execution)
	if len(indexesRule1) == 1 {
		bestPrice = book.Items[indexesRule1[0]].Price
		return
	}

	indexesRule2, imbalance := execRule2(buyAmountSum, sellAmountSum, indexesRule1)
	if len(indexesRule2) == 1 {
		bestPrice = book.Items[indexesRule2[0]].Price
		return
	}

	bestPrice = execRule3(book, indexesRule1[0], refPrice, pricePrecision, indexesRule2, imbalance)

	return
}

// get best price from reference price
// if min < ref < max, choose ref; else choose the closest price to ref price
func bestPriceFromRefPrice(minPrice, maxPrice, refPrice sdk.Dec) sdk.Dec {
	if minPrice.LTE(refPrice) {
		return minPrice
	}
	if maxPrice.GTE(refPrice) {
		return maxPrice
	}
	return refPrice
}

func markCurBlockToFutureExpireBlockList(ctx sdk.Context, keeper keeper.Keeper) {
	curBlockHeight := ctx.BlockHeight()
	feeParams := keeper.GetParams(ctx)

	// Add current blockHeight to future Height
	// which will solve expire orders in current block.
	futureHeight := curBlockHeight + feeParams.OrderExpireBlocks

	// the feeParams.OrderExpireBlocks param can be change during the blockchain running,
	// so we use an array to record the expire blocks in the feature block height
	futureExpireHeightList := keeper.GetExpireBlockHeight(ctx, futureHeight)
	futureExpireHeightList = append(futureExpireHeightList, curBlockHeight)
	keeper.SetExpireBlockHeight(ctx, futureHeight, futureExpireHeightList)
}

func cleanLastBlockClosedOrders(ctx sdk.Context, keeper keeper.Keeper) {
	// drop expired data
	lastClosedOrderIDs := keeper.GetLastClosedOrderIDs(ctx)
	for _, orderID := range lastClosedOrderIDs {
		keeper.DropOrder(ctx, orderID)
	}

	keeper.GetDiskCache().DecreaseStoreOrderNum(int64(len(lastClosedOrderIDs)))
}

// Deal the block from create to current height which is Expired
func cacheExpiredBlockToCurrentHeight(ctx sdk.Context, keeper keeper.Keeper) {
	logger := ctx.Logger().With("module", "order")
	curBlockHeight := ctx.BlockHeight()

	lastExpiredBlockHeight := keeper.GetLastExpiredBlockHeight(ctx)
	if lastExpiredBlockHeight == 0 {
		lastExpiredBlockHeight = curBlockHeight - 1
	}

	// check orders in expired blocks, remove expired orders by order id
	for height := lastExpiredBlockHeight + 1; height <= curBlockHeight; height++ {
		var expiredHeight int64
		expiredBlocks := keeper.GetExpireBlockHeight(ctx, height)
		for _, expiredHeight = range expiredBlocks {
			keeper.DropExpiredOrdersByBlockHeight(ctx, expiredHeight)
			logger.Info(fmt.Sprintf("currentHeight(%d), expire orders at blockHeight(%d)",
				curBlockHeight, expiredHeight))
		}
	}

	if !keeper.AnyProductLocked(ctx) {
		height := lastExpiredBlockHeight
		if curBlockHeight > 1 {
			for ; height < curBlockHeight; height++ {
				var expiredHeight int64
				for _, expiredHeight = range keeper.GetExpireBlockHeight(ctx, height) {
					keeper.DropBlockOrderNum(ctx, expiredHeight)
					logger.Info(fmt.Sprintf("currentHeight(%d), drop Data at blockHeight(%d)",
						curBlockHeight, expiredHeight))
				}
				keeper.DropExpireBlockHeight(ctx, height)
			}
		}
		keeper.SetLastExpiredBlockHeight(ctx, height)
	}
}

func cleanupOrdersWhoseTokenPairHaveBeenDelisted(ctx sdk.Context, keeper keeper.Keeper) {
	products := keeper.GetProductsFromDepthBookMap()
	for _, product := range products {
		tokenPair := keeper.GetDexKeeper().GetTokenPair(ctx, product)
		if tokenPair == nil {
			cleanupOrdersByProduct(ctx, keeper, product)
		}
	}
}

func cleanupOrdersByProduct(ctx sdk.Context, keeper keeper.Keeper, product string) {
	depthBook := keeper.GetDepthBookCopy(product)
	for _, item := range depthBook.Items {
		buyKey := types.FormatOrderIDsKey(product, item.Price, types.BuyOrder)
		orderIDList := keeper.GetProductPriceOrderIDs(buyKey)
		sellKey := types.FormatOrderIDsKey(product, item.Price, types.SellOrder)
		orderIDList = append(orderIDList, keeper.GetProductPriceOrderIDs(sellKey)...)
		cleanOrdersByOrderIDList(ctx, keeper, orderIDList)
	}
}

func cleanOrdersByOrderIDList(ctx sdk.Context, keeper keeper.Keeper, orderIDList []string) {
	logger := ctx.Logger()
	for _, orderID := range orderIDList {
		order := keeper.GetOrder(ctx, orderID)
		keeper.CancelOrder(ctx, order, logger)
	}
}

func cleanupExpiredOrders(ctx sdk.Context, keeper keeper.Keeper) {

	// Look forward to see what height will this block expired
	markCurBlockToFutureExpireBlockList(ctx, keeper)

	// Clean the expired orders which is collected by the last block
	cleanLastBlockClosedOrders(ctx, keeper)

	// Look backward to see who is expired and cache the expired orders
	cacheExpiredBlockToCurrentHeight(ctx, keeper)
}

func matchOrders(ctx sdk.Context, keeper keeper.Keeper) {
	blockHeight := ctx.BlockHeight()
	orderNum := keeper.GetBlockOrderNum(ctx, blockHeight)
	// no new orders in this block & no product lock in previous blocks, skip match
	if orderNum == 0 && !keeper.AnyProductLocked(ctx) {
		return
	}

	// step0: get active products
	products := keeper.GetDiskCache().GetNewDepthbookKeys()
	products = keeper.FilterDelistedProducts(ctx, products)
	keeper.GetDexKeeper().SortProducts(ctx, products) // sort products

	// step1: calc best price and max execution for every active product, save latest price
	//updatedProductsBaseprice := make(map[string]types.MatchResult)
	updatedProductsBasePrice := calcMatchPriceAndExecution(ctx, keeper, products)

	// step1.1: recover locked depth book
	lockMap := keeper.GetDexKeeper().GetLockedProductsCopy(ctx)
	for product := range lockMap.Data {
		products = append(products, product)
	}
	keeper.GetDexKeeper().SortProducts(ctx, products) // sort products

	// step2: execute match results, fill orders in match results, transfer tokens and collect fees
	executeMatch(ctx, keeper, products, updatedProductsBasePrice, lockMap)

	// step3: save match results for querying
	if len(updatedProductsBasePrice) > 0 {
		blockMatchResult := &types.BlockMatchResult{
			BlockHeight: blockHeight,
			ResultMap:   updatedProductsBasePrice,
			TimeStamp:   ctx.BlockHeader().Time.Unix(),
		}
		keeper.SetBlockMatchResult(blockMatchResult)
	}
}

func calcMatchPriceAndExecution(ctx sdk.Context, k keeper.Keeper, products []string) map[string]types.MatchResult {
	resultMap := make(map[string]types.MatchResult)

	for _, product := range products {
		tokenPair := k.GetDexKeeper().GetTokenPair(ctx, product)
		if tokenPair == nil {
			continue
		}
		book := k.GetDepthBookCopy(product)
		bestPrice, maxExecution := periodicAuctionMatchPrice(book, tokenPair.MaxPriceDigit,
			k.GetLastPrice(ctx, product))
		if maxExecution.IsPositive() {
			k.SetLastPrice(ctx, product, bestPrice)
			resultMap[product] = types.MatchResult{BlockHeight: ctx.BlockHeight(), Price: bestPrice,
				Quantity: maxExecution, Deals: []types.Deal{}}
		}
	}

	return resultMap
}

func lockProduct(ctx sdk.Context, k keeper.Keeper, logger log.Logger, product string, matchResult types.MatchResult,
	buyExecutedCnt, sellExecutedCnt sdk.Dec) {
	blockHeight := ctx.BlockHeight()

	lock := &types.ProductLock{
		BlockHeight:  matchResult.BlockHeight,
		Price:        matchResult.Price,
		Quantity:     matchResult.Quantity,
		BuyExecuted:  buyExecutedCnt,
		SellExecuted: sellExecutedCnt,
	}
	k.SetProductLock(ctx, product, lock)
	logger.Info(fmt.Sprintf("BlockHeight<%d> lock product(%s)", blockHeight, product))
}

func executeMatchedUpdatedProduct(ctx sdk.Context, k keeper.Keeper,
	updatedProductsBasePrice map[string]types.MatchResult, feeParams *types.Params, blockRemainDeals int64,
	product string, logger log.Logger) int64 {

	matchResult := updatedProductsBasePrice[product]
	buyExecutedCnt := sdk.ZeroDec()
	sellExecutedCnt := sdk.ZeroDec()

	if blockRemainDeals <= 0 {
		lockProduct(ctx, k, logger, product, matchResult, buyExecutedCnt, sellExecutedCnt)

		return blockRemainDeals
	}

	deals, blockRemainDeals := fillDepthBook(ctx, k, product,
		matchResult.Price, matchResult.Quantity, &buyExecutedCnt, &sellExecutedCnt, blockRemainDeals, feeParams)
	matchResult.Deals = deals
	updatedProductsBasePrice[product] = matchResult

	logger.Info(fmt.Sprintf("matchResult(%d-%s): price: %v, quantity: %v, buyExecuted: %v"+
		", sellExecuted: %v, dealsNum: %d", matchResult.BlockHeight, product,
		matchResult.Price, matchResult.Quantity, buyExecutedCnt, sellExecutedCnt, len(deals)))

	// if filling not done, lock the product
	if buyExecutedCnt.LT(matchResult.Quantity) || sellExecutedCnt.LT(matchResult.Quantity) {
		lockProduct(ctx, k, logger, product, matchResult, buyExecutedCnt, sellExecutedCnt)
	}

	return blockRemainDeals
}

func executeLockedProduct(ctx sdk.Context, k keeper.Keeper,
	updatedProductsBasePrice map[string]types.MatchResult, lockMap *types.ProductLockMap,
	feeParams *types.Params, blockRemainDeals int64, product string,
	logger log.Logger) int64 {

	if blockRemainDeals <= 0 {
		return blockRemainDeals
	}

	// fill locked product
	blockHeight := ctx.BlockHeight()
	lock := lockMap.Data[product]

	buyExecuted := lock.BuyExecuted
	sellExecuted := lock.SellExecuted
	deals, blockRemainDeals := fillDepthBook(ctx, k, product,
		lock.Price, lock.Quantity, &buyExecuted, &sellExecuted, blockRemainDeals, feeParams)

	// if deals not empty, add match result
	if len(deals) > 0 {
		updatedProductsBasePrice[product] = types.MatchResult{
			BlockHeight: lock.BlockHeight,
			Price:       lock.Price,
			Quantity:    lock.Quantity,
			Deals:       deals,
		}
	}

	logger.Info(fmt.Sprintf("BlockHeight<%d> execute locked product(%s<%d>): price: %v, "+
		"quantity: %v, buyExecuted: %v, sellExecuted: %v, dealsNum: %d",
		blockHeight, product, lock.BlockHeight, lock.Price, lock.Quantity, buyExecuted,
		sellExecuted, len(deals)))

	lock.BuyExecuted = buyExecuted
	lock.SellExecuted = sellExecuted
	// if execution is done, unlock product
	if buyExecuted.GTE(lock.Quantity) && sellExecuted.GTE(lock.Quantity) {
		k.UnlockProduct(ctx, product)
		logger.Info(fmt.Sprintf("BlockHeight<%d> unlock product(%s<%d>)", blockHeight,
			product, lock.BlockHeight))
	} else {
		// update product lock
		k.SetProductLock(ctx, product, lock)
	}

	return blockRemainDeals
}

func executeMatch(ctx sdk.Context, k keeper.Keeper, products []string,
	updatedProductsBasePrice map[string]types.MatchResult, lockMap *types.ProductLockMap) {
	logger := ctx.Logger().With("module", "order")
	feeParams := k.GetParams(ctx)
	blockRemainDeals := feeParams.MaxDealsPerBlock

	for _, product := range products {
		if _, ok := updatedProductsBasePrice[product]; ok {
			blockRemainDeals = executeMatchedUpdatedProduct(ctx, k, updatedProductsBasePrice, feeParams,
				blockRemainDeals, product, logger)
		} else if _, ok := lockMap.Data[product]; ok {
			blockRemainDeals = executeLockedProduct(ctx, k, updatedProductsBasePrice, lockMap, feeParams,
				blockRemainDeals, product, logger)
		}
	}
}
