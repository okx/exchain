package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

// errors
var (
	ErrInternalFeeSplit              = sdkerrors.Register(ModuleName, 2, "internal feesplit error")
	ErrFeeSplitDisabled              = sdkerrors.Register(ModuleName, 3, "feesplit module is disabled by governance")
	ErrFeeSplitAlreadyRegistered     = sdkerrors.Register(ModuleName, 4, "feesplit already exists for given contract")
	ErrFeeSplitNoContractDeployed    = sdkerrors.Register(ModuleName, 5, "no contract deployed")
	ErrFeeSplitContractNotRegistered = sdkerrors.Register(ModuleName, 6, "no feesplit registered for contract")
	ErrFeeSplitDeployerIsNotEOA      = sdkerrors.Register(ModuleName, 7, "deployer is not EOA")
	ErrFeeAccountNotFound            = sdkerrors.Register(ModuleName, 8, "account not found")
	ErrDerivedNotMatched             = sdkerrors.Register(ModuleName, 9, "derived address not matched")
)
