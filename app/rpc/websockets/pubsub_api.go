package websockets

import (
	"fmt"
	"sync"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"

	rpcfilters "github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// PubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
type PubSubAPI struct {
	clientCtx context.CLIContext
	events    *rpcfilters.EventSystem
	filtersMu *sync.RWMutex
	filters   map[rpc.ID]*wsSubscription
	logger    log.Logger
}

// NewAPI creates an instance of the ethereum PubSub API.
func NewAPI(clientCtx context.CLIContext, log log.Logger) *PubSubAPI {
	return &PubSubAPI{
		clientCtx: clientCtx,
		events:    rpcfilters.NewEventSystem(clientCtx.Client),
		filtersMu: new(sync.RWMutex),
		filters:   make(map[rpc.ID]*wsSubscription),
		logger:    log.With("module", "websocket-client"),
	}
}

func (api *PubSubAPI) subscribe(conn *wsConn, params []interface{}) (rpc.ID, error) {
	method, ok := params[0].(string)
	if !ok {
		return "0", fmt.Errorf("invalid parameters")
	}

	switch method {
	case "newHeads":
		// TODO: handle extra params
		return api.subscribeNewHeads(conn)
	case "logs":
		if len(params) > 1 {
			return api.subscribeLogs(conn, params[1])
		}

		return api.subscribeLogs(conn, nil)
	case "newPendingTransactions":
		return api.subscribePendingTransactions(conn)
	case "syncing":
		return api.subscribeSyncing(conn)
	default:
		return "0", fmt.Errorf("unsupported method %s", method)
	}
}

func (api *PubSubAPI) unsubscribe(id rpc.ID) bool {
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

func (api *PubSubAPI) subscribeNewHeads(conn *wsConn) (rpc.ID, error) {
	sub, _, err := api.events.SubscribeNewHeads()
	if err != nil {
		return "", fmt.Errorf("error creating block filter: %s", err.Error())
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	go func(headersCh <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case event := <-headersCh:
				data, ok := event.Data.(tmtypes.EventDataNewBlockHeader)
				if !ok {
					api.logger.Error(fmt.Sprintf("invalid data type %T, expected EventDataTx", event.Data), "ID", sub.ID())
					continue
				}
				headerWithBlockHash, err := rpctypes.EthHeaderWithBlockHashFromTendermint(&data.Header)
				if err != nil {
					api.logger.Error("failed to get header with block hash", "error", err)
					continue
				}

				api.filtersMu.RLock()
				if f, found := api.filters[sub.ID()]; found {
					// write to ws conn
					res := &SubscriptionNotification{
						Jsonrpc: "2.0",
						Method:  "eth_subscription",
						Params: &SubscriptionResult{
							Subscription: sub.ID(),
							Result:       headerWithBlockHash,
						},
					}

					err = f.conn.WriteJSON(res)
					if err != nil {
						api.logger.Error("failed to write header", "ID", sub.ID(), "blocknumber", headerWithBlockHash.Number, "error", err)
					} else {
						api.logger.Debug("successfully write header", "ID", sub.ID(), "blocknumber", headerWithBlockHash.Number)
					}
				}
				api.filtersMu.RUnlock()

				if err != nil {
					api.unsubscribe(sub.ID())
				}
			case err := <-errCh:
				if err != nil {
					api.unsubscribe(sub.ID())
					api.logger.Error("websocket recv error, close the conn", "ID", sub.ID(), "error", err)
				}
				return
			case <-unsubscribed:
				api.logger.Debug("NewHeads channel is closed", "ID", sub.ID())
				return
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}

func (api *PubSubAPI) subscribeLogs(conn *wsConn, extra interface{}) (rpc.ID, error) {
	crit := filters.FilterCriteria{}

	if extra != nil {
		params, ok := extra.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid criteria")
		}

		if params["address"] != nil {
			address, ok := params["address"].(string)
			addresses, sok := params["address"].([]interface{})
			if !ok && !sok {
				return "", fmt.Errorf("invalid address; must be address or array of addresses")
			}

			if ok {
				if !common.IsHexAddress(address) {
					return "", fmt.Errorf("invalid address")
				}
				crit.Addresses = []common.Address{common.HexToAddress(address)}
			} else if sok {
				crit.Addresses = []common.Address{}
				for _, addr := range addresses {
					address, ok := addr.(string)
					if !ok || !common.IsHexAddress(address) {
						return "", fmt.Errorf("invalid address")
					}

					crit.Addresses = append(crit.Addresses, common.HexToAddress(address))
				}
			}
		}

		if params["topics"] != nil {
			topics, ok := params["topics"].([]interface{})
			if !ok {
				return "", fmt.Errorf("invalid topics")
			}

			topicFilterLists, err := resolveTopicList(topics)
			if err != nil {
				return "", fmt.Errorf("invalid topics")
			}
			crit.Topics = topicFilterLists
		}
	}

	sub, _, err := api.events.SubscribeLogs(crit)
	if err != nil {
		return rpc.ID(""), err
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	go func(ch <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case event := <-ch:
				go func(event coretypes.ResultEvent) {
					dataTx, ok := event.Data.(tmtypes.EventDataTx)
					if !ok {
						api.logger.Error(fmt.Sprintf("invalid event data %T, expected EventDataTx", event.Data))
						return
					}

					var resultData evmtypes.ResultData
					resultData, err = evmtypes.DecodeResultData(dataTx.TxResult.Result.Data)
					if err != nil {
						api.logger.Error("failed to decode result data", "error", err)
						return
					}

					logs := rpcfilters.FilterLogs(resultData.Logs, crit.FromBlock, crit.ToBlock, crit.Addresses, crit.Topics)
					if len(logs) == 0 {
						api.logger.Debug("no matched logs", "ID", sub.ID(), "txhash", resultData.TxHash)
						return
					}

					api.filtersMu.RLock()
					if f, found := api.filters[sub.ID()]; found {
						// write to ws conn
						res := &SubscriptionNotification{
							Jsonrpc: "2.0",
							Method:  "eth_subscription",
							Params: &SubscriptionResult{
								Subscription: sub.ID(),
							},
						}
						for _, singleLog := range logs {
							res.Params.Result = singleLog
							err = f.conn.WriteJSON(res)
							if err != nil {
								api.logger.Error("failed to write log", "ID", sub.ID(), "height", singleLog.BlockNumber, "txhash", singleLog.TxHash, "error", err)
								break
							}
							api.logger.Debug("successfully write log", "ID", sub.ID(), "height", singleLog.BlockNumber, "txhash", singleLog.TxHash)
						}
					}
					api.filtersMu.RUnlock()

					if err != nil {
						api.unsubscribe(sub.ID())
					}
				}(event)
			case err := <-errCh:
				if err != nil {
					api.unsubscribe(sub.ID())
					api.logger.Error("websocket recv error, close the conn", "ID", sub.ID(), "error", err)
				}
				return
			case <-unsubscribed:
				api.logger.Debug("Logs channel is closed", "ID", sub.ID())
				return
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}

func resolveTopicList(params []interface{}) ([][]common.Hash, error) {
	topicFilterLists := make([][]common.Hash, len(params))
	for i, param := range params { // eg: ["0xddf252......f523b3ef", null, ["0x000000......32fea9e4", "0x000000......ab14dc5d"]]
		if param == nil {
			// 1.1 if the topic is null
			topicFilterLists[i] = nil
		} else {
			// 2.1 judge if the param is the type of string or not
			topicStr, ok := param.(string)
			// 2.1 judge if the param is the type of string slice or not
			topicSlices, sok := param.([]interface{})
			if !ok && !sok {
				// if both judgement are false, return invalid topics
				return topicFilterLists, fmt.Errorf("invalid topics")
			}

			if ok {
				// 2.2 This is string
				// 2.3 judge the topic is a valid hex hash or not
				if !IsHexHash(topicStr) {
					return topicFilterLists, fmt.Errorf("invalid topics")
				}
				// 2.4 add this topic to topic-hash-lists
				topicHash := common.HexToHash(topicStr)
				topicFilterLists[i] = []common.Hash{topicHash}
			} else if sok {
				// 2.2 This is slice of string
				topicHashes := make([]common.Hash, len(topicSlices))
				for n, topicStr := range topicSlices {
					//2.3 judge every topic
					topicHash, ok := topicStr.(string)
					if !ok || !IsHexHash(topicHash) {
						return topicFilterLists, fmt.Errorf("invalid topics")
					}
					topicHashes[n] = common.HexToHash(topicHash)
				}
				// 2.4 add this topic slice to topic-hash-lists
				topicFilterLists[i] = topicHashes
			}
		}
	}
	return topicFilterLists, nil
}

func IsHexHash(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*common.HashLength && isHex(s)
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

func (api *PubSubAPI) subscribePendingTransactions(conn *wsConn) (rpc.ID, error) {
	sub, _, err := api.events.SubscribePendingTxs()
	if err != nil {
		return "", fmt.Errorf("error creating block filter: %s", err.Error())
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	go func(txsCh <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case ev := <-txsCh:
				data, ok := ev.Data.(tmtypes.EventDataTx)
				if !ok {
					api.logger.Error(fmt.Sprintf("invalid data type %T, expected EventDataTx", ev.Data), "ID", sub.ID())
					continue
				}
				txHash := common.BytesToHash(data.Tx.Hash())

				api.filtersMu.RLock()
				if f, found := api.filters[sub.ID()]; found {
					// write to ws conn
					res := &SubscriptionNotification{
						Jsonrpc: "2.0",
						Method:  "eth_subscription",
						Params: &SubscriptionResult{
							Subscription: sub.ID(),
							Result:       txHash,
						},
					}

					err = f.conn.WriteJSON(res)
					if err != nil {
						api.logger.Error("failed to write pending tx", "ID", sub.ID(), "error", err)
					} else {
						api.logger.Debug("successfully write pending tx", "ID", sub.ID(), "txhash", txHash)
					}
				}
				api.filtersMu.RUnlock()

				if err != nil {
					api.unsubscribe(sub.ID())
				}
			case err := <-errCh:
				if err != nil {
					api.unsubscribe(sub.ID())
					api.logger.Error("websocket recv error, close the conn", "ID", sub.ID(), "error", err)
				}
				return
			case <-unsubscribed:
				api.logger.Debug("PendingTransactions channel is closed", "ID", sub.ID())
				return
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}

func (api *PubSubAPI) subscribeSyncing(conn *wsConn) (rpc.ID, error) {
	sub, _, err := api.events.SubscribeNewHeads()
	if err != nil {
		return "", fmt.Errorf("error creating block filter: %s", err.Error())
	}

	unsubscribed := make(chan struct{})
	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:          sub,
		conn:         conn,
		unsubscribed: unsubscribed,
	}
	api.filtersMu.Unlock()

	status, err := api.clientCtx.Client.Status()
	if err != nil {
		return "", fmt.Errorf("error get sync status: %s", err.Error())
	}
	startingBlock := hexutil.Uint64(status.SyncInfo.EarliestBlockHeight)
	highestBlock := hexutil.Uint64(0)

	var result interface{}

	go func(headersCh <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case <-headersCh:

				newStatus, err := api.clientCtx.Client.Status()
				if err != nil {
					api.logger.Error(fmt.Sprintf("error get sync status: %s", err.Error()))
					continue
				}

				if !newStatus.SyncInfo.CatchingUp {
					result = false
				} else {
					result = map[string]interface{}{
						"startingBlock": startingBlock,
						"currentBlock":  hexutil.Uint64(newStatus.SyncInfo.LatestBlockHeight),
						"highestBlock":  highestBlock,
					}
				}

				api.filtersMu.RLock()
				if f, found := api.filters[sub.ID()]; found {
					// write to ws conn
					res := &SubscriptionNotification{
						Jsonrpc: "2.0",
						Method:  "eth_subscription",
						Params: &SubscriptionResult{
							Subscription: sub.ID(),
							Result:       result,
						},
					}

					err = f.conn.WriteJSON(res)
					if err != nil {
						api.logger.Error("failed to write syncing status", "ID", sub.ID(), "error", err)
					} else {
						api.logger.Debug("successfully write syncing status", "ID", sub.ID())
					}
				}
				api.filtersMu.RUnlock()

				if err != nil {
					api.unsubscribe(sub.ID())
				}

			case err := <-errCh:
				if err != nil {
					api.unsubscribe(sub.ID())
					api.logger.Error("websocket recv error, close the conn", "ID", sub.ID(), "error", err)
				}
				return
			case <-unsubscribed:
				api.logger.Debug("Syncing channel is closed", "ID", sub.ID())
				return
			}
		}
	}(sub.Event(), sub.Err())

	return sub.ID(), nil
}
