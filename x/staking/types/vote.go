package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Votes is the alias of sdk.Dec to represent the number of votes for multi-voting function
type Votes = sdk.Dec

// MustUnmarshalVote unmarshals the vote bytes and return it
func MustUnmarshalVote(cdc *codec.Codec, bytes []byte) Votes {
	var vote Votes
	cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &vote)
	return vote
}

// VoteResponse is the struct for query all the votes on a validator
type VoteResponse struct {
	VoterAddr sdk.AccAddress `json:"voter_address"`
	Votes     sdk.Dec        `json:"votes"`
}

// NewVoteResponse creates a new instance of VoteResponse
func NewVoteResponse(voterAddr sdk.AccAddress, votes Votes) VoteResponse {
	return VoteResponse{
		voterAddr,
		votes,
	}
}

// String returns a human readable string representation of VoteResponse
func (vr VoteResponse) String() string {
	return fmt.Sprintf("%s\n  Votes:   %s", vr.VoterAddr.String(), vr.Votes)
}

// VoteResponses is the type alias of VoteResponse slice
type VoteResponses []VoteResponse

// String returns a human readable string representation of VoteResponses
func (vrs VoteResponses) String() (strFormat string) {
	for _, vr := range vrs {
		strFormat = fmt.Sprintf("%s%s\n", strFormat, vr.String())
	}

	return strings.TrimSpace(strFormat)
}
