package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

// Param module codespace constants
const (
	DefaultCodespace string = "params"
	BaseParamsError = 4001

	CodeInvalidMaxProposalNum uint32 = BaseParamsError+4
)

// ErrInvalidMaxProposalNum returns error when the number of params to change are out of limit
var RegisteredErrInvalidParamsNum = sdkerrors.Register(params.ModuleName, CodeInvalidMaxProposalNum, "invalid param number")

// ErrInvalidParamsNum returns error when the number of params to change are out of limit
func ErrInvalidParamsNum(codespace string, msg string) sdk.EnvelopedErr {
	return sdk.EnvelopedErr{sdkerrors.Wrap(RegisteredErrInvalidParamsNum, msg)}
}
