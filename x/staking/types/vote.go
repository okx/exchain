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

// SharesResponse is the struct for query all the shares added to a validator
type SharesResponse struct {
	DelAddr sdk.AccAddress `json:"delegator_address"`
	Shares  sdk.Dec        `json:"shares"`
}

// NewSharesResponse creates a new instance of SharesResponse
func NewSharesResponse(delAddr sdk.AccAddress, shares Votes) SharesResponse {
	return SharesResponse{
		delAddr,
		shares,
	}
}

// String returns a human readable string representation of SharesResponse
func (sr SharesResponse) String() string {
	return fmt.Sprintf("%s\n  Shares:   %s", sr.DelAddr.String(), sr.Shares)
}

// SharesResponses is the type alias of SharesResponse slice
type SharesResponses []SharesResponse

// String returns a human readable string representation of SharesResponses
func (srs SharesResponses) String() (strFormat string) {
	for _, sr := range srs {
		strFormat = fmt.Sprintf("%s%s\n", strFormat, sr.String())
	}

	return strings.TrimSpace(strFormat)
}
