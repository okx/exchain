package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/exchain/x/distribution/types"
	"github.com/okex/exchain/x/staking/exported"
)

// initialize rewards for a new validator
func (k Keeper) initializeValidator(ctx sdk.Context, val exported.ValidatorI) {
	// set accumulated commissions
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), types.InitialValidatorAccumulatedCommission())
}
