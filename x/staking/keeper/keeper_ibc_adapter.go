package keeper

import (
	logrusplugin "github.com/itsfunny/go-cell/sdk/log/logrus"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	types2 "github.com/okex/exchain/x/staking/types"
)


func (k Keeper) TrackHistoricalInfo(ctx sdk.Context) {
	entryNum := k.HistoricalEntries(ctx)

	// Prune store to ensure we only have parameter-defined historical entries.
	// In most cases, this will involve removing a single historical entry.
	// In the rare scenario when the historical entries gets reduced to a lower value k'
	// from the original value k. k - k' entries must be deleted from the store.
	// Since the entries to be deleted are always in a continuous range, we can iterate
	// over the historical entries starting from the most recent version to be pruned
	// and then return at the first empty entry.
	for i := ctx.BlockHeight() - int64(entryNum); i >= 0; i-- {
		_, found := k.GetHistoricalInfo(ctx, i)
		if found {
			k.DeleteHistoricalInfo(ctx, i)
		} else {
			break
		}
	}

	// if there is no need to persist historicalInfo, return
	if entryNum == 0 {
		return
	}

	// Create HistoricalInfo struct
	lastVals := k.GetLastValidators(ctx)
	historicalEntry := types2.NewHistoricalInfo(ctx.BlockHeader(), lastVals)

	// Set latest HistoricalInfo at current height
	k.SetHistoricalInfo(ctx, ctx.BlockHeight(), historicalEntry)
}

// SetHistoricalInfo sets the historical info at a given height
func (k Keeper) SetHistoricalInfo(ctx sdk.Context, height int64, hi types2.HistoricalInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	value := types2.MustMarshalHistoricalInfo(k.cdc, hi)
	store.Set(key, value)
}

func (k Keeper) HistoricalEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.GetIfExists(ctx, types.KeyHistoricalEntries, &res)
	if res == 0 {
		res = 10000
		k.paramstore.Set(ctx, types.KeyHistoricalEntries, &res)
	}
	return
}

// DeleteHistoricalInfo deletes the historical info at a given height
func (k Keeper) DeleteHistoricalInfo(ctx sdk.Context, height int64) {
	logrusplugin.Info("清楚historical", "height", height)
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	store.Delete(key)
}

// get the group of the bonded validators
func (k Keeper) GetLastValidators(ctx sdk.Context) (validators []types2.Validator) {
	store := ctx.KVStore(k.storeKey)

	// add the actual validator power sorted store
	maxValidators := k.MaxValidators(ctx)
	validators = make([]types2.Validator, maxValidators)

	iterator := sdk.KVStorePrefixIterator(store, types.LastValidatorPowerKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {

		// sanity check
		if i >= int(maxValidators) {
			panic("more validators than maxValidators found")
		}
		address := types.AddressFromLastValidatorPowerKey(iterator.Key())
		validator := k.mustGetValidator(ctx, address)

		validators[i] = validator
		i++
	}
	return validators[:i] // trim
}
