package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
)

// GetVote gets the vote entity
func (k Keeper) GetVote(ctx sdk.Context, voterAddr sdk.AccAddress, valAddr sdk.ValAddress) (types.Votes, bool) {
	store := ctx.KVStore(k.storeKey)
	votesBytes := store.Get(types.GetVoteKey(valAddr, voterAddr))
	var votes types.Votes
	// the voter never votes to this val before
	if votesBytes == nil {
		return votes, false
	}

	votes = types.MustUnmarshalVote(k.cdc, votesBytes)
	return votes, true
}

// SetVote sets votes to store
func (k Keeper) SetVote(ctx sdk.Context, voterAddr sdk.AccAddress, valAddr sdk.ValAddress, votes types.Votes) {
	key := types.GetVoteKey(valAddr, voterAddr)
	voteBytes := k.cdc.MustMarshalBinaryLengthPrefixed(votes)
	ctx.KVStore(k.storeKey).Set(key, voteBytes)
}

// DeleteVote deletes votes entire from store
func (k Keeper) DeleteVote(ctx sdk.Context, valAddr sdk.ValAddress, voterAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetVoteKey(valAddr, voterAddr))
}

// GetValidatorVotes returns all votes made to a specific validator and it's useful for querier
func (k Keeper) GetValidatorVotes(ctx sdk.Context, valAddr sdk.ValAddress) types.SharesResponses {
	store := ctx.KVStore(k.storeKey)

	var sharesResps types.SharesResponses
	iterator := sdk.KVStorePrefixIterator(store, types.GetVotesToValidatorsKey(valAddr))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// 1.get the delegator address
		delAddr := sdk.AccAddress(iterator.Key()[1+sdk.AddrLen:])

		// 2.get the shares
		shares := types.MustUnmarshalVote(k.cdc, iterator.Value())

		// 3.assemble the result
		sharesResps = append(sharesResps, types.NewSharesResponse(delAddr, shares))
	}

	return sharesResps
}

// IterateVotes iterates through all of the votes from store
func (k Keeper) IterateVotes(ctx sdk.Context, fn func(index int64, voterAddr sdk.AccAddress, valAddr sdk.ValAddress,
	votes types.Votes) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VoteKey)
	defer iterator.Close()

	boundIndex := sdk.AddrLen + 1
	for i := int64(0); iterator.Valid(); iterator.Next() {
		// 1.get voter/validator address from the key
		key := iterator.Key()
		valAddr, voterAddr := sdk.ValAddress(key[1:boundIndex]), sdk.AccAddress(key[boundIndex:])

		// 2.get the votes
		var vote types.Votes
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		// 3.call back the function
		if stop := fn(i, voterAddr, valAddr, vote); stop {
			break
		}
		i++
	}
}

// GetDelegatorsByProxy returns all delegator addresses binding a proxy and it's useful for querier
func (k Keeper) GetDelegatorsByProxy(ctx sdk.Context, proxyAddr sdk.AccAddress) (delAddrs []sdk.AccAddress) {
	k.IterateProxy(ctx, proxyAddr, false, func(_ int64, delAddr, _ sdk.AccAddress) (stop bool) {
		delAddrs = append(delAddrs, delAddr)
		return false
	})

	return
}
