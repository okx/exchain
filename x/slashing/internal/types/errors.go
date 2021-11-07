package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const BaseSlashingError = 5001

// x/slashing module sentinel errors
var (
	ErrNoValidatorForAddress        = sdkerrors.Register(ModuleName, BaseSlashingError+1, "address is not associated with any known validator")
	ErrBadValidatorAddr             = sdkerrors.Register(ModuleName, BaseSlashingError+2, "validator does not exist for that address")
	ErrValidatorJailed              = sdkerrors.Register(ModuleName, BaseSlashingError+3, "validator still jailed; cannot be unjailed")
	ErrValidatorNotJailed           = sdkerrors.Register(ModuleName, BaseSlashingError+4, "validator not jailed; cannot be unjailed")
	ErrMissingSelfDelegation        = sdkerrors.Register(ModuleName, BaseSlashingError+5, "validator has no self-delegation; cannot be unjailed")
	ErrSelfDelegationTooLowToUnjail = sdkerrors.Register(ModuleName, BaseSlashingError+6, "validator's self delegation less than minimum; cannot be unjailed")
	ErrNoSigningInfoFound           = sdkerrors.Register(ModuleName, BaseSlashingError+7, "no validator signing info found")
	ErrValidatorDestroying          = sdkerrors.Register(ModuleName, BaseSlashingError+8, "can not unjail it when the validator is destroying")
)
