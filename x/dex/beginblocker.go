package dex

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/perf"
)

// BeginBlocker called every block, reset cache.
func BeginBlocker(ctx sdk.Context, keeper IKeeper) {
	seq := perf.GetPerf().OnBeginBlockEnter(ctx, ModuleName)
	defer perf.GetPerf().OnBeginBlockExit(ctx, ModuleName, seq)
	keeper.ResetCache(ctx)
}
