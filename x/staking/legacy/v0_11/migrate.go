package v0_11

import "github.com/okex/okchain/x/staking/legacy/v0_10"

func Migrate(oldGenState v0_10.GenesisState) GenesisState {
	params := Params{
		UnbondingTime:      oldGenState.Params.UnbondingTime,
		MaxValidators:      oldGenState.Params.MaxValidators,
		Epoch:              oldGenState.Params.Epoch,
		MaxValsToAddShares: oldGenState.Params.MaxValsToVote,
		BondDenom:          oldGenState.Params.BondDenom,
		MinDelegation:      oldGenState.Params.MinDelegation,
		MinSelfDelegation:  DefaultMinSelfDelegation,
	}

	allShares := make([]SharesExported, len(oldGenState.Votes))
	for i, vote := range oldGenState.Votes {
		allShares[i] = SharesExported{
			DelAddress:       vote.VoterAddress,
			ValidatorAddress: vote.ValidatorAddress,
			Shares:           vote.Votes,
		}
	}

	return GenesisState{
		Params:               params,
		LastTotalPower:       oldGenState.LastTotalPower,
		LastValidatorPowers:  oldGenState.LastValidatorPowers,
		Validators:           oldGenState.Validators,
		Delegators:           oldGenState.Delegators,
		UnbondingDelegations: oldGenState.UnbondingDelegations,
		AllShares:            allShares,
		ProxyDelegatorKeys:   oldGenState.ProxyDelegatorKeys,
		Exported:             oldGenState.Exported,
	}
}
