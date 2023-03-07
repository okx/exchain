package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"

	"github.com/okx/okbchain/x/feesplit/types"
)

// GetParamsWithCache returns the total set of fees parameters from cache。
func (k Keeper) GetParamsWithCache(ctx sdk.Context) (params types.Params) {
	if ctx.UseParamCache() {
		if types.GetParamsCache().IsNeedParamsUpdate() {
			params = k.GetParams(ctx)
			types.GetParamsCache().UpdateParams(params, ctx.IsCheckTx())
		} else {
			params = types.GetParamsCache().GetParams()
		}
	} else {
		params = k.GetParams(ctx)
	}

	return params
}

// GetParams returns the total set of fees parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the fees parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
	types.GetParamsCache().SetNeedParamsUpdate()
}
