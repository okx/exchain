package feesplit

import (
	"github.com/okex/exchain/x/feesplit/keeper"
	"github.com/okex/exchain/x/feesplit/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
