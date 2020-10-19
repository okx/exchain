package farm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker allocates the native token to the pools in PoolsYieldNativeToken
// according to the value of locked token in pool
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	logger := k.Logger(ctx)

	moduleAcc := k.SupplyKeeper().GetModuleAccount(ctx, MintFarmingAccount)
	yieldedNativeTokenAmt := moduleAcc.GetCoins().AmountOf(sdk.DefaultBondDenom)
	logger.Debug(fmt.Sprintf("amount of yielded native token: %s", yieldedNativeTokenAmt))
	if yieldedNativeTokenAmt.LTE(sdk.ZeroDec()) {
		return
	}

	// 1. gets all pools in PoolsYieldNativeToken
	lockedPoolValueMap, pools, totalPoolsValue := calculateAllocateInfo(ctx, k)

	if totalPoolsValue.LTE(sdk.ZeroDec()) {
		return
	}

	// 2. allocate native token to pools according to the value
	remainingNativeTokenAmt := yieldedNativeTokenAmt
	for i, pool := range pools {
		poolName := pool.Name
		var allocatedAmt sdk.Dec
		if i == len(pools)-1 {
			allocatedAmt = remainingNativeTokenAmt
		} else {
			allocatedAmt = lockedPoolValueMap[poolName].MulTruncate(yieldedNativeTokenAmt).QuoTruncate(totalPoolsValue)
		}
		logger.Debug(fmt.Sprintf("Pool %s allocate %s yielded native token", poolName, allocatedAmt.String()))
		current := k.GetPoolCurrentRewards(ctx, pool.Name)
		allocatedCoins := sdk.NewDecCoinsFromDec(sdk.DefaultBondDenom, allocatedAmt)
		current.Rewards = current.Rewards.Add(allocatedCoins)
		pool.RemainingRewards = pool.RemainingRewards.Add(allocatedCoins)
		k.SetPoolCurrentRewards(ctx, pool.Name, current)
		remainingNativeTokenAmt = remainingNativeTokenAmt.Sub(allocatedAmt)
		logger.Debug(fmt.Sprintf("Pool %s rewards are %s", poolName, current.Rewards))
		k.SetFarmPool(ctx, pool)
	}
	logger.Debug(fmt.Sprintf("RemainingNativeTokenAmt: %s", remainingNativeTokenAmt))

	// 3.liquidate native token minted at current block for yield farming
	err := k.SupplyKeeper().SendCoinsFromModuleToModule(
		ctx, MintFarmingAccount, YieldFarmingAccount,
		sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, yieldedNativeTokenAmt)},
	)
	if err != nil {
		panic("should not happen")
	}
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

}

// calculateAllocateInfo gets all pools in PoolsYieldNativeToken
func calculateAllocateInfo(ctx sdk.Context, k keeper.Keeper) (map[string]sdk.Dec, []types.FarmPool, sdk.Dec) {
	lockedPoolValue := make(map[string]sdk.Dec)
	var pools types.FarmPools
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
		poolValue := k.GetLockedPoolValue(ctx, pool)
		if poolValue.LTE(sdk.ZeroDec()) {
			continue
		}
		pools = append(pools, pool)
		lockedPoolValue[poolName] = poolValue
		totalPoolsValue = totalPoolsValue.Add(poolValue)
	}
	return lockedPoolValue, pools, totalPoolsValue
}
