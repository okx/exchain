package transfer

import (
	"github.com/okex/exchain/libs/ibc-go/modules/application/transfer/keeper"
	"github.com/okex/exchain/libs/ibc-go/modules/application/transfer/types"
)

var (
	NewKeeper            = keeper.NewKeeper
	ModuleCdc    = types.ModuleCdc
)
