package match

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okexchain/x/order/keeper"
	"github.com/okex/okexchain/x/order/match/continuousauction"
	"github.com/okex/okexchain/x/order/match/periodicauction"
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
