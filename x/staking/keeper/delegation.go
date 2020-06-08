package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
)

// UpdateProxy updates the shares by the total delegated and self delegated tokens of a proxy
func (k Keeper) UpdateProxy(ctx sdk.Context, delegator types.Delegator, tokens sdk.Dec) (err sdk.Error) {
	if !delegator.HasProxy() {
		return nil
	}
	// delegator has bound a proxy, need update proxy's shares
	if proxy, found := k.GetDelegator(ctx, delegator.ProxyAddress); found {
		// tokens might be negative
		proxy.TotalDelegatedTokens = proxy.TotalDelegatedTokens.Add(tokens)
		if proxy.TotalDelegatedTokens.LT(sdk.ZeroDec()) {
			return types.ErrInvalidProxyUpdating(types.DefaultCodespace)
		}

		finalTokens := proxy.TotalDelegatedTokens.Add(proxy.Tokens)
		k.SetDelegator(ctx, proxy)
		return k.UpdateShares(ctx, proxy.DelegatorAddress, finalTokens)
	}
	return sdk.ErrInvalidAddress(delegator.ProxyAddress.String())
}

// Delegate handles the process of delegating
func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, token sdk.DecCoin) sdk.Error {

	delQuantity, minDelLimit := token.Amount, k.ParamsMinDelegation(ctx)
	if delQuantity.LT(minDelLimit) {
		return types.ErrInsufficientQuantity(types.DefaultCodespace, delQuantity.String(), minDelLimit.String())
	}

	// 1.transfer account's okt into bondPool
	coins := token.ToCoins()
	if err := k.supplyKeeper.DelegateCoinsFromAccountToModule(ctx, delAddr, types.BondedPoolName, coins); err != nil {
		return err
	}

	// 2.get
	delegator, found := k.GetDelegator(ctx, delAddr)
	if !found {
		delegator = types.NewDelegator(delAddr)
	}

	// 3.update delegator
	delegator.Tokens = delegator.Tokens.Add(delQuantity)
	k.SetDelegator(ctx, delegator)

	if delegator.HasProxy() {
		//delegator have bound with some proxy, need update proxy's shares
		return k.UpdateProxy(ctx, delegator, delQuantity)

	}
	// 4.update shares when delAddr has added already
	return k.UpdateShares(ctx, delegator.DelegatorAddress, delegator.Tokens)
}

// Withdraw handles the process of withdrawing token from deposit account
func (k Keeper) Withdraw(ctx sdk.Context, delAddr sdk.AccAddress, token sdk.DecCoin) (time.Time, sdk.Error) {
	delegator, found := k.GetDelegator(ctx, delAddr)
	if !found {
		return time.Time{}, types.ErrNoDelegationToAddShares(types.DefaultCodespace, delAddr.String())
	}
	quantity, minDelLimit := token.Amount, k.ParamsMinDelegation(ctx)
	if quantity.LT(minDelLimit) {
		return time.Time{}, types.ErrInsufficientQuantity(types.DefaultCodespace, quantity.String(), minDelLimit.String())
	} else if delegator.Tokens.LT(quantity) {
		return time.Time{}, types.ErrInsufficientDelegation(types.DefaultCodespace, quantity.String(), delegator.Tokens.String())
	}

	// 1.some okt transfer bondPool into unbondPool
	k.bondedTokensToNotBonded(ctx, token)

	// 2.delete delegator in store, or set back
	leftTokens := delegator.Tokens.Sub(quantity)
	if delegator.HasProxy() {
		if sdkErr := k.UpdateProxy(ctx, delegator, quantity.Mul(sdk.NewDec(-1))); sdkErr != nil {
			return time.Time{}, sdkErr
		}
	}
	if leftTokens.IsZero() {
		// withdraw all shares
		lastVals, lastShares := k.GetLastValsAddedSharesExisted(ctx, delAddr)
		k.WithdrawLastShares(ctx, delAddr, lastVals, lastShares)
		if delegator.HasProxy() {
			k.SetProxyBinding(ctx, delegator.ProxyAddress, delAddr, true)
		}
		k.DeleteDelegator(ctx, delAddr)
	} else {
		delegator.Tokens = leftTokens
		k.SetDelegator(ctx, delegator)
		if !delegator.HasProxy() {
			if err := k.UpdateShares(ctx, delegator.DelegatorAddress, delegator.Tokens); err != nil {
				return time.Time{}, err
			}
		}
	}

	// 3.set undelegation and into store
	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	undelegation, found := k.GetUndelegating(ctx, delAddr)
	if !found {
		undelegation = types.NewUndelegationInfo(delAddr, quantity, completionTime)
	} else {
		k.DeleteAddrByTimeKey(ctx, undelegation.CompletionTime, delAddr)
		undelegation.Quantity = undelegation.Quantity.Add(quantity)
		undelegation.CompletionTime = completionTime
	}
	k.SetUndelegating(ctx, undelegation)
	k.SetAddrByTimeKeyWithNilValue(ctx, completionTime, delAddr)

	return completionTime, nil
}

// GetUndelegating gets UndelegationInfo entity from store
func (k Keeper) GetUndelegating(ctx sdk.Context, delAddr sdk.AccAddress) (undelegationInfo types.UndelegationInfo,
	found bool) {
	bytes := ctx.KVStore(k.storeKey).Get(types.GetUndelegationInfoKey(delAddr))
	if bytes == nil {
		return undelegationInfo, false
	}

	undelegationInfo = types.MustUnMarshalUndelegationInfo(k.cdc, bytes)
	return undelegationInfo, true
}

// SetUndelegating sets UndelegationInfo entity to store
func (k Keeper) SetUndelegating(ctx sdk.Context, undelegationInfo types.UndelegationInfo) {
	key := types.GetUndelegationInfoKey(undelegationInfo.DelegatorAddress)
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(undelegationInfo)
	ctx.KVStore(k.storeKey).Set(key, bytes)
}

// DeleteUndelegating deletes UndelegationInfo from store
func (k Keeper) DeleteUndelegating(ctx sdk.Context, delAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetUndelegationInfoKey(delAddr))
}

// CompleteUndelegation handles the final process when the undelegation is completed
func (k Keeper) CompleteUndelegation(ctx sdk.Context, delAddr sdk.AccAddress) (sdk.Dec, sdk.Error) {
	ud, found := k.GetUndelegating(ctx, delAddr)
	if !found {
		return sdk.NewDec(0), types.ErrNotInDelegating(k.Codespace(), delAddr.String())
	}

	coin := sdk.NewDecCoinsFromDec(k.GetParams(ctx).BondDenom, ud.Quantity)

	err := k.supplyKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.NotBondedPoolName, ud.DelegatorAddress, coin)
	if err != nil {
		return sdk.NewDec(0), err
	}

	k.DeleteUndelegating(ctx, delAddr)
	return ud.Quantity, nil
}

// IterateUndelegationInfo iterates through all of the undelegation info
func (k Keeper) IterateUndelegationInfo(ctx sdk.Context,
	fn func(index int64, undelegationInfo types.UndelegationInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnDelegationInfoKey)
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		var undelegationInfo types.UndelegationInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &undelegationInfo)
		if stop := fn(i, undelegationInfo); stop {
			break
		}
		i++
	}
}

// SetAddrByTimeKeyWithNilValue sets the time+delAddr key into store with an empty value
func (k Keeper) SetAddrByTimeKeyWithNilValue(ctx sdk.Context, timestamp time.Time, delAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetCompleteTimeWithAddrKey(timestamp, delAddr), []byte{})
}

// DeleteAddrByTimeKey deletes the time+delAddr key from store
func (k Keeper) DeleteAddrByTimeKey(ctx sdk.Context, timestamp time.Time, delAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetCompleteTimeWithAddrKey(timestamp, delAddr))
}

// IterateKeysBeforeCurrentTime iterates for all keys of (time+delAddr) from time 0 until the current Blockheader time
func (k Keeper) IterateKeysBeforeCurrentTime(ctx sdk.Context, currentTime time.Time,
	fn func(index int64, key []byte) (stop bool)) {

	timeKeyIterator := k.getAddrByTimeKeyIterator(ctx, currentTime)
	defer timeKeyIterator.Close()

	for i := int64(0); timeKeyIterator.Valid(); timeKeyIterator.Next() {
		key := timeKeyIterator.Key()
		if stop := fn(i, key); stop {
			break
		}
		i++
	}
}

// getAddrByTimeKeyIterator gets the iterator of keys from time 0 until endTime
func (k Keeper) getAddrByTimeKeyIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCompleteTimeKey(endTime)
	return store.Iterator(types.UnDelegateQueueKey, sdk.PrefixEndBytes(key))
}
