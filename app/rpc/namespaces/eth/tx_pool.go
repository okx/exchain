package eth

import (
	"errors"
	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"sync"
)

const FlagEnableTxPool = "enable-tx-pool"

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
	if tx.Data.AccountNonce < currentNonce {
		return errors.New("AccountNonce of tx is less than currentNonce in memPool")
	}

	// map need lock
	pool.mu.Lock()
	if err := pool.insertTx(address, tx); err != nil {
		pool.mu.Unlock()
		return err
	}

	if err := pool.continueBroadcast(clientCtx, currentNonce, address); err != nil {
		pool.mu.Unlock()
		return err
	}
	pool.mu.Unlock()

	return nil
}

func (pool *TxPool) updateTxPool(index int, address common.Address, tx *evmtypes.MsgEthereumTx) {
	if index >= len(pool.addressTxsPool[address]) {
		pool.addressTxsPool[address] = append(pool.addressTxsPool[address], tx)
	} else {
		tmpTx := make([]*evmtypes.MsgEthereumTx, len(pool.addressTxsPool[address][index:]))
		copy(tmpTx, pool.addressTxsPool[address][index:])

		pool.addressTxsPool[address] =
			append(append(pool.addressTxsPool[address][:index], tx), tmpTx...)
	}
}

// insert the tx into the txPool
func (pool *TxPool) insertTx(address common.Address, tx *evmtypes.MsgEthereumTx) error {
	index := 0
	for index < len(pool.addressTxsPool[address]) {
		// the tx nonce has in txPool, drop duplicate tx
		if tx.Data.AccountNonce == pool.addressTxsPool[address][index].Data.AccountNonce {
			return errors.New("duplicate tx, this AccountNonce of tx has been send")
		}
		// find the index to insert
		if tx.Data.AccountNonce < pool.addressTxsPool[address][index].Data.AccountNonce {
			break
		}
		index++
	}

	// update txPool
	pool.updateTxPool(index, address, tx)

	return nil
}

// iterate through the txPool map, check if need to continue broadcast tx and do it
func (pool *TxPool) continueBroadcast(clientCtx clientcontext.CLIContext, currentNonce uint64, address common.Address) error {
	i := 0
	for i < len(pool.addressTxsPool[address]) {
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
