package dex

import (
	"fmt"

	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/dex/types"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
)

// EndBlocker called every block
func EndBlocker(ctx sdk.Context, k IKeeper) {
	seq := perf.GetPerf().OnEndBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnEndBlockExit(ctx, types.ModuleName, seq)
	// complete withdraw
	currentTime := ctx.BlockHeader().Time
	k.IterateWithdrawAddress(ctx, currentTime,
		func(_ int64, key []byte) (stop bool) {
			oldTime, addr := types.SplitWithdrawTimeKey(key)
			err := k.CompleteWithdraw(ctx, addr)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("complete undelegate failed: %s, addr:%s", err.Error(), addr.String()))
			} else {
				ctx.Logger().Debug(fmt.Sprintf("complete undelegate successful, addr: %s", addr.String()))
				k.DeleteWithdrawCompleteTimeAddress(ctx, oldTime, addr)
			}
			return false
		})
}
