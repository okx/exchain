package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/types"
)

// UndelegateMinSelfDelegation unbonds the msd from validator
func (k Keeper) UndelegateMinSelfDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator types.Validator,
) (completionTime time.Time, err sdk.Error) {
	// 0.check the msd on validator
	if validator.MinSelfDelegation.IsZero() {
		return completionTime, types.ErrNoMinSelfDelegation(types.DefaultCodespace, validator.OperatorAddress.String())
	}

	// 1.check the remained vote from validator
	remainVotes := validator.GetDelegatorShares().Sub(k.getVotesFromFixedMinSelfDelegation())
	if remainVotes.LT(sdk.ZeroDec()) {
		return completionTime, types.ErrMoreMinSelfDelegation(types.DefaultCodespace, validator.OperatorAddress.String())
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
		// if there is no vote on the validator, remove it
		if remainVotes.IsZero() && validator.GetMinSelfDelegation().IsZero() {
			k.RemoveValidator(ctx, validator.OperatorAddress)
			return
		}
	}
	// kick out the val from the vals-set
	k.DeleteValidatorByPowerIndex(ctx, validator)
	// ATTENTION:update DelegatorShares must go after DeleteValidatorByPowerIndex
	validator.DelegatorShares = remainVotes
	k.SetValidator(ctx, validator)

	return
}

// VoteMinSelfDelegation votes fixed msd (0.001okt) to validator itself during the creation
func (k Keeper) VoteMinSelfDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator *types.Validator,
	fixedMSDToken sdk.DecCoin) (err sdk.Error) {
	// 0. transfer account's okt (0.001okt) into bondPool
	coins := fixedMSDToken.ToCoins()
	err = k.supplyKeeper.DelegateCoinsFromAccountToModule(ctx, delAddr, types.BondedPoolName, coins)
	if err != nil {
		return err
	}

	// 1. vote to validator itself
	k.voteMinSelfDelegation(ctx, validator)

	return nil
}

func (k Keeper) voteMinSelfDelegation(ctx sdk.Context, pValidator *types.Validator) {
	k.DeleteValidatorByPowerIndex(ctx, *pValidator)
	//TODO: current rule: any msd -> 1 votes
	voteDec := k.getVotesFromFixedMinSelfDelegation()
	pValidator.DelegatorShares = pValidator.GetDelegatorShares().Add(voteDec)
	k.SetValidator(ctx, *pValidator)
	k.SetValidatorByPowerIndex(ctx, *pValidator)
}

// RULES: any msd -> 1 votes
func (k Keeper) getVotesFromFixedMinSelfDelegation() sdk.Dec {
	return sdk.OneDec()
}