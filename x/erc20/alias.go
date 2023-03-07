package erc20

import (
	"github.com/okx/okbchain/x/erc20/keeper"
	"github.com/okx/okbchain/x/erc20/types"
)

// nolint
const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

// nolint
var (
	NewKeeper                = keeper.NewKeeper
	NewIBCTransferHooks      = keeper.NewIBCTransferHooks
	NewSendToIbcEventHandler = keeper.NewSendToIbcEventHandler

	NewSendNative20ToIbcEventHandler = keeper.NewSendNative20ToIbcEventHandler
)

//nolint
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
