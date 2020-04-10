package keeper

import (
	sdkGov "github.com/cosmos/cosmos-sdk/x/gov"
)

var (
	// NewRouter is alias of cm gov NewRouter
	NewRouter = sdkGov.NewRouter
)
