package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/distribution/types"
)

func (k Keeper) GetDistributionType(ctx sdk.Context) (distrType uint32) {
	distrType = types.DistributionTypeOffChain
	if k.paramSpace.Has(ctx, types.ParamStoreKeyDistributionType) {
		k.paramSpace.Get(ctx, types.ParamStoreKeyDistributionType, &distrType)
	}

	return distrType
}

func (k Keeper) SetDistributionType(ctx sdk.Context, distrType uint32) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyDistributionType, &distrType)
}

func (k Keeper) GetWithdrawRewardEnabled(ctx sdk.Context) (enabled bool) {
	enabled = true
	if k.paramSpace.Has(ctx, types.ParamStoreKeyWithdrawRewardEnabled) {
		k.paramSpace.Get(ctx, types.ParamStoreKeyWithdrawRewardEnabled, &enabled)
	}

	return enabled
}

func (k Keeper) SetWithdrawRewardEnabled(ctx sdk.Context, enabled bool) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyWithdrawRewardEnabled, &enabled)
}
