package stream

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// BeginBlocker runs the logic of BeginBlocker with version 0.
// BeginBlocker resets keeper cache111.
func BeginBlocker(ctx sdk.Context, keeper Keeper) {
	keeper.stream.Cache.Reset()
}
