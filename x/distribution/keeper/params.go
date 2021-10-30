package keeper

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"

	"github.com/okex/exchain/x/distribution/types"
)

// GetParams returns the total set of distribution parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the distribution parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
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
