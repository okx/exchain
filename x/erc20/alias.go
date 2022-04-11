package erc20

import (
	"github.com/okex/exchain/x/erc20/keeper"
	"github.com/okex/exchain/x/erc20/types"
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
)

//nolint
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
