package v0_36

import (
	v034distr "github.com/okex/okchain/x/distribution/legacy/v0_34"
)

// Migrate accepts exported genesis state from v0.34 and migrates it to v0.36
// genesis state. All entries are identical except for validator slashing events
// which now include the period.
func Migrate(oldGenState v034distr.GenesisState) GenesisState {
	// migrate slash events which now have the period included
	return NewGenesisState(oldGenState.WithdrawAddrEnabled, oldGenState.DelegatorWithdrawInfos,
		oldGenState.PreviousProposer, oldGenState.ValidatorAccumulatedCommissions)
}
