package types

import (
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// OVPair is designed for querying validators by rest
type OVPair struct {
	OperAddr sdk.ValAddress `json:"operator_address"`
	ValAddr  string         `json:"validator_address"`
}

// OVPairs is the alias of the OVPair slice
type OVPairs []OVPair

// String returns a human readable string representation of OVPairs
func (ovPairs OVPairs) String() (out string) {
	for _, ovPair := range ovPairs {
		out = fmt.Sprintf("%s%s:%s\n", out, ovPair.OperAddr.String(), ovPair.ValAddr)
	}
	return
}
