package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/staking/types"
)

// UpdateValidatorCommission attempts to update a validator's commission rate.
// An error is returned if the new commission rate is invalid.
func (k Keeper) UpdateValidatorCommission(ctx sdk.Context,
	validator types.Validator, newRate sdk.Dec) (types.Commission, error) {

	commission := validator.Commission
	blockTime := ctx.BlockHeader().Time
	if err := commission.ValidateNewRate(newRate, blockTime); err != nil {
		return commission, err
	}

	commission.Rate = newRate
	commission.UpdateTime = blockTime
	return commission, nil
}
