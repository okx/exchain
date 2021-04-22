package v0_17

import "github.com/okex/exchain/x/staking/legacy/v0_11"

const (
	ModuleName = "staking"
)

// Migrate removes BondDenom
func Migrate(oldGenState v0_11.GenesisState) GenesisState {
	params := Params{
		UnbondingTime:      oldGenState.Params.UnbondingTime,
		MaxValidators:      oldGenState.Params.MaxValidators,
		Epoch:              oldGenState.Params.Epoch,
		MaxValsToAddShares: oldGenState.Params.MaxValsToAddShares,
		MinDelegation:      oldGenState.Params.MinDelegation,
		MinSelfDelegation:  DefaultMinSelfDelegation,
	}

	return GenesisState{
		Params:               params,
		LastTotalPower:       oldGenState.LastTotalPower,
		LastValidatorPowers:  oldGenState.LastValidatorPowers,
		Validators:           oldGenState.Validators,
		Delegators:           oldGenState.Delegators,
		UnbondingDelegations: oldGenState.UnbondingDelegations,
		AllShares:            oldGenState.AllShares,
		ProxyDelegatorKeys:   oldGenState.ProxyDelegatorKeys,
		Exported:             oldGenState.Exported,
	}
}
