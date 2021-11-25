package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	if data, gas := k.configCache.params, k.configCache.gasParam; gas != 0 {
		ctx.GasMeter().ConsumeGas(gas, "evm.keeper.GetParams")
		return data
	}
	startGas := ctx.GasMeter().GasConsumed()
	k.paramSpace.GetParamSet(ctx, &params)
	k.configCache.setParams(params, ctx.GasMeter().GasConsumed()-startGas)
	return
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
