package pendingtx

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	rpcfilters "github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/watcher"
)

type Watcher struct {
	clientCtx context.CLIContext
	events    *rpcfilters.EventSystem
	logger    log.Logger

	sender Sender
}

type Sender interface {
	Send(hash []byte, tx *watcher.Transaction) error
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
	pendingSub, _, err := w.events.SubscribePendingTxs()
	if err != nil {
		w.logger.Error("error creating block filter", "error", err.Error())
	}

	confirmedSub, _, err := w.events.SubscribeConfirmedTx()
	if err != nil {
		w.logger.Error("error creating block filter", "error", err.Error())
	}

	go func(pendingCh <-chan coretypes.ResultEvent, confirmedCh <-chan coretypes.ResultEvent) {
		for {
			select {
			case re := <-pendingCh:
				tx, txHash, err := w.newTransactionByEvent(re, "pending")
				if err != nil {
					continue
				}

				go func(hash []byte, tx *watcher.Transaction) {
					w.logger.Debug("push pending tx to MQ", "txHash=", txHash.String())
					err = w.sender.Send(hash, tx)
					if err != nil {
						w.logger.Error("failed to send pending tx", "hash", txHash.String(), "error", err)
					}
				}(txHash.Bytes(), tx)
			case re := <-confirmedCh:
				tx, txHash, err := w.newTransactionByEvent(re, "confirmed")
				if err != nil {
					continue
				}

				go func(hash []byte, tx *watcher.Transaction) {
					w.logger.Debug("push confirmed tx to MQ", "txHash=", txHash.String())
					err = w.sender.Send(hash, tx)
					if err != nil {
						w.logger.Error("failed to send confirmed tx", "hash", txHash.String(), "error", err)
					}
				}(txHash.Bytes(), tx)
			}
		}
	}(pendingSub.Event(), confirmedSub.Event())
}

func (w *Watcher) newTransactionByEvent(re coretypes.ResultEvent, txType string) (*watcher.Transaction, common.Hash, error) {
	data, ok := re.Data.(tmtypes.EventDataTx)
	if !ok {
		w.logger.Error(fmt.Sprintf("invalid %s tx data type %T, expected EventDataTx", txType, re.Data))
		return nil, common.Hash{}, errors.New("invalid tx data type")
	}

	txHash := common.BytesToHash(data.Tx.Hash(data.Height))
	w.logger.Debug(fmt.Sprintf("receive %s tx", txType), "txHash=", txHash.String())

	// only watch evm tx
	ethTx, err := rpctypes.RawTxToEthTx(w.clientCtx, data.Tx, data.Height)
	if err != nil {
		w.logger.Error("failed to decode raw tx to eth tx", "hash", txHash.String(), "error", err)
		return nil, common.Hash{}, err
	}

	tx, err := watcher.NewTransaction(ethTx, txHash, common.Hash{}, uint64(data.Height), uint64(data.Index))
	if err != nil {
		w.logger.Error("failed to new transaction", "hash", txHash.String(), "error", err)
		return nil, common.Hash{}, err
	}

	return tx, txHash, nil
}
