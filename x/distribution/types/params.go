package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params - structure for params query
type Params struct {
	CommunityTax        sdk.Dec `json:"community_tax"`
	WithdrawAddrEnabled bool    `json:"withdraw_addr_enabled"`
}

// NewParams creates a new instance of Params
func NewParams(communityTax sdk.Dec, withdrawAddrEnabled bool) Params {
	return Params{
		CommunityTax:        communityTax,
		WithdrawAddrEnabled: withdrawAddrEnabled,
	}
}

// String returns a human readable string representation of Params
func (p Params) String() string {
	return fmt.Sprintf(`Distribution Params:
  Community Tax:          %s
  Withdraw Addr Enabled:  %t`,
		p.CommunityTax, p.WithdrawAddrEnabled)
}

// MarshalYAML implements the text format for yaml marshaling
func (p Params) MarshalYAML() (interface{}, error) {
	return p.String(), nil
}
