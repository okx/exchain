package watcher

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okx/okbchain/x/evm/types"
)

type evmTx struct {
	msgEvmTx  *types.MsgEthereumTx
	txHash    ethcmn.Hash
	blockHash ethcmn.Hash
	height    uint64
	index     uint64
}

func NewEvmTx(msgEvmTx *types.MsgEthereumTx, txHash ethcmn.Hash, blockHash ethcmn.Hash, height, index uint64) *evmTx {
	return &evmTx{
		msgEvmTx:  msgEvmTx,
		txHash:    txHash,
		blockHash: blockHash,
		height:    height,
		index:     index,
	}
}

func (etx *evmTx) GetTxHash() ethcmn.Hash {
	if etx == nil {
		return ethcmn.Hash{}
	}
	return etx.txHash
}

func (etx *evmTx) GetTransaction() *Transaction {
	if etx == nil || etx.msgEvmTx == nil {
		return nil
	}
	ethTx, e := NewTransaction(etx.msgEvmTx, etx.txHash, etx.blockHash, etx.height, etx.index)
	if e != nil {
		return nil
	}
	return ethTx
}

func (etx *evmTx) GetFailedReceipts(cumulativeGas, gasUsed uint64) *TransactionReceipt {
	if etx == nil {
		return nil
	}
	tr := newTransactionReceipt(TransactionFailed, etx.msgEvmTx, etx.txHash, etx.blockHash, etx.index, etx.height, &types.ResultData{}, cumulativeGas, gasUsed)
	return &tr
}

func (etx *evmTx) GetIndex() uint64 {
	return etx.index
}

type MsgEthTx struct {
	*Transaction
	Key []byte
}

func (m MsgEthTx) GetType() uint32 {
	return TypeOthers
}

func (m MsgEthTx) GetKey() []byte {
	return append(prefixTx, m.Key...)
}

func (etx *evmTx) GetTxWatchMessage() WatchMessage {
	if etx == nil || etx.msgEvmTx == nil {
		return nil
	}

	return newMsgEthTx(etx.msgEvmTx, etx.txHash, etx.blockHash, etx.height, etx.index)
}

func newTransaction(tx *types.MsgEthereumTx, txHash, blockHash ethcmn.Hash, blockNumber, index uint64) *Transaction {
	return &Transaction{
		Hash:              txHash,
		tx:                tx,
		originBlockHash:   &blockHash,
		originBlockNumber: blockNumber,
		originIndex:       index,
	}
}

func newMsgEthTx(tx *types.MsgEthereumTx, txHash, blockHash ethcmn.Hash, height, index uint64) *MsgEthTx {
	ethTx := newTransaction(tx, txHash, blockHash, height, index)

	return &MsgEthTx{
		Transaction: ethTx,
		Key:         txHash.Bytes(),
	}
}
