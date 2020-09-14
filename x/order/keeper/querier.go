package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/order/types"
)

// nolint
const (
	DefaultBookSize = 200
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryOrderDetail:
			return queryOrder(ctx, path[1:], req, keeper)
		case types.QueryDepthBook:
			return queryDepthBook(ctx, path[1:], req, keeper)
		case types.QueryStore:
			return queryStore(ctx, path[1:], req, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)

		case types.QueryDepthBookV2:
			return queryDepthBookV2(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown order query endpoint")
		}
	}
}

// nolint: unparam
func queryOrder(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte,
	err sdk.Error) {
	order := keeper.GetOrder(ctx, path[0])
	if order == nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("order(%v) does not exist", path[0]))
	}
	bz := keeper.cdc.MustMarshalJSON(order)
	return bz, nil
}

// QueryDepthBookParams as input parameters when querying the depthBook
type QueryDepthBookParams struct {
	Product string
	Size    uint
}

// NewQueryDepthBookParams creates a new instance of QueryProposalParams
func NewQueryDepthBookParams(product string, size uint) QueryDepthBookParams {
	if size == 0 {
		size = DefaultBookSize
	}
	return QueryDepthBookParams{
		Product: product,
		Size:    size,
	}
}

// nolint
type BookResItem struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

// BookRes is used to return the result of queryDepthBook
type BookRes struct {
	Asks []BookResItem `json:"asks"`
	Bids []BookResItem `json:"bids"`
}

// nolint: unparam
func queryDepthBook(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte,
	sdk.Error) {
	var params QueryDepthBookParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(
			sdk.AppendMsgToErr("incorrectly formatted request Data", err.Error()))
	}
	if params.Size == 0 {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("invalid param: size= %d", params.Size))
	}
	tokenPair := keeper.GetDexKeeper().GetTokenPair(ctx, params.Product)
	if tokenPair == nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Non-exist product: %s", params.Product))
	}
	depthBook := keeper.GetDepthBookFromDB(ctx, params.Product)

	var asks []BookResItem
	var bids []BookResItem
	for _, item := range depthBook.Items {
		if item.SellQuantity.IsPositive() {
			asks = append([]BookResItem{{item.Price.String(), item.SellQuantity.String()}}, asks...)
		}
		if item.BuyQuantity.IsPositive() {
			bids = append(bids, BookResItem{item.Price.String(), item.BuyQuantity.String()})
		}
	}
	if uint(len(asks)) > params.Size {
		asks = asks[:params.Size]
	}
	if uint(len(bids)) > params.Size {
		bids = bids[:params.Size]
	}

	bookRes := BookRes{
		Asks: asks,
		Bids: bids,
	}
	bz := keeper.cdc.MustMarshalJSON(bookRes)
	return bz, nil
}

// StoreStatistic is used to store the state of depthBook
type StoreStatistic struct {
	StoreOrderNum   int64
	DepthBookNum    map[string]int64
	BookOrderIDsNum map[string]int64
}

func getStoreStatistic(ctx sdk.Context, keeper Keeper) *StoreStatistic {
	storeOrderNum := keeper.GetStoreOrderNum(ctx)
	ss := &StoreStatistic{
		StoreOrderNum: storeOrderNum,
	}

	depthBookMap := make(map[string]types.DepthBook)

	depthStore := ctx.KVStore(keeper.orderStoreKey)
	iter := sdk.KVStorePrefixIterator(depthStore, types.DepthBookKey)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var depthBook types.DepthBook
		bz := iter.Value()
		keeper.cdc.MustUnmarshalBinaryBare(bz, &depthBook)
		depthBookMap[types.GetKey(iter)] = depthBook
	}
	ss.DepthBookNum = make(map[string]int64, len(depthBookMap))
	ss.BookOrderIDsNum = make(map[string]int64, len(depthBookMap)*500)
	for product, depthBook := range depthBookMap {
		ss.DepthBookNum[product] = int64(len(depthBook.Items))
		for _, item := range depthBook.Items {
			if item.BuyQuantity.IsPositive() {
				key := types.FormatOrderIDsKey(product, item.Price, types.BuyOrder)
				orderIDs := keeper.GetProductPriceOrderIDs(key)
				ss.BookOrderIDsNum[key] = int64(len(orderIDs))
			}
			if item.SellQuantity.IsPositive() {
				key := types.FormatOrderIDsKey(product, item.Price, types.SellOrder)
				orderIDs := keeper.GetProductPriceOrderIDs(key)
				ss.BookOrderIDsNum[key] = int64(len(orderIDs))
			}
		}
	}
	return ss
}

func queryStore(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte,
	sdk.Error) {
	ss := getStoreStatistic(ctx, keeper)
	bz := keeper.cdc.MustMarshalJSON(ss)
	return bz, nil
}

func queryParameters(ctx sdk.Context, keeper Keeper) (res []byte, err sdk.Error) {
	params := keeper.GetParams(ctx)
	res, errRes := codec.MarshalJSONIndent(keeper.cdc, params)
	if errRes != nil {
		return nil, sdk.ErrInternal(
			sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDepthBookV2(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params QueryDepthBookParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}
	if params.Size == 0 {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("invalid param: size= %d", params.Size))
	}
	depthBook := keeper.GetDepthBookFromDB(ctx, params.Product)

	var asks []BookResItem
	var bids []BookResItem
	for _, item := range depthBook.Items {
		if item.SellQuantity.IsPositive() {
			asks = append([]BookResItem{{item.Price.String(), item.SellQuantity.String()}}, asks...)
		}
		if item.BuyQuantity.IsPositive() {
			bids = append(bids, BookResItem{item.Price.String(), item.BuyQuantity.String()})
		}
	}
	if uint(len(asks)) > params.Size {
		asks = asks[:params.Size]
	}
	if uint(len(bids)) > params.Size {
		bids = bids[:params.Size]
	}

	bookRes := BookRes{
		Asks: asks,
		Bids: bids,
	}

	res, err := common.JSONMarshalV2(bookRes)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return res, nil
}
