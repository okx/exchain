package crisis

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

// check all registered invariants
func EndBlocker(ctx sdk.Context, k Keeper) {
	return // TODO: add back when fix mtp store iterator bug
	if k.InvCheckPeriod() == 0 || ctx.BlockHeight()%int64(k.InvCheckPeriod()) != 0 {
		// skip running the invariant check
		return
	}
	k.AssertInvariants(ctx)
}
