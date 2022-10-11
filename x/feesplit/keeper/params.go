package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/x/feesplit/types"
)

// GetParams returns the total set of fees parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if ctx.UseParamCache() {
		if types.GetParamsCache().IsNeedParamsUpdate() {
			k.paramSpace.GetParamSet(ctx, &params)
			ctx.Logger().Error("从db中取得params1")
			types.GetParamsCache().UpdateParams(params, ctx.IsCheckTx())
		} else {
			ctx.Logger().Error("从cache中取得params")
			params = types.GetParamsCache().GetParams()
		}
	} else {
		ctx.Logger().Error("从db中取得params2")
		k.paramSpace.GetParamSet(ctx, &params)
	}

	return params
}

// SetParams sets the fees parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
	types.GetParamsCache().SetNeedParamsUpdate()
}
