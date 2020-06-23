package order

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
)

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Params     types.Params   `json:"params"`
	OpenOrders []*types.Order `json:"open_orders"`
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:     types.DefaultParams(),
		OpenOrders: nil,
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	return nil
}

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data GenesisState) {
	keeper.SetParams(ctx, &data.Params)

	// reset open order& depth book
	for _, order := range data.OpenOrders {
		if order == nil {
			panic("the nil pointer is not expected")
		}
		height := types.GetBlockHeightFromOrderID(order.OrderID)

		futureHeight := height + data.Params.OrderExpireBlocks
		futureExpireHeightList := keeper.GetExpireBlockHeight(ctx, futureHeight)
		futureExpireHeightList = append(futureExpireHeightList, height)
		keeper.SetExpireBlockHeight(ctx, futureHeight, futureExpireHeightList)

		orderNum := keeper.GetBlockOrderNum(ctx, height)
		keeper.SetBlockOrderNum(ctx, height, orderNum+1)
		keeper.SetOrder(ctx, order.OrderID, order)

		// update depth book and orderIDsMap in cache
		keeper.InsertOrderIntoDepthBook(order)
	}
	if len(data.OpenOrders) > 0 {
		keeper.Cache2Disk(ctx)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	tokenPairs := keeper.GetDexKeeper().GetTokenPairs(ctx)

	var openOrders []*types.Order
	var num int64 = 1
	for _, pair := range tokenPairs {
		if pair == nil {
			panic("the nil pointer is not expected")
		}
		product := fmt.Sprintf("%s_%s", pair.BaseAssetSymbol, pair.QuoteAssetSymbol)
		// update token pairs price
		pair.InitPrice = keeper.GetLastPrice(ctx, product)
		keeper.GetDexKeeper().UpdateTokenPair(ctx, product, pair)

		// get open orders
		depthBook := keeper.GetDepthBookFromDB(ctx, product)
		var openIDs []string
		for _, item := range depthBook.Items {
			if item.SellQuantity.IsPositive() {
				key := types.FormatOrderIDsKey(product, item.Price, types.SellOrder)
				ids := keeper.GetProductPriceOrderIDsFromDB(ctx, key)

				openIDs = append(openIDs, ids...)
			}
			if item.BuyQuantity.IsPositive() {
				key := types.FormatOrderIDsKey(product, item.Price, types.BuyOrder)
				ids := keeper.GetProductPriceOrderIDsFromDB(ctx, key)
				openIDs = append(openIDs, ids...)
			}
		}

		for _, orderID := range openIDs {
			order := keeper.GetOrder(ctx, orderID)
			if order.Status == types.OrderStatusFilled ||
				order.Status == types.OrderStatusCancelled ||
				order.Status == types.OrderStatusExpired {
				continue
			}
			// change orderID for order expire and order id on new chain
			// ID+genesisBlockHeight+1~n
			//order.OrderID = common.FormatOrderID(1, num)
			//order.Quantity = order.RemainQuantity
			//order.Status = types.OrderStatusOpen
			//order.TxHash=??
			openOrders = append(openOrders, order)
			num++
		}
	}

	return GenesisState{
		Params:     *params,
		OpenOrders: openOrders,
	}
}
