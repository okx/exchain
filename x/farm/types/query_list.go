package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"strings"
)

// PoolNameList is the type alias for []string
type PoolNameList []string

// String returns a human readable string representation of PoolNameList
func (pnl PoolNameList) String() string {
	out := "Pool Name List:\n"
	for _, poolName := range pnl {
		out = fmt.Sprintf("%s  %s\n", out, poolName)
	}
	return strings.TrimSpace(out)
}

// AccAddrList is the type alias for []sdk.AccAddress
type AccAddrList []sdk.AccAddress

// String returns a human readable string representation of AccAddrList
func (aal AccAddrList) String() string {
	out := "Account Address List:\n"
	for _, accAddr := range aal {
		out = fmt.Sprintf("%s  %s\n", out, accAddr.String())
	}
	return strings.TrimSpace(out)
}
