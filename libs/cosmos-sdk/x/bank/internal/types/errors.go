package types

import (
	"fmt"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

// x/bank module sentinel errors
var (
	ErrNoInputs            = sdkerrors.Register(ModuleName, 1, "no inputs to send transaction")
	ErrNoOutputs           = sdkerrors.Register(ModuleName, 2, "no outputs to send transaction")
	ErrInputOutputMismatch = sdkerrors.Register(ModuleName, 3, "sum inputs != sum outputs")
	ErrSendDisabled        = sdkerrors.Register(ModuleName, 4, "send transactions are disabled")
)

func ErrUnSupportQueryType(data string) *sdkerrors.Error {
	return sdkerrors.Register(ModuleName, 5, fmt.Sprintf("%s is not support", data))
}
