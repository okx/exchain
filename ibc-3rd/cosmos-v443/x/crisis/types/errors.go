package types

import (
	sdkerrors "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/errors"
)

// x/crisis module sentinel errors
var (
	ErrNoSender         = sdkerrors.Register(ModuleName, 2, "sender address is empty")
	ErrUnknownInvariant = sdkerrors.Register(ModuleName, 3, "unknown invariant")
)
