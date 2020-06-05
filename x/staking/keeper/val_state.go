package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

// KickOutAndReturnValidatorSetUpdates shows the main logic when a validator is kicked out of validator-set in an epoch
func (k Keeper) KickOutAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate) {
	logger := k.Logger(ctx)
	// 1.get the last validator set
	lastBondedVals := k.getLastValidatorsByAddr(ctx)
	logMap(logger, lastBondedVals, "LastBondedValAddrs")

	// 2.get the abandoned validator addrs
	abandonedValAddrs := k.getAbandonedValidatorAddrs(ctx)
	logSlice(logger, abandonedValAddrs, "AbandonedValAddrs")

	store := ctx.KVStore(k.storeKey)
	totalPower := k.GetLastTotalPower(ctx)

	// 3.look for the ahead candidate and promote it
	iterator := sdk.KVStoreReversePrefixIterator(store, types.ValidatorsByPowerIndexKey)
	defer iterator.Close()
	for abandonedNum := len(abandonedValAddrs); iterator.Valid() && abandonedNum > 0; iterator.Next() {
		valAddr := iterator.Value()
		// get the key of map
		valKey := getLastValidatorsMapKey(valAddr)

		// look for the ahead candidate
		_, found := lastBondedVals[valKey]
		if found {
			// not the val to promote
			continue
		}
		// promote the candidate
		validator := k.mustGetValidator(ctx, valAddr)
		// if we get to a zero-power validator without votes, just pass
		if validator.PotentialConsensusPowerByVotes() == 0 {
			continue
		}

		switch {
		case validator.IsUnbonded():
			validator = k.unbondedToBonded(ctx, validator)
		case validator.IsUnbonding():
			validator = k.unbondingToBonded(ctx, validator)
		case validator.IsBonded():
			panic("Panic. Candidate validator is not allowed to be in bonded status")
		default:
			panic("unexpected validator status")
		}

		// calculate the new power of candidate validator
		newPower := validator.ConsensusPowerByVotes()
		// update the validator to tendermint
		updates = append(updates, validator.ABCIValidatorUpdateByVotes())
		// set validator power on lookup index
		k.SetLastValidatorPower(ctx, valAddr, newPower)
		// cumsum the total power
		totalPower = totalPower.Add(sdk.NewInt(newPower))

		abandonedNum--
	}

	// 4.discharge the abandoned validators
	for _, valAddr := range abandonedValAddrs {
		validator := k.mustGetValidator(ctx, valAddr)
		switch {
		case validator.IsUnbonded():
			logger.Debug(fmt.Sprintf("validator %s is already in the unboned status", validator.OperatorAddress.String()))
		case validator.IsUnbonding():
			logger.Debug(fmt.Sprintf("validator %s is already in the unbonding status", validator.OperatorAddress.String()))
		case validator.IsBonded(): // bonded to unbonding
			k.bondedToUnbonding(ctx, validator)
			// delete from the bonded validator index
			k.DeleteLastValidatorPower(ctx, validator.GetOperator())
			// update the validator set
			updates = append(updates, validator.ABCIValidatorUpdateZero())
			// reduce the total power
			valKey := getLastValidatorsMapKey(valAddr)
			oldPowerBytes, found := lastBondedVals[valKey]
			if !found {
				panic("Never occur")
			}
			var oldPower int64
			k.cdc.MustUnmarshalBinaryLengthPrefixed(oldPowerBytes, &oldPower)
			totalPower = totalPower.Sub(sdk.NewInt(oldPower))
		default:
			panic("unexpected validator status")
		}
	}

	// 5. update the total power of this block to store
	k.SetLastTotalPower(ctx, totalPower)

	return updates
}

// getLastValidatorsMapKey gets the map key of last validator-set from val address
func getLastValidatorsMapKey(valAddr sdk.ValAddress) (key [sdk.AddrLen]byte) {
	copy(key[:], valAddr[:])
	return
}

func logMap(logger log.Logger, valMap validatorsByAddr, title string) {
	logger.Debug(title)
	for i := range valMap {
		logger.Debug(sdk.ValAddress(i[:]).String())
	}
}

func logSlice(logger log.Logger, valAddrs []sdk.ValAddress, title string) {
	logger.Debug(title)
	for _, addr := range valAddrs {
		logger.Debug(addr.String())
	}
}

// AppendAbandonedValidatorAddrs appends validator addresses to kick out
func (k Keeper) AppendAbandonedValidatorAddrs(ctx sdk.Context, ConsAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, ConsAddr)
	abandonedValAddr := k.getAbandonedValidatorAddrs(ctx)
	// if there are several validators to destroy in one block
	abandonedValAddr = append(abandonedValAddr, validator.OperatorAddress)
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(abandonedValAddr)
	ctx.KVStore(k.storeKey).Set(types.ValidatorAbandonedKey, bytes)
}

// getAbandonedValidatorAddrs gets the abandoned validator addresses
func (k Keeper) getAbandonedValidatorAddrs(ctx sdk.Context) (abandonedValAddr []sdk.ValAddress) {
	bytes := ctx.KVStore(k.storeKey).Get(types.ValidatorAbandonedKey)
	if len(bytes) == 0 {
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &abandonedValAddr)
	return
}

// DeleteAbandonedValidatorAddrs deletes the abandoned validator addresses
func (k Keeper) DeleteAbandonedValidatorAddrs(ctx sdk.Context) {
	ctx.KVStore(k.storeKey).Delete(types.ValidatorAbandonedKey)
}

// IsKickedOut tells whether there're jailed validators to kick out in an epoch
func (k Keeper) IsKickedOut(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.ValidatorAbandonedKey)
	return len(bytes) != 0
}
