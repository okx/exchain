package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/gov/types"
	"github.com/okex/okexchain/x/staking/exported"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Address             sdk.ValAddress // address of the validator operator
	BondedTokens        sdk.Int        // Power of a Validator
	DelegatorShares     sdk.Dec        // Total outstanding delegator shares
	DelegatorDeductions sdk.Dec        // Delegator deductions from validator's delegators voting independently
	Vote                types.VoteOption     // Vote of the validator
}

func newValidatorGovInfo(address sdk.ValAddress, bondedTokens sdk.Int, delegatorShares,
	delegatorDeductions sdk.Dec, vote types.VoteOption) validatorGovInfo {

	return validatorGovInfo{
		Address:             address,
		BondedTokens:        bondedTokens,
		DelegatorShares:     delegatorShares,
		DelegatorDeductions: delegatorDeductions,
		Vote:                vote,
	}
}

func tallyDelegatorVotes(
	ctx sdk.Context, keeper Keeper, currValidators map[string]validatorGovInfo, proposalID uint64,
	voteP *types.Vote, voterPower, totalVotedPower *sdk.Dec, results map[types.VoteOption]sdk.Dec,
) {
	// iterate over all the votes
	votesIterator := keeper.GetVotes(ctx, proposalID)
	if voteP != nil {
		votesIterator = append(votesIterator, *voteP)
	}
	for i := 0; i < len(votesIterator); i++ {
		vote := votesIterator[i]

		// if validator, just record it in the map
		// if delegator tally voting power
		valAddrStr := sdk.ValAddress(vote.Voter).String()
		if val, ok := currValidators[valAddrStr]; ok {
			val.Vote = vote.Option
			currValidators[valAddrStr] = val
		} else {
			// iterate over all delegations from voter, deduct from any delegated-to validators
			delegation := keeper.sk.Delegator(ctx, vote.Voter)
			if delegation == nil {
				continue
			}
			for _, val := range delegation.GetShareAddedValidatorAddresses() {
				valAddrStr := val.String()
				if valInfo, ok := currValidators[valAddrStr]; ok {
					valInfo.DelegatorDeductions = valInfo.DelegatorDeductions.Add(delegation.GetLastAddedShares())
					currValidators[valAddrStr] = valInfo

					votedPower := delegation.GetLastAddedShares()
					// calculate vote power of delegator for voterPowerRate
					if voteP != nil && vote.Voter.Equals(voteP.Voter) {
						voterPower.Add(votedPower)
					}
					results[vote.Option] = results[vote.Option].Add(votedPower)
					*totalVotedPower = totalVotedPower.Add(votedPower)
				}
			}
		}
	}
}

func tallyValidatorVotes(
	currValidators map[string]validatorGovInfo, voteP *types.Vote, voterPower,
	totalPower, totalVotedPower *sdk.Dec, results map[types.VoteOption]sdk.Dec,
) {
	// iterate over the validators again to tally their voting power
	for key, val := range currValidators {
		// calculate all vote power of current validators including delegated for voterPowerRate
		*totalPower = totalPower.Add(val.DelegatorShares)
		if val.Vote == types.OptionEmpty {
			continue
		}

		valValidVotedPower := val.DelegatorShares.Sub(val.DelegatorDeductions)
		if voteP != nil && sdk.ValAddress(voteP.Voter).String() == key {
			// calculate vote power of validator after deduction for voterPowerRate
			*voterPower = voterPower.Add(valValidVotedPower)
		}
		results[val.Vote] = results[val.Vote].Add(valValidVotedPower)
		*totalVotedPower = totalVotedPower.Add(valValidVotedPower)
	}
}

func preTally(
	ctx sdk.Context, keeper Keeper, proposal types.Proposal, voteP *types.Vote,
) (results map[types.VoteOption]sdk.Dec, totalVotedPower sdk.Dec, voterPowerRate sdk.Dec) {
	results = make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	totalVotedPower = sdk.ZeroDec()
	totalPower := sdk.ZeroDec()
	voterPower := sdk.ZeroDec()
	currValidators := make(map[string]validatorGovInfo)

	// fetch all the current validators except candidate, insert them into currValidators
	keeper.sk.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		currValidators[validator.GetOperator().String()] = newValidatorGovInfo(
			validator.GetOperator(),
			validator.GetBondedTokens(),
			validator.GetDelegatorShares(),
			sdk.ZeroDec(),
			types.OptionEmpty,
		)

		return false
	})

	tallyDelegatorVotes(ctx, keeper, currValidators, proposal.ProposalID,
		voteP, &voterPower, &totalVotedPower, results)

	tallyValidatorVotes(currValidators, voteP, &voterPower, &totalPower, &totalVotedPower, results)
	if totalPower.GT(sdk.ZeroDec()) {
		voterPowerRate = voterPower.Quo(totalPower)
	} else {
		voterPowerRate = sdk.ZeroDec()
	}

	return results, totalVotedPower, voterPowerRate
}

// tally and return status before voting period end time
func tallyStatusInVotePeriod(
	ctx sdk.Context, keeper Keeper, tallyResults types.TallyResult,
) (types.ProposalStatus, bool) {
	tallyParams := keeper.GetTallyParams(ctx)
	totalPower := tallyResults.TotalPower
	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if totalPower.IsZero() {
		return types.StatusRejected, false
	}
	// If no one votes (everyone abstains), proposal fails
	if totalPower.Sub(tallyResults.Abstain).Equal(sdk.ZeroDec()) {
		return types.StatusRejected, false
	}
	// If more than 1/3 of voters veto, proposal fails
	if tallyResults.NoWithVeto.Quo(totalPower).GT(tallyParams.Veto) {
		return types.StatusRejected, true
	}
	// If more than or equal to 1/2 of non-abstain vote not Yes, proposal fails
	if tallyResults.NoWithVeto.Add(tallyResults.No).Quo(totalPower.Sub(tallyResults.Abstain)).
		GTE(tallyParams.Threshold) {
		return types.StatusRejected, false
	}
	// If more than 2/3 of totalPower vote Yes, proposal passes
	if tallyResults.Yes.Quo(totalPower).GT(tallyParams.YesInVotePeriod) {
		return types.StatusPassed, false
	}

	return types.StatusVotingPeriod, false
}

// tally and return status expire voting period end time
func tallyStatusExpireVotePeriod(
	ctx sdk.Context, keeper Keeper, tallyResults types.TallyResult,
) (types.ProposalStatus, bool) {
	tallyParams := keeper.GetTallyParams(ctx)
	totalVoted := tallyResults.TotalVotedPower
	totalPower := tallyResults.TotalPower
	// TODO: Upgrade the spec to cover all of these cases & remove pseudo code.
	// If there is no staked coins, the proposal fails
	if totalPower.IsZero() {
		return types.StatusRejected, false
	}
	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVoted.Quo(totalPower)
	if percentVoting.LT(tallyParams.Quorum) {
		return types.StatusRejected, true
	}
	// If no one votes (everyone abstains), proposal fails
	if totalVoted.Sub(tallyResults.Abstain).Equal(sdk.ZeroDec()) {
		return types.StatusRejected, false
	}
	// If more than 1/3 of voters veto, proposal fails
	if tallyResults.NoWithVeto.Quo(totalVoted).GT(tallyParams.Veto) {
		return types.StatusRejected, true
	}
	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if tallyResults.Yes.Quo(totalVoted.Sub(tallyResults.Abstain)).GT(tallyParams.Threshold) {
		return types.StatusPassed, false
	}
	// If more than 1/2 of non-abstaining voters vote No, proposal fails

	return types.StatusRejected, false
}

// Tally counts the votes for proposal
func Tally(ctx sdk.Context, keeper Keeper, proposal types.Proposal, isExpireVoteEndTime bool,
) (types.ProposalStatus, bool, types.TallyResult) {
	results, totalVotedPower, _ := preTally(ctx, keeper, proposal, nil)
	tallyResults := types.NewTallyResultFromMap(results)
	tallyResults.TotalPower = keeper.totalPower(ctx)
	tallyResults.TotalVotedPower = totalVotedPower

	if isExpireVoteEndTime {
		status, distribute := tallyStatusExpireVotePeriod(ctx, keeper, tallyResults)
		return status, distribute, tallyResults
	}
	status, distribute := tallyStatusInVotePeriod(ctx, keeper, tallyResults)
	return status, distribute, tallyResults
}
