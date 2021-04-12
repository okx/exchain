package periodicauction

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/exchain/x/order/keeper"
)

// PaEngine is the periodic auction match engine
type PaEngine struct {
}

// nolint
func (e *PaEngine) Run(ctx sdk.Context, keeper keeper.Keeper) {
	cleanupExpiredOrders(ctx, keeper)
	cleanupOrdersWhoseTokenPairHaveBeenDelisted(ctx, keeper)
	matchOrders(ctx, keeper)
}
