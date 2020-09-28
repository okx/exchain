package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// GetWhitelist gets the whole whitelist currently
func (k Keeper) GetWhitelist(ctx sdk.Context) (whitelist types.PoolNameList) {
	store := ctx.KVStore(k.StoreKey())
	iterator := sdk.KVStorePrefixIterator(store, types.PoolsYieldNativeTokenPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		whitelist = append(whitelist, types.SplitPoolsYieldNativeTokenKey(iterator.Key()))
	}

	return
}

// GetWhitelist sets the pool name as a member into whitelist
func (k Keeper) SetWhitelist(ctx sdk.Context, poolName string) {
	ctx.KVStore(k.storeKey).Set(types.GetWhitelistMemberKey(poolName), []byte(""))
}
