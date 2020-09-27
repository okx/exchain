package types

import (
	"fmt"
	"strings"
)

// Whitelist is the type alias for []string
type Whitelist []string

// String returns a human readable string representation of Whitelist
func (w Whitelist) String() string {
	out := "Whitelist:\n"
	for _, poolName := range w {
		out = fmt.Sprintf("%s  %s\n", out, poolName)
	}
	return strings.TrimSpace(out)
}
