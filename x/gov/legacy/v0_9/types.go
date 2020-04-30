// nolint
package v0_9

import (
	"time"

	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
	"github.com/okex/okchain/x/gov/types"
	"github.com/okex/okchain/x/params"
	upgradeTypes "github.com/okex/okchain/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
)

const (
	ModuleName = "gov"
)

type (
	Proposal struct {
		// Proposal content interface
		sdkGovTypes.Content `json:"content" yaml:"content"`

		ProposalID uint64 `json:"id" yaml:"id"`
		// Status of the Proposal {Pending, Active, Passed, Rejected}
		Status v08gov.ProposalStatus `json:"proposal_status" yaml:"proposal_status"`
		// Result of Tallys
		FinalTallyResult v08gov.TallyResult `json:"final_tally_result" yaml:"final_tally_result"`

		// Time of the block where TxGovSubmitProposal was included
		SubmitTime time.Time `json:"submit_time" yaml:"submit_time"`
		// Time that the Proposal would expire if deposit amount isn't met
		DepositEndTime time.Time `json:"deposit_end_time" yaml:"deposit_end_time"`
		// Current deposit on this proposal. Initial value is set at InitialDeposit
		TotalDeposit sdk.DecCoins `json:"total_deposit" yaml:"total_deposit"`

		// Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
		VotingStartTime time.Time `json:"voting_start_time" yaml:"voting_start_time"`
		// Time that the VotingPeriod for this proposal will end and votes will be tallied
		VotingEndTime time.Time `json:"voting_end_time" yaml:"voting_end_time"`
	}

	ParameterChangeProposal struct {
		sdkparams.ParameterChangeProposal
		Height uint64 `json:"height" yaml:"height"`
	}

	GenesisState struct {
		StartingProposalID uint64                       `json:"starting_proposal_id" yaml:"starting_proposal_id"`
		Deposits           types.Deposits               `json:"deposits" yaml:"deposits"`
		Votes              types.Votes                  `json:"votes" yaml:"votes"`
		Proposals          []Proposal                   `json:"proposals" yaml:"proposals"`
		DepositParams      sdkGovTypes.DepositParams    `json:"deposit_params" yaml:"deposit_params"`
		VotingParams       sdkGovTypes.VotingParams     `json:"voting_params" yaml:"voting_params"`
		TallyParams        sdkGovTypes.TallyParams      `json:"tally_params" yaml:"tally_params"`
		TendermintParams   sdkGovTypes.TendermintParams `json:"tendermint_params" yaml:"tendermint_params"`
	}
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*sdkGovTypes.Content)(nil), nil)

	cdc.RegisterConcrete(sdkGovTypes.TextProposal{}, "cosmos-sdk/TextProposal", nil)
	cdc.RegisterConcrete(params.ParameterChangeProposal{}, "params/ParameterChangeProposal", nil)
	cdc.RegisterConcrete(upgradeTypes.AppUpgradeProposal{}, "okchain/AppUpgradeProposal", nil)
}
