package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
)

func (k Keeper) GetUnbondingDelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valAddr sdk.ValAddress) (ubd types.UnbondingDelegation, found bool) {

	store := ctx.KVStore(k.storeKey)
	key := types.GetUBDKey(delAddr, valAddr)
	value := store.Get(key)
	if value == nil {
		return ubd, false
	}

	ubd = types.MustUnmarshalUBD(k.cdcMarshl.GetCdc(), value)
	return ubd, true
}

func (k Keeper) GetDelegatorUnbondingDelegations(ctx sdk.Context,
	delAddr sdk.AccAddress, page *query.PageRequest) (types.UnbondingDelegations, *query.PageResponse, error) {

	var unbondingDelegations types.UnbondingDelegations

	store := ctx.KVStore(k.storeKey)
	unbStore := prefix.NewStore(store, types.GetUBDsKey(delAddr))
	pageRes, err := query.Paginate(unbStore, page, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUBD(k.cdcMarshl.GetCdc(), value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, nil, nil
	}

	return unbondingDelegations, pageRes, nil
}
