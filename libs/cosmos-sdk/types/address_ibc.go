package types

import (
	"errors"
	"strings"
)

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func IBCAccAddressFromBech32(address string) (AccAddress, error) {
	if len(strings.TrimSpace(address)) == 0 {
		return AccAddress{}, errors.New("empty address string is not allowed")
	}
	return AccAddressFromBech32ByPrefix(address, GetConfig().GetBech32AccountAddrPrefix())
}
