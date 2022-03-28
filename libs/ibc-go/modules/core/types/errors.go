package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
)

var ErrIbcDisabled = sdkerrors.Register(host.ModuleName, 1, "IBC are disabled")
