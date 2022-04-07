package ibc

import (
	"github.com/okex/exchain/libs/ibc-go/modules/core/keeper"
	"github.com/okex/exchain/libs/ibc-go/modules/core/types"
)

type (
	Keeper           = keeper.Keeper
)
const(
)
var (
	NewKeeper  = keeper.NewKeeper
	ModuleCdc    = types.ModuleCdc
	DefaultGenesisState  = types.DefaultGenesisState
)
