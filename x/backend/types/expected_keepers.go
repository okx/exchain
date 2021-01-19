package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/ammswap"
	ammswaptypes "github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/dex"
	dextypes "github.com/okex/okexchain/x/dex/types"
	farmtypes "github.com/okex/okexchain/x/farm/types"
	"github.com/okex/okexchain/x/order"
	ordertypes "github.com/okex/okexchain/x/order/types"
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
	SetObserverKeeper(keeper dex.StreamKeeper)
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
	SetObserverKeeper(k ammswaptypes.BackendKeeper)
}

// FarmKeeper expected farm keeper
type FarmKeeper interface {
	SetObserverKeeper(k farmtypes.BackendKeeper)
	GetFarmPools(ctx sdk.Context) (pools farmtypes.FarmPools)
	GetWhitelist(ctx sdk.Context) (whitelist farmtypes.PoolNameList)
	GetParams(ctx sdk.Context) (params farmtypes.Params)
	GetPoolLockedValue(ctx sdk.Context, pool farmtypes.FarmPool) sdk.Dec
	CalculateAmountYieldedBetween(ctx sdk.Context, pool farmtypes.FarmPool) (farmtypes.FarmPool, sdk.SysCoins)
	SupplyKeeper() supply.Keeper
	GetFarmPoolNamesForAccount(ctx sdk.Context, accAddr sdk.AccAddress) (poolNames farmtypes.PoolNameList)
	GetFarmPool(ctx sdk.Context, poolName string) (pool farmtypes.FarmPool, found bool)
	GetLockInfo(ctx sdk.Context, addr sdk.AccAddress, poolName string) (info farmtypes.LockInfo, found bool)
	GetEarnings(ctx sdk.Context, poolName string, accAddr sdk.AccAddress) (farmtypes.Earnings, sdk.Error)
}

// MintKeeper expected mint keeper
type MintKeeper interface {
	GetParams(ctx sdk.Context) (params mint.Params)
}
