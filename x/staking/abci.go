package staking

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k Keeper) {
	k.TrackHistoricalInfo(ctx)
}
