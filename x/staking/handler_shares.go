package staking

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/keeper"
	"github.com/okex/okchain/x/staking/types"
)

func handleMsgBindProxy(ctx sdk.Context, msg types.MsgBindProxy, k keeper.Keeper) sdk.Result {
	delegator, found := k.GetDelegator(ctx, msg.DelAddr)

	if !found || delegator.Tokens.IsZero() {
		return types.ErrNoDelegationToAddShares(types.DefaultCodespace, msg.DelAddr.String()).Result()
	}

	if !delegator.Shares.Equal(sdk.ZeroDec()) {
		return types.ErrAlreadyAddedShares(types.DefaultCodespace, delegator.DelegatorAddress.String()).Result()
	}

	// proxy must delegated
	proxyDelegator, found := k.GetDelegator(ctx, msg.ProxyAddress)
	if !found || proxyDelegator.Tokens.IsZero() {
		return types.ErrNotFoundProxy(types.DefaultCodespace, msg.ProxyAddress.String()).Result()
	}

	// target delegator must reg as a proxy
	if !proxyDelegator.IsProxy {
		return types.ErrDelegatorNotAProxy(types.DefaultCodespace, msg.ProxyAddress.String()).Result()
	}

	// double proxy is denied on okchain
	if delegator.IsProxy {
		return types.ErrDoubleProxy(types.DefaultCodespace, delegator.DelegatorAddress.String()).Result()
	}

	// same proxy, only update shares
	if delegator.HasProxy() && delegator.ProxyAddress.Equals(proxyDelegator.DelegatorAddress) {
		updateTokens := proxyDelegator.TotalDelegatedTokens.Add(proxyDelegator.Tokens)
		if err := k.UpdateShares(ctx, proxyDelegator.DelegatorAddress, updateTokens); err != nil {
			return types.ErrInvalidDelegation(types.DefaultCodespace, proxyDelegator.DelegatorAddress.String()).Result()
		}
		return sdk.Result{Events: ctx.EventManager().Events()}
	}

	// unbind from the original proxy
	if len(delegator.ProxyAddress) != 0 {
		if sdkErr := unbindProxy(ctx, delegator.DelegatorAddress, k); sdkErr != nil {
			return sdkErr.Result()
		}
	}

	// bind proxy relationship
	delegator.BindProxy(msg.ProxyAddress)

	// update proxy's shares weight
	proxyDelegator.TotalDelegatedTokens = proxyDelegator.TotalDelegatedTokens.Add(delegator.Tokens)

	k.SetDelegator(ctx, delegator)
	k.SetDelegator(ctx, proxyDelegator)
	k.SetProxyBinding(ctx, proxyDelegator.DelegatorAddress, delegator.DelegatorAddress, false)

	finalTokens := proxyDelegator.TotalDelegatedTokens.Add(proxyDelegator.Tokens)

	if err := k.UpdateShares(ctx, proxyDelegator.DelegatorAddress, finalTokens); err != nil {
		return types.ErrInvalidDelegation(types.DefaultCodespace, proxyDelegator.DelegatorAddress.String()).Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func unbindProxy(ctx sdk.Context, delAddr sdk.AccAddress, k keeper.Keeper) sdk.Error {
	delegator, found := k.GetDelegator(ctx, delAddr)
	if !found {
		return types.ErrNoDelegationToAddShares(types.DefaultCodespace, delAddr.String())
	}

	proxyDelegator, found := k.GetDelegator(ctx, delegator.ProxyAddress)
	if !found {
		return types.ErrNotFoundProxy(types.DefaultCodespace, delegator.ProxyAddress.String())
	}

	// update proxy's shares weight
	if k.UpdateProxy(ctx, delegator, delegator.Tokens.Mul(sdk.NewDec(-1))) != nil {
		return types.ErrInvalidDelegation(types.DefaultCodespace, proxyDelegator.DelegatorAddress.String())
	}
	// unbind proxy relationship
	delegator.UnbindProxy()
	k.SetDelegator(ctx, delegator)
	k.SetProxyBinding(ctx, proxyDelegator.DelegatorAddress, delegator.DelegatorAddress, true)

	return nil
}

func handleMsgUnbindProxy(ctx sdk.Context, msg types.MsgUnbindProxy, k keeper.Keeper) sdk.Result {
	if sdkErr := unbindProxy(ctx, msg.DelAddr, k); sdkErr != nil {
		return sdkErr.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func regProxy(ctx sdk.Context, proxyAddr sdk.AccAddress, k keeper.Keeper) sdk.Result {
	// check status
	proxy, found := k.GetDelegator(ctx, proxyAddr)
	if !found {
		return types.ErrNoDelegationToAddShares(types.DefaultCodespace, proxyAddr.String()).Result()
	}
	if proxy.IsProxy {
		return types.ErrAlreadyProxied(types.DefaultCodespace, proxyAddr.String()).Result()
	}
	if len(proxy.ProxyAddress) != 0 {
		return types.ErrAlreadyBound(types.DefaultCodespace, proxy.DelegatorAddress.String()).Result()
	}

	proxy.RegProxy(true)
	k.SetDelegator(ctx, proxy)

	if k.UpdateShares(ctx, proxy.DelegatorAddress, proxy.Tokens) != nil {
		return types.ErrInvalidDelegation(types.DefaultCodespace, proxy.DelegatorAddress.String()).Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}

}

func unregProxy(ctx sdk.Context, proxyAddr sdk.AccAddress, k keeper.Keeper) sdk.Result {
	// check status
	proxy, found := k.GetDelegator(ctx, proxyAddr)
	if !found {
		return types.ErrNotFoundProxy(types.DefaultCodespace, proxyAddr.String()).Result()
	}

	if !proxy.IsProxy {
		return types.ErrNeverProxied(types.DefaultCodespace, proxyAddr.String()).Result()
	}

	proxy.RegProxy(false)
	// unreg action, we need to erase all proxy relationship
	proxy.TotalDelegatedTokens = sdk.ZeroDec()
	k.ClearProxy(ctx, proxy.DelegatorAddress)
	k.SetDelegator(ctx, proxy)

	if k.UpdateShares(ctx, proxy.DelegatorAddress, proxy.Tokens) != nil {
		return types.ErrInvalidDelegation(types.DefaultCodespace, proxy.DelegatorAddress.String()).Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}

}

func handleRegProxy(ctx sdk.Context, msg types.MsgRegProxy, k keeper.Keeper) sdk.Result {
	if msg.Reg {
		return regProxy(ctx, msg.ProxyAddress, k)
	}

	return unregProxy(ctx, msg.ProxyAddress, k)
}

func handleMsgAddShares(ctx sdk.Context, msg types.MsgAddShares, k keeper.Keeper) sdk.Result {
	maxValsToAddShares := int(k.ParamsMaxValsToAddShares(ctx))
	if len(msg.ValAddrs) == 0 {
		return types.ErrNilValidatorAddrs(DefaultCodespace).Result()
	} else if len(msg.ValAddrs) > maxValsToAddShares {
		return types.ErrExceedValidatorAddrs(DefaultCodespace, maxValsToAddShares).Result()
	}

	// 0. check whether the delegator has delegation
	delegator, found := k.GetDelegator(ctx, msg.DelAddr)
	if !found || delegator.Tokens.IsZero() {
		return types.ErrNoDelegationToAddShares(types.DefaultCodespace, msg.DelAddr.String()).Result()
	}
	if delegator.HasProxy() {
		return types.ErrAddSharesDuringProxy(types.DefaultCodespace, delegator.DelegatorAddress.String(),
			delegator.ProxyAddress.String()).Result()
	}

	// 1. get last validators which were added shares to and existed in the store
	lastVals, lastShares := k.GetLastValsAddedSharesExisted(ctx, msg.DelAddr)

	// 2. withdraw the shares last time
	k.WithdrawLastShares(ctx, msg.DelAddr, lastVals, lastShares)

	// 3. get validators to add shares this time (if the validator doesn't exist, return error)
	vals, sdkErr := k.GetValidatorsToAddShares(ctx, msg.ValAddrs)
	if sdkErr != nil {
		return sdkErr.Result()
	}
	if sdkErr = validateSharesAdding(vals); sdkErr != nil {
		return sdkErr.Result()
	}

	// 4. get the total amount of self token and delegated token
	totalTokens := delegator.Tokens.Add(delegator.TotalDelegatedTokens)

	// 5. add shares to the vals this time
	shares, sdkErr := k.AddSharesToValidators(ctx, msg.DelAddr, vals, totalTokens)
	if sdkErr != nil {
		return sdkErr.Result()
	}

	// 6. update the delegator entity for this time
	delegator.ValidatorAddresses = getValsAddrs(vals)
	delegator.Shares = shares
	k.SetDelegator(ctx, delegator)

	ctx.EventManager().EmitEvent(buildEventForHandlerAddShares(delegator))
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// validateSharesAdding gives a quick validity of target validators before shares adding
func validateSharesAdding(vals types.Validators) sdk.Error {
	if len(vals) == 0 {
		return types.ErrNoAvailableValsToAddShares(types.DefaultCodespace)
	}

	if valAddr, ok := isDismissed(vals); ok {
		return types.ErrAddSharesToDismission(types.DefaultCodespace, valAddr.String())
	}

	return nil
}

// isDismissed tells whether validator with zero-msd is among the shares adding targets and returns the first dismissed
// validator address
func isDismissed(vals types.Validators) (sdk.ValAddress, bool) {
	valsLen := len(vals)
	for i := 0; i < valsLen; i++ {
		if vals[i].MinSelfDelegation.IsZero() {
			return vals[i].OperatorAddress, true
		}
	}

	return nil, false
}

// getValsAddrs gets validator addresses from a set of validator's entities
func getValsAddrs(vals types.Validators) []sdk.ValAddress {
	lenVals := len(vals)
	valAddrs := make([]sdk.ValAddress, lenVals)
	for i := 0; i < lenVals; i++ {
		valAddrs[i] = vals[i].OperatorAddress
	}
	return valAddrs
}

func buildEventForHandlerAddShares(delegator types.Delegator) sdk.Event {
	lenAttributes := len(delegator.ValidatorAddresses) + 2
	attributes := make([]sdk.Attribute, lenAttributes)
	attributes[0] = sdk.NewAttribute(types.AttributeKeyDelegator, delegator.DelegatorAddress.String())
	attributes[1] = sdk.NewAttribute(types.AttributeKeyShares, delegator.Shares.String())
	for i := 2; i < lenAttributes; i++ {
		attributes[i] = sdk.NewAttribute(types.AttributeKeyValidatorToAddShares, delegator.ValidatorAddresses[i-2].String())
	}

	return sdk.NewEvent(types.EventTypeAddShares, attributes...)
}

func handleMsgDeposit(ctx sdk.Context, msg types.MsgDeposit, k keeper.Keeper) sdk.Result {

	if msg.Amount.Denom != k.BondDenom(ctx) {
		return ErrBadDenom(k.Codespace()).Result()
	}

	err := k.Delegate(ctx, msg.DelegatorAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.DelegatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgWithdraw(ctx sdk.Context, msg types.MsgWithdraw, k keeper.Keeper) sdk.Result {
	if msg.Amount.Denom != k.BondDenom(ctx) {
		return ErrBadDenom(k.Codespace()).Result()
	}

	completionTime, err := k.Withdraw(ctx, msg.DelegatorAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})
	completionTimeBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	return sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().Events()}
}

func handleMsgDestroyValidator(ctx sdk.Context, msg types.MsgDestroyValidator, k keeper.Keeper) (result sdk.Result) {
	valAddr := sdk.ValAddress(msg.DelAddr)
	// 0.check to see if the validator which belongs to the delegator is existed
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return ErrNoValidatorFound(types.DefaultCodespace, valAddr.String()).Result()
	}

	completionTime, sdkErr := k.WithdrawMinSelfDelegation(ctx, msg.DelAddr, validator)
	if sdkErr != nil {
		return sdkErr.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, validator.GetMinSelfDelegation().String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})

	completionTimeBytes := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	return sdk.Result{Data: completionTimeBytes, Events: ctx.EventManager().Events()}

}
