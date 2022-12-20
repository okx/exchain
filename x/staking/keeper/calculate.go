package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking/exported"
	"github.com/okex/exchain/x/staking/types"
)

func (k Keeper) CheckStatistics(ctx sdk.Context) {
	logger := k.Logger(ctx)
	valFilter := getFilterFromAddress(k.monitoredValidators)
	delegatorFilter := getFilterFromAddress(k.monitoredDelegators)
	tarValsTotalShares, totalValsTotalShares := sdk.ZeroDec(), sdk.ZeroDec()
	officialValidatorStakingOKT, officialDelegatorStakingOKT := sdk.ZeroDec(), sdk.ZeroDec()
	communityValidatorStakingOKT, communityDelegatorStakingOKT := sdk.ZeroDec(), sdk.ZeroDec()
	officialValidatorOutstandingOKT, communityValidatorOutstandingOKT := sdk.ZeroDec(), sdk.ZeroDec()
	totalStakingOKT := sdk.ZeroDec()
	// iterate validators
	k.IterateValidators(ctx, func(index int64, val exported.ValidatorI) (stop bool) {
		totalValsTotalShares = totalValsTotalShares.Add(val.GetDelegatorShares())
		minSelf := val.GetMinSelfDelegation()
		outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		if _, ok := valFilter[val.GetOperator().String()]; ok {
			tarValsTotalShares = tarValsTotalShares.Add(val.GetDelegatorShares())
			officialValidatorStakingOKT = officialValidatorStakingOKT.Add(minSelf)
			officialValidatorOutstandingOKT = officialValidatorOutstandingOKT.Add(outstanding)
		} else {
			communityValidatorStakingOKT = communityValidatorStakingOKT.Add(minSelf)
			communityValidatorOutstandingOKT = communityValidatorOutstandingOKT.Add(outstanding)
		}

		return false
	})

	// iterate delegators
	k.IterateDelegator(ctx, func(index int64, delegator types.Delegator) bool {
		if _, ok := delegatorFilter[delegator.GetDelegatorAddress().String()]; ok {
			officialDelegatorStakingOKT = officialDelegatorStakingOKT.Add(delegator.Tokens)
		} else {
			communityDelegatorStakingOKT = communityDelegatorStakingOKT.Add(delegator.Tokens)
		}
		return false
	})

	totalStakingOKT = totalStakingOKT.Add(officialValidatorStakingOKT).Add(officialDelegatorStakingOKT).
		Add(communityValidatorStakingOKT).Add(communityDelegatorStakingOKT).
		Add(officialValidatorOutstandingOKT).Add(communityValidatorOutstandingOKT)

	molecule, denominator := sdk.ConvertDecToFloat64(tarValsTotalShares), sdk.ConvertDecToFloat64(totalValsTotalShares)
	k.metric.AllValidatorsAndCandidateShare.Set(denominator)
	k.metric.ControlledValidatorsAndCandidateShare.Set(molecule)
	k.metric.ControlledValidatorsAndCandidateShareRatio.Set(molecule / denominator)

	k.metric.OfficialValidatorStakingOKT.Set(sdk.ConvertDecToFloat64(officialValidatorStakingOKT))
	k.metric.OfficialDelegatorStakingOKT.Set(sdk.ConvertDecToFloat64(officialDelegatorStakingOKT))
	k.metric.OfficialValidatorOutstandingOKT.Set(sdk.ConvertDecToFloat64(officialValidatorOutstandingOKT))
	k.metric.CommunityValidatorStakingOKT.Set(sdk.ConvertDecToFloat64(communityValidatorStakingOKT))
	k.metric.CommunityDelegatorStakingOKT.Set(sdk.ConvertDecToFloat64(communityDelegatorStakingOKT))
	k.metric.CommunityValidatorOutstandingOKT.Set(sdk.ConvertDecToFloat64(communityValidatorOutstandingOKT))
	k.metric.TotalStakingOKT.Set(sdk.ConvertDecToFloat64(totalStakingOKT))
	logger.Error("Staking okt.", "official_validator", officialValidatorStakingOKT,
		"official_delegator", officialDelegatorStakingOKT,
		"official_validator_outstanding", officialValidatorOutstandingOKT,
		"community_validator", communityValidatorStakingOKT,
		"community_delegator", communityDelegatorStakingOKT,
		"community_validator_outstanding", communityValidatorOutstandingOKT,
		"total_stakingOKT", totalStakingOKT)
}

// build a filter
func getFilterFromAddress(addrs []string) map[string]struct{} {
	valLen := len(addrs)
	valFilter := make(map[string]struct{}, valLen)
	for i := 0; i < valLen; i++ {
		valFilter[addrs[i]] = struct{}{}
	}

	return valFilter
}
