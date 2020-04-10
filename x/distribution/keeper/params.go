package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// ParamKeyTable is the type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyWithdrawAddrEnabled, true,
	)
}

// GetWithdrawAddrEnabled returns the current WithdrawAddrEnabled
// nolint: errcheck
func (k Keeper) GetWithdrawAddrEnabled(ctx sdk.Context) bool {
	var enabled bool
	k.paramSpace.Get(ctx, ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}

// SetWithdrawAddrEnabled sets the value of enabled
// nolint: errcheck
func (k Keeper) SetWithdrawAddrEnabled(ctx sdk.Context, enabled bool) {
	k.paramSpace.Set(ctx, ParamStoreKeyWithdrawAddrEnabled, &enabled)
}
