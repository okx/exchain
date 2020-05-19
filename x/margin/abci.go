package margin

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/perf"
	"github.com/okex/okchain/x/margin/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// 	TODO: fill out if your application requires beginblock, if not you can delete this function
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k Keeper) {
	seq := perf.GetPerf().OnEndBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnEndBlockExit(ctx, types.ModuleName, seq)
	// complete withdraw
	currentTime := ctx.BlockHeader().Time
	k.IterateWithdrawAddress(ctx, currentTime,
		func(_ int64, key []byte) (stop bool) {
			oldTime, addr := types.SplitWithdrawTimeKey(key)
			err := k.CompleteWithdraw(ctx, addr)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("complete withdraw failed: %s, addr:%s", err.Error(), addr.String()))
			} else {
				ctx.Logger().Debug(fmt.Sprintf("complete withdraw successful, addr: %s", addr.String()))
				k.DeleteWithdrawCompleteTimeAddress(ctx, oldTime, addr)
			}
			return false
		})
}
