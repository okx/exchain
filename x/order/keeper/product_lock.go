package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/order/types"
)

// nolint
func (k Keeper) IsProductLocked(product string) bool {
	return k.dexKeeper.IsTokenPairLocked(product)
}

// nolint
func (k Keeper) SetProductLock(ctx sdk.Context, product string, lock *types.ProductLock) {
	k.dexKeeper.LockTokenPair(ctx, product, lock)
}

// nolint
func (k Keeper) UnlockProduct(ctx sdk.Context, product string) {
	k.dexKeeper.UnlockTokenPair(ctx, product)
}

// nolint
func (k Keeper) AnyProductLocked() bool {
	return k.dexKeeper.IsAnyProductLocked()
}
