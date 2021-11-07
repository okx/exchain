package keeper

import (
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking/types"
)

// WithdrawMinSelfDelegation withdraws the msd from validator
func (k Keeper) WithdrawMinSelfDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator types.Validator,
) (completionTime time.Time, err error) {
	// 0.check the msd on validator
	if validator.MinSelfDelegation.IsZero() {
		return completionTime, types.ErrNoMinSelfDelegation(validator.OperatorAddress.String())
	}

	// 1.check the remained shares on the validator
	remainShares := validator.GetDelegatorShares().Sub(k.getSharesFromDefaultMinSelfDelegation())
	if remainShares.LT(sdk.ZeroDec()) {
		return completionTime, types.ErrMoreMinSelfDelegation(validator.OperatorAddress.String())
	}

	// 2.unbond msd
	k.bondedTokensToNotBonded(ctx, sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, validator.MinSelfDelegation))
	completionTime = ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	minSelfUndelegation := types.NewUndelegationInfo(delAddr, validator.MinSelfDelegation, completionTime)
	k.SetUndelegating(ctx, minSelfUndelegation)
	k.SetAddrByTimeKeyWithNilValue(ctx, minSelfUndelegation.CompletionTime, minSelfUndelegation.DelegatorAddress)

	// 3.clear the msd
	validator.MinSelfDelegation = sdk.ZeroDec()

	// 4.jail the validator
	validator.Jailed = true

	// 5.call the hooks of slashing module
	k.AfterValidatorDestroyed(ctx, validator.ConsAddress(), validator.OperatorAddress)

	// 6.change status of validator
	switch validator.Status {
	case sdk.Bonded:
		// set the validator info to enforce the update of validator-set
		k.AppendAbandonedValidatorAddrs(ctx, validator.ConsAddress())
	case sdk.Unbonding:
	case sdk.Unbonded:
		// if there is no shares on the validator, remove it
		if remainShares.IsZero() && validator.GetMinSelfDelegation().IsZero() {
			k.RemoveValidator(ctx, validator.OperatorAddress)
			return
		}
	}
	// kick out the val from the vals-set
	k.DeleteValidatorByPowerIndex(ctx, validator)
	// ATTENTION:update DelegatorShares must go after DeleteValidatorByPowerIndex
	validator.DelegatorShares = remainShares
	k.SetValidator(ctx, validator)

	return
}

// AddSharesAsMinSelfDelegation adds shares of equal value of default msd (0.001okt) to validator itself during the creation
func (k Keeper) AddSharesAsMinSelfDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator *types.Validator,
	defaultMSDToken sdk.SysCoin) (err error) {
	// 0. transfer account's okt (0.001okt as default) into bondPool
	coins := sdk.SysCoins{defaultMSDToken}
	err = k.supplyKeeper.DelegateCoinsFromAccountToModule(ctx, delAddr, types.BondedPoolName, coins)
	if err != nil {
		return err
	}

	// 1. add shares for default msd to validator itself
	k.addSharesAsDefaultMinSelfDelegation(ctx, validator)

	return nil
}

func (k Keeper) addSharesAsDefaultMinSelfDelegation(ctx sdk.Context, pValidator *types.Validator) {
	k.DeleteValidatorByPowerIndex(ctx, *pValidator)
	//TODO: current rule: any msd -> 1 shares
	shares := k.getSharesFromDefaultMinSelfDelegation()
	pValidator.DelegatorShares = pValidator.GetDelegatorShares().Add(shares)
	k.SetValidator(ctx, *pValidator)
	k.SetValidatorByPowerIndex(ctx, *pValidator)
}

// RULES: any msd -> 1 shares
func (k Keeper) getSharesFromDefaultMinSelfDelegation() sdk.Dec {
	return sdk.OneDec()
}
