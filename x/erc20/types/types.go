package types

import (
	"strings"

	ibctransferType "github.com/okx/okbchain/libs/ibc-go/modules/apps/transfer/types"
)

const (
	IbcDenomPrefix = "ibc/"
	ibcDenomLen    = len(IbcDenomPrefix) + 64
)

// IsValidIBCDenom returns if denom is a valid ibc denom
func IsValidIBCDenom(denom string) bool {
	if err := ibctransferType.ValidateIBCDenom(denom); err != nil {
		return false
	}
	return len(denom) == ibcDenomLen && strings.HasPrefix(denom, IbcDenomPrefix)
}
