package txindex

import (
	"context"

	"github.com/okex/exchain/libs/tendermint/libs/service"

	"github.com/okex/exchain/libs/tendermint/types"
)

const (
	subscriber = "IndexerService"
)

// IndexerService connects event bus and transaction indexer together in order
// to index transactions coming from event bus.
type IndexerService struct {
	service.BaseService

	idr      TxIndexer
	eventBus *types.EventBus
	quit     chan struct{}
}

// NewIndexerService returns a new service instance.
func NewIndexerService(idr TxIndexer, eventBus *types.EventBus) *IndexerService {
	is := &IndexerService{idr: idr, eventBus: eventBus}
	is.BaseService = *service.NewBaseService(nil, "IndexerService", is)
	is.quit = make(chan struct{})
	return is
}

// OnStart implements service.Service by subscribing for all transactions
// and indexing them by events.
func (is *IndexerService) OnStart() error {
	// Use SubscribeUnbuffered here to ensure both subscriptions does not get
	// cancelled due to not pulling messages fast enough. Cause this might
	// sometimes happen when there are no other subscribers.

	blockHeadersSub, err := is.eventBus.SubscribeUnbuffered(
		context.Background(),
		subscriber,
		types.EventQueryNewBlockHeader)
	if err != nil {
		return err
	}

	txsSub, err := is.eventBus.SubscribeUnbuffered(context.Background(), subscriber, types.EventQueryTx)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-blockHeadersSub.Out():
				eventDataHeader := msg.Data().(types.EventDataNewBlockHeader)
				height := eventDataHeader.Header.Height
				batch := NewBatch(eventDataHeader.NumTxs)
				for i := int64(0); i < eventDataHeader.NumTxs; i++ {
					msg2 := <-txsSub.Out()
					txResult := msg2.Data().(types.EventDataTx).TxResult
					if err = batch.Add(&txResult); err != nil {
						is.Logger.Error("Can't add tx to batch",
							"height", height,
							"index", txResult.Index,
							"err", err)
					}
				}
				if err = is.idr.AddBatch(batch); err != nil {
					is.Logger.Error("Failed to index block", "height", height, "err", err)
				} else {
					is.Logger.Info("Indexed block", "height", height)
				}
			case <-blockHeadersSub.Cancelled():
				close(is.quit)
				return
			}

		}
	}()
	return nil
}

// OnStop implements service.Service by unsubscribing from all transactions.
func (is *IndexerService) OnStop() {
	if is.eventBus.IsRunning() {
		_ = is.eventBus.UnsubscribeAll(context.Background(), subscriber)
	}
}

func (is *IndexerService) Wait() {
	<-is.quit
}
