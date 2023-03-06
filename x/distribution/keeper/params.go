package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"

	"github.com/okx/okbchain/x/distribution/types"
)

// GetParams returns the total set of distribution parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// GetParamsForInitGenesis returns the total set of distribution parameters.
func (k Keeper) GetParamsForInitGenesis(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSetForInitGenesis(ctx, &params, types.IgnoreInitGenesisList)
	return params
}

// SetParams sets the distribution parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// SetParamsForInitGenesis sets the distribution parameters to the param space, and ignore the target keys for additional
func (k Keeper) SetParamsForInitGenesis(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSetForInitGenesis(ctx, &params, types.IgnoreInitGenesisList)
}

// GetCommunityTax returns the current CommunityTax rate from the global param store
// nolint: errcheck
func (k Keeper) GetCommunityTax(ctx sdk.Context) (percent sdk.Dec) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyCommunityTax, &percent)
	return percent
}

// SetCommunityTax sets the value of community tax
// nolint: errcheck
func (k Keeper) SetCommunityTax(ctx sdk.Context, percent sdk.Dec) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyCommunityTax, &percent)
}

// GetWithdrawAddrEnabled returns the current WithdrawAddrEnabled
// nolint: errcheck
func (k Keeper) GetWithdrawAddrEnabled(ctx sdk.Context) (enabled bool) {
	k.paramSpace.Get(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}

// SetWithdrawAddrEnabled sets the value of enabled
// nolint: errcheck
func (k Keeper) SetWithdrawAddrEnabled(ctx sdk.Context, enabled bool) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
}
