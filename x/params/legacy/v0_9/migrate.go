package v0_9

import (
	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
)

// Migrate converts the state between the differet cm version
func Migrate(oldgovParams v08gov.GovParams) GenesisState {
	params := Params{
		MaxDepositPeriod: oldgovParams.ParamChangeMaxDepositPeriod,
		MinDeposit:       oldgovParams.ParamChangeMinDeposit,
		VotingPeriod:     oldgovParams.ParamChangeVotingPeriod,
		MaxBlockHeight:   uint64(oldgovParams.ParamChangeMaxBlockHeight),
	}

	return GenesisState{
		Params: params,
	}
}
