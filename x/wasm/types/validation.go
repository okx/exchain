package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxWasmSize = 500 * 1024

	// MaxLabelSize is the longest label that can be used when Instantiating a contract
	MaxLabelSize = 128
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
