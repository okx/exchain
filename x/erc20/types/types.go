package types

import (
	"strings"

	ibctransferType "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
)

const (
	ibcDenomPrefix = "ibc/"
	ibcDenomLen    = len(ibcDenomPrefix) + 64
)

// IsValidIBCDenom returns if denom is a valid ibc denom
func IsValidIBCDenom(denom string) bool {
	if err := ibctransferType.ValidateIBCDenom(denom); err != nil {
		return false
	}
	return len(denom) == ibcDenomLen && strings.HasPrefix(denom, ibcDenomPrefix)
}
