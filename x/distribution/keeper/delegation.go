package keeper

import (
	"fmt"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/types"
	stakingexported "github.com/okx/okbchain/x/staking/exported"
)

// initialize starting info for a new delegation
func (k Keeper) initializeDelegation(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) {
	if !k.CheckDistributionProposalValid(ctx) {
		return
	}

	logger := k.Logger(ctx)
	// period has already been incremented - we want to store the period ended by this delegation action
	previousPeriod := k.GetValidatorCurrentRewards(ctx, val).Period - 1

	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, val, previousPeriod)
	delegation := k.stakingKeeper.Delegator(ctx, del)

	k.SetDelegatorStartingInfo(ctx, val, del, types.NewDelegatorStartingInfo(previousPeriod, delegation.GetLastAddedShares(), uint64(ctx.BlockHeight())))
	logger.Debug("initializeDelegation", "ValAddress", val, "Delegator", del, "Shares", delegation.GetLastAddedShares())
}

// calculate the rewards accrued by a delegation between two periods
func (k Keeper) calculateDelegationRewardsBetween(ctx sdk.Context, val stakingexported.ValidatorI,
	startingPeriod, endingPeriod uint64, stake sdk.Dec) (rewards sdk.DecCoins) {
	logger := k.Logger(ctx)
	// sanity check
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	// sanity check
	if stake.IsNegative() {
		panic("stake should not be negative")
	}

	// return staking * (ending - starting)
	starting := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), startingPeriod)
	ending := k.GetValidatorHistoricalRewards(ctx, val.GetOperator(), endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	rewards = difference.MulDecTruncate(stake)
	logger.Debug("calculateDelegationRewardsBetween", "Validator", val.GetOperator(),
		"Start", starting.CumulativeRewardRatio, "End", ending.CumulativeRewardRatio, "Stake", stake,
		"Difference", difference, "Rewards", rewards)
	return
}

// calculate the total rewards accrued by a delegation
func (k Keeper) calculateDelegationRewards(ctx sdk.Context, val stakingexported.ValidatorI, delAddr sdk.AccAddress, endingPeriod uint64) (rewards sdk.DecCoins) {
	logger := k.Logger(ctx)
	del := k.stakingKeeper.Delegator(ctx, delAddr)

	// fetch starting info for delegation
	startingInfo := k.GetDelegatorStartingInfo(ctx, val.GetOperator(), del.GetDelegatorAddress())

	if startingInfo.Height == uint64(ctx.BlockHeight()) {
		// started this height, no rewards yet
		logger.Debug(fmt.Sprintf("calculateDelegationRewards end, error, val:%s, del:%s, height:%d",
			val.GetOperator().String(), delAddr.String(), startingInfo.Height))
		return
	}

	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake
	if stake.GT(del.GetLastAddedShares()) {
		panic(fmt.Sprintf("calculated final stake for delegator %s greater than current stake"+
			"\n\tfinal stake:\t%s"+
			"\n\tcurrent stake:\t%s",
			del.GetDelegatorAddress(), stake, del.GetLastAddedShares()))
	}

	// calculate rewards for final period
	rewards = rewards.Add(k.calculateDelegationRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake)...)

	logger.Debug("calculateDelegationRewards", "Validator", val.GetOperator(),
		"Delegator", delAddr, "Start", startingPeriod, "End", endingPeriod, "Stake", stake, "Rewards", rewards)

	return rewards
}

//withdraw rewards according to the specified validator by delegator
func (k Keeper) withdrawDelegationRewards(ctx sdk.Context, val stakingexported.ValidatorI, delAddress sdk.AccAddress) (sdk.Coins, error) {
	if !k.CheckDistributionProposalValid(ctx) {
		return nil, types.ErrCodeNotSupportWithdrawDelegationRewards()
	}

	logger := k.Logger(ctx)

	// check existence of delegator starting info
	if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), delAddress) {
		del := k.stakingKeeper.Delegator(ctx, delAddress)
		if del.GetLastAddedShares().IsZero() {
			return nil, types.ErrCodeZeroDelegationShares()
		}
		k.initExistedDelegationStartInfo(ctx, val, del)
	}

	// end current period and calculate rewards
	endingPeriod := k.incrementValidatorPeriod(ctx, val)
	rewardsRaw := k.calculateDelegationRewards(ctx, val, delAddress, endingPeriod)
	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())

	// defensive edge case may happen on the very final digits
	// of the decCoins due to operation order of the distribution mechanism.
	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logger.Info(fmt.Sprintf("missing rewards rounding error, delegator %v"+
			"withdrawing rewards from validator %v, should have received %v, got %v",
			val.GetOperator(), delAddress, rewardsRaw, rewards))
	}

	// truncate coins, return remainder to community pool
	coins, remainder := rewards.TruncateWithPrec(k.GetRewardTruncatePrecision(ctx))

	// add coins to user account
	if !coins.IsZero() {
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, delAddress)
		err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins)
		logger.Debug("SendCoinsFromModuleToAccount", "From", types.ModuleName,
			"To", withdrawAddr, "Coins", coins)
		if err != nil {
			return nil, err
		}
	}

	// update the outstanding rewards and the community pool only if the
	// transaction was successful
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding.Sub(rewards))
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
	k.SetFeePool(ctx, feePool)

	// decrement reference count of starting period
	startingInfo := k.GetDelegatorStartingInfo(ctx, val.GetOperator(), delAddress)
	startingPeriod := startingInfo.PreviousPeriod
	k.decrementReferenceCount(ctx, val.GetOperator(), startingPeriod)

	// remove delegator starting info
	k.DeleteDelegatorStartingInfo(ctx, val.GetOperator(), delAddress)

	logger.Debug("withdrawDelegationRewards", "Validator", val.GetOperator(), "Delegator", delAddress,
		"Stake", startingInfo.Stake, "StartingPeriod", startingPeriod, "EndingPeriod", endingPeriod,
		"RewardsRaw", rewardsRaw, "Rewards", rewards, "Coins", coins, "Remainder", remainder)
	return coins, nil
}

//initExistedDelegationStartInfo If the delegator existed but no start info, it add shares before distribution proposal, and need to set a new start info
func (k Keeper) initExistedDelegationStartInfo(ctx sdk.Context, val stakingexported.ValidatorI, del stakingexported.DelegatorI) {
	if !k.CheckDistributionProposalValid(ctx) {
		return
	}

	logger := k.Logger(ctx)
	//set previous validator period 0
	previousPeriod := uint64(0)
	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, val.GetOperator(), previousPeriod)

	k.SetDelegatorStartingInfo(ctx, val.GetOperator(), del.GetDelegatorAddress(),
		types.NewDelegatorStartingInfo(previousPeriod, del.GetLastAddedShares(), 0))

	logger.Debug("initExistedDelegationStartInfo", "Validator", val.GetOperator(),
		"Delegator", del.GetDelegatorAddress(), "Shares", del.GetLastAddedShares())
	return
}
