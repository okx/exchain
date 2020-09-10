package stream

import (
	"fmt"

	"github.com/okex/okchain/x/dex"

	"github.com/okex/okchain/x/stream/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/common/monitor"
	"github.com/tendermint/tendermint/libs/log"
)

// nolint
type Keeper struct {
	metric *monitor.StreamMetrics
	stream *Stream
}

// nolint
func NewKeeper(orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, accountKeeper types.AccountKeeper, cdc *codec.Codec, logger log.Logger, cfg *config.Config, metrics *monitor.StreamMetrics) Keeper {
	logger = logger.With("module", "stream")
	k := Keeper{
		metric: metrics,
		stream: NewStream(orderKeeper, tokenKeeper, dexKeeper, cdc, logger, cfg),
	}
	dexKeeper.SetObserverKeeper(k)
	accountKeeper.SetObserverKeeper(k)
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
