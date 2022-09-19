package infura

import (
	"log"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// BeginBlocker runs the logic of BeginBlocker with version 0.
// BeginBlocker resets keeper cache.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	log.Println("lcm infura BeginBlocker")
	if !k.stream.enable {
		return
	}
	k.stream.cache.Reset()
}
