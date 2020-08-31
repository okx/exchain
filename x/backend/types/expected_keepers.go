package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dextypes "github.com/okex/okchain/x/dex/types"
	"github.com/okex/okchain/x/order"
	ordertypes "github.com/okex/okchain/x/order/types"
	"github.com/okex/okchain/x/stream/exported"
	"github.com/okex/okchain/x/token"
	"github.com/willf/bitset"
)

//OrderKeeper expected order keeper
type OrderKeeper interface {
	GetOrder(ctx sdk.Context, orderID string) *order.Order
	GetUpdatedOrderIDs() []string
	GetTxHandlerMsgResult() []bitset.BitSet
	GetBlockOrderNum(ctx sdk.Context, blockHeight int64) int64
	GetBlockMatchResult() *ordertypes.BlockMatchResult
	GetLastPrice(ctx sdk.Context, product string) sdk.Dec
	GetBestBidAndAsk(ctx sdk.Context, product string) (sdk.Dec, sdk.Dec)
}

// TokenKeeper expected token keeper
type TokenKeeper interface {
	GetFeeDetailList() []*token.FeeDetail
	GetParams(ctx sdk.Context) (params token.Params)
}

// DexKeeper expected dex keeper
type DexKeeper interface {
	GetTokenPairs(ctx sdk.Context) []*dextypes.TokenPair
	GetTokenPair(ctx sdk.Context, product string) *dextypes.TokenPair
	SetObserverKeeper(keeper exported.StreamKeeper)
}

// MarketKeeper expected market keeper which would get data from pulsar & redis
type MarketKeeper interface {
	GetTickerByProducts(products []string) ([]map[string]string, error)
	GetKlineByProductID(productID uint64, granularity, size int) ([][]string, error)
}
