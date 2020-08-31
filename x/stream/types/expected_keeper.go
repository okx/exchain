package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order"
	"github.com/okex/okchain/x/stream/exported"
	"github.com/okex/okchain/x/token"
	"github.com/willf/bitset"
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
	GetTxHandlerMsgResult() []bitset.BitSet
}

type TokenKeeper interface {
	GetFeeDetailList() []*token.FeeDetail
	GetCoinsInfo(ctx sdk.Context, addr sdk.AccAddress) token.CoinsInfo
}

type AccountKeeper interface {
	SetObserverKeeper(observer auth.ObserverI)
}

type DexKeeper interface {
	GetTokenPairs(ctx sdk.Context) []*dex.TokenPair
	SetObserverKeeper(keeper exported.StreamKeeper)
}
