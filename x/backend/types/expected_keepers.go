package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dextypes "github.com/okex/okchain/x/dex/types"
	"github.com/okex/okchain/x/order"
	ordertypes "github.com/okex/okchain/x/order/types"
	"github.com/okex/okchain/x/token"
)

// expected order keeper
type OrderKeeper interface {
	GetOrder(ctx sdk.Context, orderId string) *order.Order
	GetUpdatedOrderIDs() []string
	GetBlockOrderNum(ctx sdk.Context, blockHeight int64) int64
	GetBlockMatchResult() *ordertypes.BlockMatchResult
	GetLastPrice(ctx sdk.Context, product string) sdk.Dec
	GetBestBidAndAsk(ctx sdk.Context, product string) (sdk.Dec, sdk.Dec)
}

// expected token keeper
type TokenKeeper interface {
	GetFeeDetailList() []*token.FeeDetail
	GetParams(ctx sdk.Context) (params token.Params)
}

type DexKeeper interface {
	GetTokenPairs(ctx sdk.Context) []*dextypes.TokenPair
}

// expected market keeper which would get data from pulsar & redis
type MarketKeeper interface {
	InitTokenPairMap(ctx sdk.Context, dk DexKeeper)
	GetTickers() ([]map[string]string, error)
	GetTickerByInstruments(instruments []string) map[string]Ticker
	GetKlineByInstrument(instrument string, granularity, size int) ([][]string, error)
}
