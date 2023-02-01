package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

type CM45DepositParams struct {
	MinDeposit       sdk.SysCoins `json:"min_deposit,omitempty" yaml:"min_deposit,omitempty"`
	MaxDepositPeriod string       `json:"max_deposit_period,omitempty" yaml:"max_deposit_period,omitempty"`
}

type CM45VotingParams struct {
	VotingPeriod string `json:"voting_period,omitempty" yaml:"voting_period,omitempty"` //  Length of the voting period.
}

type CM45TallyParams struct {
	Quorum          sdk.Dec `json:"quorum,omitempty" yaml:"quorum,omitempty"`                         //  Minimum percentage of total stake needed to vote for a result to be considered valid
	Threshold       sdk.Dec `json:"threshold,omitempty" yaml:"threshold,omitempty"`                   //  Minimum proportion of Yes votes for proposal to pass. Initial value: 0.5
	Veto            sdk.Dec `json:"veto_threshold,omitempty" yaml:"veto_threshold,omitempty"`         //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	YesInVotePeriod sdk.Dec `json:"yes_in_vote_period,omitempty" yaml:"yes_in_vote_period,omitempty"` //
}

type CM45Params struct {
	VotingParams  CM45VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams   CM45TallyParams   `json:"tally_params" yaml:"tally_params"`
	DepositParams CM45DepositParams `json:"deposit_params" yaml:"deposit_parmas"`
}

func NewCM45Params(vp CM45VotingParams, tp CM45TallyParams, dp CM45DepositParams) CM45Params {
	return CM45Params{
		VotingParams:  vp,
		TallyParams:   tp,
		DepositParams: dp,
	}
}
