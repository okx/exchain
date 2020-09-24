package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker allocates the native token to the pools in PoolsYieldNativeToken
// according to the value of locked token in pool
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, types.ModuleName)
	yieldedNativeTokenAmt := moduleAcc.GetCoins().AmountOf(sdk.DefaultBondDenom)
	if yieldedNativeTokenAmt.LTE(sdk.ZeroDec()) {
		return
	}

	// 1. gets all pools in PoolsYieldNativeToken
	lockedPoolValue, pools, totalPoolsValue := calculateAllocateInfo(ctx, k)

	// 2. allocate native token to pools according to the value
	for poolName, pool := range pools {
		yieldAmt := lockedPoolValue[poolName].MulTruncate(yieldedNativeTokenAmt).QuoTruncate(totalPoolsValue)
		yieldNativeToken := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, yieldAmt)}
		pool.YieldedCoins.Add(yieldNativeToken)
		err := k.SupplyKeeper().BurnCoins(ctx, types.ModuleName, yieldNativeToken)
		if err != nil {
			panic("should not happen")
		}
		k.SetFarmPool(ctx, pool)
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

}

// calculateAllocateInfo gets all pools in PoolsYieldNativeToken
func calculateAllocateInfo(ctx sdk.Context, k keeper.Keeper) (map[string]sdk.Dec, map[string]types.FarmPool, sdk.Dec) {
	lockedPoolValue := make(map[string]sdk.Dec)
	pools := make(map[string]types.FarmPool)
	totalPoolsValue := sdk.ZeroDec()

	store := ctx.KVStore(k.StoreKey())
	iterator := sdk.KVStorePrefixIterator(store, types.PoolsYieldNativeTokenPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		poolName := types.SplitPoolsYieldNativeTokenKey(iterator.Key())
		pool, found := k.GetFarmPool(ctx, poolName)
		if !found {
			panic("should not happen")
		}
		pools[poolName] = pool
		poolValue := k.GetLockedPoolValue(ctx, pool)
		lockedPoolValue[poolName] = poolValue
		totalPoolsValue = totalPoolsValue.Add(poolValue)
	}
	return lockedPoolValue, pools, totalPoolsValue
}
