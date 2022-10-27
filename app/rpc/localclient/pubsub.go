package localclient

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	rpcfilters "github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	rpcclient "github.com/okex/exchain/libs/tendermint/rpc/client"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// PubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
type PubSubAPI struct {
	clientCtx context.CLIContext
	events    *rpcfilters.EventSystem
	filtersMu *sync.RWMutex
	filters   map[rpc.ID]*localSubscription
	logger    log.Logger
}

// NewAPI creates an instance of the ethereum PubSub API.
func NewAPI(client rpcclient.Client, log log.Logger) *PubSubAPI {
	return &PubSubAPI{
		events:    rpcfilters.NewEventSystem(client),
		filtersMu: new(sync.RWMutex),
		filters:   make(map[rpc.ID]*localSubscription),
		logger:    log.With("module", "local-client-pubsub"),
	}
}

func (api *PubSubAPI) Unsubscribe(id rpc.ID) bool {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	if api.filters[id] == nil {
		api.logger.Debug("client doesn't exist in filters", "ID", id)
		return false
	}
	if api.filters[id].sub != nil {
		api.filters[id].sub.Unsubscribe(api.events)
	}
	close(api.filters[id].unsubscribed)
	delete(api.filters, id)
	api.logger.Debug("close client channel & delete client from filters", "ID", id)
	return true
}

func (api *PubSubAPI) SubscribeLogs(conn chan<- *ethtypes.Log, query ethereum.FilterQuery) (rpc.ID, error) {
	crit := filters.FilterCriteria{
		Addresses: query.Addresses,
		Topics:    query.Topics,
	}

	sub, _, err := api.events.SubscribeLogsBatch(crit)
	if err != nil {
		return rpc.ID(""), err
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &localSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	go func(ch <-chan coretypes.ResultEvent, errCh <-chan error) {
		quit := false
		for {
			select {
			case event := <-ch:
				go func(event coretypes.ResultEvent) {
					//batch receive txResult
					txs, ok := event.Data.(tmtypes.EventDataTxs)
					if !ok {
						api.logger.Error(fmt.Sprintf("invalid event data %T, expected EventDataTxs", event.Data))
						return
					}

					for _, txResult := range txs.Results {
						if quit {
							return
						}

						//check evm type event
						if !evmtypes.IsEvmEvent(txResult) {
							continue
						}

						//decode txResult data
						var resultData evmtypes.ResultData
						resultData, err = evmtypes.DecodeResultData(txResult.Data)
						if err != nil {
							api.logger.Error("failed to decode result data", "error", err)
							return
						}

						//filter logs
						logs := rpcfilters.FilterLogs(resultData.Logs, crit.FromBlock, crit.ToBlock, crit.Addresses, crit.Topics)
						if len(logs) == 0 {
							continue
						}

						//write log to client by each tx
						api.filtersMu.RLock()
						if f, found := api.filters[sub.ID()]; found {
							for _, singleLog := range logs {
								f.conn <- singleLog
								api.logger.Info("successfully write log", "ID", sub.ID(), "height", singleLog.BlockNumber, "txHash", singleLog.TxHash)
							}
						}
						api.filtersMu.RUnlock()

						if err != nil {
							//unsubscribe and quit current routine
							api.Unsubscribe(sub.ID())
							return
						}
					}
				}(event)
			case err := <-errCh:
				quit = true
				if err != nil {
					api.Unsubscribe(sub.ID())
					api.logger.Error("websocket recv error, close the conn", "ID", sub.ID(), "error", err)
				}
				return
			case <-unsubscribed:
				quit = true
				api.logger.Debug("Logs channel is closed", "ID", sub.ID())
				return
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}

type localSubscription struct {
	sub          *rpcfilters.Subscription
	unsubscribed chan struct{} // closed when unsubscribing
	conn         chan<- *ethtypes.Log
}
