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
func (k Keeper) getAccountsLockedTo(ctx sdk.Context, poolName string) (lockerAddrList types.AccAddrList) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, append(types.Pool2AddressPrefix, []byte(poolName)...))
	defer iterator.Close()

	splitIndex := 1 + len(poolName)
	for ; iterator.Valid(); iterator.Next() {
		lockerAddrList = append(lockerAddrList, iterator.Key()[splitIndex:])
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

// GetFarmPools gets all pools that exist currently in the store
func (k Keeper) GetFarmPools(ctx sdk.Context) (pools types.FarmPools) {
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

func (k Keeper) SetAddressInFarmPool(ctx sdk.Context, poolName string, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAddressInFarmPoolKey(poolName, addr), []byte(""))
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

// GetPoolLockedValue gets the value of locked tokens in pool priced in quote symbol
func (k Keeper) GetPoolLockedValue(ctx sdk.Context, pool types.FarmPool) sdk.Dec {
	if pool.TotalValueLocked.Amount.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}

	poolValue := sdk.ZeroDec()
	params := k.GetParams(ctx)
	quoteSymbol := params.QuoteSymbol
	swapParams := k.swapKeeper.GetParams(ctx)
	// calculate locked lpt value
	if swaptypes.IsPoolToken(pool.LockedSymbol) {
		poolValue = k.calculateLockedLPTValue(ctx, pool, quoteSymbol, swapParams)
	} else {
		poolValue = k.calculateQuoteTokenAmtToBuy(ctx, pool.TotalValueLocked, quoteSymbol, swapParams)
	}
	return poolValue
}

func (k Keeper) calculateLockedLPTValue(
	ctx sdk.Context, pool types.FarmPool, quoteSymbol string, swapParams swaptypes.Params,
) (poolValue sdk.Dec) {
	token0Symbol, token1Symbol := swaptypes.SplitPoolToken(pool.LockedSymbol)

	// calculate how much assets the TotalValueLocked can redeem
	token0Amount, token1Amount, err := k.swapKeeper.GetRedeemableAssets(ctx, token0Symbol, token1Symbol,
		pool.TotalValueLocked.Amount)
	if err != nil {
		panic("should not happen")
	}
	
	if token0Symbol == quoteSymbol || token1Symbol == quoteSymbol {
		return k.calculateLPTValueWithQuote(ctx, token0Amount, token1Amount, quoteSymbol, swapParams)
	}
	return k.calculateLPTValueWithoutQuote(ctx, token0Amount, token1Amount, quoteSymbol, swapParams)
}

// calculateLPTValueWithQuote calculates the value of LPT which represents token pair containing quote symbol
func (k Keeper) calculateLPTValueWithQuote(
	ctx sdk.Context, token0Amount, token1Amount sdk.DecCoin, quoteSymbol string, swapParams swaptypes.Params,
) sdk.Dec {
	var baseTokenAmount, quoteTokenAmount sdk.DecCoin
	if token0Amount.Denom == quoteSymbol {
		baseTokenAmount = token1Amount
		quoteTokenAmount = token0Amount
	} else {
		baseTokenAmount = token0Amount
		quoteTokenAmount = token1Amount
	}
	swappedQuoteTokenAmt := k.calculateQuoteTokenAmtToBuy(ctx, baseTokenAmount, quoteSymbol, swapParams)
	poolValue := swappedQuoteTokenAmt.Add(quoteTokenAmount.Amount)
	return poolValue
}

// calculateLPTValueWithoutQuote calculates the value of LPT which represents token pair not containing quote symbol
func (k Keeper) calculateLPTValueWithoutQuote(
	ctx sdk.Context, token0Amount, token1Amount sdk.DecCoin, quoteSymbol string, swapParams swaptypes.Params,
) sdk.Dec {
	// calculate how much quote token the base token can swap
	quote0TokenAmt := k.calculateQuoteTokenAmtToBuy(ctx, token0Amount, quoteSymbol, swapParams)
	quote1TokenAmt := k.calculateQuoteTokenAmtToBuy(ctx, token1Amount, quoteSymbol, swapParams)
	return quote0TokenAmt.Add(quote1TokenAmt)
}

func (k Keeper) calculateQuoteTokenAmtToBuy(
	ctx sdk.Context, coin sdk.DecCoin, quoteSymbol string, params swaptypes.Params,
) sdk.Dec {
	// calculate how much quote token the base token can swap
	tokenPair, err := k.swapKeeper.GetSwapTokenPair(ctx, swaptypes.GetSwapTokenPairName(coin.Denom, quoteSymbol))
	if err != nil {
		return sdk.ZeroDec()
	}
	swappedCoin := swapkeeper.CalculateTokenToBuy(tokenPair, coin, quoteSymbol, params)
	return swappedCoin.Amount
}

// Iterate over all lock infos
func (k Keeper) IterateAllLockInfos(
	ctx sdk.Context, handler func(lockInfo types.LockInfo) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.Address2PoolPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var lockInfo types.LockInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &lockInfo)
		if handler(lockInfo) {
			break
		}
	}
}
