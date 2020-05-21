//+build !stream

package stream

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/dex"
	"github.com/okex/okchain/x/order"
	"github.com/okex/okchain/x/token"
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

type AccountKeeper interface {
	GetUpdatedAccAddress(ctx sdk.Context) (accs []sdk.AccAddress)
}

type DexKeeper interface {
	IsTokenPairChanged(ctx sdk.Context) bool
	GetTokenPairs(ctx sdk.Context) []*dex.TokenPair
	GetNewTokenPairs(ctx sdk.Context) []*dex.TokenPair
}
