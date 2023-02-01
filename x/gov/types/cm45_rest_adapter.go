package types

import (
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

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

// CM45Proposal is constructed to be compatible with the REST API of cosmos v0.45.1
type CM45Proposal struct {
	Content `json:"content" yaml:"content"` // Proposal content interface

	ProposalID       uint64      `json:"proposal_id" yaml:"proposal_id"`               //  ID of the proposal
	Status           string      `json:"status" yaml:"status"`                         // Status of the Proposal {Pending, Active, Passed, Rejected}
	FinalTallyResult TallyResult `json:"final_tally_result" yaml:"final_tally_result"` // Result of Tallys

	SubmitTime     time.Time    `json:"submit_time" yaml:"submit_time"`           // Time of the block where TxGovSubmitProposal was included
	DepositEndTime time.Time    `json:"deposit_end_time" yaml:"deposit_end_time"` // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit   sdk.SysCoins `json:"total_deposit" yaml:"total_deposit"`       // Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime time.Time `json:"voting_start_time" yaml:"voting_start_time"` // Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time `json:"voting_end_time" yaml:"voting_end_time"`     // Time that the VotingPeriod for this proposal will end and votes will be tallied
}

// WrapProposalForCosmosAPI is for compatibility with the standard cosmos REST API
func WrapProposalForCosmosAPI(proposal Proposal, content Content) Proposal {
	return Proposal{
		Content:          content,
		ProposalID:       proposal.ProposalID,
		Status:           proposal.Status,
		FinalTallyResult: proposal.FinalTallyResult,
		SubmitTime:       proposal.SubmitTime,
		DepositEndTime:   proposal.DepositEndTime,
		TotalDeposit:     proposal.TotalDeposit,
		VotingStartTime:  proposal.VotingStartTime,
		VotingEndTime:    proposal.VotingEndTime,
	}
}

func (p Proposal) ToCM45Proposal() CM45Proposal {
	cm45p := CM45Proposal{
		Content:          p.Content,
		ProposalID:       p.ProposalID,
		Status:           p.Status.ToCM45Status(),
		FinalTallyResult: p.FinalTallyResult,
		SubmitTime:       p.SubmitTime,
		DepositEndTime:   p.DepositEndTime,
		TotalDeposit:     p.TotalDeposit,
		VotingStartTime:  p.VotingStartTime,
		VotingEndTime:    p.VotingEndTime,
	}
	return cm45p
}

func (status ProposalStatus) ToCM45Status() string {
	switch status {
	case StatusDepositPeriod:
		return "PROPOSAL_STATUS_DEPOSIT_PERIOD"

	case StatusVotingPeriod:
		return "PROPOSAL_STATUS_VOTING_PERIOD"

	case StatusPassed:
		return "PROPOSAL_STATUS_PASSED"

	case StatusRejected:
		return "PROPOSAL_STATUS_REJECTED"

	case StatusFailed:
		return "PROPOSAL_STATUS_FAILED"

	default:
		return ""
	}
}

type WrappedProposal struct {
	P CM45Proposal `json:"proposal" yaml:"result"`
}

func NewWrappedProposal(p CM45Proposal) WrappedProposal {
	return WrappedProposal{
		P: p,
	}
}

type WrappedProposals struct {
	Ps []CM45Proposal `json:"proposals" yaml:"result"`
}

func NewWrappedProposals(ps []CM45Proposal) WrappedProposals {
	return WrappedProposals{
		Ps: ps,
	}
}

type WrappedTallyResult struct {
	TR TallyResult `json:"tally"`
}

func NewWrappedTallyResult(tr TallyResult) WrappedTallyResult {
	return WrappedTallyResult{
		TR: tr,
	}
}
