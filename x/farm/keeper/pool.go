package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swapkeeper "github.com/okex/okexchain/x/ammswap/keeper"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/farm/types"
)

func (k Keeper) SetFarmPool(ctx sdk.Context, pool types.FarmPool) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetFarmPoolKey(pool.PoolName), k.cdc.MustMarshalBinaryLengthPrefixed(pool))
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

func (k Keeper) SetLockInfo(ctx sdk.Context, lockInfo types.LockInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetLockInfoKey(lockInfo.Address, lockInfo.PoolName), k.cdc.MustMarshalBinaryLengthPrefixed(lockInfo))
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

func (k Keeper) GetLockedPoolValue(ctx sdk.Context, pool types.FarmPool) sdk.Dec {
	if pool.TotalLockedCoin.Amount.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}

	poolValue := sdk.ZeroDec()
	params := k.GetParams(ctx)
	quoteToken := params.QuoteToken
	swapParams := k.swapKeeper.GetParams(ctx)
	// calculate locked lpt value
	if swaptypes.IsPoolToken(pool.LockedTokenSymbol) {
		token0, token1 := swaptypes.SplitPoolToken(pool.LockedTokenSymbol)
		if token0 == quoteToken || token1 == quoteToken {
			// calculate how much assets the TotalLockedCoin can redeem
			token0Coin, token1Coin, err := k.swapKeeper.GetRedeemableAssets(ctx, token0, token1,
				pool.TotalLockedCoin.Amount)
			if err != nil {
				panic("should not happen")
			}
			var baseCoin, quoteCoin sdk.DecCoin
			if token0Coin.Denom == quoteToken {
				baseCoin = token1Coin
				quoteCoin = token0Coin
			} else {
				baseCoin = token0Coin
				quoteCoin = token1Coin
			}
			quoteTokenAmt := k.GetSwappedQuoteTokenAmt(ctx, baseCoin, quoteToken, swapParams)
			poolValue = quoteTokenAmt.Add(quoteCoin.Amount)
		} else {
			// calculate how much assets the TotalLockedCoin can redeem
			token0Coin, token1Coin, err := k.swapKeeper.GetRedeemableAssets(ctx, token0, token1,
				pool.TotalLockedCoin.Amount)
			if err != nil {
				panic("should not happen")
			}
			// calculate how much quote token the base token can swap
			quote0TokenAmt := k.GetSwappedQuoteTokenAmt(ctx, token0Coin, quoteToken, swapParams)
			quote1TokenAmt := k.GetSwappedQuoteTokenAmt(ctx, token1Coin, quoteToken, swapParams)
			poolValue = quote0TokenAmt.Add(quote1TokenAmt)
		}
	} else {
		poolValue = k.GetSwappedQuoteTokenAmt(ctx, pool.TotalLockedCoin, quoteToken, swapParams)
	}
	return poolValue
}

func (k Keeper) GetSwappedQuoteTokenAmt(
	ctx sdk.Context, coin sdk.DecCoin, quoteToken string, params swaptypes.Params,
) sdk.Dec {
	// calculate how much quote token the base token can swap
	tokenPair, err := k.swapKeeper.GetSwapTokenPair(ctx, swaptypes.GetSwapTokenPairName(coin.Denom, quoteToken))
	if err != nil {
		panic("should not happen")
	}
	swappedCoin := swapkeeper.CalculateTokenToBuy(tokenPair, coin, quoteToken, params)
	return swappedCoin.Amount
}
