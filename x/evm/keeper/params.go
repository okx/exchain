package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if data, gas := k.ConfigCache.GetParams(); gas != 0 {
		ctx.GasMeter().ConsumeGas(gas, "evm.keeper.GetParams")
		return data
	}
	startGas := ctx.GasMeter().GasConsumed()
	k.paramSpace.GetParamSet(ctx, &params)
	k.ConfigCache.setParams(params, ctx.GasMeter().GasConsumed()-startGas)
	return
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
