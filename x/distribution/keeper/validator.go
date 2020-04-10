package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/distribution/types"
	"github.com/okex/okchain/x/staking/exported"
)

// initialize rewards for a new validator
func (k Keeper) initializeValidator(ctx sdk.Context, val exported.ValidatorI) {
	// set accumulated commission
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), types.InitialValidatorAccumulatedCommission())
}
