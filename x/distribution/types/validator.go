package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorAccumulatedCommission is the accumulated commission for a validator
// kept as a running counter, can be withdrawn at any time
type ValidatorAccumulatedCommission = sdk.DecCoins

// InitialValidatorAccumulatedCommission returns the initial accumulated commission (zero)
func InitialValidatorAccumulatedCommission() ValidatorAccumulatedCommission {
	return ValidatorAccumulatedCommission{}
}
