package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ContractDeploymentWhitelist is the type alias for []sdk.AccAddress
type ContractDeploymentWhitelist []sdk.AccAddress

// String returns a human readable string representation of ContractDeploymentWhitelist
func (cdw ContractDeploymentWhitelist) String() string {
	var b strings.Builder
	b.WriteString("Contract Deployment Whitelist:\n")
	for _, addr := range cdw {
		b.WriteString(addr.String())
		b.WriteByte('\n')
	}

	return strings.TrimSpace(b.String())
}
