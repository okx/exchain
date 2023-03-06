package types

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

// GetDelegatorAddress gets delegator address
func (d Delegator) GetDelegatorAddress() sdk.AccAddress {
	return d.DelegatorAddress
}
