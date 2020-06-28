package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/types"
	"github.com/okex/okchain/x/staking/exported"
	stakingexported "github.com/okex/okchain/x/staking/exported"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	valPortion  = sdk.NewDecWithPrec(25, 2)
	votePortion = sdk.NewDecWithPrec(75, 2)
)

// AllocateTokens allocates fees from fee_collector
//1. 25% rewards to validators, equally.
//2. 75% rewards to validators and candidates, by shares' weight
func (k Keeper) AllocateTokens(ctx sdk.Context, totalPreviousPower int64,
	previousProposer sdk.ConsAddress, previousVotes []abci.VoteInfo) {
	logger := k.Logger(ctx)
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.supplyKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollected := feeCollector.GetCoins()

	if feesCollected.Empty() {
		logger.Debug("No fee to distributed.")
		return
	}
	logger.Debug("AllocateTokens", "TotalFee", feesCollected.String())

	// transfer collected fees to the distribution module account
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollected)
	if err != nil {
		panic(err)
	}

	feePool := k.GetFeePool(ctx)
	if totalPreviousPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected)
		k.SetFeePool(ctx, feePool)
		logger.Debug("totalPreviousPower is zero, send fees to community pool", "fees", feesCollected)
		return
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

	feesToVals := feesCollected.MulDecTruncate(sdk.OneDec().Sub(k.GetCommunityTax(ctx)))
	feeByEqual, feeByVote := feesToVals.MulDecTruncate(valPortion), feesToVals.MulDecTruncate(votePortion)
	feesToCommunity := feesCollected.Sub(feeByEqual.Add(feeByVote))
	remainByEqual := k.allocateByEqual(ctx, feeByEqual, previousVotes) //allocate rewards equally between validators
	remainByShare := k.allocateByShares(ctx, feeByVote)                //allocate rewards by shares
	feesToCommunity = feesToCommunity.Add(remainByEqual.Add(remainByShare))

	// allocate community funding
	if !feesToCommunity.IsZero() {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesToCommunity)
		k.SetFeePool(ctx, feePool)
		logger.Debug("Send fees to community pool", "community_pool", feesToCommunity)
	}
}

func (k Keeper) allocateByEqual(ctx sdk.Context, rewards sdk.DecCoins, previousVotes []abci.VoteInfo) sdk.DecCoins {
	logger := k.Logger(ctx)

	//count the total sum of the unJailed val
	var validators []stakingexported.ValidatorI
	for _, vote := range previousVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
		if validator != nil {
			if validator.IsJailed() {
				logger.Debug(fmt.Sprintf("validator %s is jailed, not allowed to get reward by equal", validator.GetOperator()))
			} else {
				validators = append(validators, validator)
			}
		}
	}

	//calculate the proportion of every valid validator
	powerFraction := sdk.NewDec(1).QuoTruncate(sdk.NewDec(int64(len(validators))))

	//beginning allocating rewards equally
	remaining := rewards
	reward := rewards.MulDecTruncate(powerFraction)
	for _, val := range validators {
		k.AllocateTokensToValidator(ctx, val, reward)
		logger.Debug("allocate by equal", val.GetOperator(), reward.String())
		remaining = remaining.Sub(reward)
	}
	return remaining
}

func (k Keeper) allocateByShares(ctx sdk.Context, rewards sdk.DecCoins) sdk.DecCoins {
	logger := k.Logger(ctx)

	//allocate tokens proportionally by votes to validators and candidates
	var validators []stakingexported.ValidatorI
	k.stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingexported.ValidatorI) (stop bool) {
		if validator != nil {
			if validator.IsJailed() {
				logger.Debug(fmt.Sprintf("validator %s is jailed, not allowed to get reward by shares weight",
					validator.GetOperator()))
			} else {
				validators = append(validators, validator)
			}
		}
		return false
	})

	//calculate total Shares-Weight
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
		logger.Debug("allocate by shares", val.GetOperator(), reward.String())
		remaining = remaining.Sub(reward)
	}
	return remaining
}

// AllocateTokensToValidator allocate tokens to a particular validator, splitting according to commissions
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val exported.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commissions
	// commissions is always 1.0, so tokens.MulDec(val.GetCommission()) = tokens
	// only update current commissions
	commission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	commission = commission.Add(tokens)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), commission)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
}
