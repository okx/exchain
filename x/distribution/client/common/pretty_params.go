package common

import (
	"encoding/json"
	"fmt"
)

// PrettyParams is the struct for CLI output
type PrettyParams struct {
	WithdrawAddrEnabled json.RawMessage `json:"withdraw_addr_enabled"`
}

// NewPrettyParams creates a new PrettyParams
func NewPrettyParams(withdrawAddrEnabled json.RawMessage) PrettyParams {
	return PrettyParams{
		WithdrawAddrEnabled: withdrawAddrEnabled,
	}
}

// String returns the params string
func (pp PrettyParams) String() string {
	return fmt.Sprintf(`Distribution Params:
  Withdraw Addr Enabled:  %s`, pp.WithdrawAddrEnabled)
}
