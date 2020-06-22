package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Shares is the alias of sdk.Dec to represent the amount of shares for adding shares to validators
type Shares = sdk.Dec

// MustUnmarshalShares unmarshals the shares bytes and return it
func MustUnmarshalShares(cdc *codec.Codec, bytes []byte) Shares {
	var shares Shares
	cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &shares)
	return shares
}

// SharesResponse is the struct for query all the shares on a validator
type SharesResponse struct {
	DelAddr sdk.AccAddress `json:"delegator_address"`
	Shares  sdk.Dec        `json:"shares"`
}

// NewSharesResponse creates a new instance of sharesResponse
func NewSharesResponse(delAddr sdk.AccAddress, shares Shares) SharesResponse {
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
