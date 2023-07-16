package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

// MaxSaltSize is the longest salt that can be used when instantiating a contract
const MaxSaltSize = 64

var (
	// MaxLabelSize is the longest label that can be used when Instantiating a contract
	MaxLabelSize = 128 // extension point for chains to customize via compile flag.

	// MaxWasmSize is the largest a compiled contract code can be when storing code on chain
	MaxWasmSize = 800 * 1024 // extension point for chains to customize via compile flag.
)

func validateWasmCode(s []byte) error {
	if len(s) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "is required")
	}
	if len(s) > MaxWasmSize {
		return sdkerrors.Wrapf(ErrLimit, "cannot be longer than %d bytes", MaxWasmSize)
	}
	return nil
}

func validateLabel(label string) error {
	if label == "" {
		return sdkerrors.Wrap(ErrEmpty, "is required")
	}
	if len(label) > MaxLabelSize {
		return sdkerrors.Wrap(ErrLimit, "cannot be longer than 128 characters")
	}
	return nil
}

// ValidateSalt ensure salt constraints
func ValidateSalt(salt []byte) error {
	switch n := len(salt); {
	case n == 0:
		return sdkerrors.Wrap(ErrEmpty, "is required")
	case n > MaxSaltSize:
		return sdkerrors.Wrapf(ErrLimit, "cannot be longer than %d characters", MaxSaltSize)
	}
	return nil
}
