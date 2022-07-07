package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// GetLastAddedShares gets the last shares added to validators of a delegator for other module
func (d Delegator) GetDelegatorAddress() sdk.AccAddress {
	return d.DelegatorAddress
}
