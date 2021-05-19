package eth

import (
	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"sync"
)

type TxPool struct {
	addressTxsPool map[common.Address][]*evmtypes.MsgEthereumTx // All currently processable transactions
	mu             sync.Mutex
}

func NewTxPool() *TxPool {
	pool := &TxPool{
		addressTxsPool: make(map[common.Address][]*evmtypes.MsgEthereumTx),
	}

	return pool
}

func (pool *TxPool) CacheAndBroadcastTx(clientCtx clientcontext.CLIContext, address common.Address,
	currentNonce uint64, tx *evmtypes.MsgEthereumTx) error {
	needInsert := true
	txNonce := tx.Data.AccountNonce
	if txNonce == currentNonce {
		needInsert = false
		// do broadcast
		if err := pool.doBroadcast(clientCtx, tx); err != nil {
			return err
		}
		currentNonce++
	}
	// map need lock
	pool.mu.Lock()
	if needInsert {
		pool.insertTx(txNonce, address, tx)
	}
	err := pool.continueBroadcast(clientCtx, currentNonce, address)
	pool.mu.Unlock()

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

// insert the tx into the txPool
func (pool *TxPool) insertTx(txNonce uint64, address common.Address, tx *evmtypes.MsgEthereumTx) {
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

// iterate through the txPool map, check if need to continue broadcast tx and do it
func (pool *TxPool) continueBroadcast(clientCtx clientcontext.CLIContext, currentNonce uint64, address common.Address) error {
	i := 0
	txsLen := len(pool.addressTxsPool[address])
	for i < txsLen {
		if pool.addressTxsPool[address][i].Data.AccountNonce != currentNonce {
			break
		}
		// do broadcast
		if err := pool.doBroadcast(clientCtx, pool.addressTxsPool[address][i]); err != nil {
			return err
		}
		// update currentNonce
		currentNonce++
		i++
	}

	// update txPool
	if i != 0 {
		pool.addressTxsPool[address] = pool.addressTxsPool[address][i:]
	}

	return nil
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