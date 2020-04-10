package periodicauction

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/order/keeper"
)

// PaEngine is the periodic auction match engine
type PaEngine struct {
}

// Run
func (e *PaEngine) Run(ctx sdk.Context, keeper keeper.Keeper) {
}
