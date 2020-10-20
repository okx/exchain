package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
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

func (k Keeper) isPoolNameExistedInWhiteList(ctx sdk.Context, poolName string) bool {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolsYieldNativeTokenPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if poolName == types.SplitPoolsYieldNativeTokenKey(iterator.Key()) {
			return true
		}
	}

	return false
}

func (k Keeper) satisfyWhiteListAdmittance(ctx sdk.Context, pool types.FarmPool) sdk.Error {
	quoteTokenSymbol := k.GetParams(ctx).QuoteSymbol
	if !swaptypes.IsPoolToken(pool.SymbolLocked) {
		// locked token is common token
		// check the existence of locked token with default quoteTokenSymbol in Params
		if !k.isSwapTokenPairExisted(ctx, pool.SymbolLocked, quoteTokenSymbol) {
			return types.ErrTokenNotExist(types.DefaultParamspace, swaptypes.GetSwapTokenPairName(pool.SymbolLocked, quoteTokenSymbol))
		}

		return nil
	}

	// locked token is lpt
	tokenSymbol0, tokenSymbol1 := swaptypes.SplitPoolToken(pool.SymbolLocked)
	if tokenSymbol0 == quoteTokenSymbol || tokenSymbol1 == quoteTokenSymbol {
		// base or quote token contains default quoteTokenSymbol in Params
		// check the existence of locked token
		if !k.isSwapTokenPairExisted(ctx, "", "", pool.SymbolLocked) {
			return types.ErrTokenNotExist(types.DefaultParamspace, pool.SymbolLocked)
		}

		return nil
	}

	// base or quote token don't contain default quoteTokenSymbol in Params
	// check the existence of locked token both with default quoteTokenSymbol in Params
	if !k.isSwapTokenPairExisted(ctx, tokenSymbol0, quoteTokenSymbol) {
		return types.ErrTokenNotExist(types.DefaultParamspace, swaptypes.GetSwapTokenPairName(tokenSymbol0, quoteTokenSymbol))
	}

	if !k.isSwapTokenPairExisted(ctx, tokenSymbol1, quoteTokenSymbol) {
		return types.ErrTokenNotExist(types.DefaultParamspace, swaptypes.GetSwapTokenPairName(tokenSymbol1, quoteTokenSymbol))
	}

	return nil
}

func (k Keeper) isSwapTokenPairExisted(ctx sdk.Context, baseTokenSymbol, quoteTokenSymbol string, fullSwapTokenName ...string) bool {
	var swapTokenPairName string
	if len(fullSwapTokenName) == 0 {
		swapTokenPairName = swaptypes.GetSwapTokenPairName(baseTokenSymbol, quoteTokenSymbol)
	} else if len(fullSwapTokenName) == 1 {
		swapTokenPairName = fullSwapTokenName[0]
	} else {
		// over one swap token name to check
		return false
	}
	_, err := k.swapKeeper.GetSwapTokenPair(ctx, swapTokenPairName)
	if err != nil {
		return false
	}

	return true
}
