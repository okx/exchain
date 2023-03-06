package pendingtx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	rpcfilters "github.com/okx/okbchain/app/rpc/namespaces/eth/filters"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	coretypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

type Watcher struct {
	clientCtx context.CLIContext
	events    *rpcfilters.EventSystem
	logger    log.Logger

	sender Sender
}

type Sender interface {
	SendPending(hash []byte, tx *PendingTx) error
	SendRmPending(hash []byte, tx *RmPendingTx) error
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

	rmPendingSub, _, err := w.events.SubscribeRmPendingTx()
	if err != nil {
		w.logger.Error("error creating block filter", "error", err.Error())
	}

	go func(pendingCh <-chan coretypes.ResultEvent, rmPendingdCh <-chan coretypes.ResultEvent) {
		for {
			select {
			case re := <-pendingCh:
				data, ok := re.Data.(tmtypes.EventDataTx)
				if !ok {
					w.logger.Error(fmt.Sprintf("invalid pending tx data type %T, expected EventDataTx", re.Data))
					continue
				}
				txHash := common.BytesToHash(data.Tx.Hash())
				w.logger.Debug("receive pending tx", "txHash=", txHash.String())

				tx, err := evmtypes.TxDecoder(w.clientCtx.Codec)(data.Tx, data.Height)
				if err != nil {
					w.logger.Error("failed to decode raw tx", "hash", txHash.String(), "error", err)
					continue
				}

				var input string
				var value *big.Int
				var to *common.Address
				ethTx, ok := tx.(*evmtypes.MsgEthereumTx)
				if ok {
					input = hexutil.Bytes(ethTx.Data.Payload).String()
					value = ethTx.Data.Amount
					to = ethTx.Data.Recipient
				} else {
					b, err := w.clientCtx.Codec.MarshalJSON(tx)
					if err != nil {
						w.logger.Error("failed to Marshal tx", "hash", txHash.String(), "error", err)
						continue
					}
					input = string(b)
				}

				pendingTx := &PendingTx{
					From:     tx.GetFrom(),
					To:       to,
					Hash:     txHash,
					Nonce:    hexutil.Uint64(data.Nonce),
					Value:    (*hexutil.Big)(value),
					Gas:      hexutil.Uint64(tx.GetGas()),
					GasPrice: (*hexutil.Big)(tx.GetGasPrice()),
					Input:    input,
				}

				go func() {
					w.logger.Debug("push pending tx to MQ", "txHash=", pendingTx.Hash.String())
					err = w.sender.SendPending(pendingTx.Hash.Bytes(), pendingTx)
					if err != nil {
						w.logger.Error("failed to send pending tx", "hash", pendingTx.Hash.String(), "error", err)
					}
				}()
			case re := <-rmPendingdCh:
				data, ok := re.Data.(tmtypes.EventDataRmPendingTx)
				if !ok {
					w.logger.Error(fmt.Sprintf("invalid rm pending tx data type %T, expected EventDataTx", re.Data))
					continue
				}
				txHash := common.BytesToHash(data.Hash).String()
				go func() {
					w.logger.Debug("push rm pending tx to MQ", "txHash=", txHash)
					err = w.sender.SendRmPending(data.Hash, &RmPendingTx{
						From:   data.From,
						Hash:   txHash,
						Nonce:  hexutil.Uint64(data.Nonce).String(),
						Delete: true,
						Reason: int(data.Reason),
					})
					if err != nil {
						w.logger.Error("failed to send rm pending tx", "hash", txHash, "error", err)
					}
				}()
			}
		}
	}(pendingSub.Event(), rmPendingSub.Event())
}
