// nolint
package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	CodeNotSupport uint32 = 67047
)

// ErrNotSupport returns an error when not support
func ErrNotSupport() sdk.EnvelopedErr {
	return sdk.EnvelopedErr{Err: sdkerrors.New(DefaultCodespace, CodeNotSupport,
		fmt.Sprintf("failed. The current version does not support it"))}
}
