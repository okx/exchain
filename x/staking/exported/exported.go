package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegatorI expected delegator functions
type DelegatorI interface {
	GetVotedValidatorAddresses() []sdk.ValAddress
	GetLastVotes() sdk.Dec
}
