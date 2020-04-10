package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGov "github.com/cosmos/cosmos-sdk/x/gov"
)

// Vote defines the vote for proposal
type Vote struct {
	sdkGov.Vote
	VoteID uint64 `json:"vote_id" yaml:"vote_id"`
}

// NewVote creates a new Vote instance
func NewVote(proposalID uint64, voter sdk.AccAddress, option sdkGov.VoteOption, voteID uint64) Vote {
	return Vote{sdkGov.NewVote(proposalID, voter, option), voteID}
}

func (v Vote) String() string {
	return fmt.Sprintf("%s, the vote ID is %d", v.Vote.String(), v.VoteID)
}

// Votes is a collection of Vote objects
type Votes []Vote

func (v Votes) String() string {
	if len(v) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Votes for Proposal %d:", v[0].ProposalID)
	for _, vot := range v {
		out += fmt.Sprintf("\n  %s: %s", vot.Voter, vot.Option)
	}
	return out
}

// Equals returns whether two votes are equal.
func (v Vote) Equals(comp Vote) bool {
	return v.Vote.Equals(comp.Vote) && v.VoteID == comp.VoteID
}

// Empty returns whether a vote is empty.
func (v Vote) Empty() bool {
	return v.Vote.Equals(sdkGov.Vote{}) && v.VoteID == 0
}
