package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if ctx.IsDeliver() {
		if types.GetEvmParamsCache().IsNeedParamsUpdate() {
			k.paramSpace.GetParamSet(ctx, &params)
			types.GetEvmParamsCache().UpdateParams(params)
		} else {
			params = types.GetEvmParamsCache().GetParams()
		}
	} else {
		k.paramSpace.GetParamSet(ctx, &params)
	}

	return
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.EvmStateDb.WithContext(ctx).SetParams(params)
	k.paramSpace.SetParamSet(ctx, &params)
	types.GetEvmParamsCache().SetNeedParamsUpdate()
}
