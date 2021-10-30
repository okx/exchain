package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"

	"github.com/okex/exchain/x/distribution/types"
)

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetDistributionAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
