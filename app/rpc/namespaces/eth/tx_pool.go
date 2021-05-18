package eth

import (
	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"sync"
)

type TxPool struct {
	addressTxsPool	map[common.Address][]*evmtypes.MsgEthereumTx // All currently processable transactions
	txChan			chan *ChanData
	mu				sync.RWMutex
}

// data struct for transmitting of chan to txPool
type ChanData struct {
	address      *common.Address
	tx           *evmtypes.MsgEthereumTx
	currentNonce *hexutil.Uint64
}

func NewTxPool() *TxPool {
	pool := &TxPool{
		addressTxsPool: make(map[common.Address][]*evmtypes.MsgEthereumTx),
		txChan:         make(chan *ChanData),
	}

	return pool
}

func (pool *TxPool) SetData(chanData *ChanData) {
	pool.txChan <- chanData
}

func (pool *TxPool) DoBroadcastTx(clientCtx clientcontext.CLIContext) {
	for {
		select {
		case chanData := <-pool.txChan:
			address := *(chanData.address)
			txNonce := chanData.tx.Data.AccountNonce
			currentNonce := *(chanData.currentNonce)
			needInsert := true
			if hexutil.Uint64(txNonce) == currentNonce {
				needInsert = false
				// do broadcast
				txEncoder := authclient.GetTxEncoder(clientCtx.Codec)
				txBytes, err := txEncoder(chanData.tx)
				if err != nil {
					break
				}
				_, err = clientCtx.BroadcastTx(txBytes)
				if err != nil {
					break
				}

				// update currentNonce
				currentNonce++
			}
			// map need lock
			pool.mu.RLock()
			txsLen := len(pool.addressTxsPool[address])
			if needInsert {
				index := 0
				for index < txsLen {

					/*
					// the tx nonce has in txPool, drop duplicate tx
					if txNonce == pool.addressTxsPool[address][index].Data.AccountNonce {
						return
					}
					*/

					// find the index to insert
					if txNonce < pool.addressTxsPool[address][index].Data.AccountNonce {
						break
					}
					index++
				}

				// update txPool
				if index >= txsLen {
					pool.mu.RUnlock()
					pool.mu.Lock()
					pool.addressTxsPool[address] =append(pool.addressTxsPool[address], chanData.tx)
					pool.mu.Unlock()
				} else {
					tmpTx := make([]*evmtypes.MsgEthereumTx, len(pool.addressTxsPool[address][index:]))
					copy(tmpTx, pool.addressTxsPool[address][index:])

					pool.mu.RUnlock()
					pool.mu.Lock()
					pool.addressTxsPool[address] =
						append(append(pool.addressTxsPool[address][:index], chanData.tx), tmpTx...)
					pool.mu.Unlock()
				}

			} else {
				var err error
				i := 0
				for i < txsLen {
					if hexutil.Uint64(pool.addressTxsPool[address][i].Data.AccountNonce) != currentNonce {
						break
					}
					// do broadcast
					txEncoder := authclient.GetTxEncoder(clientCtx.Codec)
					var txBytes []byte
					txBytes, err = txEncoder(pool.addressTxsPool[address][i])
					if err != nil {
						break
					}
					_, err = clientCtx.BroadcastTx(txBytes)
					if err != nil {
						break
					}

					// update currentNonce
					currentNonce++
					i++
				}

				// update txPool
				pool.mu.RUnlock()
				if err == nil && i != 0 {
					pool.mu.Lock()
					pool.addressTxsPool[address] = pool.addressTxsPool[address][i:]
					pool.mu.Unlock()
				}
			}

		} // end select
	}
}
