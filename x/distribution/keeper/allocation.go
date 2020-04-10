package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingexported "github.com/cosmos/cosmos-sdk/x/staking/exported"
	"github.com/okex/okchain/x/distribution/types"
	"github.com/okex/okchain/x/staking/exported"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	valPortion  = sdk.NewDecWithPrec(25, 2)
	votePortion = sdk.NewDecWithPrec(75, 2)
)

// AllocateTokens allocates fees from fee_collector
//1. 25% rewards to validators, equally.
//2. 75% rewards to validators and candidators, by votes' wight
func (k Keeper) AllocateTokens(ctx sdk.Context, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, previousVotes []abci.VoteInfo) {
	logger := ctx.Logger().With("module", "distr")
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block

	// get the module account of feeCollector
	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	// get the total Coins from the module account-feeCollector
	feesCollected := feeCollector.GetCoins()

	if feesCollected.Empty() {
		logger.Debug("No fee to distributed.")
		return
	}
	logger.Debug("AllocateTokens", "TotalFee", feesCollected.String())

	if totalPreviousPower == 0 {
		// if the total previous power is zero, just return without allocate the fees util the power recovers
		logger.Error("totalPreviousPower is 0, skip this allocation of fees")
		return
	}

	// transfer collected fees to the distribution module account
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollected)
	if err != nil {
		panic(err)
	}

	preProposerVal := k.stakingKeeper.ValidatorByConsAddr(ctx, previousProposer)
	if preProposerVal == nil {
		// previous proposer can be unknown if say, the unbonding period is 1 block, so
		// e.g. a validator undelegates at block X, it's removed entirely by
		// block X+1's endblock, then X+2 we need to refer to the previous
		// proposer for X+1, but we've forgotten about them.
		logger.Error(fmt.Sprintf(
			"WARNING: Cannot find the entity of previous proposer validator %s.\n"+
				"This should happen only if the proposer unbonded completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.", previousProposer.String()))
	}

	fee1, fee2 := feesCollected.MulDecTruncate(valPortion), feesCollected.MulDecTruncate(votePortion)
	remaining := feesCollected.Sub(fee1.Add(fee2))
	remain2 := k.allocateByVal(ctx, fee1, previousVotes) //allocate rewards equally between validators
	remain1 := k.allocateByVotePower(ctx, fee2)          //allocate rewards by votes
	remaining = remaining.Add(remain1.Add(remain2))

	// if it remains some coins, allocate to proposer
	if !remaining.IsZero() {
		// if we can't find previous proposer validator from store or being jailed
		// then transfer the remaining to fee module account back
		if preProposerVal == nil || preProposerVal.IsJailed() {
			err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, remaining)
			if err != nil {
				panic(err)
			}
			logger.Debug("No Proposer to receive remaining", "Remainder to feeCollector", remaining)
			return
		}

		k.AllocateTokensToValidator(ctx, preProposerVal, remaining)
		logger.Debug("Send remaining to previous proposer",
			"previous proposer", preProposerVal.GetOperator().String(),
			"remaining coins", remaining,
		)
	}
}

func (k Keeper) allocateByVal(ctx sdk.Context, rewards sdk.DecCoins, previousVotes []abci.VoteInfo) sdk.DecCoins {
	logger := ctx.Logger().With("module", "distr")

	//count the total sum of the unJailed val
	validators := make([]stakingexported.ValidatorI, 0)
	for _, vote := range previousVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		if validator == nil {
			// previous validator can be unknown if say, the unbonding period is 1 block, so
			// e.g. a validator undelegates at block X, it's removed entirely by
			// block X+1's endblock, then X+2 we need to refer to the previous
			// validator for X+1, but we've forgotten about them.
			continue
		} else if validator.IsJailed() {
			logger.Debug(fmt.Sprintf("validator %s is jailed, not allowed to get reward by equal", validator.GetOperator()))
		} else {
			validators = append(validators, validator)
		}
	}

	//calculate the proportion of every valid validator
	powerFraction := sdk.NewDec(1).QuoTruncate(sdk.NewDec(int64(len(validators))))

	//beginning allocating rewards
	remaining := rewards
	for _, val := range validators {
		reward := rewards.MulDecTruncate(powerFraction)
		k.AllocateTokensToValidator(ctx, val, reward)
		logger.Debug("allocate by equal", val.GetOperator(), reward.String())
		remaining = remaining.Sub(reward)
	}
	return remaining
}

func (k Keeper) allocateByVotePower(ctx sdk.Context, rewards sdk.DecCoins) sdk.DecCoins {
	logger := ctx.Logger().With("module", "distr")

	//allocate tokens proportionally by votes to validators and candidators
	var validators []stakingexported.ValidatorI
	k.stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingexported.ValidatorI) (stop bool) {
		if validator != nil {
			if validator.IsJailed() {
				logger.Debug(fmt.Sprintf("validator %s is jailed, not allowed to get reward by votes weight",
					validator.GetOperator()))
			} else {
				validators = append(validators, validator)
			}
		}
		return false
	})

	//calculate total Votes-Weight
	var totalVotes = sdk.NewDec(0)
	sum := len(validators)
	for i := 0; i < sum; i++ {
		totalVotes = totalVotes.Add(validators[i].GetDelegatorShares())
	}

	//beginning allocating rewards
	remaining := rewards
	for _, val := range validators {
		powerFraction := val.GetDelegatorShares().QuoTruncate(totalVotes)
		reward := rewards.MulDecTruncate(powerFraction)
		k.AllocateTokensToValidator(ctx, val, reward)
		logger.Debug("allocate by votes", val.GetOperator(), reward.String())
		remaining = remaining.Sub(reward)
	}
	return remaining
}

// AllocateTokensToValidator allocate tokens to a particular validator, splitting according to commission
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val exported.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission
	// commission is always 1.0, so tokens.MulDec(val.GetCommission()) = tokens
	// only update current commission
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	currentCommission = currentCommission.Add(tokens)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), currentCommission)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
}
