package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegatorI expected delegator functions
type DelegatorI interface {
	GetShareAddedValidatorAddresses() []sdk.ValAddress
	GetLastAddedShares() sdk.Dec
}
