package authz

import (
	sdkerrors "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/errors"
)

// x/authz module sentinel errors
var (
	ErrInvalidExpirationTime = sdkerrors.Register(ModuleName, 3, "expiration time of authorization should be more than current time")
)
