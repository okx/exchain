package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/types"
)

// GetIbcEnabled retrieves the ibc enabled boolean from the param store
func (k Keeper) GetIbcEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramSpace.Get(ctx, types.KeyIbcEnabled, &res)
	return res
}

// GetParams returns the total set of ibc parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.GetIbcEnabled(ctx))
}

// SetParams sets the total set of ibc parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
