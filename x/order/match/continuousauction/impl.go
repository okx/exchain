package continuousauction

import (
	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"

	"github.com/okx/exchain/x/order/keeper"
)

// nolint
type CaEngine struct {
}

// nolint
func (e *CaEngine) Run(ctx sdk.Context, keeper keeper.Keeper) {
	// TODO
}
