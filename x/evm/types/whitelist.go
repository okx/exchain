package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"strings"
)

// ContractDeploymentWhitelist is the type alias for []ethcmn.Address
type ContractDeploymentWhitelist []ethcmn.Address

// String returns a human readable string representation of ContractDeploymentWhitelist
func (cdw ContractDeploymentWhitelist) String() string {
	var b strings.Builder
	b.WriteString("Contract Deployment Whitelist:\n")
	for _, addr := range cdw {
		b.WriteString(addr.Hex())
		b.WriteByte('\n')
	}

	return strings.TrimSpace(b.String())
}
