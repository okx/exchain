package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	outtypes "github.com/okex/exchain/x/staking/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	historicalEntry := outtypes.NewHistoricalInfo(ctx.BlockHeader(), lastVals)

	// Set latest HistoricalInfo at current height
	k.SetHistoricalInfo(ctx, ctx.BlockHeight(), historicalEntry)
}

// SetHistoricalInfo sets the historical info at a given height
func (k Keeper) SetHistoricalInfo(ctx sdk.Context, height int64, hi outtypes.HistoricalInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	value := outtypes.MustMarshalHistoricalInfo(k.cdcMarshl.GetCdc(), hi)
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
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	store.Delete(key)
}

// get the group of the bonded validators
func (k Keeper) GetLastValidators(ctx sdk.Context) (validators []outtypes.Validator) {
	store := ctx.KVStore(k.storeKey)

	// add the actual validator power sorted store
	maxValidators := k.MaxValidators(ctx)
	validators = make([]outtypes.Validator, maxValidators)

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

// GetHistoricalInfo gets the historical info at a given height
func (k Keeper) GetHistoricalInfo(ctx sdk.Context, height int64) (types.HistoricalInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetHistoricalInfoKey(height)

	value := store.Get(key)
	if value == nil {
		return types.HistoricalInfo{}, false
	}

	return types.MustUnmarshalHistoricalInfo(k.cdcMarshl.GetCdc(), value), true
}

func (k Keeper) GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []types.Delegation {
	delegations := make([]types.Delegation, 0)

	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delegator)
	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) //smallest to largest
	defer iterator.Close()

	i := 0
	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdcMarshl.GetCdc(), iterator.Value())
		delegations = append(delegations, delegation)
		i++
	}

	return delegations
}

func (k Keeper) DelegatorDelegations(ctx sdk.Context, req *outtypes.QueryDelegatorDelegationsRequest) (*outtypes.QueryDelegatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var delegations types.Delegations

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddr)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.GetDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdcMarshl.GetCdc(), value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delegationResps, err := DelegationsToDelegationResponses(ctx, k, delegations)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &outtypes.QueryDelegatorDelegationsResponse{DelegationResponses: delegationResps, Pagination: pageRes}, nil

}
