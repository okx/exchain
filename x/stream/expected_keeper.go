//+build !stream

package stream

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/order"
	"github.com/okex/okexchain/x/stream/exported"
	"github.com/okex/okexchain/x/token"
)

type OrderKeeper interface {
	GetOrder(ctx sdk.Context, orderID string) *order.Order
	GetUpdatedOrderIDs() []string
	GetBlockOrderNum(ctx sdk.Context, blockHeight int64) int64
	GetBlockMatchResult() *order.BlockMatchResult
	GetLastPrice(ctx sdk.Context, product string) sdk.Dec
	GetBestBidAndAsk(ctx sdk.Context, product string) (sdk.Dec, sdk.Dec)
	GetUpdatedDepthbookKeys() []string
	GetDepthBookCopy(product string) *order.DepthBook
	GetProductPriceOrderIDs(key string) []string
}

type TokenKeeper interface {
	GetFeeDetailList() []*token.FeeDetail
	GetCoinsInfo(ctx sdk.Context, addr sdk.AccAddress) token.CoinsInfo
}

type DexKeeper interface {
	GetTokenPairs(ctx sdk.Context) []*dex.TokenPair
	SetObserverKeeper(sk exported.StreamKeeper)
}

type AccountKeeper interface {
	SetObserverKeeper(observer auth.ObserverI)
}
