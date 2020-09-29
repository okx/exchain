package types

import (
	"fmt"
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
