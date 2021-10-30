package stream

import (
	"fmt"

	"github.com/okex/exchain/x/ammswap"
	backend "github.com/okex/exchain/x/backend/types"
	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/stream/types"

	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
	"github.com/okex/exchain/dependence/cosmos-sdk/server/config"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/auth"
	"github.com/okex/exchain/x/common/monitor"
	"github.com/okex/exchain/dependence/tendermint/libs/log"
)

// nolint
type Keeper struct {
	metric *monitor.StreamMetrics
	stream *Stream
}

// nolint
func NewKeeper(orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, accountKeeper types.AccountKeeper,
	swapKeeper types.SwapKeeper, farmKeeper types.FarmKeeper, cdc *codec.Codec, logger log.Logger, cfg *config.Config, metrics *monitor.StreamMetrics) Keeper {
	logger = logger.With("module", "stream")
	k := Keeper{
		metric: metrics,
		stream: NewStream(orderKeeper, tokenKeeper, dexKeeper, swapKeeper, farmKeeper, cdc, logger, cfg),
	}
	if k.stream.engines != nil {
		dexKeeper.SetObserverKeeper(k)
		accountKeeper.SetObserverKeeper(k)
		swapKeeper.SetObserverKeeper(k)
		farmKeeper.SetObserverKeeper(k)
	}
	return k
}

// nolint
func (k Keeper) SyncTx(ctx sdk.Context, tx *auth.StdTx, txHash string, timestamp int64) {
	if k.stream.engines[EngineAnalysisKind] != nil {
		k.stream.logger.Debug(fmt.Sprintf("[stream engine] get new tx, txHash: %s", txHash))
		txs := backend.GenerateTx(tx, txHash, ctx, k.stream.orderKeeper, timestamp)
		for _, tx := range txs {
			k.stream.Cache.AddTransaction(tx)
		}
	}
}

// GetMarketKeeper returns market keeper
func (k Keeper) GetMarketKeeper() MarketKeeper {
	return k.stream.marketKeeper
}

// AnalysisEnable returns true when analysis is enable
func (k Keeper) AnalysisEnable() bool {
	return k.stream.AnalysisEnable
}

// OnAddNewTokenPair called by dex when new token pair listed
func (k Keeper) OnAddNewTokenPair(ctx sdk.Context, tokenPair *dex.TokenPair) {
	k.stream.logger.Debug(fmt.Sprintf("OnAddNewTokenPair:%s", tokenPair.Name()))
	k.stream.Cache.AddNewTokenPair(tokenPair)
	k.stream.Cache.SetTokenPairChanged(true)
}

// OnTokenPairUpdated called by dex when token pair updated
func (k Keeper) OnTokenPairUpdated(ctx sdk.Context) {
	k.stream.logger.Debug("OnTokenPairUpdated:true")
	k.stream.Cache.SetTokenPairChanged(true)
}

// OnAccountUpdated called by auth when account updated
func (k Keeper) OnAccountUpdated(acc auth.Account) {
	k.stream.logger.Debug(fmt.Sprintf("OnAccountUpdated:%s", acc.GetAddress()))
	k.stream.Cache.AddUpdatedAccount(acc)
}

// OnSwapToken called by swap
func (k Keeper) OnSwapToken(ctx sdk.Context, address sdk.AccAddress, swapTokenPair ammswap.SwapTokenPair, sellAmount sdk.SysCoin, buyAmount sdk.SysCoin) {
	swapInfo := &backend.SwapInfo{
		Address:          address.String(),
		TokenPairName:    swapTokenPair.TokenPairName(),
		BaseTokenAmount:  swapTokenPair.BasePooledCoin.String(),
		QuoteTokenAmount: swapTokenPair.QuotePooledCoin.String(),
		SellAmount:       sellAmount.String(),
		BuysAmount:       buyAmount.String(),
		Price:            swapTokenPair.BasePooledCoin.Amount.Quo(swapTokenPair.QuotePooledCoin.Amount).String(),
		Timestamp:        ctx.BlockTime().Unix(),
	}
	k.stream.Cache.AddSwapInfo(swapInfo)
}

func (k Keeper) OnSwapCreateExchange(ctx sdk.Context, swapTokenPair ammswap.SwapTokenPair) {
	k.stream.Cache.AddNewSwapTokenPair(&swapTokenPair)
}

func (k Keeper) OnFarmClaim(ctx sdk.Context, address sdk.AccAddress, poolName string, claimedCoins sdk.SysCoins) {
	if claimedCoins.IsZero() {
		return
	}
	claimInfo := &backend.ClaimInfo{
		Address:   address.String(),
		PoolName:  poolName,
		Claimed:   claimedCoins.String(),
		Timestamp: ctx.BlockTime().Unix(),
	}
	k.stream.Cache.AddClaimInfo(claimInfo)
}
