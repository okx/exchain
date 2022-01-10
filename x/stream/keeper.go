package stream

import (
	"fmt"

	"github.com/okex/exchain/x/ammswap"
	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/stream/types"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server/config"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/common/monitor"
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
	}
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
}

func (k Keeper) OnSwapCreateExchange(ctx sdk.Context, swapTokenPair ammswap.SwapTokenPair) {
	k.stream.Cache.AddNewSwapTokenPair(&swapTokenPair)
}

func (k Keeper) OnFarmClaim(ctx sdk.Context, address sdk.AccAddress, poolName string, claimedCoins sdk.SysCoins) {
	if claimedCoins.IsZero() {
		return
	}

}
