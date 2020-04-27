package common

import (
	"encoding/json"
	"fmt"
)

// PrettyParams is the struct for CLI output
type PrettyParams struct {
	CommunityTax        json.RawMessage `json:"community_tax"`
	WithdrawAddrEnabled json.RawMessage `json:"withdraw_addr_enabled"`
}

// newPrettyParams creates a new PrettyParams
func newPrettyParams(communityTax, withdrawAddrEnabled json.RawMessage) PrettyParams {
	return PrettyParams{
		CommunityTax:        communityTax,
		WithdrawAddrEnabled: withdrawAddrEnabled,
	}
}

// String returns the params string
func (pp PrettyParams) String() string {
	return fmt.Sprintf(`Distribution Params:
  Community Tax:          %s
  Withdraw Addr Enabled:  %s`,
		pp.CommunityTax, pp.WithdrawAddrEnabled)
}
