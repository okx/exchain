package farm

import (
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

const (
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	ModuleName        = types.ModuleName
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
