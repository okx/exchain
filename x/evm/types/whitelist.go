package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddressList is the type alias for []sdk.AccAddress
type AddressList []sdk.AccAddress

// String returns a human readable string representation of AddressList
func (al AddressList) String() string {
	var b strings.Builder
	b.WriteString("Address List:\n")
	for i := 0; i < len(al); i++ {
		b.WriteString(al[i].String())
		b.WriteByte('\n')
	}

	return strings.TrimSpace(b.String())
}
