package pendingtx

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	rpcfilters "github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

type Watcher struct {
	clientCtx context.CLIContext
	events    *rpcfilters.EventSystem
	logger    log.Logger

	sender Sender
}

type Sender interface {
	Send(hash []byte, tx *rpctypes.Transaction) error
}

func NewWatcher(clientCtx context.CLIContext, log log.Logger, sender Sender) *Watcher {
	return &Watcher{
		clientCtx: clientCtx,
		events:    rpcfilters.NewEventSystem(clientCtx.Client),
		logger:    log.With("module", "pendingtx-watcher"),

		sender: sender,
	}
}

func (w *Watcher) Start() {
	sub, _, err := w.events.SubscribePendingTxs()
	if err != nil {
		w.logger.Error("error creating block filter", "error", err.Error())
	}

	go func(txsCh <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case ev := <-txsCh:
				data, ok := ev.Data.(tmtypes.EventDataTx)
				if !ok {
					w.logger.Error(fmt.Sprintf("invalid data type %T, expected EventDataTx", ev.Data), "ID", sub.ID())
					continue
				}
				txHash := common.BytesToHash(data.Tx.Hash(data.Height))
				w.logger.Debug("receive tx from mempool", "txHash=", txHash.String())

				ethTx, err := rpctypes.RawTxToEthTx(w.clientCtx, data.Tx)
				if err != nil {
					w.logger.Error("failed to decode raw tx to eth tx", "hash", txHash.String(), "error", err)
					continue
				}

				tx, err := rpctypes.NewTransaction(ethTx, txHash, common.Hash{}, uint64(data.Height), uint64(data.Index))
				if err != nil {
					w.logger.Error("failed to new transaction", "hash", txHash.String(), "error", err)
					continue
				}

				go func(hash []byte, tx *rpctypes.Transaction) {
					w.logger.Debug("push pending tx to MQ", "txHash=", txHash.String())
					err = w.sender.Send(hash, tx)
					if err != nil {
						w.logger.Error("failed to send pending tx", "hash", txHash.String(), "error", err)
					}
				}(txHash.Bytes(), tx)
			}
		}
	}(sub.Event(), sub.Err())
}
