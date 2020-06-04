package keeper

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/exported"

	"github.com/okex/okchain/x/staking/types"
)

// RegisterInvariantsCustom registers all staking invariants for okchain
func RegisterInvariantsCustom(ir sdk.InvariantRegistry, k Keeper) {

	ir.RegisterRoute(types.ModuleName, "module-accounts",
		ModuleAccountInvariantsCustom(k))
	ir.RegisterRoute(types.ModuleName, "nonnegative-power",
		NonNegativePowerInvariantCustom(k))
	ir.RegisterRoute(types.ModuleName, "positive-delegator",
		PositiveDelegatorInvariant(k))
	ir.RegisterRoute(types.ModuleName, "delegator-votes",
		DelegatorVotesInvariant(k))
}

// DelegatorVotesInvariant checks whether all the votes which persist
// in the store add up to the correct total votes amount stored in each existed validator
//TODO:if the self-votes based on msd is related with time-calculating, this DelegatorVotesInvariant will not pass
// because we can't calculate the votes number base on msd of a validator afterwards
func DelegatorVotesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		validators := k.GetAllValidators(ctx)
		for _, validator := range validators {

			valTotalVotes := validator.GetDelegatorShares()

			totalVotes := k.getVotesFromDefaultMinSelfDelegation()

			votes := k.GetValidatorVotes(ctx, validator.GetOperator())
			for _, vote := range votes {
				totalVotes = totalVotes.Add(vote.Votes)
			}

			if !valTotalVotes.Equal(totalVotes) {
				broken = true
				msg += fmt.Sprintf("broken delegator votes invariance:\n"+
					"\tvalidator.DelegatorShares: %v\n"+
					"\tsum of Vote.Votes and min self delegation: %v\n", valTotalVotes, totalVotes)
			}
		}
		return sdk.FormatInvariant(types.ModuleName, "delegator votes", msg), broken
	}
}

// PositiveDelegatorInvariant checks that all tokens delegated by delegator are greater than 0
func PositiveDelegatorInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		k.IterateDelegator(ctx, func(index int64, delegator types.Delegator) bool {
			if delegator.Tokens.IsNegative() {
				count++
				msg += fmt.Sprintf("\tdelegation with negative tokens: %+v\n", delegator)
			}

			if delegator.Tokens.IsZero() {
				count++
				msg += fmt.Sprintf("\tdelegation with zero tokens: %+v\n", delegator)
			}

			return false
		})

		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "positive tokens of delegator", fmt.Sprintf(
			"%d invalid tokens of delegator found\n%s", count, msg)), broken
	}
}

// NonNegativePowerInvariantCustom checks that all stored validators have nonnegative power
func NonNegativePowerInvariantCustom(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		iterator := k.ValidatorsPowerStoreIterator(ctx)

		for ; iterator.Valid(); iterator.Next() {
			validator, found := k.GetValidator(ctx, iterator.Value())
			if !found {
				panic(fmt.Sprintf("validator record not found for address: %X\n", iterator.Value()))
			}

			powerKey := types.GetValidatorsByPowerIndexKey(validator)

			if !bytes.Equal(iterator.Key(), powerKey) {
				broken = true
				msg += fmt.Sprintf("power store invariance:\n\tvalidator.Power: %v"+
					"\n\tkey should be: %v\n\tkey in store: %v\n",
					validator.ConsensusPowerByVotes(), powerKey, iterator.Key())
			}

			if validator.DelegatorShares.IsNegative() {
				broken = true
				msg += fmt.Sprintf("\tnegative votes for validator: %v\n", validator)
			}
		}
		iterator.Close()
		return sdk.FormatInvariant(types.ModuleName, "nonnegative power",
			fmt.Sprintf("found invalid validator powers\n%s", msg)), broken
	}
}

// ModuleAccountInvariantsCustom check invariants for module account
func ModuleAccountInvariantsCustom(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		bonded := sdk.ZeroDec()
		notBonded := sdk.ZeroDec()
		bondedPool := k.GetBondedPool(ctx)
		notBondedPool := k.GetNotBondedPool(ctx)
		bondDenom := k.BondDenom(ctx)

		k.IterateValidators(ctx, func(index int64, validator exported.ValidatorI) bool {
			bonded = bonded.Add(validator.GetMinSelfDelegation())
			return false
		})

		k.IterateDelegator(ctx, func(index int64, delegator types.Delegator) bool {
			bonded = bonded.Add(delegator.Tokens)
			return false
		})

		k.IterateUndelegationInfo(ctx, func(_ int64, undelegationInfo types.UndelegationInfo) bool {
			notBonded = notBonded.Add(undelegationInfo.Quantity)
			return false
		})

		poolBonded := bondedPool.GetCoins().AmountOf(bondDenom)
		poolNotBonded := notBondedPool.GetCoins().AmountOf(bondDenom)
		broken := !poolBonded.Equal(bonded) || !poolNotBonded.Equal(notBonded)

		// Bonded tokens should be equal to the sum of delegators' tokens
		// Not-bonded tokens should be equal to the sum of undelegation infos' tokens
		return sdk.FormatInvariant(types.ModuleName, "bonded and not bonded module account coins", fmt.Sprintf(
			"\tPool's bonded tokens: %v\n"+
				"\tsum of bonded tokens: %v\n"+
				"not bonded token invariance:\n"+
				"\tPool's not bonded tokens: %v\n"+
				"\tsum of not bonded tokens: %v\n"+
				"module accounts total (bonded + not bonded):\n"+
				"\tModule Accounts' tokens: %v\n"+
				"\tsum tokens:              %v\n",
			poolBonded, bonded, poolNotBonded, notBonded, poolBonded.Add(poolNotBonded), bonded.Add(notBonded))), broken
	}
}
