package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/backend/types"
)

func GetNewOrders(ctx sdk.Context, orderKeeper types.OrderKeeper) []*backend.Order {
	orders, err := backend.GetNewOrdersAtEndBlock(ctx, orderKeeper)
	if err != nil {
		return nil
	}
	return orders
}

func GetDealsAndMatchResult(ctx sdk.Context, orderKeeper types.OrderKeeper) ([]*backend.Deal, []*backend.MatchResult, error) {
	deals, matchResults, err := backend.GetNewDealsAndMatchResultsAtEndBlock(ctx, orderKeeper)
	return deals, matchResults, err
}

func GetMatchResults(ctx sdk.Context, orderKeeper types.OrderKeeper) []*backend.MatchResult {
	_, matchResults, err := GetDealsAndMatchResult(ctx, orderKeeper)
	if err != nil {
		return nil
	}
	return matchResults
}
