package keeper

import (
	"fmt"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SanityCheck checks validator shares
func (k Keeper) SanityCheck(ctx sdk.Context) error {
	k.Logger(ctx).Error("start sanity check in module staking")
	validators := k.GetAllValidators(ctx)
	for _, validator := range validators {

		valTotalShares := validator.GetDelegatorShares()

		var totalShares sdk.Dec
		if validator.MinSelfDelegation.Equal(sdk.ZeroDec()) && validator.Jailed {
			totalShares = sdk.ZeroDec()
		} else {
			//TODO:if the self-votes based on msd is related with time-calculating, this DelegatorVotesInvariant will not pass
			// because we can't calculate the votes number base on msd of a validator afterwards
			totalShares = sdk.OneDec()
		}

		votes := k.GetValidatorAllShares(ctx, validator.GetOperator())
		for _, vote := range votes {
			totalShares = totalShares.Add(vote.Shares)
		}

		if !valTotalShares.Equal(totalShares) {
			msg := fmt.Sprintf("validator address:%s, broken delegator votes invariance:\n"+
				"\tvalidator.DelegatorShares: %v\n"+
				"\tsum of Vote.Votes and min self delegation: %v\n", validator.OperatorAddress, valTotalShares, totalShares)
			return errors.New(msg)
		}
	}
	return nil
}
