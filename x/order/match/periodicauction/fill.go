package periodicauction

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	orderkeeper "github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
	token "github.com/okex/okchain/x/token/types"
)

func fillBuyOrders(ctx sdk.Context, keeper orderkeeper.Keeper, product string,
	bestPrice, maxExecution sdk.Dec, buyExecuted *sdk.Dec,
	blockRemainDeals int64, feeParams *types.Params) ([]types.Deal, int64) {

	var buyDeals []types.Deal
	book := keeper.GetDepthBookCopy(product)

	// Fill buy orders, prices from high to low
	index := 0
	for index < len(book.Items) {
		if !(book.Items[index].Price.GTE(bestPrice) && buyExecuted.LT(maxExecution)) {
			break
		}
		// item.InitPrice >= bestPrice, fill buy orders
		fillAmount := sdk.MinDec(book.Items[index].BuyQuantity, maxExecution.Sub(*buyExecuted))
		if fillAmount.IsZero() { // no buyer at this price
			index++
			continue
		}

		// Fill buy orders at this price
		key := types.FormatOrderIDsKey(product, book.Items[index].Price, types.BuyOrder)
		filledBuyDeals, filledBuyAmount, filledDealsCnt := fillOrderByKey(ctx, keeper, key,
			fillAmount, bestPrice, feeParams, blockRemainDeals)
		blockRemainDeals -= filledDealsCnt

		buyDeals = append(buyDeals, filledBuyDeals...)
		*buyExecuted = buyExecuted.Add(filledBuyAmount)

		book.Sub(index, filledBuyAmount, types.BuyOrder)

		res := book.RemoveIfEmpty(index)
		if !res {
			index++
		}

		if blockRemainDeals <= 0 {
			break
		}
	}
	keeper.SetDepthBook(product, book)

	return buyDeals, blockRemainDeals
}

func fillSellOrders(ctx sdk.Context, keeper orderkeeper.Keeper, product string,
	bestPrice, maxExecution sdk.Dec, sellExecuted *sdk.Dec,
	blockRemainDeals int64, feeParams *types.Params) ([]types.Deal, int64) {

	var sellDeals []types.Deal
	book := keeper.GetDepthBookCopy(product)

	// Fill sell orders, prices from low to high
	index := len(book.Items) - 1
	for index >= 0 {

		if !(book.Items[index].Price.LTE(bestPrice) && sellExecuted.LT(maxExecution)) {
			break
		}

		// item.InitPrice <= bestPrice, fill sell orders
		fillAmount := sdk.MinDec(book.Items[index].SellQuantity, maxExecution.Sub(*sellExecuted))

		if fillAmount.IsZero() {
			index--
			continue
		}
		// Fill sell orders at this price
		key := types.FormatOrderIDsKey(product, book.Items[index].Price, types.SellOrder)

		filledSellDeals, filledSellAmount, filledDealsCnt := fillOrderByKey(ctx, keeper,
			key, fillAmount, bestPrice, feeParams, blockRemainDeals)

		blockRemainDeals -= filledDealsCnt
		sellDeals = append(sellDeals, filledSellDeals...)
		*sellExecuted = sellExecuted.Add(filledSellAmount)

		book.Sub(index, filledSellAmount, types.SellOrder)
		book.RemoveIfEmpty(index)

		if blockRemainDeals <= 0 {
			break
		}
		index--
	}
	keeper.SetDepthBook(product, book) // update depthbook on filled

	return sellDeals, blockRemainDeals
}

// fillDepthBook will fill orders in depth book with bestPrice.
// It will update book and orderIDsMap, also update orders, charge fees, and transfer tokens,
// then return all deals.
func fillDepthBook(ctx sdk.Context,
	keeper orderkeeper.Keeper,
	product string,
	bestPrice,
	maxExecution sdk.Dec,
	buyExecutedCnt,
	sellExecutedCnt *sdk.Dec,
	blockRemainDeals int64,
	feeParams *types.Params) ([]types.Deal, int64) {

	var deals []types.Deal
	if maxExecution.IsZero() {
		return deals, blockRemainDeals
	}

	buyDeals, blockRemainDeals := fillBuyOrders(ctx, keeper, product, bestPrice, maxExecution,
		buyExecutedCnt, blockRemainDeals, feeParams)
	deals = append(deals, buyDeals...)
	if blockRemainDeals <= 0 {
		return deals, blockRemainDeals
	}

	sellDeals, blockRemainDeals := fillSellOrders(ctx, keeper, product, bestPrice, maxExecution,
		sellExecutedCnt, blockRemainDeals, feeParams)
	deals = append(deals, sellDeals...)

	return deals, blockRemainDeals
}

// Fill orders in orderIDsMap at specific key
func fillOrderByKey(ctx sdk.Context, keeper orderkeeper.Keeper, key string,
	needFillAmount sdk.Dec, fillPrice sdk.Dec, feeParams *types.Params,
	remainDeals int64) ([]types.Deal, sdk.Dec, int64) {

	deals := []types.Deal{}
	filledAmount := sdk.ZeroDec()
	orderIDsMap := keeper.GetDiskCache().GetOrderIDsMapCopy()
	filledDealsCnt := int64(0)

	orderIDs, ok := orderIDsMap.Data[key]
	// if key not found in orderIDsMap, return
	if !ok {
		return deals, filledAmount, filledDealsCnt
	}

	index := 0
	for filledDealsCnt < remainDeals && filledAmount.LT(needFillAmount) {
		order := keeper.GetOrder(ctx, orderIDs[index])
		if filledAmount.Add(order.RemainQuantity).LTE(needFillAmount) {
			filledAmount = filledAmount.Add(order.RemainQuantity)
			deal := fillOrder(order, ctx, keeper, fillPrice, order.RemainQuantity, feeParams)
			deals = append(deals, *deal)

			filledDealsCnt++
			index++
		} else {
			deal := fillOrder(order, ctx, keeper, fillPrice, needFillAmount.Sub(filledAmount), feeParams)
			deals = append(deals, *deal)
			filledAmount = needFillAmount

			break
		}

	}

	unFilledOrderIDs := orderIDs[index:] // update orderIDs, remove filled orderIDs
	// Note: orderIDs cannot be nil, we will use empty slice to remove Data on keeper
	if len(unFilledOrderIDs) == 0 {
		unFilledOrderIDs = []string{}
	}
	keeper.SetOrderIDs(key, unFilledOrderIDs) // update orderIDsMap on filled

	return deals, filledAmount, filledDealsCnt
}

func balanceAccount(order *types.Order, ctx sdk.Context, keeper orderkeeper.Keeper,
	fillPrice, fillQuantity sdk.Dec) {

	symbols := strings.Split(order.Product, "_")
	// transfer tokens
	var outputCoins, inputCoins sdk.DecCoins
	if order.Side == types.BuyOrder {
		outputCoins = sdk.DecCoins{{Denom: symbols[1], Amount: fillPrice.Mul(fillQuantity)}}
		inputCoins = sdk.DecCoins{{Denom: symbols[0], Amount: fillQuantity}}
	} else {
		outputCoins = sdk.DecCoins{{Denom: symbols[0], Amount: fillQuantity}}
		inputCoins = sdk.DecCoins{{Denom: symbols[1], Amount: fillPrice.Mul(fillQuantity)}}
	}
	keeper.BalanceAccount(ctx, order.Sender, outputCoins, inputCoins)
}

func chargeFee(order *types.Order, ctx sdk.Context, keeper orderkeeper.Keeper, fillQuantity sdk.Dec,
	feeParams *types.Params) sdk.DecCoins {
	// charge fee
	fee := orderkeeper.GetZeroFee()
	if order.Status == types.OrderStatusFilled {
		lockedFee := orderkeeper.GetOrderNewFee(order)
		fee = orderkeeper.GetOrderCostFee(order, ctx)
		receiveFee := lockedFee.Sub(fee)

		keeper.UnlockCoins(ctx, order.Sender, lockedFee, token.LockCoinsTypeFee)
		keeper.AddFeeDetail(ctx, order.Sender, receiveFee, types.FeeTypeOrderReceive)
		order.RecordOrderReceiveFee(receiveFee)

		err := keeper.AddCollectedFees(ctx, fee, order.Sender, types.FeeTypeOrderNew, false)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("Send fee failed:%s\n", err.Error()))
		}
	}
	dealFee := orderkeeper.GetDealFee(order, fillQuantity, ctx, keeper, feeParams)
	err := keeper.SendFeesToProductOwner(ctx, dealFee, order.Sender, types.FeeTypeOrderDeal, order.Product)
	if err == nil {
		order.RecordOrderDealFee(fee)
	}

	return dealFee
}

// Fill an order. Update order, charge fee and transfer tokens. Return a deal.
// If an order is fully filled but still lock some coins, unlock it.
func fillOrder(order *types.Order, ctx sdk.Context, keeper orderkeeper.Keeper,
	fillPrice, fillQuantity sdk.Dec, feeParams *types.Params) *types.Deal {

	// update order
	order.Fill(fillPrice, fillQuantity)

	balanceAccount(order, ctx, keeper, fillPrice, fillQuantity)
	// if fully filled and still need unlock coins
	if order.Status == types.OrderStatusFilled && order.RemainLocked.IsPositive() {
		needUnlockCoins := order.NeedUnlockCoins()
		keeper.UnlockCoins(ctx, order.Sender, needUnlockCoins, token.LockCoinsTypeQuantity)
		order.Unlock()
	}

	dealFee := chargeFee(order, ctx, keeper, fillQuantity, feeParams)

	keeper.UpdateOrder(order, ctx) // update order info on filled

	return &types.Deal{OrderID: order.OrderID, Side: order.Side, Quantity: fillQuantity, Fee: dealFee.String()}
}
