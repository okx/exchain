package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/ammswap"
	"github.com/okex/okexchain/x/backend/exported"
	dextypes "github.com/okex/okexchain/x/dex/types"
	"github.com/okex/okexchain/x/order"
	ordertypes "github.com/okex/okexchain/x/order/types"
	streamexported "github.com/okex/okexchain/x/stream/exported"
	"github.com/okex/okexchain/x/token"
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
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.SysCoins
	GetTokensInfo(ctx sdk.Context) (tokens []token.Token)
}

// DexKeeper expected dex keeper
type DexKeeper interface {
	GetTokenPairs(ctx sdk.Context) []*dextypes.TokenPair
	GetTokenPair(ctx sdk.Context, product string) *dextypes.TokenPair
	SetObserverKeeper(keeper streamexported.StreamKeeper)
}

// MarketKeeper expected market keeper which would get data from pulsar & redis
type MarketKeeper interface {
	GetTickerByProducts(products []string) ([]map[string]string, error)
	GetKlineByProductID(productID uint64, granularity, size int) ([][]string, error)
}

// SwapKeeper expected swap keeper
type SwapKeeper interface {
	GetSwapTokenPairs(ctx sdk.Context) []ammswap.SwapTokenPair
	GetSwapTokenPair(ctx sdk.Context, tokenPairName string) (ammswap.SwapTokenPair, error)
	GetParams(ctx sdk.Context) (params ammswap.Params)
	GetPoolTokenAmount(ctx sdk.Context, poolTokenName string) sdk.Dec
	SetObserverKeeper(k exported.BackendKeeper)
}
