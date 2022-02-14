package ibc

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/types"
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
