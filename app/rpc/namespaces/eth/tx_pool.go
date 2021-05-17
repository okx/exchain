package eth

import (
	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

type TxPool struct {
	addressTxsPool 		map[common.Address][]*evmtypes.MsgEthereumTx   // All currently processable transactions
	txChan     			chan *ChanData
}

// data struct for transmitting of chan to txPool
type ChanData struct {
	address			*common.Address
	tx				*evmtypes.MsgEthereumTx
	currentNonce	*hexutil.Uint64
}

func NewTxPool() *TxPool {
	pool := &TxPool{
		addressTxsPool: 	make(map[common.Address][]*evmtypes.MsgEthereumTx),
		txChan:				make(chan *ChanData),
	}

	return pool
}

func (pool *TxPool) SetData(chanData *ChanData) {
	pool.txChan <- chanData
}

func (pool *TxPool) DoBroadcastTx(clientCtx clientcontext.CLIContext) {
	for {
		select {
		case chanData := <- pool.txChan:
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
			if needInsert {
				index := 0
				for index< len(pool.addressTxsPool[address]) {
					/*
						// the tx nonce has in txPool, drop duplicate tx
						if txNonce == api.txPool.addressTxsPool[address][index].Data.AccountNonce {
							return
						}
					*/
					// find the index to insert
					if txNonce < pool.addressTxsPool[address][index].Data.AccountNonce {
						break
					}
					index++
				}
				tmpTx := make([]*evmtypes.MsgEthereumTx, len(pool.addressTxsPool[address][index:]))
				copy(tmpTx, pool.addressTxsPool[address][index:])
				pool.addressTxsPool[address] =
					append(append(pool.addressTxsPool[address][:index], chanData.tx), tmpTx...)
			} else {
				for i:=0; i< len(pool.addressTxsPool[address]); i++ {
					if hexutil.Uint64(pool.addressTxsPool[address][i].Data.AccountNonce) == currentNonce {
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
						continue
					}

					// update txPool
					pool.addressTxsPool[address]= pool.addressTxsPool[address][i:]
					break
				}
			}

		} // end select
	}
}
