// nolint
package types

import (
	"fmt"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
)

const (
	CodeInvalidCommissionRate                 uint32 = 67047
	CodeNotSupportEditValidatorCommissionRate uint32 = 67048
	CodeDisabledOperate                       uint32 = 67049
	CodeNoDelegatorValidator                  uint32 = 67050
)

// ErrInvalidCommissionRate returns an error when commission rate not be between 0 and 1 (inclusive)
func ErrInvalidCommissionRate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeInvalidCommissionRate,
		"commission rate must be between 0 and 1 (inclusive)")
}

func ErrCodeDisabledOperate() sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeDisabledOperate, "disable operate")
}

func ErrCodeNoDelegatorValidator(delegator string, validator string) sdk.Error {
	return sdkerrors.New(DefaultCodespace, CodeNoDelegatorValidator, fmt.Sprintf("delegator %s not vote validator %s", delegator, validator))
}
