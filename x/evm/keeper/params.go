package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if types.EvmParamsCache.IsNeedParamsUpdate() && ctx.IsDeliverorAsync() {
		k.paramSpace.GetParamSet(ctx, &params)
		types.EvmParamsCache.UpdateParams(params)
	} else {
		params = types.EvmParamsCache.GetParams()
	}
	return
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
	types.EvmParamsCache.SetNeedParamsUpdate()
}
