package v0_9

import (
	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
	v08upgrade "github.com/okex/okchain/x/upgrade/legacy/v0_8"
	"github.com/okex/okchain/x/upgrade/types"
)

// Migrate converts the app state from an old version to a new one
func Migrate(oldGenesisState v08upgrade.GenesisState, oldgovParams v08gov.GovParams) GenesisState {
	params := types.UpgradeParams{
		AppUpgradeMaxDepositPeriod: oldgovParams.AppUpgradeMaxDepositPeriod,
		AppUpgradeMinDeposit:       oldgovParams.AppUpgradeMinDeposit,
		AppUpgradeVotingPeriod:     oldgovParams.AppUpgradeVotingPeriod,
	}

	return GenesisState{
		GenesisVersion: oldGenesisState.GenesisVersion,
		Params:         params,
	}
}
