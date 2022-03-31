package infura

import (
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/common/monitor"
	evm "github.com/okex/exchain/x/evm/watcher"
)

// nolint
type Keeper struct {
	metric *monitor.StreamMetrics
	stream *Stream
}

// nolint
func NewKeeper(evmKeeper EvmKeeper, logger log.Logger, metrics *monitor.StreamMetrics) Keeper {
	logger = logger.With("module", "infura")
	k := Keeper{
		metric: metrics,
		stream: NewStream(logger),
	}
	if k.stream.enable {
		evmKeeper.SetObserverKeeper(k)
	}
	return k
}

func (k Keeper) OnSaveTransactionReceipt(tr evm.TransactionReceipt) {
	k.stream.cache.AddTransactionReceipt(tr)
}
