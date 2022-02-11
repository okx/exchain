package feemarket

import (
	"github.com/okex/exchain/x/feemarket/keeper"
	"github.com/okex/exchain/x/feemarket/types"
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
	NewKeeper = keeper.NewKeeper
)

//nolint
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
