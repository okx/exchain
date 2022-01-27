package types

import sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

// NOTE: We can't use 1 since that error code is reserved for internal errors.
const (
	DefaultCodespace string = ModuleName
)

var (
	// ErrInvalidState returns an error resulting from an invalid Storage State.
	ErrEmptyTreasures = sdkerrors.Register(ModuleName, 2, "treasures is not empty")

	ErrDuplicatedTreasure     = sdkerrors.Register(ModuleName, 3, "treasures can not be duplicate")
	ErrUnexpectedProposalType = sdkerrors.Register(ModuleName, 3, "Unsupported proposal type of mint module")
)
