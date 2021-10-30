package order

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"

	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/order/keeper"
	"github.com/okex/exchain/x/order/types"
	//"github.com/okex/exchain/x/common/version"
)

// BeginBlocker runs the logic of BeginBlocker with version 0.
// BeginBlocker resets keeper cache.
func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	seq := perf.GetPerf().OnBeginBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnBeginBlockExit(ctx, types.ModuleName, seq)

	keeper.ResetCache(ctx)
}
