package order

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/common/perf"
	"github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
	//"github.com/okex/okchain/x/common/version"
)

// BeginBlocker runs the logic of BeginBlocker with version 0.
// BeginBlocker resets keeper cache.
func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	seq := perf.GetPerf().OnBeginBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnBeginBlockExit(ctx, types.ModuleName, seq)

	keeper.ResetCache(ctx)
}
