//+build !stream

package stream

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/common/monitor"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct{}

func NewKeeper(ok OrderKeeper, tk TokenKeeper, dk DexKeeper, ak AccountKeeper,
	cdc *codec.Codec, logger log.Logger,
	cfg *config.Config, metrics *monitor.StreamMetrics) Keeper {

	k := Keeper{}
	// do this only if stream module enabled
	// dk.SetObserverKeeper(k)
	// ak.SetObserverKeeper(k)
	return k
}

func (k Keeper) SyncTx(ctx sdk.Context, tx *auth.StdTx, txHash string, timestamp int64) {}

func (k Keeper) GetMarketKeeper() backend.MarketKeeper { return nil }

func (k Keeper) AnalysisEnable() bool { return false }
