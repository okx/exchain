package types

import (
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
)

const DefaultCodespace string = ModuleName

// errors
var (
	ErrInternalFeeSplit              = sdkerrors.Register(DefaultCodespace, 1, "internal feesplit error")
	ErrFeeSplitDisabled              = sdkerrors.Register(DefaultCodespace, 2, "feesplit module is disabled by governance")
	ErrFeeSplitAlreadyRegistered     = sdkerrors.Register(DefaultCodespace, 3, "feesplit already exists for given contract")
	ErrFeeSplitNoContractDeployed    = sdkerrors.Register(DefaultCodespace, 4, "no contract deployed")
	ErrFeeSplitContractNotRegistered = sdkerrors.Register(DefaultCodespace, 5, "no feesplit registered for contract")
	ErrFeeSplitDeployerIsNotEOA      = sdkerrors.Register(DefaultCodespace, 6, "deployer is not EOA")
	ErrFeeAccountNotFound            = sdkerrors.Register(DefaultCodespace, 7, "account not found")
	ErrDerivedNotMatched             = sdkerrors.Register(DefaultCodespace, 8, "derived address not matched")
)
