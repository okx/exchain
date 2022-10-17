package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const DefaultCodespace string = ModuleName

// errors
var (
	ErrNotFeesplitHeight             = sdkerrors.Register(DefaultCodespace, 1, "not feesplit height")
	ErrInternalFeeSplit              = sdkerrors.Register(DefaultCodespace, 2, "internal feesplit error")
	ErrFeeSplitDisabled              = sdkerrors.Register(DefaultCodespace, 3, "feesplit module is disabled by governance")
	ErrFeeSplitAlreadyRegistered     = sdkerrors.Register(DefaultCodespace, 4, "feesplit already exists for given contract")
	ErrFeeSplitNoContractDeployed    = sdkerrors.Register(DefaultCodespace, 5, "no contract deployed")
	ErrFeeSplitContractNotRegistered = sdkerrors.Register(DefaultCodespace, 6, "no feesplit registered for contract")
	ErrFeeSplitDeployerIsNotEOA      = sdkerrors.Register(DefaultCodespace, 7, "deployer is not EOA")
	ErrFeeAccountNotFound            = sdkerrors.Register(DefaultCodespace, 8, "account not found")
	ErrDerivedNotMatched             = sdkerrors.Register(DefaultCodespace, 9, "derived address not matched")
)
