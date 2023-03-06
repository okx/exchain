package types

import (
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
)

// ICA Host sentinel errors
var (
	ErrHostSubModuleDisabled = sdkerrors.Register(SubModuleName, 2, "host submodule is disabled")
)
