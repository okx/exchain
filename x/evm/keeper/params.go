package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
	"sync"
)

var (
	once sync.Once
)

// GetParams returns the total set of evm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	once.Do(func() {
		k.paramSpace.GetParamSet(ctx, &params)
		ctx.Cache().UpdateEvmParams(sdk.EvmParamsCopy{IsUpdate: false,
			EnableCreate:                      params.EnableCreate,
			EnableCall:                        params.EnableCall,
			ExtraEIPs:                         params.ExtraEIPs,
			EnableContractDeploymentWhitelist: params.EnableContractDeploymentWhitelist,
			EnableContractBlockedList:         params.EnableContractBlockedList,
			MaxGasLimitPerTx:                  params.MaxGasLimitPerTx})
	})
	evmParams := ctx.Cache().GetEvmParam()
	if evmParams.IsUpdate || evmParams.MaxGasLimitPerTx == 0 {
		k.paramSpace.GetParamSet(ctx, &params)
		ctx.Cache().UpdateEvmParams(sdk.EvmParamsCopy{IsUpdate: false,
			EnableCreate:                      params.EnableCreate,
			EnableCall:                        params.EnableCall,
			ExtraEIPs:                         params.ExtraEIPs,
			EnableContractDeploymentWhitelist: params.EnableContractDeploymentWhitelist,
			EnableContractBlockedList:         params.EnableContractBlockedList,
			MaxGasLimitPerTx:                  params.MaxGasLimitPerTx})
	}

	params = types.NewParams(
		evmParams.EnableCreate,
		evmParams.EnableCall,
		evmParams.EnableContractDeploymentWhitelist,
		evmParams.EnableContractBlockedList,
		evmParams.MaxGasLimitPerTx,
		evmParams.ExtraEIPs...)
	return
}

// SetParams sets the evm parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
	if ctx.Cache() != nil {
		evmParams := ctx.Cache().GetEvmParam()
		evmParams.IsUpdate = true
		ctx.Cache().UpdateEvmParams(evmParams)
	}
}
