package transfer

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/application/transfer/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/application/transfer/types"
)

var (
	NewKeeper            = keeper.NewKeeper
	ModuleCdc    = types.ModuleCdc
)
