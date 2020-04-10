// nolint
package v0_9

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"

	"github.com/okex/okchain/x/common/proto"
	v08gov "github.com/okex/okchain/x/gov/legacy/v0_8"
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

	AppUpgradeProposal struct {
		Title              string                   `json:"title" yaml:"titile"`
		Description        string                   `json:"description" yaml:"description"`
		ProtocolDefinition proto.ProtocolDefinition `json:"protocol_definition" yaml:"protocol_definition"`
	}

	Deposit struct {
		ProposalID uint64         `json:"proposal_id" yaml:"proposal_id"`
		Depositor  sdk.AccAddress `json:"depositor" yaml:"depositor"`
		Amount     sdk.DecCoins   `json:"amount" yaml:"amount"`
		DepositID  uint64         `json:"deposit_id" yaml:"deposit_id"`
	}

	Deposits []Deposit

	// Vote
	Vote struct {
		ProposalID uint64            `json:"proposal_id" yaml:"proposal_id"`
		Voter      sdk.AccAddress    `json:"voter" yaml:"voter"`
		Option     v08gov.VoteOption `json:"option" yaml:"option"`
		VoteID     uint64            `json:"vote_id" yaml:"vote_id"`
	}

	Votes []Vote

	GenesisState struct {
		StartingProposalID uint64                       `json:"starting_proposal_id" yaml:"starting_proposal_id"`
		Deposits           Deposits                     `json:"deposits" yaml:"deposits"`
		Votes              Votes                        `json:"votes" yaml:"votes"`
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
	cdc.RegisterConcrete(ParameterChangeProposal{}, "params/ParameterChangeProposal", nil)
	cdc.RegisterConcrete(AppUpgradeProposal{}, "okchain/AppUpgradeProposal", nil)
}

// GetTitle returns the title of a dex list proposal.
func (apu AppUpgradeProposal) GetTitle() string { return apu.Title }

// GetDescription returns the description of a dex list proposal.
func (apu AppUpgradeProposal) GetDescription() string { return apu.Description }

// ProposalRoute returns the routing key of a dex list proposal.
func (apu AppUpgradeProposal) ProposalRoute() string { return "" }

// ProposalType returns the type of a dex list proposal.
func (apu AppUpgradeProposal) ProposalType() string { return "" }

// ValidateBasic validates the parameter change proposal
func (apu AppUpgradeProposal) ValidateBasic() sdk.Error {
	return nil
}

func (apu AppUpgradeProposal) String() string {
	return fmt.Sprintf(`Proposal:
  title:              %s
  proposalType:               %s
  Version:            %d
  Software:           %s
  Switch Height:      %d
  Threshold:          %s`, apu.Title, apu.ProposalType(), apu.ProtocolDefinition.Version,
		apu.ProtocolDefinition.Software, apu.ProtocolDefinition.Height, apu.ProtocolDefinition.Threshold.String())
}
