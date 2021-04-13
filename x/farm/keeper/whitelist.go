package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swaptypes "github.com/okex/exchain/x/ammswap/types"
	"github.com/okex/exchain/x/farm/types"
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

// DeleteWhiteList removes the pool name from whitelist
func (k Keeper) DeleteWhiteList(ctx sdk.Context, poolName string) {
	ctx.KVStore(k.storeKey).Delete(types.GetWhitelistMemberKey(poolName))
}

func (k Keeper) isPoolNameExistedInWhiteList(ctx sdk.Context, poolName string) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetWhitelistMemberKey(poolName))
}

func (k Keeper) satisfyWhiteListAdmittance(ctx sdk.Context, pool types.FarmPool) sdk.Error {
	quoteTokenSymbol := k.GetParams(ctx).QuoteSymbol
	// lock token is quote symbol
	if pool.MinLockAmount.Denom == quoteTokenSymbol {
		return nil
	}
	if !swaptypes.IsPoolToken(pool.MinLockAmount.Denom) {
		// locked token is common token
		// check the existence of locked token with default quoteTokenSymbol in Params
		if !k.isSwapTokenPairExisted(ctx, pool.MinLockAmount.Denom, quoteTokenSymbol) {
			return types.ErrSwapTokenPairNotExist(swaptypes.GetSwapTokenPairName(pool.MinLockAmount.Denom, quoteTokenSymbol))

		}

		return nil
	}

	// locked token is lpt
	tokenSymbol0, tokenSymbol1 := swaptypes.SplitPoolToken(pool.MinLockAmount.Denom)
	if tokenSymbol0 == quoteTokenSymbol || tokenSymbol1 == quoteTokenSymbol {
		// base or quote token contains default quoteTokenSymbol in Params
		// check the existence of locked token
		if !k.isSwapTokenPairExisted(ctx, tokenSymbol0, tokenSymbol1) {
			return types.ErrTokenNotExist(pool.MinLockAmount.Denom)
		}

		return nil
	}

	// base or quote token don't contain default quoteTokenSymbol in Params
	// check the existence of locked token both with default quoteTokenSymbol in Params
	if !k.isSwapTokenPairExisted(ctx, tokenSymbol0, quoteTokenSymbol) {
		return types.ErrSwapTokenPairNotExist(swaptypes.GetSwapTokenPairName(tokenSymbol0, quoteTokenSymbol))
	}

	if !k.isSwapTokenPairExisted(ctx, tokenSymbol1, quoteTokenSymbol) {
		return types.ErrSwapTokenPairNotExist(swaptypes.GetSwapTokenPairName(tokenSymbol1, quoteTokenSymbol))
	}

	return nil
}

func (k Keeper) isSwapTokenPairExisted(ctx sdk.Context, baseTokenSymbol, quoteTokenSymbol string) bool {
	_, err := k.swapKeeper.GetSwapTokenPair(ctx, swaptypes.GetSwapTokenPairName(baseTokenSymbol, quoteTokenSymbol))
	if err != nil {
		return false
	}

	return true
}
