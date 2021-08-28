package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/staking/exported"
)

func (k Keeper) CheckStatistics(ctx sdk.Context) {
	filter := getValFilter(k.monitoredValidators)
	tarValsTotalShares, totalValsTotalShares := sdk.ZeroDec(), sdk.ZeroDec()
	// sum shares
	k.IterateValidators(ctx, func(index int64, val exported.ValidatorI) (stop bool) {
		totalValsTotalShares = totalValsTotalShares.Add(val.GetDelegatorShares())
		if _, ok := filter[val.GetOperator().String()]; ok {
			tarValsTotalShares = tarValsTotalShares.Add(val.GetDelegatorShares())
		}

		return false
	})

	molecule, denominator := common.ConvertDecToFloat64(tarValsTotalShares), common.ConvertDecToFloat64(totalValsTotalShares)
	k.metric.AllValidatorsAndCandidateShare.Set(denominator)
	k.metric.ControlledValidatorsAndCandidateShare.Set(molecule)
	k.metric.ControlledValidatorsAndCandidateShareRatio.Set(molecule / denominator)
}

// build a filter
func getValFilter(valAddrs []string) map[string]struct{} {
	valLen := len(valAddrs)
	valFilter := make(map[string]struct{}, valLen)
	for i := 0; i < valLen; i++ {
		valFilter[valAddrs[i]] = struct{}{}
	}

	return valFilter
}
