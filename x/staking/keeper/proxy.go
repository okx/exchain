package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
)

// ClearProxy clears the ProxyAddress on the delegator who has bound
func (k Keeper) ClearProxy(ctx sdk.Context, proxyAddr sdk.AccAddress) {
	k.IterateProxy(ctx, proxyAddr, true, func(_ int64, delAddr, _ sdk.AccAddress) (stop bool) {
		delegator, found := k.GetDelegator(ctx, delAddr)
		if found {
			delegator.UnbindProxy()
			k.SetDelegator(ctx, delegator)
		}
		return false
	})
}

// SetProxyBinding sets or deletes the key of proxy relationship
func (k Keeper) SetProxyBinding(ctx sdk.Context, proxyAddress, delegatorAddress sdk.AccAddress, isRemove bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetProxyDelegatorKey(proxyAddress, delegatorAddress)

	if isRemove {
		store.Delete(key)
	} else {
		store.Set(key, []byte(""))
	}
}

// IterateProxy iterates all the info between delegator and its proxy
func (k Keeper) IterateProxy(ctx sdk.Context, proxyAddr sdk.AccAddress, isClear bool,
	fn func(index int64, delAddr, proxyAddr sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetProxyDelegatorKey(proxyAddr, []byte{}))
	defer iterator.Close()

	index := sdk.AddrLen + 1
	for i := int64(0); iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if stop := fn(i, key[index:], key[1:index]); stop {
			break
		}
		if isClear {
			store.Delete(key)
		}
		i++
	}
}

// UpdateShares withdraws and adds shares continuously on the same validator set with different amount of shares
func (k Keeper) UpdateShares(ctx sdk.Context, delAddr sdk.AccAddress, tokens sdk.Dec) sdk.Error {
	// get last validators that were added shares to and existing in the store
	vals, lastShares := k.GetLastValsAddedSharesExisted(ctx, delAddr)
	if vals == nil {
		// if the delegator never adds shares, just pass
		return nil
	}

	lenVals := len(vals)
	shares, sdkErr := calculateWeight(ctx.BlockTime().Unix(), tokens)
	if sdkErr != nil {
		return sdkErr
	}

	delegator, found := k.GetDelegator(ctx, delAddr)
	if !found {
		return types.ErrNoDelegatorExisted(types.DefaultCodespace, delAddr.String())
	}

	logger := k.Logger(ctx)
	for i := 0; i < lenVals; i++ {
		if vals[i].MinSelfDelegation.IsZero() {
			return types.ErrAddSharesToDismission(types.DefaultCodespace, vals[i].OperatorAddress.String())
		}

		// 1.delete related store
		k.DeleteValidatorByPowerIndex(ctx, vals[i])

		// 2.update shares
		k.SetShares(ctx, delAddr, vals[i].OperatorAddress, shares)

		// 3.update validator
		logger.Debug("update shares", vals[i].OperatorAddress.String(), vals[i].DelegatorShares.String())
		vals[i].DelegatorShares = vals[i].DelegatorShares.Sub(lastShares).Add(shares)
		k.SetValidator(ctx, vals[i])
		k.SetValidatorByPowerIndex(ctx, vals[i])
	}

	// update the delegator struct
	delegator.Shares = shares
	k.SetDelegator(ctx, delegator)

	return nil
}

// AddSharesToValidators adds shares to validators and return the amount of the shares
func (k Keeper) AddSharesToValidators(ctx sdk.Context, delAddr sdk.AccAddress, vals types.Validators, tokens sdk.Dec) (
	shares types.Shares, sdkErr sdk.Error) {
	lenVals := len(vals)
	shares, sdkErr = calculateWeight(ctx.BlockTime().Unix(), tokens)
	if sdkErr != nil {
		return
	}
	for i := 0; i < lenVals; i++ {
		k.addShares(ctx, delAddr, vals[i], shares)
	}
	return
}

// WithdrawLastShares withdraws the shares last time from the validators
func (k Keeper) WithdrawLastShares(ctx sdk.Context, delAddr sdk.AccAddress, lastValsAddedSharesTo types.Validators,
	lastShares types.Shares) {
	lenLastVals := len(lastValsAddedSharesTo)
	for i := 0; i < lenLastVals; i++ {
		k.withdrawShares(ctx, delAddr, lastValsAddedSharesTo[i], lastShares)
	}
}

func (k Keeper) withdrawShares(ctx sdk.Context, delAddr sdk.AccAddress, val types.Validator, shares types.Shares) {
	// 1.delete shares entity
	k.DeleteShares(ctx, val.OperatorAddress, delAddr)

	// 2.update validator entity
	k.DeleteValidatorByPowerIndex(ctx, val)

	// 3.update validator's shares
	val.DelegatorShares = val.GetDelegatorShares().Sub(shares)

	// 3.check whether the validator should be removed
	if val.IsUnbonded() && val.GetMinSelfDelegation().IsZero() && val.GetDelegatorShares().IsZero() {
		k.RemoveValidator(ctx, val.OperatorAddress)
		return
	}

	k.SetValidator(ctx, val)
	k.SetValidatorByPowerIndex(ctx, val)
}

func (k Keeper) addShares(ctx sdk.Context, delAddr sdk.AccAddress, val types.Validator, shares types.Shares) {
	// 1.update shares entity
	k.SetShares(ctx, delAddr, val.OperatorAddress, shares)

	// 2.update validator entity
	k.DeleteValidatorByPowerIndex(ctx, val)
	val.DelegatorShares = val.GetDelegatorShares().Add(shares)
	k.SetValidator(ctx, val)
	k.SetValidatorByPowerIndex(ctx, val)
}

// GetLastValsAddedSharesExisted gets last validators that the delegator added shares to last time
func (k Keeper) GetLastValsAddedSharesExisted(ctx sdk.Context, delAddr sdk.AccAddress) (types.Validators, types.Shares) {
	// 1.get delegator entity
	delegator, found := k.GetDelegator(ctx, delAddr)

	// if not found
	if !found {
		return nil, sdk.ZeroDec()
	}

	// 2.get validators that were added shares to and existing in the store
	lenVals := len(delegator.ValidatorAddresses)
	var vals types.Validators
	for i := 0; i < lenVals; i++ {
		val, found := k.GetValidator(ctx, delegator.ValidatorAddresses[i])
		if found {
			// the validator that were added shares to hasn't been removed
			vals = append(vals, val)
		}
	}

	return vals, delegator.Shares
}

// GetValidatorsToAddShares gets the validators from their validator addresses
func (k Keeper) GetValidatorsToAddShares(ctx sdk.Context, valAddrs []sdk.ValAddress) (types.Validators, sdk.Error) {
	lenVals := len(valAddrs)
	vals := make(types.Validators, lenVals)
	for i := 0; i < lenVals; i++ {
		val, found := k.GetValidator(ctx, valAddrs[i])
		if found {
			// get the validator hasn't been removed
			vals[i] = val
		} else {
			return nil, types.ErrNoValidatorFound(types.DefaultCodespace, valAddrs[i].String())
		}
	}

	return vals, nil
}
