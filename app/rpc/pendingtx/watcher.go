package pendingtx

import (
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/common"
	rpcfilters "github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/tendermint/tendermint/libs/log"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Watcher struct {
	clientCtx context.CLIContext
	events    *rpcfilters.EventSystem
	logger    log.Logger

	sender Sender
}

type Sender interface {
	Send(hash, tx []byte) error
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
				txHash := common.BytesToHash(data.Tx.Hash())
				fmt.Println("txHash=", txHash.String())
				ethTx, err := rpctypes.RawTxToEthTx(w.clientCtx, data.Tx)
				if err != nil {
					w.logger.Debug("invalid tx", "error", err)
					continue
				}

				tx, err := rpctypes.NewTransaction(ethTx, txHash, common.Hash{}, uint64(data.Height), uint64(data.Index))
				if err != nil {
					w.logger.Error("invalid tx", "error", err)
					continue
				}
				txBytes, err := json.Marshal(tx)
				if err != nil {
					w.logger.Error("failed to marshal result to JSON", "error", err)
					continue
				}
				fmt.Println(string(txBytes))

				go func(hash, tx []byte) {
					err = w.sender.Send(hash, tx)
					if err != nil {
						debug.PrintStack()
						w.logger.Error("failed to send pendingtx", "hash", txHash.String(), "error", err)
					}
				}(txHash.Bytes(), txBytes)
			}
		}
	}(sub.Event(), sub.Err())
}
