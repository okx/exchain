package match

import (
	"sync"

	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"

	"github.com/okx/exchain/x/order/keeper"
	"github.com/okx/exchain/x/order/match/continuousauction"
	"github.com/okx/exchain/x/order/match/periodicauction"
)

// nolint
const DefaultAuctionType = "periodicauction"

// nolint
var (
	once        sync.Once
	engine      Engine
	auctionType = DefaultAuctionType
)

// GetEngine : periodic auction only today
func GetEngine() Engine {
	once.Do(func() {
		if auctionType == DefaultAuctionType {
			engine = &periodicauction.PaEngine{}
		} else {
			engine = &continuousauction.CaEngine{}
		}
	})
	return engine
}

// nolint
type Engine interface {
	Run(ctx sdk.Context, keeper keeper.Keeper)
}
