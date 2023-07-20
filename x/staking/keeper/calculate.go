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
	officialValidatorStakingOKT, officialDelegatorStakingOKT, officialDelegatorUnStakingOKT, officialDelegatorAmountOKT :=
		sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()
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

	for d, _ := range delegatorFilter {
		address, _ := sdk.AccAddressFromBech32(d)
		undelegation, found := k.GetUndelegating(ctx, address)
		if found {
			officialDelegatorUnStakingOKT = officialDelegatorUnStakingOKT.Add(undelegation.Quantity)
		}

		account := k.accountKeeper.GetAccount(ctx, address)
		if account != nil {
			coins := account.GetCoins()
			officialDelegatorAmountOKT = officialDelegatorAmountOKT.Add(coins.AmountOf(sdk.DefaultBondDenom))
		}
	}

	officeRewards := k.GetOfficeRewards()
	officialTotal := sdk.ConvertDecToFloat64(officialValidatorStakingOKT) + sdk.ConvertDecToFloat64(officialDelegatorStakingOKT) + officeRewards + sdk.ConvertDecToFloat64(officialDelegatorUnStakingOKT) + sdk.ConvertDecToFloat64(officialDelegatorAmountOKT)

	molecule, denominator := sdk.ConvertDecToFloat64(tarValsTotalShares), sdk.ConvertDecToFloat64(totalValsTotalShares)
	k.metric.AllValidatorsAndCandidateShare.Set(denominator)
	k.metric.ControlledValidatorsAndCandidateShare.Set(molecule)
	k.metric.ControlledValidatorsAndCandidateShareRatio.Set(molecule / denominator)

	k.metric.OfficialValidatorStakingOKT.Set(sdk.ConvertDecToFloat64(officialValidatorStakingOKT))
	k.metric.OfficialDelegatorStakingOKT.Set(sdk.ConvertDecToFloat64(officialDelegatorStakingOKT))
	k.metric.OfficialDelegatorUnStakingOKT.Set(sdk.ConvertDecToFloat64(officialDelegatorUnStakingOKT))
	k.metric.OfficialDelegatorAmountOKT.Set(sdk.ConvertDecToFloat64(officialDelegatorAmountOKT))
	k.metric.OfficialRewards.Set(officeRewards)
	k.metric.OfficialTotal.Set(officialTotal)
	k.metric.CommunityValidatorStakingOKT.Set(sdk.ConvertDecToFloat64(communityValidatorStakingOKT))
	k.metric.CommunityDelegatorStakingOKT.Set(sdk.ConvertDecToFloat64(communityDelegatorStakingOKT))

	totalSupplyOKT := k.supplyKeeper.GetSupplyByDenom(ctx, "okt")
	k.metric.TotalSupplyOKT.Set(sdk.ConvertDecToFloat64(totalSupplyOKT))
	logger.Error("Staking okt.", "official_validator", officialValidatorStakingOKT,
		"official_delegator", officialDelegatorStakingOKT,
		"official_un_delegator", officialDelegatorUnStakingOKT,
		"official_delegator_amount", officialDelegatorAmountOKT,
		"official_rewards", officeRewards,
		"official_total", officialTotal,
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
