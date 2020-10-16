package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swapkeeper "github.com/okex/okexchain/x/ammswap/keeper"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/farm/types"
)

func (k Keeper) SetFarmPool(ctx sdk.Context, pool types.FarmPool) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetFarmPoolKey(pool.Name), k.cdc.MustMarshalBinaryLengthPrefixed(pool))
}

func (k Keeper) GetFarmPool(ctx sdk.Context, poolName string) (pool types.FarmPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetFarmPoolKey(poolName))
	if bz == nil {
		return pool, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pool)
	return pool, true
}

func (k Keeper) DeleteFarmPool(ctx sdk.Context, poolName string) {
	store := ctx.KVStore(k.storeKey)
	// delete pool from whitelist
	store.Delete(types.GetWhitelistMemberKey(poolName))
	// delete pool key
	store.Delete(types.GetFarmPoolKey(poolName))
}

// getFarmPoolNamesForAccount gets all pool names that an account has locked coins in from the store
func (k Keeper) getFarmPoolNamesForAccount(ctx sdk.Context, accAddr sdk.AccAddress) (poolNames types.PoolNameList) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.Address2PoolPrefix, accAddr...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		poolNames = append(poolNames, types.SplitPoolNameFromLockInfoKey(iterator.Key()))
	}

	return
}

// getAccountsLockedTo gets all addresses of accounts that have locked coins in a pool
func (k Keeper) getAccountsLockedTo(ctx sdk.Context, poolName string) (accAddrList types.AccAddrList) {
	//TODO:
	tmp := []string{
		"okexchain13hxtcclwfhks2vzrm3hlx0kp4cuxggg26p9kqf",
		"okexchain1gnq40xzh75fl8yzlk32hwx3xeew7nycndhj843",
		"okexchain1nh4mtjpstnvtxf24vjzgnf5slqksja9fsdt7fl",
	}

	for _, s := range tmp {
		accAddr, _ := sdk.AccAddressFromBech32(s)
		accAddrList = append(accAddrList, accAddr)
	}
	return
}

// getPoolNum gets the number of pools that already exist
func (k Keeper) getPoolNum(ctx sdk.Context) types.PoolNum {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FarmPoolPrefix)
	defer iterator.Close()
	var num uint
	for ; iterator.Valid(); iterator.Next() {
		num++
	}

	return types.NewPoolNum(num)
}

// getFarmPools gets all pools that exist currently in the store
func (k Keeper) getFarmPools(ctx sdk.Context) (pools types.FarmPools) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.FarmPoolPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.FarmPool
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &pool)
		pools = append(pools, pool)
	}

	return
}

func (k Keeper) SetAddressInFarmPool(ctx sdk.Context, poolName string, addr sdk.AccAddress, amount sdk.DecCoin) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAddressInFarmPoolKey(poolName, addr),  k.cdc.MustMarshalBinaryLengthPrefixed(amount))
}

func (k Keeper) GetAddressInFarmPool(ctx sdk.Context, poolName string, addr sdk.AccAddress) (amount sdk.DecCoin) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAddressInFarmPoolKey(poolName, addr))
	if bz == nil {
		return amount
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &amount)
	return amount
}

// HasAddressInFarmPool check existence of the pool associated with a address
func (k Keeper) HasAddressInFarmPool(ctx sdk.Context, poolName string, addr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetAddressInFarmPoolKey(poolName, addr))
}

func (k Keeper) DeleteAddressInFarmPool(ctx sdk.Context, poolName string, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAddressInFarmPoolKey(poolName, addr))
}

// getFarmPools gets all addresses that locked tokens in a pool
func (k Keeper) getAddressesInFarmPool(ctx sdk.Context, poolName string) (addrs []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.Pool2AddressPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		// TODO

	}

	return
}

func (k Keeper) SetLockInfo(ctx sdk.Context, lockInfo types.LockInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLockInfoKey(lockInfo.Owner, lockInfo.PoolName), k.cdc.MustMarshalBinaryLengthPrefixed(lockInfo))
}

func (k Keeper) GetLockInfo(ctx sdk.Context, addr sdk.AccAddress, poolName string) (info types.LockInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetLockInfoKey(addr, poolName))
	if bz == nil {
		return info, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
	return info, true
}

// HasLockInfo check existence of the address associated with a pool
func (k Keeper) HasLockInfo(ctx sdk.Context, addr sdk.AccAddress, poolName string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetLockInfoKey(addr, poolName))
}

func (k Keeper) DeleteLockInfo(ctx sdk.Context, addr sdk.AccAddress, poolName string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLockInfoKey(addr, poolName))
}

func (k Keeper) GetLockedPoolValue(ctx sdk.Context, pool types.FarmPool) sdk.Dec {
	if pool.TotalValueLocked.Amount.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}

	poolValue := sdk.ZeroDec()
	params := k.GetParams(ctx)
	quoteToken := params.QuoteSymbol
	swapParams := k.swapKeeper.GetParams(ctx)
	// calculate locked lpt value
	if swaptypes.IsPoolToken(pool.SymbolLocked) {
		token0, token1 := swaptypes.SplitPoolToken(pool.SymbolLocked)
		if token0 == quoteToken || token1 == quoteToken {
			// calculate how much assets the TotalValueLocked can redeem
			token0Amount, token1Amount, err := k.swapKeeper.GetRedeemableAssets(ctx, token0, token1,
				pool.TotalValueLocked.Amount)
			if err != nil {
				panic("should not happen")
			}
			var baseTokenAmount, quoteTokenAmount sdk.DecCoin
			if token0Amount.Denom == quoteToken {
				baseTokenAmount = token1Amount
				quoteTokenAmount = token0Amount
			} else {
				baseTokenAmount = token0Amount
				quoteTokenAmount = token1Amount
			}
			swappedQuoteTokenAmt := k.GetSwappedQuoteTokenAmt(ctx, baseTokenAmount, quoteToken, swapParams)
			poolValue = swappedQuoteTokenAmt.Add(quoteTokenAmount.Amount)
		} else {
			// calculate how much assets the TotalValueLocked can redeem
			token0Amount, token1Amount, err := k.swapKeeper.GetRedeemableAssets(ctx, token0, token1,
				pool.TotalValueLocked.Amount)
			if err != nil {
				panic("should not happen")
			}
			// calculate how much quote token the base token can swap
			quote0TokenAmt := k.GetSwappedQuoteTokenAmt(ctx, token0Amount, quoteToken, swapParams)
			quote1TokenAmt := k.GetSwappedQuoteTokenAmt(ctx, token1Amount, quoteToken, swapParams)
			poolValue = quote0TokenAmt.Add(quote1TokenAmt)
		}
	} else {
		poolValue = k.GetSwappedQuoteTokenAmt(ctx, pool.TotalValueLocked, quoteToken, swapParams)
	}
	return poolValue
}

func (k Keeper) GetSwappedQuoteTokenAmt(
	ctx sdk.Context, coin sdk.DecCoin, quoteToken string, params swaptypes.Params,
) sdk.Dec {
	// calculate how much quote token the base token can swap
	tokenPair, err := k.swapKeeper.GetSwapTokenPair(ctx, swaptypes.GetSwapTokenPairName(coin.Denom, quoteToken))
	if err != nil {
		return sdk.ZeroDec()
	}
	swappedCoin := swapkeeper.CalculateTokenToBuy(tokenPair, coin, quoteToken, params)
	return swappedCoin.Amount
}
