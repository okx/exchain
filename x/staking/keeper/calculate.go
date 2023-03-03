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
	// iterate validators
	k.IterateValidators(ctx, func(index int64, val exported.ValidatorI) (stop bool) {
		totalValsTotalShares = totalValsTotalShares.Add(val.GetDelegatorShares())
		minSelf := val.GetMinSelfDelegation()
		if _, ok := valFilter[val.GetOperator().String()]; ok {
			tarValsTotalShares = tarValsTotalShares.Add(val.GetDelegatorShares())
			officialValidatorStakingOKT = officialValidatorStakingOKT.Add(minSelf)
		} else {
			communityValidatorStakingOKT = communityValidatorStakingOKT.Add(minSelf)
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

	officeRewards := k.GetOfficeRewards()
	officialStakingAndRewards := sdk.ConvertDecToFloat64(officialValidatorStakingOKT) + sdk.ConvertDecToFloat64(officialDelegatorStakingOKT) + officeRewards

	molecule, denominator := sdk.ConvertDecToFloat64(tarValsTotalShares), sdk.ConvertDecToFloat64(totalValsTotalShares)
	k.metric.AllValidatorsAndCandidateShare.Set(denominator)
	k.metric.ControlledValidatorsAndCandidateShare.Set(molecule)
	k.metric.ControlledValidatorsAndCandidateShareRatio.Set(molecule / denominator)

	k.metric.OfficialValidatorStakingOKT.Set(sdk.ConvertDecToFloat64(officialValidatorStakingOKT))
	k.metric.OfficialDelegatorStakingOKT.Set(sdk.ConvertDecToFloat64(officialDelegatorStakingOKT))
	k.metric.OfficialRewards.Set(officeRewards)
	k.metric.OfficialStakingAndRewards.Set(officialStakingAndRewards)
	k.metric.CommunityValidatorStakingOKT.Set(sdk.ConvertDecToFloat64(communityValidatorStakingOKT))
	k.metric.CommunityDelegatorStakingOKT.Set(sdk.ConvertDecToFloat64(communityDelegatorStakingOKT))

	totalSupplyOKT := k.supplyKeeper.GetSupplyByDenom(ctx, "okt")
	k.metric.TotalSupplyOKT.Set(sdk.ConvertDecToFloat64(totalSupplyOKT))
	logger.Error("Staking okt.", "official_validator", officialValidatorStakingOKT,
		"official_delegator", officialDelegatorStakingOKT,
		"official_rewards", officeRewards,
		"official_staking_and_rewards", officialStakingAndRewards,
		"community_validator", communityValidatorStakingOKT,
		"community_delegator", communityDelegatorStakingOKT,
		"total_supply_okt", totalSupplyOKT)
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
