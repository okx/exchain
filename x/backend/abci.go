package backend

import (
	"fmt"
	"strconv"

	"github.com/okex/okchain/x/backend/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	orderTypes "github.com/okex/okchain/x/order/types"
)

// Called every block, check expired orders
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	if keeper.Config.EnableBackend && keeper.Config.EnableMktCompute {
		keeper.Logger.Debug(fmt.Sprintf("begin backend endblocker: block---%d", ctx.BlockHeight()))
		storeNewOrders(ctx, keeper)
		updateOrders(ctx, keeper)
		storeDealAndMatchResult(ctx, keeper)
		storeFeeDetails(keeper)
		storeTransactions(keeper)
		keeper.Flush()
		keeper.Logger.Debug(fmt.Sprintf("end backend endblocker: block---%d", ctx.BlockHeight()))
	}
}

func storeTransactions(keeper Keeper) {
	defer types.PrintStackIfPanic()

	txs := keeper.Cache.GetTransactions()
	txsLen := len(txs)

	cnt, err := keeper.Orm.AddTransactions(txs)
	if err != nil {
		keeper.Logger.Error(fmt.Sprintf("[backend] Expect to insert %d txs, inserted Count %d, err: %+v", txsLen, cnt, err))
	} else {
		keeper.Logger.Debug(fmt.Sprintf("[backend] Expect to insert %d txs, inserted Count %d", txsLen, cnt))
	}
}

func storeDealAndMatchResult(ctx sdk.Context, keeper Keeper) {
	timestamp := ctx.BlockHeader().Time.Unix()
	keeper.Orm.MaxBlockTimestamp = timestamp
	deals, results, err := GetNewDealsAndMatchResultsAtEndBlock(ctx, keeper.OrderKeeper)

	if err != nil {
		keeper.Logger.Error(fmt.Sprintf("[backend] failed to GetNewDealsAndMatchResultsAtEndBlock, error: %s", err.Error()))
	}

	if len(results) > 0 {
		cnt, err := keeper.Orm.AddMatchResults(results)
		if err != nil {
			keeper.Logger.Error(fmt.Sprintf("[backend] Expect to insert %d matchResults, inserted Count %d, err: %+v", len(results), cnt, err))
		} else {
			keeper.Logger.Debug(fmt.Sprintf("[backend] Expect to insert %d matchResults, inserted Count %d", len(results), cnt))
		}
	}
	if len(deals) > 0 {
		cnt, err := keeper.Orm.AddDeals(deals)
		if err != nil {
			keeper.Logger.Error(fmt.Sprintf("[backend] Expect to insert %d deals, inserted Count %d, err: %+v", len(deals), cnt, err))
		} else {
			keeper.Logger.Debug(fmt.Sprintf("[backend] Expect to insert %d deals, inserted Count %d", len(deals), cnt))
		}
	}

	ts := keeper.Orm.MaxBlockTimestamp
	keeper.UpdateTickersBuffer(ts-types.SecondsInADay, ts+1, keeper.Cache.ProductsBuf)
}

func storeFeeDetails(keeper Keeper) {
	feeDetails := keeper.TokenKeeper.GetFeeDetailList()
	if len(feeDetails) > 0 {
		cnt, err := keeper.Orm.AddFeeDetails(feeDetails)
		if err != nil {
			keeper.Logger.Error(fmt.Sprintf("[backend] Expect to insert %d feeDetails, inserted Count %d, err: %+v", len(feeDetails), cnt, err))
		} else {
			keeper.Logger.Debug(fmt.Sprintf("[backend] Expect to insert %d feeDetails, inserted Count %d", len(feeDetails), cnt))
		}
	}
}

func storeNewOrders(ctx sdk.Context, keeper Keeper) {
	orders, err := GetNewOrdersAtEndBlock(ctx, keeper.OrderKeeper)
	if err != nil {
		keeper.Logger.Error(fmt.Sprintf("[backend] failed to GetNewOrdersAtEndBlock, error: %s", err.Error()))
	}

	if len(orders) > 0 {
		cnt, err := keeper.Orm.AddOrders(orders)
		if err != nil {
			keeper.Logger.Error(fmt.Sprintf("[backend] Expect to insert %d orders, inserted Count %d, err: %+v", len(orders), cnt, err))
		} else {
			keeper.Logger.Debug(fmt.Sprintf("[backend] Expect to insert %d orders, inserted Count %d", len(orders), cnt))
		}
	}
}

func updateOrders(ctx sdk.Context, keeper Keeper) {
	orders := GetUpdatedOrdersAtEndBlock(ctx, keeper.OrderKeeper)
	if len(orders) > 0 {
		cnt, err := keeper.Orm.UpdateOrders(orders)
		if err != nil {
			keeper.Logger.Error(fmt.Sprintf("[backend] Expect to update %d orders, updated Count %d, err: %+v", len(orders), cnt, err))
		} else {
			keeper.Logger.Debug(fmt.Sprintf("[backend] Expect to update %d orders, updated Count %d", len(orders), cnt))
		}
	}
}

func GetNewDealsAndMatchResultsAtEndBlock(ctx sdk.Context, orderKeeper types.OrderKeeper) ([]*types.Deal, []*types.MatchResult, error) {
	result := orderKeeper.GetBlockMatchResult()
	if result == nil {
		return []*types.Deal{}, []*types.MatchResult{}, nil
	}

	blockHeight := ctx.BlockHeight()
	totalDeals := 0
	for _, matchResult := range result.ResultMap {
		totalDeals += len(matchResult.Deals)
	}
	deals := make([]*types.Deal, 0, totalDeals)
	results := make([]*types.MatchResult, 0, len(result.ResultMap))
	for product, matchResult := range result.ResultMap {
		price, err := strconv.ParseFloat(matchResult.Price.String(), 64)
		if err == nil && matchResult.BlockHeight == blockHeight {
			if total, err := strconv.ParseFloat(matchResult.Quantity.String(), 64); err == nil {
				results = append(results, &types.MatchResult{
					BlockHeight: blockHeight,
					Product:     product,
					Price:       price,
					Quantity:    total,
					Timestamp:   ctx.BlockHeader().Time.Unix(),
				})
			}
		} else {
			return deals, results, err
		}

		for _, record := range matchResult.Deals {
			order := orderKeeper.GetOrder(ctx, record.OrderID)
			if quantity, err := strconv.ParseFloat(record.Quantity.String(), 64); err == nil {

				deal := &types.Deal{
					BlockHeight: blockHeight,
					OrderId:     record.OrderID,
					Side:        record.Side,
					Sender:      order.Sender.String(),
					Product:     product,
					Price:       price,
					Quantity:    quantity,
					Fee:         record.Fee,
					Timestamp:   ctx.BlockHeader().Time.Unix(),
				}
				deals = append(deals, deal)

			}
		}
	}
	return deals, results, nil
}

func GetNewOrdersAtEndBlock(ctx sdk.Context, orderKeeper types.OrderKeeper) ([]*types.Order, error) {
	blockHeight := ctx.BlockHeight()
	orderNum := orderKeeper.GetBlockOrderNum(ctx, blockHeight)
	orders := make([]*types.Order, 0, orderNum)
	var index int64 = 0
	for ; index < orderNum; index++ {
		orderId := orderTypes.FormatOrderID(blockHeight, index+1)
		order := orderKeeper.GetOrder(ctx, orderId)
		if order != nil {
			orderDb := &types.Order{
				TxHash:         order.TxHash,
				OrderId:        order.OrderID,
				Sender:         order.Sender.String(),
				Product:        order.Product,
				Side:           order.Side,
				Price:          order.Price.String(),
				Quantity:       order.Quantity.String(),
				Status:         order.Status,
				FilledAvgPrice: order.FilledAvgPrice.String(),
				RemainQuantity: order.RemainQuantity.String(),
				Timestamp:      order.Timestamp,
			}
			orders = append(orders, orderDb)
		} else {
			return nil, fmt.Errorf("failed to get order with orderId: %+v at blockHeight: %d", orderId, blockHeight)
		}
	}
	return orders, nil
}

func GetUpdatedOrdersAtEndBlock(ctx sdk.Context, orderKeeper types.OrderKeeper) []*types.Order {
	orderIds := orderKeeper.GetUpdatedOrderIDs()
	orders := make([]*types.Order, 0, len(orderIds))
	for _, orderId := range orderIds {
		order := orderKeeper.GetOrder(ctx, orderId)
		if order != nil {
			orderDb := &types.Order{
				TxHash:         order.TxHash,
				OrderId:        order.OrderID,
				Sender:         order.Sender.String(),
				Product:        order.Product,
				Side:           order.Side,
				Price:          order.Price.String(),
				Quantity:       order.Quantity.String(),
				Status:         order.Status,
				FilledAvgPrice: order.FilledAvgPrice.String(),
				RemainQuantity: order.RemainQuantity.String(),
				Timestamp:      order.Timestamp,
			}
			orders = append(orders, orderDb)
		}
	}
	return orders
}
