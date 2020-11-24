// DONTCOVER
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const BaseEvidenceError = 9001


// x/evidence module sentinel errors
var (
	ErrNoEvidenceHandlerExists = sdkerrors.Register(ModuleName, BaseEvidenceError + 1, "unregistered handler for evidence type")
	ErrInvalidEvidence         = sdkerrors.Register(ModuleName, BaseEvidenceError + 2, "invalid evidence")
	ErrNoEvidenceExists        = sdkerrors.Register(ModuleName, BaseEvidenceError + 3, "evidence does not exist")
	ErrEvidenceExists          = sdkerrors.Register(ModuleName, BaseEvidenceError + 4, "evidence already exists")
)
