// nolint
package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	CodeInvalidCommissionRate                 uint32 = 67047
	CodeNotSupportEditValidatorCommissionRate uint32 = 67048
	CodeDisabledOperate                       uint32 = 67049
)

// ErrInvalidCommissionRate returns an error when commission rate not be between 0 and 1 (inclusive)
func ErrInvalidCommissionRate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCommissionRate,
		"commission rate must be between 0 and 1 (inclusive)")
}

func ErrCodeNotSupportEditValidatorCommissionRate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNotSupportEditValidatorCommissionRate,
		"not support edit validator commission rate")
}

func ErrCodeDisabledOperate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeDisabledOperate, "disable operate")
}
