package common

import sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

// IBC port sentinel errors
var (
	ErrDisableProxyBeforeHeight = sdkerrors.Register(ModuleProxy, 1, "this feature is disable")
)
