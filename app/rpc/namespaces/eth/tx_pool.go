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
	addressTxsPool map[common.Address][]*evmtypes.MsgEthereumTx // All currently processable transactions
	txChan         chan *ChanData
	mu             sync.Mutex
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
				if err := pool.doBroadcast(clientCtx, chanData.tx); err != nil {
					break
				}
				currentNonce++
			}
			// map need lock
			pool.mu.Lock()
			if needInsert {
				pool.doInsert(txNonce, address, chanData.tx)
			} else {
				pool.doNoInsert(clientCtx, currentNonce, address, chanData.tx)
			}
			pool.mu.Unlock()

		} // end select
	}
}

func (pool *TxPool) doBroadcast(clientCtx clientcontext.CLIContext, tx *evmtypes.MsgEthereumTx) error {
	txEncoder := authclient.GetTxEncoder(clientCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return err
	}
	_, err = clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}
	return nil
}

func (pool *TxPool) updateTxPool(index int, address common.Address, tx *evmtypes.MsgEthereumTx) {
	txsLen := len(pool.addressTxsPool[address])
	if index >= txsLen {
		pool.addressTxsPool[address] = append(pool.addressTxsPool[address], tx)
	} else {
		tmpTx := make([]*evmtypes.MsgEthereumTx, len(pool.addressTxsPool[address][index:]))
		copy(tmpTx, pool.addressTxsPool[address][index:])

		pool.addressTxsPool[address] =
			append(append(pool.addressTxsPool[address][:index], tx), tmpTx...)
	}
}

func (pool *TxPool) doInsert(txNonce uint64, address common.Address, tx *evmtypes.MsgEthereumTx) {
	index := 0
	txsLen := len(pool.addressTxsPool[address])
	for index < txsLen {
		// the tx nonce has in txPool, drop duplicate tx
		if txNonce == pool.addressTxsPool[address][index].Data.AccountNonce {
			return
		}
		// find the index to insert
		if txNonce < pool.addressTxsPool[address][index].Data.AccountNonce {
			break
		}
		index++
	}

	// update txPool
	pool.updateTxPool(index, address, tx)
}

func (pool *TxPool) doNoInsert(clientCtx clientcontext.CLIContext, currentNonce hexutil.Uint64,
	address common.Address, tx *evmtypes.MsgEthereumTx) {
	i := 0
	txsLen := len(pool.addressTxsPool[address])
	for i < txsLen {
		if hexutil.Uint64(pool.addressTxsPool[address][i].Data.AccountNonce) != currentNonce {
			break
		}
		// do broadcast
		if err := pool.doBroadcast(clientCtx, tx); err != nil {
			return
		}
		// update currentNonce
		currentNonce++
		i++
	}

	// update txPool
	if i != 0 {
		pool.addressTxsPool[address] = pool.addressTxsPool[address][i:]
	}
}