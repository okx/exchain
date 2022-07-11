package watcher

import (
	"encoding/json"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

type stdWatchTx struct {
	stdTx     *authtypes.StdTx
	txHash    ethcmn.Hash
	blockHash ethcmn.Hash
	height    uint64
}

func NewStdWatchTx(tx sdk.Tx, txHash ethcmn.Hash, blockHash ethcmn.Hash, height uint64) *stdWatchTx {
	realTx, ok := tx.(*authtypes.StdTx)
	if !ok {
		return nil
	}
	return &stdWatchTx{
		stdTx:     realTx,
		txHash:    txHash,
		blockHash: blockHash,
		height:    height,
	}
}

func (tx *stdWatchTx) GetTxHash() ethcmn.Hash {
	if tx == nil {
		return ethcmn.Hash{}
	}
	return tx.txHash
}

func (tx *stdWatchTx) GetTransaction() *Transaction {
	if tx == nil {
		return nil
	}

	return &Transaction{
		BlockHash:   &tx.blockHash,
		BlockNumber: (*hexutil.Big)(new(big.Int).SetUint64(tx.height)),
		Hash:        tx.txHash,
	}
}

func (tx *stdWatchTx) GetFailedReceipts(cumulativeGas, gasUsed uint64) *TransactionReceipt {
	return nil
}

func (tx *stdWatchTx) GetIndex() uint64 {
	return 0
}

func (tx *stdWatchTx) GetTxWatchMessage() WatchMessage {
	if tx == nil || tx.stdTx == nil {
		return nil
	}

	return newMsgStdTx(tx.txHash, tx.blockHash, tx.height)
}

type MsgStdTx struct {
	tr  string
	Key []byte
}

func (m MsgStdTx) GetType() uint32 {
	return TypeOthers
}

func (m MsgStdTx) GetKey() []byte {
	return append(prefixTx, m.Key...)
}

func (m MsgStdTx) GetValue() string {
	return m.tr

}

func newMsgStdTx(txHash, blockHash ethcmn.Hash, height uint64) *MsgStdTx {
	stdTx := &Transaction{
		BlockHash:   &blockHash,
		BlockNumber: (*hexutil.Big)(new(big.Int).SetUint64(height)),
		Hash:        txHash,
	}

	js, err := json.Marshal(stdTx)
	if err != nil {
		return nil
	}
	return &MsgStdTx{
		tr:  string(js),
		Key: txHash.Bytes(),
	}
}
